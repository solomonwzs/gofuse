package simplefs

import (
	"sync/atomic"

	"github.com/solomonwzs/gofuse/fuse"
)

type File interface {
	Name() (name string)
	ReadAt(b []byte, off int64) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
}

type FileNode struct {
	Attr     fuse.FuseAttr
	parent   *FileNode
	children []*FileNode
	file     File
}

type FileTree struct {
	root       *FileNode
	curInodeID uint64
	nodeIndex  map[uint64]*FileNode
}

func NewFileTree(attr fuse.FuseAttr) *FileTree {
	node := &FileNode{
		Attr:     attr,
		parent:   nil,
		children: []*FileNode{},
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

func (ft *FileTree) GetNode(ino uint64) *FileNode {
	if n, exist := ft.nodeIndex[ino]; exist {
		return n
	} else {
		return nil
	}
}

func (ft *FileTree) NewNode() (n *FileNode) {
	n = &FileNode{
		Attr: fuse.FuseAttr{
			Ino: atomic.AddUint64(&ft.curInodeID, 1),
		},
	}
	return
}
