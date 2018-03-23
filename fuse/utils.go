package fuse

func padding64bits(size int) int {
	if size&7 == 0 {
		return size
	} else {
		return ((size >> 3) + 1) << 3
	}
}
