package gofuse

import (
	"os"
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
