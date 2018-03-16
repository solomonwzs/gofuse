package gofuse

import (
	"fmt"
	"unsafe"
)

// fuse opcode
const (
	_FUSE_LOOKUP       = 1
	_FUSE_FORGET       = 2 /* no reply */
	_FUSE_GETATTR      = 3
	_FUSE_SETATTR      = 4
	_FUSE_READLINK     = 5
	_FUSE_SYMLINK      = 6
	_FUSE_MKNOD        = 8
	_FUSE_MKDIR        = 9
	_FUSE_UNLINK       = 10
	_FUSE_RMDIR        = 11
	_FUSE_RENAME       = 12
	_FUSE_LINK         = 13
	_FUSE_OPEN         = 14
	_FUSE_READ         = 15
	_FUSE_WRITE        = 16
	_FUSE_STATFS       = 17
	_FUSE_RELEASE      = 18
	_FUSE_FSYNC        = 20
	_FUSE_SETXATTR     = 21
	_FUSE_GETXATTR     = 22
	_FUSE_LISTXATTR    = 23
	_FUSE_REMOVEXATTR  = 24
	_FUSE_FLUSH        = 25
	_FUSE_INIT         = 26
	_FUSE_OPENDIR      = 27
	_FUSE_READDIR      = 28
	_FUSE_RELEASEDIR   = 29
	_FUSE_FSYNCDIR     = 30
	_FUSE_GETLK        = 31
	_FUSE_SETLK        = 32
	_FUSE_SETLKW       = 33
	_FUSE_ACCESS       = 34
	_FUSE_CREATE       = 35
	_FUSE_INTERRUPT    = 36
	_FUSE_BMAP         = 37
	_FUSE_DESTROY      = 38
	_FUSE_IOCTL        = 39
	_FUSE_POLL         = 40
	_FUSE_NOTIFY_REPLY = 41
	_FUSE_BATCH_FORGET = 42
	_FUSE_FALLOCATE    = 43

	// CUSE specific operations
	_CUSE_INIT = 4096
)

type fuseInHeader struct {
	len     uint32
	opcode  uint32
	unique  uint64
	nodeid  uint64
	uid     uint32
	gid     uint32
	pid     uint32
	padding uint32
}

type fuseOutHeader struct {
	len    uint32
	err    int32
	unique uint64
}

type fuseInitIn struct {
	major        uint32
	minor        uint32
	maxReadahead uint32
	flags        uint32
}

type fuseInitOut struct {
	major    uint32
	minor    uint32
	unused   uint32
	flags    uint32
	maxRead  uint32
	maxWrite uint32
	devMajor uint32
	devMinor uint32
	spare    [10]uint32
}

func parseFuseInHeader(buf []byte) {
	header := (*fuseInHeader)(unsafe.Pointer(&buf[0]))
	fmt.Println(header, unsafe.Sizeof(*header))
}
