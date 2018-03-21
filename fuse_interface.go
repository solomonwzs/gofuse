package gofuse

/*
#include "fuse_kernel_7_26.h"
*/
import "C"

type FuseOperations interface {
	GetAttr(
		ctx *FuseRequestContext,
		path string,
		attr *C.struct_fuse_attr,
	) (err error)
}
