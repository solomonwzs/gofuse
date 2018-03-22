package gofuse

type FuseOperations interface {
	GetAttr(
		ctx *FuseRequestContext,
		path string,
		attr *FuseAttr,
	) (err error)
}
