package gofuse

/*
#include <errno.h>
#include <sys/stat.h>
#include "fuse_kernel_7_26.h"

#define SIZEOF_FUSE_IN_HEADER sizeof(struct fuse_in_header)
#define SIZEOF_FUSE_OUT_HEADER sizeof(struct fuse_out_header)
#define SIZEOF_FUSE_INIT_OUT sizeof(struct fuse_init_out)
#define SIZEOF_FUSE_INIT_IN sizeof(struct fuse_init_in)
#define SIZEOF_FUSE_ATTR_OUT sizeof(struct fuse_attr_out)
*/
import "C"

const (
	FUSE_LOOKUP       = C.FUSE_LOOKUP
	FUSE_FORGET       = C.FUSE_FORGET
	FUSE_GETATTR      = C.FUSE_GETATTR
	FUSE_SETATTR      = C.FUSE_SETATTR
	FUSE_READLINK     = C.FUSE_READLINK
	FUSE_SYMLINK      = C.FUSE_SYMLINK
	FUSE_MKNOD        = C.FUSE_MKNOD
	FUSE_MKDIR        = C.FUSE_MKDIR
	FUSE_UNLINK       = C.FUSE_UNLINK
	FUSE_RMDIR        = C.FUSE_RMDIR
	FUSE_RENAME       = C.FUSE_RENAME
	FUSE_LINK         = C.FUSE_LINK
	FUSE_OPEN         = C.FUSE_OPEN
	FUSE_READ         = C.FUSE_READ
	FUSE_WRITE        = C.FUSE_WRITE
	FUSE_STATFS       = C.FUSE_STATFS
	FUSE_RELEASE      = C.FUSE_RELEASE
	FUSE_FSYNC        = C.FUSE_FSYNC
	FUSE_SETXATTR     = C.FUSE_SETXATTR
	FUSE_GETXATTR     = C.FUSE_GETXATTR
	FUSE_LISTXATTR    = C.FUSE_LISTXATTR
	FUSE_REMOVEXATTR  = C.FUSE_REMOVEXATTR
	FUSE_FLUSH        = C.FUSE_FLUSH
	FUSE_INIT         = C.FUSE_INIT
	FUSE_OPENDIR      = C.FUSE_OPENDIR
	FUSE_READDIR      = C.FUSE_READDIR
	FUSE_RELEASEDIR   = C.FUSE_RELEASEDIR
	FUSE_FSYNCDIR     = C.FUSE_FSYNCDIR
	FUSE_GETLK        = C.FUSE_GETLK
	FUSE_SETLK        = C.FUSE_SETLK
	FUSE_SETLKW       = C.FUSE_SETLKW
	FUSE_ACCESS       = C.FUSE_ACCESS
	FUSE_CREATE       = C.FUSE_CREATE
	FUSE_INTERRUPT    = C.FUSE_INTERRUPT
	FUSE_BMAP         = C.FUSE_BMAP
	FUSE_DESTROY      = C.FUSE_DESTROY
	FUSE_IOCTL        = C.FUSE_IOCTL
	FUSE_POLL         = C.FUSE_POLL
	FUSE_NOTIFY_REPLY = C.FUSE_NOTIFY_REPLY
	FUSE_BATCH_FORGET = C.FUSE_BATCH_FORGET
	FUSE_FALLOCATE    = C.FUSE_FALLOCATE
	FUSE_READDIRPLUS  = C.FUSE_READDIRPLUS
	FUSE_RENAME2      = C.FUSE_RENAME2
	FUSE_LSEEK        = C.FUSE_LSEEK
)

const (
	S_IFDIR = C.S_IFDIR
)

const (
	_SIZEOF_FUSE_IN_HEADER  = C.SIZEOF_FUSE_IN_HEADER
	_SIZEOF_FUSE_OUT_HEADER = C.SIZEOF_FUSE_OUT_HEADER

	_SIZEOF_FUSE_INIT_IN  = C.SIZEOF_FUSE_INIT_IN
	_SIZEOF_FUSE_INIT_OUT = C.SIZEOF_FUSE_INIT_OUT

	_SIZEOF_FUSE_ATTR_OUT = C.SIZEOF_FUSE_ATTR_OUT
)

type (
	FuseInHeader    = C.struct_fuse_in_header
	FuseOutHeader   = C.struct_fuse_out_header
	FuseInitIn      = C.struct_fuse_init_in
	FuseInitOut     = C.struct_fuse_init_out
	FuseInterruptIn = C.struct_fuse_interrupt_in
	FuseGetattrIn   = C.struct_fuse_getattr_in
	FuseAttrOut     = C.struct_fuse_attr_out
	FuseAttr        = C.struct_fuse_attr
)
