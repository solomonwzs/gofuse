package gofuse

/*
#cgo CFLAGS: -I${SRCDIR}/c

#include <errno.h>
#include <sys/stat.h>
#include "fuse_kernel_7_26.h"

#define SIZEOF_FUSE_IN_HEADER sizeof(struct fuse_in_header)
#define SIZEOF_FUSE_OUT_HEADER sizeof(struct fuse_out_header)
#define SIZEOF_FUSE_INIT_OUT sizeof(struct fuse_init_out)
#define SIZEOF_FUSE_ATTR_OUT sizeof(struct fuse_attr_out)
*/
import "C"
import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

func handleFuseRequest(buf []byte, w io.Writer) (err error) {

	header := (*C.struct_fuse_in_header)(unsafe.Pointer(&buf[0]))
	bodyRaw := buf[C.SIZEOF_FUSE_IN_HEADER:]
	switch header.opcode {
	case C.FUSE_INIT:
		reqInit := (*C.struct_fuse_init_in)(unsafe.Pointer(&bodyRaw[0]))
		err = handleFuseInit(header, reqInit, w)
	case C.FUSE_GETATTR:
		reqGetAttr := (*C.struct_fuse_getattr_in)(
			unsafe.Pointer(&bodyRaw[0]))
		err = handleFuseGetAttr(header, reqGetAttr, w)
	case C.FUSE_INTERRUPT:
		reqInterrupt := (*C.struct_fuse_interrupt_in)(
			unsafe.Pointer(&bodyRaw[0]))
		err = handlerFuseInterrupt(header, reqInterrupt, w)
	default:
		err = handlerFuseUnimplemented(header, bodyRaw, w)
	}

	return
}

func handlerFuseUnimplemented(header *C.struct_fuse_in_header, bodyRaw []byte,
	w io.Writer) (err error) {
	replyRaw := make([]byte, C.SIZEOF_FUSE_OUT_HEADER)
	rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&replyRaw[0]))

	rheader.len = C.uint32_t(len(replyRaw))
	rheader.error = -C.ENOSYS
	rheader.unique = header.unique

	_, err = w.Write(replyRaw)
	return
}

func handlerFuseInterrupt(header *C.struct_fuse_in_header,
	req *C.struct_fuse_interrupt_in, w io.Writer) (err error) {
	return
}

func handleFuseGetAttr(header *C.struct_fuse_in_header,
	req *C.struct_fuse_getattr_in, w io.Writer) (err error) {
	replyRaw := make([]byte, C.SIZEOF_FUSE_OUT_HEADER+C.SIZEOF_FUSE_ATTR_OUT)
	rheader := (*C.struct_fuse_out_header)(unsafe.Pointer(&replyRaw[0]))
	rbody := (*C.struct_fuse_attr_out)(unsafe.Pointer(
		&replyRaw[C.SIZEOF_FUSE_OUT_HEADER]))
	attr := &rbody.attr

	rheader.len = C.uint32_t(len(replyRaw))
	rheader.error = 0
	rheader.unique = header.unique

	rbody.attr_valid = 0
	rbody.attr_valid_nsec = 0
	rbody.dummy = req.dummy

	attr.ino = 0
	attr.size = 4096
	attr.blocks = 0
	attr.atime = C.uint64_t(time.Now().Unix())
	attr.mtime = C.uint64_t(time.Now().Unix())
	attr.ctime = C.uint64_t(time.Now().Unix())
	attr.mode = C.S_IFDIR | 0755
	attr.nlink = 1
	attr.blksize = 4096

	_, err = w.Write(replyRaw)
	return
}

func handleFuseInit(header *C.struct_fuse_in_header,
	req *C.struct_fuse_init_in, w io.Writer) (err error) {
	if req.major != _FUSE_KERNEL_VERSION || req.minor < _FUSE_KERNEL_VERSION {
		return fmt.Errorf(
			"gofuse: error fuse kernel version, expect %d.%d, got: %d.%d",
			_FUSE_KERNEL_VERSION, _FUSE_KERNEL_MINOR_VERSION,
			req.major, req.minor)
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

	_, err = w.Write(replyRaw)
	return
}
