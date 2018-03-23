package fuse

import (
	"bytes"
	"errors"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

func getReplyBodySize(header *FuseInHeader) (size int) {
	switch header.Opcode {
	case FUSE_GETATTR:
		size = _SIZEOF_FUSE_ATTR_OUT
	case FUSE_OPENDIR, FUSE_OPEN:
		size = _SIZEOF_FUSE_OPEN_OUT
	case FUSE_LOOKUP:
		size = _SIZEOF_FUSE_ENTRY_OUT
	default:
		size = 0
	}
	return
}

type FuseRequestContext struct {
	kv     map[interface{}]interface{}
	kvLock *sync.RWMutex

	deadline atomic.Value
	header   *FuseInHeader
	raw      []byte

	extBuffer    *bytes.Buffer
	extSizeLimit atomic.Value
	extLock      *sync.Mutex

	done       chan struct{}
	doneReason error
	doneLock   *sync.Mutex
}

func newFuseRequestContext(header *FuseInHeader) (ctx *FuseRequestContext) {
	size := _SIZEOF_FUSE_OUT_HEADER + getReplyBodySize(header)
	raw := make([]byte, size, size)

	ctx = &FuseRequestContext{
		kv:     make(map[interface{}]interface{}),
		kvLock: &sync.RWMutex{},

		header: header,
		raw:    raw,

		extBuffer: new(bytes.Buffer),
		extLock:   &sync.Mutex{},

		done:       make(chan struct{}),
		doneReason: nil,
		doneLock:   &sync.Mutex{},
	}
	rheader := ctx.outHeader()
	rheader.Unique = header.Unique

	return
}

func (ctx *FuseRequestContext) setExtBufferSizeLimit(size uint32) {
	ctx.extSizeLimit.Store(size)
}

func (ctx *FuseRequestContext) Header() *FuseInHeader {
	return ctx.header
}

func (ctx *FuseRequestContext) IsDone() bool {
	select {
	case <-ctx.done:
		return true
	default:
		return false
	}
}

func (ctx *FuseRequestContext) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = ctx.deadline.Load().(time.Time)
	return
}

func (ctx *FuseRequestContext) setDeadline(t time.Time) {
	ctx.deadline.Store(t)
}

func (ctx *FuseRequestContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *FuseRequestContext) Err() error {
	select {
	case <-ctx.done:
		return ctx.doneReason
	default:
		return nil
	}
}

func (ctx *FuseRequestContext) Value(key interface{}) interface{} {
	ctx.kvLock.RLock()
	defer ctx.kvLock.RUnlock()

	return ctx.kv[key]
}

func (ctx *FuseRequestContext) SetValue(key, value interface{}) (
	oldValue interface{}) {
	ctx.kvLock.Lock()
	defer ctx.kvLock.Unlock()

	oldValue, _ = ctx.kv[key]
	ctx.kv[key] = value
	return
}

func (ctx *FuseRequestContext) setDone(reason error) error {
	ctx.doneLock.Lock()
	defer ctx.doneLock.Unlock()

	if ctx.IsDone() {
		return errors.New("gofuse: context was closed")
	}
	ctx.doneReason = reason
	close(ctx.done)

	return nil
}

func (ctx *FuseRequestContext) outHeader() *FuseOutHeader {
	return (*FuseOutHeader)(unsafe.Pointer(&ctx.raw[0]))
}

func (ctx *FuseRequestContext) outBody() unsafe.Pointer {
	return unsafe.Pointer(&ctx.raw[_SIZEOF_FUSE_OUT_HEADER])
}

func (ctx *FuseRequestContext) Write(p []byte) (n int, err error) {
	if ctx.IsDone() {
		return 0, errors.New("gofuse: context was closed")
	}

	ctx.extLock.Lock()
	defer ctx.extLock.Unlock()

	size, ok := ctx.extSizeLimit.Load().(uint32)
	extLen := uint32(ctx.extBuffer.Len())
	if ok && extLen < size {
		if uint32(len(p)) < size-extLen {
			return ctx.extBuffer.Write(p)
		} else {
			return ctx.extBuffer.Write(p[:size-extLen])
		}
	} else {
		return 0, errors.New("gofuse: buffer full")
	}
}

func (ctx *FuseRequestContext) replyRaw() []byte {
	if !ctx.IsDone() {
		return nil
	}

	rheader := ctx.outHeader()
	if ctx.doneReason != nil {
		rheader.Len = _SIZEOF_FUSE_OUT_HEADER
		if errno, ok := ctx.doneReason.(syscall.Errno); ok {
			rheader.Error = -int32(errno)
		} else {
			rheader.Error = -int32(EIO)
		}
		return ctx.raw[:_SIZEOF_FUSE_OUT_HEADER]
	} else {
		rheader.Error = 0
		ctx.extLock.Lock()
		defer ctx.extLock.Unlock()

		if extLen := ctx.extBuffer.Len(); extLen == 0 {
			rheader.Len = uint32(len(ctx.raw))
			return ctx.raw
		} else {
			rheader.Len = uint32(len(ctx.raw) + extLen)
			buf := make([]byte, len(ctx.raw)+extLen)
			copy(buf, ctx.raw)
			copy(buf[len(ctx.raw):], ctx.extBuffer.Bytes())
			return buf
		}
	}
}
