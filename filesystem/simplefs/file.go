package simplefs

import (
	"sync/atomic"
	"time"

	"github.com/solomonwzs/gofuse/fuse"
)

const (
	_BLOCK_SIZE = 4096
)

type File interface {
	Name() (name string)
	ReadAt(b []byte, off int64) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	Mode() fuse.FileModeType
	Size() uint64
}

type FileNode struct {
	Attr     fuse.FuseAttr
	parent   *FileNode
	children map[uint64]*FileNode
	file     File
}

type FileNodeList []*FileNode

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
		Attr:     attr,
		parent:   nil,
		children: map[uint64]*FileNode{},
		file:     nil,
	}
	return &FileTree{
		root:       node,
		curInodeID: fuse.ROOT_INODE_ID,
		nodeIndex: map[uint64]*FileNode{
			fuse.ROOT_INODE_ID: node,
		},
	}
}

func (fl FileNodeList) Len() int {
	return len(fl)
}

func (fl FileNodeList) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}

func (fl FileNodeList) Less(i, j int) bool {
	return fl[i].file.Name() < fl[j].file.Name()
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
			parent: parent,
			Attr:   attr,
		}
		parent.children[n.Attr.Ino] = n
		ft.nodeIndex[n.Attr.Ino] = n
		return
	}
}

func (ft *FileTree) GetChildren(ino uint64) FileNodeList {
	if parent := ft.GetNode(ino); parent == nil {
		return nil
	} else {
		fl := FileNodeList(make([]*FileNode, len(parent.children)))
		return fl
	}
}
