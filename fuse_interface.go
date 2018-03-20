package gofuse

/*
#include "fuse_kernel_7_26.h"
*/
import "C"
import "context"

type FuseOperations interface {
	GetAttr(
		ctx context.Context,
		path string,
		attr *C.struct_fuse_attr,
	) (err error)
}
