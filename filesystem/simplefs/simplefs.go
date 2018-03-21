package simplefs

/*
#cgo CFLAGS: -I${SRCDIR}/../../c

#include <sys/stat.h>
#include "fuse_kernel_7_26.h"
*/
import "C"
import (
	"github.com/solomonwzs/gofuse"
)

type SimpleFS struct {
	gofuse.FileSystemUnimplemented
}

// func (fs SimpleFS) GetAttr(
// 	ctx *gofuse.FuseRequestContext,
// 	path string,
// 	attr *C.struct_fuse_attr,
// ) (err error) {
// 	attr.ino = 0
// 	attr.size = 4096
// 	attr.blocks = 0
// 	attr.atime = C.uint64_t(time.Now().Unix())
// 	attr.mtime = C.uint64_t(time.Now().Unix())
// 	attr.ctime = C.uint64_t(time.Now().Unix())
// 	attr.mode = C.S_IFDIR | 0755
// 	attr.nlink = 1
// 	attr.blksize = 4096
// 	return
// }
