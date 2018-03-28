package fuse

import "bytes"

type FuseOperations interface {
	GetAttr(
		ctx *FuseRequestContext,
		in *FuseGetattrIn,
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

	Lookup(
		ctx *FuseRequestContext,
		name []byte,
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
}
