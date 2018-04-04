package simplefs

import (
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

type FileSystem struct {
	root      *FileNode
	nodeIndex map[uint64]*FileNode
}

func NewFileSystem(attr fuse.FuseAttr) *FileSystem {
	node := &FileNode{
		Attr:     attr,
		parent:   nil,
		children: []*FileNode{},
		file:     nil,
	}
	return &FileSystem{
		root: node,
		nodeIndex: map[uint64]*FileNode{
			fuse.ROOT_INODE_ID: node,
		},
	}
}
