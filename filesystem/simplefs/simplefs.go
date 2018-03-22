package simplefs

import (
	"time"

	"github.com/solomonwzs/gofuse/fuse"
)

type SimpleFS struct {
	fuse.FileSystemUnimplemented
}

func (fs SimpleFS) GetAttr(
	ctx *fuse.FuseRequestContext,
	attr *fuse.FuseAttr,
) (err error) {
	attr.Ino = 0
	attr.Size = 4096
	attr.Blocks = 0
	attr.Atime = uint64(time.Now().Unix())
	attr.Mtime = uint64(time.Now().Unix())
	attr.Ctime = uint64(time.Now().Unix())
	attr.Mode = fuse.S_IFDIR | 0755
	attr.Nlink = 1
	attr.Blksize = 4096
	return
}

func (fs SimpleFS) Open(
	ctx *fuse.FuseRequestContext,
	open *fuse.FuseOpenOut,
) (err error) {
	open.Fh = 1
	open.Flags = fuse.FOPEN_DIRECT_IO | fuse.FOPEN_NONSEEKABLE
	return
}
