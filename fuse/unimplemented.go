package fuse

type FileSystemUnimplemented struct{}

func (fs FileSystemUnimplemented) GetAttr(
	ctx *FuseRequestContext,
	in *FuseGetattrIn,
	out *FuseAttrOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Open(
	ctx *FuseRequestContext,
	in *FuseOpenIn,
	out *FuseOpenOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Read(
	ctx *FuseRequestContext,
	in *FuseReadIn,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Lookup(
	ctx *FuseRequestContext,
	name []byte,
	out *FuseEntryOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Release(
	ctx *FuseRequestContext,
	in *FuseReleaseIn,
) (err error) {
	return nil
}
