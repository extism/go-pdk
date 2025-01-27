package memory

// extismStoreU8 stores the byte `v` at location `offset` in the host memory block.
//
//go:wasmimport extism:host/env store_u8
func extismStoreU8_(ExtismPointer, uint32)
func ExtismStoreU8(offset ExtismPointer, v uint8) {
	extismStoreU8_(offset, uint32(v))
}

// extismLoadU8 returns the byte located at `offset` in the host memory block.
//
//go:wasmimport extism:host/env load_u8
func extismLoadU8_(offset ExtismPointer) uint32
func ExtismLoadU8(offset ExtismPointer) uint8 {
	return uint8(extismLoadU8_(offset))
}

// extismStoreU64 stores the 64-bit unsigned integer value `v` at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env store_u64
func ExtismStoreU64(offset ExtismPointer, v uint64)

// extismLoadU64 returns the 64-bit unsigned integer at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env load_u64
func ExtismLoadU64(offset ExtismPointer) uint64

//go:wasmimport extism:host/env length_unsafe
func ExtismLengthUnsafe(ExtismPointer) uint64

// extismLength returns the number of bytes associated with the block of host memory
// located at `offset`.
//
//go:wasmimport extism:host/env length
func ExtismLength(offset ExtismPointer) uint64

// extismAlloc allocates `length` bytes of data with host memory for use by the plugin
// and returns its offset within the host memory block.
//
//go:wasmimport extism:host/env alloc
func ExtismAlloc(length uint64) ExtismPointer

// extismFree releases the bytes previously allocated with `extism_alloc` at the given `offset`.
//
//go:wasmimport extism:host/env free
func ExtismFree(offset ExtismPointer)
