package simplefs

import (
	"fmt"
	"time"

	"github.com/solomonwzs/gofuse/fuse"
)

type SimpleFS struct {
	fuse.FileSystemUnimplemented
}

func (fs SimpleFS) GetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseGetattrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	attr := &out.Attr
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
	in *fuse.FuseOpenIn,
	out *fuse.FuseOpenOut,
) (err error) {
	out.Fh = 1
	out.Flags = fuse.FOPEN_DIRECT_IO | fuse.FOPEN_NONSEEKABLE
	return
}

func (fs SimpleFS) Read(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
) (err error) {
	time.Sleep(1 * time.Second)
	ctx.Write(fuse.NewFuseDirentRaw(1, 0, fuse.DT_REG, []byte("file-0")))
	fmt.Printf("%+v\n", in)
	return
}

func (fs SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	name []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	fmt.Println(string(name))
	return
}
