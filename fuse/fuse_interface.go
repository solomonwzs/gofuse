package fuse

type FuseOperations interface {
	GetAttr(
		ctx *FuseRequestContext,
		attr *FuseAttr,
	) (err error)

	Open(
		ctx *FuseRequestContext,
		open *FuseOpenOut,
	) (err error)
}
