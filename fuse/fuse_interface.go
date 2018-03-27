package fuse

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

	Read(
		ctx *FuseRequestContext,
		in *FuseReadIn,
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
}
