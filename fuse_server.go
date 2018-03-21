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
	"syscall"
	"unsafe"
)

type interrupNotice struct {
	ch     chan struct{}
	unique uint64
	next   *interrupNotice
}

func newInterrupNotice() *interrupNotice {
	return &interrupNotice{
		ch:   make(chan struct{}),
		next: nil,
	}
}

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
	buf := make([]byte, C.SIZEOF_FUSE_IN_HEADER+C.SIZEOF_FUSE_INIT_IN)
	n, err := f.Read(buf)
	if err != nil {
		return
	} else if n != len(buf) {
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

	f    *os.File
	send chan []byte
	ops  FuseOperations

	end     chan struct{}
	endLock *sync.Mutex

	requests    map[uint64]*FuseRequestContext
	requestLock *sync.RWMutex

	intrN *interrupNotice
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

		f:    f,
		send: make(chan []byte),
		ops:  ops,

		end:     make(chan struct{}),
		endLock: &sync.Mutex{},

		requests:    make(map[uint64]*FuseRequestContext),
		requestLock: &sync.RWMutex{},

		intrN: newInterrupNotice(),
	}
	go fs.readLoop()
	go fs.sendLoop()
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

func (fs *FuseServer) readLoop() {
	buf := make([]byte, _FUSE_MAX_BUFFER_SIZE)
	for {
		n, err := fs.f.Read(buf)
		if err != nil {
			fs.Close()
			return
		}

		header := (*C.struct_fuse_in_header)(unsafe.Pointer(&buf[0]))
		switch header.opcode {
		case C.FUSE_INTERRUPT:
			reqIntr := (*C.struct_fuse_interrupt_in)(
				unsafe.Pointer(&buf[C.SIZEOF_FUSE_IN_HEADER]))
			fs.handlerFuseInterrupt(header, reqIntr)
		default:
			b := make([]byte, n, n)
			copy(b, buf[:n])
			go fs.handlerFuseMessage(b, fs.intrN)
		}
	}
}

func (fs *FuseServer) sendLoop() {
	for {
		select {
		case raw := <-fs.send:
			if len(raw) == 0 {
				break
			}
			if _, err := fs.f.Write(raw); err != nil {
				fs.Close()
				return
			}
		case <-fs.end:
			return
		}
	}
}

func (fs *FuseServer) handlerFuseInterrupt(
	header *C.struct_fuse_in_header, reqIntr *C.struct_fuse_interrupt_in) {

	intrN := fs.intrN
	intrN.unique = uint64(reqIntr.unique)

	fs.intrN = newInterrupNotice()
	intrN.next = fs.intrN

	close(intrN.ch)
}

func (fs *FuseServer) handlerFuseMessage(buf []byte, intrN *interrupNotice) {
	header := (*C.struct_fuse_in_header)(unsafe.Pointer(&buf[0]))
	bodyRaw := buf[C.SIZEOF_FUSE_IN_HEADER:]
	ctx := newFuseRequestContext()
	switch header.opcode {
	case C.FUSE_GETATTR:
		body := (*C.struct_fuse_getattr_in)(unsafe.Pointer(&bodyRaw[0]))
		go fs.hFuseGetAttr(ctx, header, body)
	default:
		replyRaw := make([]byte, C.SIZEOF_FUSE_OUT_HEADER)

		rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&replyRaw[0]))
		rheader.len = C.uint32_t(len(replyRaw))
		rheader.error = -C.ENOSYS
		rheader.unique = header.unique

		fs.send <- replyRaw
		return
	}

	for {
		select {
		case <-intrN.ch:
			if intrN.unique == uint64(header.unique) {
				ctx.setDone(EINTR)
			} else {
				intrN = intrN.next
			}
		case <-ctx.done:
			fs.send <- ctx.raw
			return
		}
	}
}

func (fs *FuseServer) hFuseGetAttr(ctx *FuseRequestContext,
	header *C.struct_fuse_in_header, body *C.struct_fuse_getattr_in) {
	ctx.raw = make([]byte, C.SIZEOF_FUSE_OUT_HEADER+C.SIZEOF_FUSE_ATTR_OUT)
	rbody := (*C.struct_fuse_attr_out)(
		unsafe.Pointer(&ctx.raw[C.SIZEOF_FUSE_OUT_HEADER]))
	attr := &rbody.attr

	err := fs.ops.GetAttr(ctx, "", attr)
	if err != nil {
		setErrorRaw(ctx, header, err)
	} else {
		rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&ctx.raw[0]))
		rheader.len = C.uint32_t(len(ctx.raw))
		rheader.error = 0
		rheader.unique = header.unique
	}
	ctx.setDone(err)
}

func setErrorRaw(ctx *FuseRequestContext, header *C.struct_fuse_in_header,
	err error) {
	ctx.raw = ctx.raw[:C.SIZEOF_FUSE_OUT_HEADER]
	rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&ctx.raw[0]))
	rheader.len = C.SIZEOF_FUSE_OUT_HEADER
	rheader.unique = header.unique
	if errno, ok := err.(syscall.Errno); ok {
		rheader.error = -C.int32_t(errno)
	} else {
		rheader.error = -C.int32_t(EIO)
	}
}
