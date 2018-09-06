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
