package simplefs

import (
	"bytes"
	"log"
	"os"
	"os/user"
	"strconv"
	"sync"

	"github.com/solomonwzs/gofuse/fuse"
)

var (
	_DLOG *log.Logger
	_USER *user.User
	_UID  uint32
	_GID  uint32
)

const _C_CODE = `
#include <stdio.h>
int main(int argc, char **argv) {
printf("hello world\n");
return 0;
}
`

type sampleFile struct {
	buf  []byte
	name string
	mode fuse.FileModeType
	uid  uint32
	gid  uint32
}

func newSampleFile(name string, mode fuse.FileModeType,
	uid uint32, gid uint32) (s *sampleFile) {
	return &sampleFile{
		buf:  []byte{},
		name: name,
		mode: mode,
		uid:  uid,
		gid:  gid,
	}
}

func (s *sampleFile) Name() string { return s.name }

func (s *sampleFile) Mode() fuse.FileModeType { return s.mode }

func (s *sampleFile) Size() uint64 {
	if s.mode&fuse.S_IFDIR != 0 {
		return 4096
	} else {
		return uint64(len(s.buf))
	}
}

func (s *sampleFile) Owner() (uint32, uint32) {
	return s.uid, s.gid
}

func (s *sampleFile) SetOwner(uid uint32, gid uint32) error {
	s.uid = uid
	s.gid = gid
	return nil
}

func (s *sampleFile) Resize(size uint64) error {
	if s.mode&fuse.S_IFREG == 0 {
		return ERR_ILLEGAL_OPT
	}
	if size <= s.Size() {
		s.buf = s.buf[:size]
	} else {
		s.buf = append(s.buf, make([]byte, size-s.Size())...)
	}
	return nil
}

func (s *sampleFile) Rename(name string) error {
	s.name = name
	return nil
}

func (s *sampleFile) ReadAt(b []byte, off int64) (n int, err error) {
	if s.mode&fuse.S_IFREG == 0 {
		return 0, ERR_ILLEGAL_OPT
	}
	if off >= int64(len(s.buf)) {
		return 0, nil
	} else {
		l := int64(len(b))
		end := off + l
		if end > int64(len(s.buf)) {
			end = int64(len(s.buf))
		}
		n = copy(b, s.buf[off:end])
		return
	}
}

func (s *sampleFile) WriteAt(b []byte, off int64) (n int, err error) {
	if s.mode&fuse.S_IFREG == 0 {
		return 0, ERR_ILLEGAL_OPT
	}
	if off > int64(len(s.buf)) {
		zerob := make([]byte, off-int64(len(s.buf)))
		s.buf = append(s.buf, zerob...)
	}

	if off+int64(len(b)) > int64(len(s.buf)) {
		s.buf = s.buf[:off]
		s.buf = append(s.buf, b...)
	} else {
		copy(s.buf[off:], b)
	}
	return len(b), nil
}

func init() {
	_DLOG = log.New(os.Stderr, "[simpleFS] ", log.Lshortfile)
	_USER, _ = user.Current()

	id, _ := strconv.Atoi(_USER.Uid)
	_UID = uint32(id)
	id, _ = strconv.Atoi(_USER.Gid)
	_GID = uint32(id)
}

type SimpleFS struct {
	fuse.FileSystemUnimplemented
	FTree    *FileTree
	treeLock *sync.Mutex
}

func NewSimpleFS() *SimpleFS {
	return &SimpleFS{
		FTree:    NewFileTree(),
		treeLock: &sync.Mutex{},
	}
}

func NewExampleSimpleFS() *SimpleFS {
	fs := NewSimpleFS()
	dirHello := fs.FTree.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile("hello", fuse.S_IFDIR|0755, _UID, _GID))
	fs.FTree.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile("world", fuse.S_IFDIR|0755, _UID, _GID))
	f0 := fs.FTree.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile("simple.c", fuse.S_IFREG|0644,
			_UID, _GID))
	f1 := fs.FTree.NewNode(dirHello.Ino(),
		newSampleFile("file.txt", fuse.S_IFREG|0644,
			_UID, _GID))
	f0.WriteAt([]byte(_C_CODE), 0)
	f1.WriteAt([]byte("1234567890"), 0)

	return fs
}

func (fs *SimpleFS) GetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseGetattrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}
	out.Valid = 1
	out.Attr = node.Attr()
	return
}

func (fs *SimpleFS) SetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseSetAttrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	node.SetAttr(in)
	out.Attr = node.Attr()

	return
}

func (fs *SimpleFS) Open(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseOpenIn,
	out *fuse.FuseOpenOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	// out.Flags = fuse.FOPEN_DIRECT_IO | fuse.FOPEN_NONSEEKABLE
	out.Flags = fuse.FOPEN_DIRECT_IO
	return
}

func (fs *SimpleFS) ReadDir(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
	out *fuse.FuseReadDirOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	children := fs.FTree.GetChildren(node.Ino())

	if in.Offset >= uint64(len(children)) {
		return
	}

	m := uint32(len(children))
	if m > in.Size {
		m = in.Size
	}
	for dOffset, n := range children[in.Offset:m] {
		var dt fuse.DirentType
		if n.attr.Mode&fuse.S_IFDIR != 0 {
			dt = fuse.DT_DIR
		} else {
			dt = fuse.DT_REG
		}
		out.AddDirentRaw(fuse.NewFuseDirentRaw(
			n.attr.Ino, uint64(dOffset+1)+in.Offset, dt, []byte(n.Name())))
	}

	return
}

func (fs *SimpleFS) Read(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
	out *bytes.Buffer,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	raw := make([]byte, in.Size)
	n, err := node.ReadAt(raw, int64(in.Offset))
	if err != nil {
		return
	}
	_, err = out.Write(raw[:n])
	return
}

func (fs *SimpleFS) Write(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseWriteIn,
	inRaw []byte,
	out *fuse.FuseWriteOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}
	n, err := node.WriteAt(inRaw, int64(in.Offset))
	if err != nil {
		return
	}
	out.Size = uint32(n)
	return
}

func (fs *SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	inName []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	for _, cIno := range fs.FTree.GetChildren(node.Ino()) {
		if string(inName) == cIno.Name() {
			out.Nodeid = cIno.Ino()
			out.Attr = cIno.Attr()
			out.Entry_valid = 1
			out.Attr_valid = 1
			return
		}
	}
	return fuse.ENOENT
}

func (fs *SimpleFS) Mknod(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseMknodIn,
	inName []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	f := newSampleFile(string(inName), in.Mode,
		header.Uid, header.Gid)
	n := fs.FTree.NewNode(header.Nodeid, f)

	out.Nodeid = n.Ino()
	out.Attr = n.Attr()

	return
}

func (fs *SimpleFS) Mkdir(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseMkdirIn,
	inName []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	f := newSampleFile(string(inName), in.Mode|fuse.S_IFDIR,
		header.Uid, header.Gid)
	n := fs.FTree.NewNode(header.Nodeid, f)

	out.Nodeid = n.Ino()
	out.Attr = n.Attr()

	return
}

func (fs *SimpleFS) Unlink(
	ctx *fuse.FuseRequestContext,
	inName []byte,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}
	_DLOG.Println(string(inName))
	for _, cIno := range fs.FTree.GetChildren(node.Ino()) {
		if string(inName) == cIno.Name() {
			if nlink, err := cIno.AddLink(-1); err != nil {
				return fuse.EPERM
			} else if nlink == 0 {
				fs.FTree.DelNode(cIno.Ino())
			}
			return nil
		}
	}
	return fuse.EPERM
}

func (fs *SimpleFS) Rmdir(
	ctx *fuse.FuseRequestContext,
	inName []byte,
) (err error) {
	header := ctx.Header()
	node := fs.FTree.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}
	_DLOG.Println(string(inName))
	for _, cIno := range fs.FTree.GetChildren(node.Ino()) {
		if string(inName) == cIno.Name() {
			fs.rmdir(cIno)
			return nil
		}
	}
	return fuse.EPERM
}

func (fs *SimpleFS) rmdir(node *FileNode) error {
	for _, n := range fs.FTree.GetChildren(node.Ino()) {
		if err := fs.rmdir(n); err != nil {
			return err
		}
		if nlink, err := n.AddLink(-1); err != nil {
			return err
		} else if nlink == 0 {
			fs.FTree.DelNode(n.Ino())
		}
	}
	fs.FTree.DelNode(node.Ino())
	return nil
}
