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
	out.Valid = 1
	out.Attr = ino.attr
	return
}

func (fs SimpleFS) SetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseSetAttrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	header := ctx.Header()
	_DLOG.Printf("%+v\n", header)
	ino, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}

	if in.Valid&fuse.FATTR_MODE != 0 {
		ino.attr.Mode = in.Mode
	}
	if in.Valid&fuse.FATTR_UID != 0 {
		ino.attr.Uid = in.Uid
	}
	if in.Valid&fuse.FATTR_GID != 0 {
		ino.attr.Gid = in.Gid
	}
	if in.Valid&fuse.FATTR_SIZE != 0 {
		ino.attr.Size = in.Size
		f := _FILES[header.Nodeid]
		if ino.attr.Size > uint64(len(f)) {
			newRaw := make([]byte, ino.attr.Size)
			copy(newRaw, f)
			_FILES[header.Nodeid] = newRaw
		} else {
			_FILES[header.Nodeid] = f[:ino.attr.Size]
		}
	}
	if in.Valid&fuse.FATTR_ATIME != 0 {
		ino.attr.Atime = in.Atime
		ino.attr.Atimensec = in.Atimensec
	}
	if in.Valid&fuse.FATTR_MTIME != 0 {
		ino.attr.Mtime = in.Mtime
		ino.attr.Mtimensec = in.Mtimensec
	}
	if in.Valid&fuse.FATTR_CTIME != 0 {
		ino.attr.Ctime = in.Ctime
		ino.attr.Ctimensec = in.Ctimensec
	}
	out.Attr = ino.attr

	_DLOG.Printf("%+v\n", in)
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
	ino.attr.Atime = uint64(time.Now().Unix())

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
	ino.attr.Atime = uint64(time.Now().Unix())

	n := uint64(in.Size)
	raw := _FILES[header.Nodeid]
	if n > ino.attr.Size {
		n = ino.attr.Size
	}
	_, err = out.Write(raw[in.Offset:n])
	return
}

func (fs SimpleFS) Write(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseWriteIn,
	inRaw []byte,
	out *fuse.FuseWriteOut,
) (err error) {
	header := ctx.Header()
	ino, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}
	if ino.attr.Mode&fuse.S_IFDIR != 0 {
		return
	}
	if in.Offset > ino.attr.Size {
		return fuse.EPERM
	}
	ino.attr.Mtime = uint64(time.Now().Unix())

	f := _FILES[header.Nodeid]
	if in.Offset+uint64(in.Size) > ino.attr.Size {
		f = f[:in.Offset]
		f = append(f, inRaw...)
		_FILES[header.Nodeid] = f
	} else {
		copy(f[in.Offset:], inRaw)
	}
	ino.attr.Size = uint64(len(f))

	out.Size = in.Size

	return
}

func (fs SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	inName []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	_DLOG.Println(string(inName))
	header := ctx.Header()
	pIno, exist := _INODES[header.Nodeid]
	if !exist {
		return fuse.ENOENT
	}

	var cIno *simpleInode
	for _, cIno = range pIno.children {
		if string(inName) == cIno.name {
			out.Nodeid = cIno.attr.Ino
			out.Attr = cIno.attr
			return
		}
	}
	return fuse.ENOENT
}
