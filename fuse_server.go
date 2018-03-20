package gofuse

/*
#cgo CFLAGS: -I${SRCDIR}/c

#include <errno.h>
#include <sys/stat.h>
#include "fuse_kernel_7_26.h"

#define SIZEOF_FUSE_IN_HEADER sizeof(struct fuse_in_header)
#define SIZEOF_FUSE_OUT_HEADER sizeof(struct fuse_out_header)
#define SIZEOF_FUSE_INIT_OUT sizeof(struct fuse_init_out)
#define SIZEOF_FUSE_INIT_IN sizeof(struct fuse_init_in)
#define SIZEOF_FUSE_ATTR_OUT sizeof(struct fuse_attr_out)
*/
import "C"
import (
	"C"
	"errors"
	"os"
	"sync"
)
import (
	"fmt"
	"unsafe"
)

func checkDir(dir string) (err error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return
	} else if !stat.IsDir() {
		err = fmt.Errorf("gofuse: mount point %s is not a directory", dir)
		return
	}
	return
}

func fuseInit(f *os.File) (err error) {
	buf := make([]byte, C.SIZEOF_FUSE_IN_HEADER+C.SIZEOF_FUSE_INIT_OUT)
	n, err := f.Read(buf)
	if err != nil {
		return
	} else if n != C.SIZEOF_FUSE_IN_HEADER+C.SIZEOF_FUSE_INIT_OUT {
		return errors.New("gofuse: error fuse request format")
	}

	header := (*C.struct_fuse_in_header)(unsafe.Pointer(&buf[0]))
	if header.opcode != C.FUSE_INIT {
		return fmt.Errorf(
			"gofuse: error fuse request opcode, expect: %d, got: %d",
			C.FUSE_INIT, header.opcode)
	}

	bodyRaw := buf[C.SIZEOF_FUSE_IN_HEADER:]
	body := (*C.struct_fuse_init_in)(unsafe.Pointer(&bodyRaw[0]))
	if body.major != _FUSE_KERNEL_VERSION ||
		body.minor < _FUSE_KERNEL_VERSION {
		return fmt.Errorf(
			"gofuse: error fuse kernel version, expect %d.%d, got: %d.%d",
			_FUSE_KERNEL_VERSION, _FUSE_KERNEL_MINOR_VERSION,
			body.major, body.minor)
	}

	replyRaw := make([]byte, C.SIZEOF_FUSE_OUT_HEADER+C.SIZEOF_FUSE_INIT_OUT)
	rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&replyRaw[0]))
	rbody := (*C.struct_fuse_init_out)(unsafe.Pointer(
		&replyRaw[C.SIZEOF_FUSE_OUT_HEADER]))

	rheader.len = C.uint32_t(len(replyRaw))
	rheader.error = 0
	rheader.unique = header.unique

	rbody.major = _FUSE_KERNEL_VERSION
	rbody.minor = _FUSE_KERNEL_MINOR_VERSION
	rbody.max_readahead = _MAX_BUFFER_SIZE
	rbody.flags = 0
	rbody.max_background = 0
	rbody.congestion_threshold = 0
	rbody.max_write = _MAX_BUFFER_SIZE
	rbody.time_gran = 0

	_, err = f.Write(replyRaw)
	return
}

type FuseServer struct {
	dir  string
	conf *MountConfig

	f   *os.File
	ops FuseOperations

	end     chan struct{}
	endLock *sync.Mutex
}

func NewFuseServer(dir string, conf *MountConfig, ops FuseOperations) (
	fs *FuseServer, err error) {
	if err = checkDir(dir); err != nil {
		return
	}

	f, err := mount(dir, conf)
	if err != nil {
		return
	}

	if err = fuseInit(f); err != nil {
		return
	}

	fs = &FuseServer{
		dir:  dir,
		conf: conf,

		f:   f,
		ops: ops,

		end:     make(chan struct{}),
		endLock: &sync.Mutex{},
	}
	return
}

func (fs *FuseServer) IsClosed() bool {
	select {
	case <-fs.end:
		return true
	default:
		return false
	}
}

func (fs *FuseServer) Close() error {
	fs.endLock.Lock()
	defer fs.endLock.Unlock()

	if fs.IsClosed() {
		return errors.New("gofuse: fuse server was closed")
	}
	close(fs.end)
	fs.f.Close()
	umount(fs.dir)

	return nil
}

func (fs *FuseServer) serv() {
	buf := make([]byte, _FUSE_MAX_BUFFER_SIZE)
	for {
		n, err := fs.f.Read(buf)
		if err != nil {
			fs.Close()
			return
		}
		fmt.Println(n)
	}
}
