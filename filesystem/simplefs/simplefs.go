package simplefs

import (
	"bytes"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/solomonwzs/gofuse/fuse"
)

var (
	_DLOG *log.Logger
	_USER *user.User
	_UID  uint32
	_GID  uint32
	_FT   *FileTree
)

type sampleFile struct {
	buf  []byte
	name string
	mode fuse.FileModeType
}

func newSampleFile(buf []byte, name string, mode fuse.FileModeType) (
	s *sampleFile) {
	return &sampleFile{
		buf:  buf,
		name: name,
		mode: mode,
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
		return 0, ERR_ILLEGAL_OPT
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

	_FT = NewFileTree()
	dirHello := _FT.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile(nil, "hello", fuse.S_IFDIR|0755))
	_FT.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile(nil, "world", fuse.S_IFDIR|0755))
	_FT.NewNode(fuse.ROOT_INODE_ID,
		newSampleFile([]byte("1234567"), "file0.txt", fuse.S_IFREG|0755))
	_FT.NewNode(dirHello.Ino(),
		newSampleFile([]byte("qwertyu"), "file1.txt", fuse.S_IFREG|0755))
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
	node := _FT.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}
	out.Valid = 1
	out.Attr = node.Attr()
	return
}

func (fs SimpleFS) SetAttr(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseSetAttrIn,
	out *fuse.FuseAttrOut,
) (err error) {
	header := ctx.Header()
	_DLOG.Printf("%+v\n", header)
	node := _FT.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	node.SetAttr(in)
	out.Attr = node.Attr()

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
	node := _FT.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	children := _FT.GetChildren(node.Ino())

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
		_DLOG.Println(n.Name())
		out.AddDirentRaw(fuse.NewFuseDirentRaw(
			n.attr.Ino, uint64(dOffset+1)+in.Offset, dt, []byte(n.Name())))
	}

	return
}

func (fs SimpleFS) Read(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseReadIn,
	out *bytes.Buffer,
) (err error) {
	header := ctx.Header()
	node := _FT.GetNode(header.Nodeid)
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

func (fs SimpleFS) Write(
	ctx *fuse.FuseRequestContext,
	in *fuse.FuseWriteIn,
	inRaw []byte,
	out *fuse.FuseWriteOut,
) (err error) {
	header := ctx.Header()
	node := _FT.GetNode(header.Nodeid)
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

func (fs SimpleFS) Lookup(
	ctx *fuse.FuseRequestContext,
	inName []byte,
	out *fuse.FuseEntryOut,
) (err error) {
	_DLOG.Println(string(inName))
	header := ctx.Header()
	node := _FT.GetNode(header.Nodeid)
	if node == nil {
		return fuse.ENOENT
	}

	for _, cIno := range _FT.GetChildren(node.Ino()) {
		if string(inName) == cIno.Name() {
			out.Nodeid = cIno.Ino()
			out.Attr = cIno.Attr()
			return
		}
	}
	return fuse.ENOENT
}
