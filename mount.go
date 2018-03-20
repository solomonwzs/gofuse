package gofuse

import (
	"errors"
	"os"
	"sync"
	"time"
)

type FileSystem struct {
	dir         string
	kFilehandle *os.File
	end         chan struct{}
	endLock     *sync.Mutex
}

func NewFileSystem(dir string, conf *MountConfig) (fs *FileSystem, err error) {
	if err = checkDir(dir); err != nil {
		return
	}
	f, err := mount(dir, conf)
	if err != nil {
		return
	}
	fs = &FileSystem{
		dir:         dir,
		kFilehandle: f,
		end:         make(chan struct{}),
		endLock:     &sync.Mutex{},
	}
	go fs.serv()
	return
}

func (fs *FileSystem) IsClosed() bool {
	select {
	case <-fs.end:
		return true
	default:
		return false
	}
}

func (fs *FileSystem) Close() error {
	fs.endLock.Lock()
	defer fs.endLock.Unlock()

	if fs.IsClosed() {
		return errors.New("gofuse: file system was closed")
	}
	close(fs.end)
	fs.kFilehandle.Close()
	umount(fs.dir)

	return nil
}

func (fs *FileSystem) serv() {
	for {
		buf := make([]byte, _FUSE_MAX_BUFFER_SIZE)
		time.Sleep(2 * time.Second)
		n, err := fs.kFilehandle.Read(buf)
		_DLOG.Println(n)
		if err != nil {
			return
		}

		go handleFuseRequest(buf[:n], fs.kFilehandle)
	}
}
