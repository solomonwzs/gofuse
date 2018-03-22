package fuse

type FileSystemUnimplemented struct{}

func (fs FileSystemUnimplemented) GetAttr(
	ctx *FuseRequestContext,
	attr *FuseAttr,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Open(
	ctx *FuseRequestContext,
	open *FuseOpenOut,
) (err error) {
	return ENOSYS
}
