package fuse

import (
	"errors"
	"fmt"
	"os"
	"sync"
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
	buf := make([]byte, _SIZEOF_FUSE_IN_HEADER+_SIZEOF_FUSE_INIT_IN)
	n, err := f.Read(buf)
	if err != nil {
		return
	} else if n != len(buf) {
		return errors.New("gofuse: error fuse request format")
	}

	header := (*FuseInHeader)(unsafe.Pointer(&buf[0]))
	if header.Opcode != FUSE_INIT {
		return fmt.Errorf(
			"gofuse: error fuse request opcode, expect: %d, got: %d",
			FUSE_INIT, header.Opcode)
	}

	bodyRaw := buf[_SIZEOF_FUSE_IN_HEADER:]
	body := (*FuseInitIn)(unsafe.Pointer(&bodyRaw[0]))
	if body.Major != _FUSE_KERNEL_VERSION ||
		body.Minor < _FUSE_KERNEL_VERSION {
		return fmt.Errorf(
			"gofuse: error fuse kernel version, expect %d.%d, got: %d.%d",
			_FUSE_KERNEL_VERSION, _FUSE_KERNEL_MINOR_VERSION,
			body.Major, body.Minor)
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

func (fs *FuseServer) Close() (err error) {
	fs.endLock.Lock()
	defer fs.endLock.Unlock()

	if fs.IsClosed() {
		return errors.New("gofuse: fuse server was closed")
	}
	if err = umount(fs.dir); err != nil {
		return
	}

	close(fs.end)
	_DLOG.Println(fs.f.Close())

	return nil
}

func (fs *FuseServer) readLoop() {
	buf := make([]byte, _FUSE_MAX_BUFFER_SIZE)
	for {
		n, err := fs.f.Read(buf)
		if err != nil {
			_DLOG.Println(err)
			fs.Close()
			return
		}

		header := (*FuseInHeader)(unsafe.Pointer(&buf[0]))
		_DLOG.Println(header.Opcode)
		switch header.Opcode {
		case FUSE_INTERRUPT:
			reqIntr := (*FuseInterruptIn)(
				unsafe.Pointer(&buf[_SIZEOF_FUSE_IN_HEADER]))
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
				_DLOG.Println(err)
				fs.Close()
				return
			}
		case <-fs.end:
			return
		}
	}
}

func (fs *FuseServer) handlerFuseInterrupt(
	header *FuseInHeader, reqIntr *FuseInterruptIn) {

	intrN := fs.intrN
	intrN.unique = uint64(reqIntr.Unique)

	fs.intrN = newInterrupNotice()
	intrN.next = fs.intrN

	close(intrN.ch)
}

func (fs *FuseServer) handlerFuseMessage(buf []byte, intrN *interrupNotice) {
	header := (*FuseInHeader)(unsafe.Pointer(&buf[0]))
	bodyRaw := buf[_SIZEOF_FUSE_IN_HEADER:]
	ctx := newFuseRequestContext(header)
	switch header.Opcode {
	case FUSE_GETATTR:
		in := (*FuseGetattrIn)(unsafe.Pointer(&bodyRaw[0]))
		out := (*FuseAttrOut)(ctx.outBody())
		go func() {
			err := fs.ops.GetAttr(ctx, in, out)
			ctx.setDone(err)
		}()
	case FUSE_OPEN, FUSE_OPENDIR:
		in := (*FuseOpenIn)(unsafe.Pointer(&bodyRaw[0]))
		out := (*FuseOpenOut)(ctx.outBody())
		go func() {
			err := fs.ops.Open(ctx, in, out)
			ctx.setDone(err)
		}()
	case FUSE_READ, FUSE_READDIR:
		in := (*FuseReadIn)(unsafe.Pointer(&bodyRaw[0]))
		ctx.setExtBufferSizeLimit(in.Size)
		go func() {
			err := fs.ops.Read(ctx, in)
			ctx.setDone(err)
		}()
	case FUSE_LOOKUP:
		var i int
		for i = 0; i < len(bodyRaw); i++ {
			if bodyRaw[i] == '\x00' {
				break
			}
		}
		name := bodyRaw[:i]
		out := (*FuseEntryOut)(ctx.outBody())
		go func() {
			err := fs.ops.Lookup(ctx, name, out)
			ctx.setDone(err)
		}()
	case FUSE_RELEASE, FUSE_RELEASEDIR:
		in := (*FuseReleaseIn)(unsafe.Pointer(&bodyRaw[0]))
		go func() {
			err := fs.ops.Release(ctx, in)
			ctx.setDone(err)
		}()
	case FUSE_DESTROY, FUSE_FORGET:
		return
	default:
		replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER)
		writeErrorRaw(replyRaw, header, syscall.ENOSYS)
		fs.send <- replyRaw
		return
	}

	for {
		select {
		case <-intrN.ch:
			if intrN.unique == uint64(header.Unique) {
				ctx.setDone(EINTR)

				replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER)
				writeErrorRaw(replyRaw, header, syscall.EINTR)
				fs.send <- replyRaw

				return
			} else {
				intrN = intrN.next
			}
		case <-ctx.Done():
			fs.send <- ctx.replyRaw()
			return
		case <-fs.end:
			replyRaw := make([]byte, _SIZEOF_FUSE_OUT_HEADER)
			writeErrorRaw(replyRaw, header, syscall.EINTR)
			_DLOG.Println(fs.f.Write(replyRaw))
		}
	}
}

func writeErrorRaw(raw []byte, header *FuseInHeader, err error) {
	rheader := (*FuseOutHeader)(unsafe.Pointer(&raw[0]))
	rheader.Len = _SIZEOF_FUSE_OUT_HEADER
	rheader.Unique = header.Unique
	if errno, ok := err.(syscall.Errno); ok {
		rheader.Error = -int32(errno)
	} else {
		rheader.Error = -int32(EIO)
	}
	return
}
