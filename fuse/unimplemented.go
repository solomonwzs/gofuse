package fuse

import "bytes"

type FileSystemUnimplemented struct{}

func (fs FileSystemUnimplemented) GetAttr(
	ctx *FuseRequestContext,
	in *FuseGetattrIn,
	out *FuseAttrOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) SetAttr(
	ctx *FuseRequestContext,
	in *FuseSetAttrIn,
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
	out *bytes.Buffer,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Write(
	ctx *FuseRequestContext,
	in *FuseWriteIn,
	inRaw []byte,
	out *FuseWriteOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) ReadDir(
	ctx *FuseRequestContext,
	in *FuseReadIn,
	out *FuseReadDirOut,
) (err error) {
	return ENOSYS
}

func (fs FileSystemUnimplemented) Lookup(
	ctx *FuseRequestContext,
	inName []byte,
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

func (fs FileSystemUnimplemented) Flush(
	ctx *FuseRequestContext,
	in *FuseFlushIn,
) (err error) {
	return nil
}

func (fs FileSystemUnimplemented) Mknod(
	ctx *FuseRequestContext,
	in *FuseMknodIn,
	inName []byte,
	out *FuseEntryOut,
) (err error) {
	return ENOSYS
}
