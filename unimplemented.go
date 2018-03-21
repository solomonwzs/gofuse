package gofuse

/*
#include "fuse_kernel_7_26.h"
*/
import "C"

type FileSystemUnimplemented struct{}

func (fs FileSystemUnimplemented) GetAttr(
	ctx *FuseRequestContext,
	path string,
	attr *C.struct_fuse_attr,
) (err error) {
	return ENOSYS
}
