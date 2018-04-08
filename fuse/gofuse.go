package fuse

import (
	"log"
	"os"
)

const (
	ROOT_INODE_ID = 1
)

const (
	_MAX_BUFFER_SIZE      = 65536
	_FUSE_MAX_BUFFER_SIZE = _MAX_BUFFER_SIZE + 100

	_FUSE_KERNEL_VERSION       = 7
	_FUSE_KERNEL_MINOR_VERSION = 26
)

var (
	_CMD_FUSERMOUNT string
	_DLOG           *log.Logger
)

func init() {
	_DLOG = log.New(os.Stderr, "[gofuse] ", log.Lshortfile)
	if cmd := os.Getenv("CMD_FUSERMOUNT"); len(cmd) != 0 {
		_CMD_FUSERMOUNT = cmd
	} else {
		_CMD_FUSERMOUNT = "/bin/fusermount"
	}
}
