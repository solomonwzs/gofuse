package gofuse

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type FileSystem struct {
	dir         string
	kFilehandle *os.File
	end         chan struct{}
	endLock     *sync.Mutex
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
	buf := make([]byte, _FUSE_MAX_BUFFER_SIZE)
	for {
		n, err := fs.kFilehandle.Read(buf)
		if err != nil {
			return
		}
		if err := handleFuseRequest(buf[:n], fs.kFilehandle); err != nil {
			fs.Close()
			return
		}
	}
}
