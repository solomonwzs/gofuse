package gofuse

import (
	"fmt"
	"io"
	"syscall"
	"time"
	"unsafe"
)

func handleFuseRequest(buf []byte, w io.Writer) (err error) {

	header := (*FuseInHeader)(unsafe.Pointer(&buf[0]))
	bodyRaw := buf[_SIZEOF_FUSE_IN_HEADER:]
	switch header.Opcode {
	case FUSE_INIT:
		reqInit := (*FuseInitIn)(unsafe.Pointer(&bodyRaw[0]))
		err = handleFuseInit(header, reqInit, w)
	case FUSE_GETATTR:
		reqGetAttr := (*FuseGetattrIn)(
			unsafe.Pointer(&bodyRaw[0]))
		err = handleFuseGetAttr(header, reqGetAttr, w)
	case FUSE_INTERRUPT:
		reqInterrupt := (*FuseInterruptIn)(
			unsafe.Pointer(&bodyRaw[0]))
		err = handlerFuseInterrupt(header, reqInterrupt, w)
	default:
		err = handlerFuseUnimplemented(header, bodyRaw, w)
	}

	return
}

func handlerFuseUnimplemented(header *FuseInHeader, bodyRaw []byte,
	w io.Writer) (err error) {
	replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER)
	rheader := (*FuseOutHeader)(unsafe.Pointer(&replyRaw[0]))

	rheader.Len = uint32(len(replyRaw))
	rheader.Error = -int32(syscall.ENOSYS)
	rheader.Unique = header.Unique

	_, err = w.Write(replyRaw)
	return
}

func handlerFuseInterrupt(header *FuseInHeader,
	req *FuseInterruptIn, w io.Writer) (err error) {
	return
}

func handleFuseGetAttr(header *FuseInHeader,
	req *FuseGetattrIn, w io.Writer) (err error) {
	replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER+_SIZEOF_FUSE_ATTR_OUT)
	rheader := (*FuseOutHeader)(unsafe.Pointer(&replyRaw[0]))
	rbody := (*FuseAttrOut)(unsafe.Pointer(
		&replyRaw[_SIZEOF_FUSE_OUT_HEADER]))
	attr := &rbody.Attr

	rheader.Len = uint32(len(replyRaw))
	rheader.Error = 0
	rheader.Unique = header.Unique

	rbody.Valid = 0
	rbody.Valid_nsec = 0
	rbody.Dummy = req.Dummy

	attr.Ino = 0
	attr.Size = 4096
	attr.Blocks = 0
	attr.Atime = uint64(time.Now().Unix())
	attr.Mtime = uint64(time.Now().Unix())
	attr.Ctime = uint64(time.Now().Unix())
	attr.Mode = S_IFDIR | 0755
	attr.Nlink = 1
	attr.Blksize = 4096

	_, err = w.Write(replyRaw)
	return
}

func handleFuseInit(header *FuseInHeader,
	req *FuseInitIn, w io.Writer) (err error) {
	if req.Major != _FUSE_KERNEL_VERSION ||
		req.Minor < _FUSE_KERNEL_VERSION {
		return fmt.Errorf(
			"gofuse: error fuse kernel version, expect %d.%d, got: %d.%d",
			_FUSE_KERNEL_VERSION, _FUSE_KERNEL_MINOR_VERSION,
			req.Major, req.Minor)
	}

	replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER+_SIZEOF_FUSE_INIT_OUT)
	rheader := (*FuseOutHeader)(unsafe.Pointer(&replyRaw[0]))
	rbody := (*FuseInitOut)(unsafe.Pointer(
		&replyRaw[_SIZEOF_FUSE_OUT_HEADER]))

	rheader.Len = uint32(len(replyRaw))
	rheader.Error = 0
	rheader.Unique = header.Unique

	rbody.Major = _FUSE_KERNEL_VERSION
	rbody.Minor = _FUSE_KERNEL_MINOR_VERSION
	rbody.Max_readahead = _MAX_BUFFER_SIZE
	rbody.Flags = 0
	rbody.Max_background = 0
	rbody.Congestion_threshold = 0
	rbody.Max_write = _MAX_BUFFER_SIZE
	rbody.Time_gran = 0

	_, err = w.Write(replyRaw)
	return
}
