package fuse

import "unsafe"

type DirentRaw []byte

func NewFuseDirentRaw(ino uint64, offset uint64, typ DirentType,
	name []byte) (raw DirentRaw) {
	size := padding64bits(_SIZEOF_FUSE_DIRENT + len(name))
	raw = make([]byte, size, size)

	dirent := (*FuseDirent)(unsafe.Pointer(&raw[0]))
	dirent.Ino = ino
	dirent.Off = offset
	dirent.Namelen = uint32(len(name))
	copy(raw[_SIZEOF_FUSE_DIRENT:], name)

	return
}
