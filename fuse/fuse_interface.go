package fuse

import "bytes"

type FuseReadDirOut []DirentRaw

func (out *FuseReadDirOut) AddDirentRaw(raw DirentRaw) {
	*out = append(*out, raw)
}

func (out *FuseReadDirOut) raw(n uint32) []byte {
	if n > uint32(len(*out)) {
		n = uint32(len(*out))
	}
	buf := new(bytes.Buffer)
	for i := uint32(0); i < n; i++ {
		buf.Write((*out)[i])
	}
	return buf.Bytes()
}

type FuseOperations interface {
	GetAttr(
		ctx *FuseRequestContext,
		in *FuseGetattrIn,
		out *FuseAttrOut,
	) (err error)

	SetAttr(
		ctx *FuseRequestContext,
		in *FuseSetAttrIn,
		out *FuseAttrOut,
	) (err error)

	Open(
		ctx *FuseRequestContext,
		in *FuseOpenIn,
		out *FuseOpenOut,
	) (err error)

	ReadDir(
		ctx *FuseRequestContext,
		in *FuseReadIn,
		out *FuseReadDirOut,
	) (err error)

	Read(
		ctx *FuseRequestContext,
		in *FuseReadIn,
		out *bytes.Buffer,
	) (err error)

	Write(
		ctx *FuseRequestContext,
		in *FuseWriteIn,
		inRaw []byte,
		out *FuseWriteOut,
	) (err error)

	Lookup(
		ctx *FuseRequestContext,
		inName []byte,
		out *FuseEntryOut,
	) (err error)

	Release(
		ctx *FuseRequestContext,
		in *FuseReleaseIn,
	) (err error)

	Flush(
		ctx *FuseRequestContext,
		in *FuseFlushIn,
	) (err error)

	Mknod(
		ctx *FuseRequestContext,
		in *FuseMknodIn,
		inName []byte,
		out *FuseEntryOut,
	) (err error)

	Mkdir(
		ctx *FuseRequestContext,
		in *FuseMkdirIn,
		inName []byte,
		out *FuseEntryOut,
	) (err error)
}
