package pdk

type extismPointer uint64

// `extismInputLength` returns the number of bytes provided by the host via its input methods.
//
//go:wasmimport extism:host/env input_length
func extismInputLength() uint64

// `extismLength` returns the number of bytes associated with the block of host memory
// located at `offset`.
//
//go:wasmimport extism:host/env length
func extismLength(offset extismPointer) uint64

//go:wasmimport extism:host/env length_unsafe
func extismLengthUnsafe(extismPointer) uint64

// `extismAlloc` allocates `length` bytes of data with host memory for use by the plugin
// and returns its offset within the host memory block.
//
//go:wasmimport extism:host/env alloc
func extismAlloc(length uint64) extismPointer

// `extismFree` releases the bytes previously allocated with `extism_alloc` at the given `offset`.
//
//go:wasmimport extism:host/env free
func extismFree(offset extismPointer)

// `extismInputLoadU8` returns the byte at location `offset` of the "input" data from the host.
//
//go:wasmimport extism:host/env input_load_u8
func extismInputLoadU8_(offset extismPointer) uint32
func extismInputLoadU8(offset extismPointer) uint8 {
	return uint8(extismInputLoadU8_(offset))
}

// `extismInputLoadU64` returns the 64-bit unsigned integer of the "input" data from the host.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env input_load_u64
func extismInputLoadU64(offset extismPointer) uint64

// `extismOutputSet` sets the "output" data from the plugin to the host to be the memory that
// has been written at `offset` with the given `length`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env output_set
func extismOutputSet(offset extismPointer, length uint64)

// `extismErrorSet` sets the "error" data from the plugin to the host to be the memory that
// has been written at `offset`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env error_set
func extismErrorSet(offset extismPointer)

// `extismConfigGet` returns the host memory block offset for the "config" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env config_get
func extismConfigGet(offset extismPointer) extismPointer

// `extismVarGet` returns the host memory block offset for the "var" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env var_get
func extismVarGet(offset extismPointer) extismPointer

// `extismVarSet` sets the host "var" memory keyed by the UTF-8 string located at `offset`
// to be the value which has been previously written at `valueOffset`.
//
// A `valueOffset` of 0 causes the old value associated with this key to be freed on the host
// and the association to be completely removed.
//
// The memory for the key can be immediately freed because the host has its own copy.
// The memory for the value, however, should not be freed, as that erases the value from the host.
//
//go:wasmimport extism:host/env var_set
func extismVarSet(offset, valueOffset extismPointer)

// `extismStoreU8` stores the byte `v` at location `offset` in the host memory block.
//
//go:wasmimport extism:host/env store_u8
func extismStoreU8_(extismPointer, uint32)
func extismStoreU8(offset extismPointer, v uint8) {
	extismStoreU8_(offset, uint32(v))
}

// `extismLoadU8` returns the byte located at `offset` in the host memory block.
//
//go:wasmimport extism:host/env load_u8
func extismLoadU8_(offset extismPointer) uint32
func extismLoadU8(offset extismPointer) uint8 {
	return uint8(extismLoadU8_(offset))
}

// `extismStoreU64` stores the 64-bit unsigned integer value `v` at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env store_u64
func extismStoreU64(offset extismPointer, v uint64)

// `extismLoadU64` returns the 64-bit unsigned integer at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env load_u64
func extismLoadU64(offset extismPointer) uint64

// `extismHTTPRequest` sends the HTTP `request` to the Extism host with the provided `body` (0 means no body)
// and returns back the memory offset to the response body.
//
//go:wasmimport extism:host/env http_request
func extismHTTPRequest(request, body extismPointer) extismPointer

// `extismHTTPStatusCode` returns the status code for the last-sent `extism_http_request` call.
//
//go:wasmimport extism:host/env http_status_code
func extismHTTPStatusCode() int32

// `extismLogInfo` logs an "info" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_info
func extismLogInfo(offset extismPointer)

// `extismLogDebug` logs a "debug" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_debug
func extismLogDebug(offset extismPointer)

// `extismLogWarn` logs a "warning" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_warn
func extismLogWarn(offset extismPointer)

// `extismLogError` logs an "error" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_error
func extismLogError(offset extismPointer)
