package gofuse

type FileSystemUnimplemented struct{}

func (fs FileSystemUnimplemented) GetAttr(
	ctx *FuseRequestContext,
	path string,
	attr *FuseAttr,
) (err error) {
	return ENOSYS
}
