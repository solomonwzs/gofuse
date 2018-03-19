package gofuse

import (
	"os"
)

const (
	_MAX_BUFFER_SIZE      = 65536
	_FUSE_MAX_BUFFER_SIZE = _MAX_BUFFER_SIZE + 100

	_FUSE_KERNEL_VERSION       = 7
	_FUSE_KERNEL_MINOR_VERSION = 26
)

var (
	_CMD_FUSERMOUNT string
)

func init() {
	if cmd := os.Getenv("CMD_FUSERMOUNT"); len(cmd) != 0 {
		_CMD_FUSERMOUNT = cmd
	} else {
		_CMD_FUSERMOUNT = "/bin/fusermount"
	}
}
