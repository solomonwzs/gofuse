package fuse

/*
#include <errno.h>
#include <sys/stat.h>
#include "fuse_kernel_7_26.h"

#define SIZEOF_FUSE_IN_HEADER sizeof(struct fuse_in_header)
#define SIZEOF_FUSE_OUT_HEADER sizeof(struct fuse_out_header)
#define SIZEOF_FUSE_INIT_OUT sizeof(struct fuse_init_out)
#define SIZEOF_FUSE_INIT_IN sizeof(struct fuse_init_in)
#define SIZEOF_FUSE_ATTR_OUT sizeof(struct fuse_attr_out)
#define SIZEOF_FUSE_OPEN_OUT sizeof(struct fuse_open_out)
#define SIZEOF_FUSE_DIRENT sizeof(struct fuse_dirent)
#define SIZEOF_FUSE_ENTRY_OUT sizeof(struct fuse_entry_out)
*/
import "C"
import "syscall"

type (
	OpcodeType      uint32
	OpenOutFlagType uint32
	FileModeType    uint32
	DirentType      uint32
)

const (
	FUSE_LOOKUP       OpcodeType = C.FUSE_LOOKUP
	FUSE_FORGET       OpcodeType = C.FUSE_FORGET
	FUSE_GETATTR      OpcodeType = C.FUSE_GETATTR
	FUSE_SETATTR      OpcodeType = C.FUSE_SETATTR
	FUSE_READLINK     OpcodeType = C.FUSE_READLINK
	FUSE_SYMLINK      OpcodeType = C.FUSE_SYMLINK
	FUSE_MKNOD        OpcodeType = C.FUSE_MKNOD
	FUSE_MKDIR        OpcodeType = C.FUSE_MKDIR
	FUSE_UNLINK       OpcodeType = C.FUSE_UNLINK
	FUSE_RMDIR        OpcodeType = C.FUSE_RMDIR
	FUSE_RENAME       OpcodeType = C.FUSE_RENAME
	FUSE_LINK         OpcodeType = C.FUSE_LINK
	FUSE_OPEN         OpcodeType = C.FUSE_OPEN
	FUSE_READ         OpcodeType = C.FUSE_READ
	FUSE_WRITE        OpcodeType = C.FUSE_WRITE
	FUSE_STATFS       OpcodeType = C.FUSE_STATFS
	FUSE_RELEASE      OpcodeType = C.FUSE_RELEASE
	FUSE_FSYNC        OpcodeType = C.FUSE_FSYNC
	FUSE_SETXATTR     OpcodeType = C.FUSE_SETXATTR
	FUSE_GETXATTR     OpcodeType = C.FUSE_GETXATTR
	FUSE_LISTXATTR    OpcodeType = C.FUSE_LISTXATTR
	FUSE_REMOVEXATTR  OpcodeType = C.FUSE_REMOVEXATTR
	FUSE_FLUSH        OpcodeType = C.FUSE_FLUSH
	FUSE_INIT         OpcodeType = C.FUSE_INIT
	FUSE_OPENDIR      OpcodeType = C.FUSE_OPENDIR
	FUSE_READDIR      OpcodeType = C.FUSE_READDIR
	FUSE_RELEASEDIR   OpcodeType = C.FUSE_RELEASEDIR
	FUSE_FSYNCDIR     OpcodeType = C.FUSE_FSYNCDIR
	FUSE_GETLK        OpcodeType = C.FUSE_GETLK
	FUSE_SETLK        OpcodeType = C.FUSE_SETLK
	FUSE_SETLKW       OpcodeType = C.FUSE_SETLKW
	FUSE_ACCESS       OpcodeType = C.FUSE_ACCESS
	FUSE_CREATE       OpcodeType = C.FUSE_CREATE
	FUSE_INTERRUPT    OpcodeType = C.FUSE_INTERRUPT
	FUSE_BMAP         OpcodeType = C.FUSE_BMAP
	FUSE_DESTROY      OpcodeType = C.FUSE_DESTROY
	FUSE_IOCTL        OpcodeType = C.FUSE_IOCTL
	FUSE_POLL         OpcodeType = C.FUSE_POLL
	FUSE_NOTIFY_REPLY OpcodeType = C.FUSE_NOTIFY_REPLY
	FUSE_BATCH_FORGET OpcodeType = C.FUSE_BATCH_FORGET
	FUSE_FALLOCATE    OpcodeType = C.FUSE_FALLOCATE
	FUSE_READDIRPLUS  OpcodeType = C.FUSE_READDIRPLUS
	FUSE_RENAME2      OpcodeType = C.FUSE_RENAME2
	FUSE_LSEEK        OpcodeType = C.FUSE_LSEEK
)

const (
	S_IFMT   FileModeType = C.S_IFMT
	S_IFDIR  FileModeType = C.S_IFDIR
	S_IFCHR  FileModeType = C.S_IFCHR
	S_IFBLK  FileModeType = C.S_IFBLK
	S_IFREG  FileModeType = C.S_IFREG
	S_IFLNK  FileModeType = C.S_IFLNK
	S_IFSOCK FileModeType = C.S_IFSOCK
)

const (
	FOPEN_DIRECT_IO   OpenOutFlagType = C.FOPEN_DIRECT_IO
	FOPEN_KEEP_CACHE  OpenOutFlagType = C.FOPEN_KEEP_CACHE
	FOPEN_NONSEEKABLE OpenOutFlagType = C.FOPEN_NONSEEKABLE
)

const (
	DT_SOCK DirentType = syscall.DT_SOCK
	DT_LNK  DirentType = syscall.DT_LNK
	DT_REG  DirentType = syscall.DT_REG
	DT_BLK  DirentType = syscall.DT_BLK
	DT_DIR  DirentType = syscall.DT_DIR
	DT_CHR  DirentType = syscall.DT_CHR
	DT_FIFO DirentType = syscall.DT_FIFO
)

const (
	_SIZEOF_FUSE_IN_HEADER  = C.SIZEOF_FUSE_IN_HEADER
	_SIZEOF_FUSE_OUT_HEADER = C.SIZEOF_FUSE_OUT_HEADER
	_SIZEOF_FUSE_INIT_IN    = C.SIZEOF_FUSE_INIT_IN
	_SIZEOF_FUSE_INIT_OUT   = C.SIZEOF_FUSE_INIT_OUT
	_SIZEOF_FUSE_ATTR_OUT   = C.SIZEOF_FUSE_ATTR_OUT
	_SIZEOF_FUSE_OPEN_OUT   = C.SIZEOF_FUSE_OPEN_OUT
	_SIZEOF_FUSE_DIRENT     = C.SIZEOF_FUSE_DIRENT
	_SIZEOF_FUSE_ENTRY_OUT  = C.SIZEOF_FUSE_ENTRY_OUT
)

type (
	FuseInHeader    C.struct_fuse_in_header
	FuseOutHeader   C.struct_fuse_out_header
	FuseInitIn      C.struct_fuse_init_in
	FuseInitOut     C.struct_fuse_init_out
	FuseInterruptIn C.struct_fuse_interrupt_in
	FuseGetattrIn   C.struct_fuse_getattr_in
	FuseAttrOut     C.struct_fuse_attr_out
	FuseAttr        C.struct_fuse_attr
	FuseOpenIn      C.struct_fuse_open_in
	FuseOpenOut     C.struct_fuse_open_out
	FuseReadIn      C.struct_fuse_read_in
	FuseDirent      C.struct_fuse_dirent
	FuseEntryOut    C.struct_fuse_entry_out
	FuseReleaseIn   C.struct_fuse_release_in
)
