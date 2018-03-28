package simplefs

import (
	"bytes"
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
	children []*simpleInode
}

var (
	_DLOG   *log.Logger
	_USER   *user.User
	_UID    uint32
	_GID    uint32
	_INODES map[uint64]*simpleInode
)

const (
	_N_ROOT = fuse.ROOT_INODE_ID + iota
	_N_FILE_1
	_N_FILE_2
	_N_DIR_1
	_N_DIR_2
)

var (
	_FILES = map[uint64][]byte{
		_N_FILE_1: []byte("0123456789"),
		_N_FILE_2: []byte("qwertyuiopasdfghjkl"),
	}
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
	_INODES[_N_ROOT] = &simpleInode{
		name: "",
		attr: fuse.FuseAttr{
			Ino:     _N_ROOT,
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
	_INODES[_N_FILE_1] = &simpleInode{
		name: "file-0.txt",
		attr: fuse.FuseAttr{
			Ino:     _N_FILE_1,
			Blocks:  1,
			Size:    uint64(len(_FILES[_N_FILE_1])),
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
	_INODES[_N_FILE_2] = &simpleInode{
		name: "file-1.txt",
		attr: fuse.FuseAttr{
			Ino:     _N_FILE_2,
			Blocks:  1,
			Size:    uint64(len(_FILES[_N_FILE_1])),
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
	_INODES[_N_DIR_1] = &simpleInode{
		name: "hello",
		attr: fuse.FuseAttr{
			Ino:     _N_DIR_1,
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
	_INODES[_N_DIR_2] = &simpleInode{
		name: "world",
		attr: fuse.FuseAttr{
			Ino:     _N_DIR_2,
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
	_INODES[_N_ROOT].children = []*simpleInode{
		_INODES[_N_FILE_1],
		_INODES[_N_DIR_1],
		_INODES[_N_DIR_2],
	}
	_INODES[_N_DIR_2].children = []*simpleInode{
		_INODES[_N_FILE_2],
	}
}

type SimpleFS struct {
	fuse.FileSystemUnimplemented
}

func (fs SimpleFS) GetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseGetattrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	header := ctx.Header()
	ino, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}
	out.Attr = ino.attr
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

func (fs SimpleFS) ReadDir(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
	out *fuse.FuseReadDirOut,
) (err error) {
	header := ctx.Header()
	ino, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}
	if in.Offset >= uint64(len(ino.children)) {
		return
	}

	m := uint32(len(ino.children))
	if m > in.Size {
		m = in.Size
	}
	for dOffset, n := range ino.children[in.Offset:m] {
		var dt fuse.DirentType
		if n.attr.Mode&fuse.S_IFDIR != 0 {
			dt = fuse.DT_DIR
		} else {
			dt = fuse.DT_REG
		}
		_DLOG.Println(n.name)
		out.AddDirentRaw(fuse.NewFuseDirentRaw(
			n.attr.Ino, uint64(dOffset+1)+in.Offset, dt, []byte(n.name)))
	}

	return
}

func (fs SimpleFS) Read(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
	out *bytes.Buffer,
) (err error) {
	header := ctx.Header()
	_DLOG.Printf("%+v\n", in)
	ino, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}
	if ino.attr.Mode&fuse.S_IFDIR != 0 {
		return
	}
	if in.Offset > ino.attr.Size {
		return
	}

	n := int(in.Size)
	raw := _FILES[header.Nodeid]
	if n > len(raw) {
		n = len(raw)
	}
	_, err = out.Write(raw[in.Offset:n])
	return
}

func (fs SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	name []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	header := ctx.Header()
	pIno, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}

	var cIno *simpleInode
	for _, cIno = range pIno.children {
		if string(name) == cIno.name {
			out.Nodeid = cIno.attr.Ino
			out.Attr = cIno.attr
			return
		}
	}
	return fuse.ENOENT
}
