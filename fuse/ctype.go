// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs ctype.go

package fuse

const (
	FUSE_LOOKUP		= 0x1
	FUSE_FORGET		= 0x2
	FUSE_GETATTR		= 0x3
	FUSE_SETATTR		= 0x4
	FUSE_READLINK		= 0x5
	FUSE_SYMLINK		= 0x6
	FUSE_MKNOD		= 0x8
	FUSE_MKDIR		= 0x9
	FUSE_UNLINK		= 0xa
	FUSE_RMDIR		= 0xb
	FUSE_RENAME		= 0xc
	FUSE_LINK		= 0xd
	FUSE_OPEN		= 0xe
	FUSE_READ		= 0xf
	FUSE_WRITE		= 0x10
	FUSE_STATFS		= 0x11
	FUSE_RELEASE		= 0x12
	FUSE_FSYNC		= 0x14
	FUSE_SETXATTR		= 0x15
	FUSE_GETXATTR		= 0x16
	FUSE_LISTXATTR		= 0x17
	FUSE_REMOVEXATTR	= 0x18
	FUSE_FLUSH		= 0x19
	FUSE_INIT		= 0x1a
	FUSE_OPENDIR		= 0x1b
	FUSE_READDIR		= 0x1c
	FUSE_RELEASEDIR		= 0x1d
	FUSE_FSYNCDIR		= 0x1e
	FUSE_GETLK		= 0x1f
	FUSE_SETLK		= 0x20
	FUSE_SETLKW		= 0x21
	FUSE_ACCESS		= 0x22
	FUSE_CREATE		= 0x23
	FUSE_INTERRUPT		= 0x24
	FUSE_BMAP		= 0x25
	FUSE_DESTROY		= 0x26
	FUSE_IOCTL		= 0x27
	FUSE_POLL		= 0x28
	FUSE_NOTIFY_REPLY	= 0x29
	FUSE_BATCH_FORGET	= 0x2a
	FUSE_FALLOCATE		= 0x2b
	FUSE_READDIRPLUS	= 0x2c
	FUSE_RENAME2		= 0x2d
	FUSE_LSEEK		= 0x2e
)

const (
	S_IFDIR	= 0x4000

	FOPEN_DIRECT_IO		= 0x1
	FOPEN_KEEP_CACHE	= 0x2
	FOPEN_NONSEEKABLE	= 0x4
)

const (
	_SIZEOF_FUSE_IN_HEADER	= 0x28
	_SIZEOF_FUSE_OUT_HEADER	= 0x10

	_SIZEOF_FUSE_INIT_IN	= 0x10
	_SIZEOF_FUSE_INIT_OUT	= 0x40

	_SIZEOF_FUSE_ATTR_OUT	= 0x68

	_SIZEOF_FUSE_OPEN_OUT	= 0x10
)

type (
	FuseInHeader	= struct {
		Len	uint32
		Opcode	uint32
		Unique	uint64
		Nodeid	uint64
		Uid	uint32
		Gid	uint32
		Pid	uint32
		Padding	uint32
	}
	FuseOutHeader	= struct {
		Len	uint32
		Error	int32
		Unique	uint64
	}
	FuseInitIn	= struct {
		Major		uint32
		Minor		uint32
		Readahead	uint32
		Flags		uint32
	}
	FuseInitOut	= struct {
		Major			uint32
		Minor			uint32
		Max_readahead		uint32
		Flags			uint32
		Max_background		uint16
		Congestion_threshold	uint16
		Max_write		uint32
		Time_gran		uint32
		Unused			[9]uint32
	}
	FuseInterruptIn	= struct {
		Unique uint64
	}
	FuseGetattrIn	= struct {
		Flags	uint32
		Dummy	uint32
		Fh	uint64
	}
	FuseAttrOut	= struct {
		Valid		uint64
		Valid_nsec	uint32
		Dummy		uint32
		Attr		FuseAttr
	}
	FuseAttr	= struct {
		Ino		uint64
		Size		uint64
		Blocks		uint64
		Atime		uint64
		Mtime		uint64
		Ctime		uint64
		Atimensec	uint32
		Mtimensec	uint32
		Ctimensec	uint32
		Mode		uint32
		Nlink		uint32
		Uid		uint32
		Gid		uint32
		Rdev		uint32
		Blksize		uint32
		Padding		uint32
	}
	FuseOpenIn	= struct {
		Flags	uint32
		Unused	uint32
	}
	FuseOpenOut	= struct {
		Fh	uint64
		Flags	uint32
		Padding	uint32
	}
)
