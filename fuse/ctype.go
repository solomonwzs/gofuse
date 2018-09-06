// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs ctype.go

package fuse

import "syscall"

const (
	FUSE_KERNEL_VERSION       = 7
	FUSE_KERNEL_MINOR_VERSION = 19
)

type (
	OpcodeType       uint32
	OpenOutFlagType  uint32
	FileModeType     uint32
	DirentType       uint32
	SetAttrValidType uint32
	WriteFlagType    uint32
)

const (
	FUSE_LOOKUP       OpcodeType = 0x1
	FUSE_FORGET       OpcodeType = 0x2
	FUSE_GETATTR      OpcodeType = 0x3
	FUSE_SETATTR      OpcodeType = 0x4
	FUSE_READLINK     OpcodeType = 0x5
	FUSE_SYMLINK      OpcodeType = 0x6
	FUSE_MKNOD        OpcodeType = 0x8
	FUSE_MKDIR        OpcodeType = 0x9
	FUSE_UNLINK       OpcodeType = 0xa
	FUSE_RMDIR        OpcodeType = 0xb
	FUSE_RENAME       OpcodeType = 0xc
	FUSE_LINK         OpcodeType = 0xd
	FUSE_OPEN         OpcodeType = 0xe
	FUSE_READ         OpcodeType = 0xf
	FUSE_WRITE        OpcodeType = 0x10
	FUSE_STATFS       OpcodeType = 0x11
	FUSE_RELEASE      OpcodeType = 0x12
	FUSE_FSYNC        OpcodeType = 0x14
	FUSE_SETXATTR     OpcodeType = 0x15
	FUSE_GETXATTR     OpcodeType = 0x16
	FUSE_LISTXATTR    OpcodeType = 0x17
	FUSE_REMOVEXATTR  OpcodeType = 0x18
	FUSE_FLUSH        OpcodeType = 0x19
	FUSE_INIT         OpcodeType = 0x1a
	FUSE_OPENDIR      OpcodeType = 0x1b
	FUSE_READDIR      OpcodeType = 0x1c
	FUSE_RELEASEDIR   OpcodeType = 0x1d
	FUSE_FSYNCDIR     OpcodeType = 0x1e
	FUSE_GETLK        OpcodeType = 0x1f
	FUSE_SETLK        OpcodeType = 0x20
	FUSE_SETLKW       OpcodeType = 0x21
	FUSE_ACCESS       OpcodeType = 0x22
	FUSE_CREATE       OpcodeType = 0x23
	FUSE_INTERRUPT    OpcodeType = 0x24
	FUSE_BMAP         OpcodeType = 0x25
	FUSE_DESTROY      OpcodeType = 0x26
	FUSE_IOCTL        OpcodeType = 0x27
	FUSE_POLL         OpcodeType = 0x28
	FUSE_NOTIFY_REPLY OpcodeType = 0x29
	FUSE_BATCH_FORGET OpcodeType = 0x2a
	FUSE_FALLOCATE    OpcodeType = 0x2b
)

const (
	S_IFMT   FileModeType = 0xf000
	S_IFDIR  FileModeType = 0x4000
	S_IFCHR  FileModeType = 0x2000
	S_IFBLK  FileModeType = 0x6000
	S_IFREG  FileModeType = 0x8000
	S_IFLNK  FileModeType = 0xa000
	S_IFSOCK FileModeType = 0xc000
)

const (
	FOPEN_DIRECT_IO   OpenOutFlagType = 0x1
	FOPEN_KEEP_CACHE  OpenOutFlagType = 0x2
	FOPEN_NONSEEKABLE OpenOutFlagType = 0x4
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
	_SIZEOF_FUSE_IN_HEADER  = 0x28
	_SIZEOF_FUSE_OUT_HEADER = 0x10
	_SIZEOF_FUSE_INIT_IN    = 0x10
	_SIZEOF_FUSE_INIT_OUT   = 0x18
	_SIZEOF_FUSE_ATTR_OUT   = 0x68
	_SIZEOF_FUSE_OPEN_OUT   = 0x10
	_SIZEOF_FUSE_DIRENT     = 0x18
	_SIZEOF_FUSE_ENTRY_OUT  = 0x80
	_SIZEOF_FUSE_WRITE_IN   = 0x28
	_SIZEOF_FUSE_WRITE_OUT  = 0x8
	_SIZEOF_FUSE_MKNOD_IN   = 0x10
	_SIZEOF_FUSE_MKDIR_IN   = 0x8
)

const (
	FATTR_MODE      SetAttrValidType = 0x1
	FATTR_UID       SetAttrValidType = 0x2
	FATTR_GID       SetAttrValidType = 0x4
	FATTR_SIZE      SetAttrValidType = 0x8
	FATTR_ATIME     SetAttrValidType = 0x10
	FATTR_MTIME     SetAttrValidType = 0x20
	FATTR_FH        SetAttrValidType = 0x40
	FATTR_ATIME_NOW SetAttrValidType = 0x80
	FATTR_MTIME_NOW SetAttrValidType = 0x100
	FATTR_LOCKOWNER SetAttrValidType = 0x200
)

const (
	FWRITE_CACHE     WriteFlagType = 0x1
	FWRITE_LOCKOWNER WriteFlagType = 0x2
)

type (
	FuseInHeader struct {
		Len     uint32
		Opcode  OpcodeType
		Unique  uint64
		Nodeid  uint64
		Uid     uint32
		Gid     uint32
		Pid     uint32
		Padding uint32
	}
	FuseOutHeader struct {
		Len    uint32
		Error  int32
		Unique uint64
	}
	FuseInitIn struct {
		Major     uint32
		Minor     uint32
		Readahead uint32
		Flags     uint32
	}
	FuseInitOut struct {
		Major                uint32
		Minor                uint32
		Max_readahead        uint32
		Flags                uint32
		Max_background       uint16
		Congestion_threshold uint16
		Max_write            uint32
	}
	FuseInterruptIn struct {
		Unique uint64
	}
	FuseGetattrIn struct {
		Flags uint32
		Dummy uint32
		Fh    uint64
	}
	FuseAttrOut struct {
		Valid      uint64
		Valid_nsec uint32
		Dummy      uint32
		Attr       FuseAttr
	}
	FuseAttr struct {
		Ino       uint64
		Size      uint64
		Blocks    uint64
		Atime     uint64
		Mtime     uint64
		Ctime     uint64
		Atimensec uint32
		Mtimensec uint32
		Ctimensec uint32
		Mode      FileModeType
		Nlink     uint32
		Uid       uint32
		Gid       uint32
		Rdev      FileModeType
		Blksize   uint32
		Padding   uint32
	}
	FuseOpenIn struct {
		Flags  uint32
		Unused uint32
	}
	FuseOpenOut struct {
		Fh      uint64
		Flags   OpenOutFlagType
		Padding uint32
	}
	FuseReadIn struct {
		Fh         uint64
		Offset     uint64
		Size       uint32
		Read_flags uint32
		Lock_owner uint64
		Flags      uint32
		Padding    uint32
	}
	FuseDirent struct {
		Ino     uint64
		Off     uint64
		Namelen uint32
		Type    DirentType
	}
	FuseEntryOut struct {
		Nodeid           uint64
		Generation       uint64
		Entry_valid      uint64
		Attr_valid       uint64
		Entry_valid_nsec uint32
		Attr_valid_nsec  uint32
		Attr             FuseAttr
	}
	FuseReleaseIn struct {
		Fh            uint64
		Flags         uint32
		Release_flags uint32
		Lock_owner    uint64
	}
	FuseFlushIn struct {
		Fh      uint64
		Unused  uint32
		Padding uint32
		Owner   uint64
	}
	FuseSetAttrIn struct {
		Valid     SetAttrValidType
		Padding   uint32
		Fh        uint64
		Size      uint64
		Owner     uint64
		Atime     uint64
		Mtime     uint64
		Unused2   uint64
		Atimensec uint32
		Mtimensec uint32
		Unused3   uint32
		Mode      FileModeType
		Unused4   uint32
		Uid       uint32
		Gid       uint32
		Unused5   uint32
	}
	FuseWriteIn struct {
		Fh          uint64
		Offset      uint64
		Size        uint32
		Write_flags WriteFlagType
		Lock_owner  uint64
		Flags       uint32
		Padding     uint32
	}
	FuseWriteOut struct {
		Size    uint32
		Padding uint32
	}
	FuseMknodIn struct {
		Mode    FileModeType
		Rdev    uint32
		Umask   uint32
		Padding uint32
	}
	FuseMkdirIn struct {
		Mode  FileModeType
		Umask uint32
	}
)
