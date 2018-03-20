package gofuse

import (
	"syscall"
)

const (
	EEXIST    = syscall.EEXIST
	EINVAL    = syscall.EINVAL
	EIO       = syscall.EIO
	ENOATTR   = syscall.ENODATA
	ENOENT    = syscall.ENOENT
	ENOSYS    = syscall.ENOSYS
	ENOTDIR   = syscall.ENOTDIR
	ENOTEMPTY = syscall.ENOTEMPTY
)
