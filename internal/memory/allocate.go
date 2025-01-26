package memory

// Allocate allocates `length` uninitialized bytes on the host.
func Allocate(length int) Memory {
	clength := uint64(length)
	offset := ExtismAlloc(clength)

	return NewMemory(offset, clength)
}

// AllocateBytes allocates and saves the `data` into Memory on the host.
func AllocateBytes(data []byte) Memory {
	clength := uint64(len(data))
	offset := ExtismAlloc(clength)

	Store(offset, data)

	return NewMemory(offset, clength)
}
