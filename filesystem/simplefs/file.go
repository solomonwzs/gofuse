package simplefs

import (
	"errors"
	"sort"
	"sync/atomic"
	"time"

	"github.com/solomonwzs/gofuse/fuse"
)

const (
	_BLOCK_SIZE = 4096
)

var (
	ERR_NOT_REG        = errors.New("it is not a regular file")
	ERR_ILLEGAL_OPT    = errors.New("illegal operation")
	ERR_NODE_NOT_EXIST = errors.New("node not exist")
)

type File interface {
	Name() (name string)
	Rename(name string) (err error)
	ReadAt(b []byte, off int64) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	Mode() fuse.FileModeType
	Size() uint64
	Resize(uint64) error
}

type rootFile struct{}

func (r rootFile) Name() string { return "." }

func (r rootFile) Rename(name string) error { return ERR_ILLEGAL_OPT }

func (r rootFile) Mode() fuse.FileModeType { return fuse.S_IFDIR | 0755 }

func (r rootFile) Size() uint64 { return 4096 }

func (r rootFile) Resize(size uint64) error { return ERR_ILLEGAL_OPT }

func (r rootFile) ReadAt(b []byte, off int64) (int, error) {
	return 0, ERR_NOT_REG
}

func (r rootFile) WriteAt(b []byte, off int64) (int, error) {
	return 0, ERR_NOT_REG
}

type FileNode struct {
	File
	attr     fuse.FuseAttr
	parent   *FileNode
	children map[uint64]*FileNode
}

func (fn *FileNode) upateSize() {
	size := fn.Size()
	blocks := size / uint64(fn.attr.Blksize)
	if size%uint64(fn.attr.Blksize) != 0 {
		blocks += 1
	}
	fn.attr.Size = size
	fn.attr.Blocks = blocks
}

func (fn *FileNode) Ino() uint64 { return fn.attr.Ino }

func (fn *FileNode) Attr() fuse.FuseAttr { return fn.attr }

func (fn *FileNode) SetAttr(in *fuse.FuseSetAttrIn) {
	if in.Valid&fuse.FATTR_MODE != 0 {
		fn.attr.Mode = in.Mode
	}
	if in.Valid&fuse.FATTR_UID != 0 {
		fn.attr.Uid = in.Uid
	}
	if in.Valid&fuse.FATTR_GID != 0 {
		fn.attr.Gid = in.Gid
	}
	if in.Valid&fuse.FATTR_SIZE != 0 {
		fn.Resize(in.Size)
		fn.attr.Size = fn.Size()
	}
	if in.Valid&fuse.FATTR_ATIME != 0 {
		fn.attr.Atime = in.Atime
		fn.attr.Atimensec = in.Atimensec
	}
	if in.Valid&fuse.FATTR_MTIME != 0 {
		fn.attr.Mtime = in.Mtime
		fn.attr.Mtimensec = in.Mtimensec
	}
	if in.Valid&fuse.FATTR_CTIME != 0 {
		fn.attr.Ctime = in.Ctime
		fn.attr.Ctimensec = in.Ctimensec
	}
}

func (fn *FileNode) ReadAt(b []byte, off int64) (int, error) {
	fn.attr.Atime = uint64(time.Now().Unix())
	return fn.File.ReadAt(b, off)
}

func (fn *FileNode) WriteAt(b []byte, off int64) (int, error) {
	now := uint64(time.Now().Unix())
	fn.attr.Atime = now
	fn.attr.Mtime = now
	n, err := fn.File.WriteAt(b, off)
	if err != nil {
		fn.upateSize()
	}
	return n, err
}

func (fn *FileNode) Resize(size uint64) error {
	err := fn.File.Resize(size)
	if err != nil {
		fn.upateSize()
	}
	return err
}

type FileNodeList []*FileNode

func (fl FileNodeList) Len() int {
	return len(fl)
}

func (fl FileNodeList) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}

func (fl FileNodeList) Less(i, j int) bool {
	return fl[i].Name() < fl[j].Name()
}

type FileTree struct {
	root       *FileNode
	curInodeID uint64
	nodeIndex  map[uint64]*FileNode
}

func NewFileTree() *FileTree {
	now := uint64(time.Now().Unix())
	attr := fuse.FuseAttr{
		Ino:     fuse.ROOT_INODE_ID,
		Blocks:  1,
		Size:    4096,
		Blksize: _BLOCK_SIZE,
		Atime:   now,
		Mtime:   now,
		Ctime:   now,
		Uid:     _UID,
		Gid:     _GID,
		Mode:    fuse.S_IFDIR | 0755,
		Nlink:   1,
	}
	node := &FileNode{
		File:     rootFile{},
		attr:     attr,
		parent:   nil,
		children: map[uint64]*FileNode{},
	}
	return &FileTree{
		root:       node,
		curInodeID: fuse.ROOT_INODE_ID,
		nodeIndex: map[uint64]*FileNode{
			fuse.ROOT_INODE_ID: node,
		},
	}
}

func (ft *FileTree) GetNode(ino uint64) *FileNode {
	if n, exist := ft.nodeIndex[ino]; exist {
		return n
	} else {
		return nil
	}
}

func (ft *FileTree) NewNode(pIno uint64, f File) (n *FileNode) {
	if parent := ft.GetNode(pIno); parent == nil {
		return nil
	} else {
		now := uint64(time.Now().Unix())

		size := f.Size()
		blocks := size / _BLOCK_SIZE
		if size%_BLOCK_SIZE != 0 {
			blocks += 1
		}

		attr := fuse.FuseAttr{
			Ino:     atomic.AddUint64(&ft.curInodeID, 1),
			Blocks:  blocks,
			Size:    size,
			Blksize: _BLOCK_SIZE,
			Atime:   now,
			Mtime:   now,
			Ctime:   now,
			Uid:     _UID,
			Gid:     _GID,
			Mode:    f.Mode(),
			Nlink:   1,
		}
		n = &FileNode{
			File:     f,
			parent:   parent,
			attr:     attr,
			children: make(map[uint64]*FileNode),
		}
		parent.children[n.attr.Ino] = n
		ft.nodeIndex[n.attr.Ino] = n
		return
	}
}

func (ft *FileTree) GetChildren(ino uint64) FileNodeList {
	if parent := ft.GetNode(ino); parent == nil {
		return nil
	} else {
		parent.attr.Atime = uint64(time.Now().Unix())
		fl := FileNodeList(make([]*FileNode, len(parent.children)))
		i := 0
		for _, node := range parent.children {
			fl[i] = node
			i += 1
		}
		sort.Sort(fl)
		return fl
	}
}
