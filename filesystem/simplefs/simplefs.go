package simplefs

import (
	"log"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/solomonwzs/gofuse/fuse"
)

type simpleInode struct {
	attr     fuse.FuseAttr
	name     string
	children map[string]*simpleInode
}

var (
	_DLOG   *log.Logger
	_USER   *user.User
	_UID    uint32
	_GID    uint32
	_INODES map[uint64]*simpleInode
)

const (
	_FILE_NAME = "only-one-file.txt"
)

func init() {
	_DLOG = log.New(os.Stderr, "[simpleFS] ", log.Lshortfile)
	_USER, _ = user.Current()
	_INODES = make(map[uint64]*simpleInode)

	id, _ := strconv.Atoi(_USER.Uid)
	_UID = uint32(id)
	id, _ = strconv.Atoi(_USER.Gid)
	_GID = uint32(id)

	now := uint64(time.Now().Unix())
	_INODES[0] = &simpleInode{
		name:     "",
		children: make(map[string]*simpleInode),
		attr: fuse.FuseAttr{
			Ino:     0,
			Blocks:  1,
			Size:    4096,
			Blksize: 4096,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    fuse.S_IFDIR | 0755,
			Nlink:   1,
		},
	}
	_INODES[1] = &simpleInode{
		name:     "file-0.txt",
		children: make(map[string]*simpleInode),
		attr: fuse.FuseAttr{
			Ino:     1,
			Blocks:  1,
			Size:    4096,
			Blksize: 4096,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    fuse.S_IFREG | 0755,
			Nlink:   1,
		},
	}
	_INODES[2] = &simpleInode{
		name:     "file-1.txt",
		children: make(map[string]*simpleInode),
		attr: fuse.FuseAttr{
			Ino:     2,
			Blocks:  1,
			Size:    4096,
			Blksize: 4096,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    fuse.S_IFREG | 0755,
			Nlink:   1,
		},
	}
	_INODES[3] = &simpleInode{
		name:     "hello",
		children: make(map[string]*simpleInode),
		attr: fuse.FuseAttr{
			Ino:     3,
			Blocks:  1,
			Size:    4096,
			Blksize: 4096,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    fuse.S_IFDIR | 0755,
			Nlink:   1,
		},
	}
	_INODES[4] = &simpleInode{
		name:     "world",
		children: make(map[string]*simpleInode),
		attr: fuse.FuseAttr{
			Ino:     4,
			Blocks:  1,
			Size:    4096,
			Blksize: 4096,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    fuse.S_IFDIR | 0755,
			Nlink:   1,
		},
	}
	_INODES[0].children[_INODES[1].name] = _INODES[1]
	_INODES[0].children[_INODES[3].name] = _INODES[3]
	_INODES[0].children[_INODES[4].name] = _INODES[4]
	_INODES[4].children[_INODES[2].name] = _INODES[2]
}

type SimpleFS struct {
	fuse.FileSystemUnimplemented
}

func (fs SimpleFS) GetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseGetattrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	_DLOG.Printf("%+v\n", ctx.Header())
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
	ctx.Write(fuse.NewFuseDirentRaw(1, 0, fuse.DT_REG, []byte(_FILE_NAME)))
	return
}

func (fs SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	name []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	_DLOG.Println(string(name))
	if string(name) != _FILE_NAME {
		return
	}
	header := ctx.Header()
	_DLOG.Printf("%+v\n", header)

	out.Nodeid = 1
	attr := &out.Attr
	attr.Ino = 1
	attr.Size = 4096
	attr.Blocks = 1
	attr.Atime = uint64(time.Now().Unix())
	attr.Mtime = uint64(time.Now().Unix())
	attr.Ctime = uint64(time.Now().Unix())
	attr.Mode = fuse.S_IFREG | 0755
	attr.Nlink = 1
	attr.Blksize = 4096
	attr.Uid = header.Uid
	attr.Gid = header.Gid
	attr.Rdev = fuse.S_IFCHR | 0755

	return
}
