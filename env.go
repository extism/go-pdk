package pdk

type extismPointer uint64

// `extism_input_length` returns the number of bytes provided by the host via its input methods.
//
//go:wasmimport extism:host/env input_length
func extism_input_length() uint64

// `extism_length` returns the number of bytes associated with the block of host memory
// located at `offset`.
//
//go:wasmimport extism:host/env length
func extism_length(offset extismPointer) uint64

//go:wasmimport extism:host/env length_unsafe
func extism_length_unsafe(extismPointer) uint64

// `extism_alloc` allocates `length` bytes of data with host memory for use by the plugin
// and returns its offset within the host memory block.
//
//go:wasmimport extism:host/env alloc
func extism_alloc(length uint64) extismPointer

// `extism_free` releases the bytes previously allocated with `extism_alloc` at the given `offset`.
//
//go:wasmimport extism:host/env free
func extism_free(offset extismPointer)

// `extism_input_load_u8` returns the byte at location `offset` of the "input" data from the host.
//
//go:wasmimport extism:host/env input_load_u8
func extism_input_load_u8_(offset extismPointer) uint32
func extism_input_load_u8(offset extismPointer) uint8 {
	return uint8(extism_input_load_u8_(offset))
}

// `extism_input_load_u64` returns the 64-bit unsigned integer of the "input" data from the host.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env input_load_u64
func extism_input_load_u64(offset extismPointer) uint64

// `extism_output_set` sets the "output" data from the plugin to the host to be the memory that
// has been written at `offset` with the given `length`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env output_set
func extism_output_set(offset extismPointer, length uint64)

// `extism_error_set` sets the "error" data from the plugin to the host to be the memory that
// has been written at `offset`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env error_set
func extism_error_set(offset extismPointer)

// `extism_config_get` returns the host memory block offset for the "config" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env config_get
func extism_config_get(offset extismPointer) extismPointer

// `extism_var_get` returns the host memory block offset for the "var" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env var_get
func extism_var_get(offset extismPointer) extismPointer

// `extism_var_set` sets the host "var" memory keyed by the UTF-8 string located at `offset`
// to be the value which has been previously written at `valueOffset`.
//
// A `valueOffset` of 0 causes the old value associated with this key to be freed on the host
// and the association to be completely removed.
//
// The memory for the key can be immediately freed because the host has its own copy.
// The memory for the value, however, should not be freed, as that erases the value from the host.
//
//go:wasmimport extism:host/env var_set
func extism_var_set(offset, valueOffset extismPointer)

// `extism_store_u8` stores the byte `v` at location `offset` in the host memory block.
//
//go:wasmimport extism:host/env store_u8
func extism_store_u8_(extismPointer, uint32)
func extism_store_u8(offset extismPointer, v uint8) {
	extism_store_u8_(offset, uint32(v))
}

// `extism_load_u8` returns the byte located at `offset` in the host memory block.
//
//go:wasmimport extism:host/env load_u8
func extism_load_u8_(offset extismPointer) uint32
func extism_load_u8(offset extismPointer) uint8 {
	return uint8(extism_load_u8_(offset))
}

// `extism_store_u64` stores the 64-bit unsigned integer value `v` at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env store_u64
func extism_store_u64(offset extismPointer, v uint64)

// `extism_load_u64` returns the 64-bit unsigned integer at location `offset` in the host memory block.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env load_u64
func extism_load_u64(offset extismPointer) uint64

// `extism_http_request` sends the HTTP `request` to the Extism host with the provided `body` (0 means no body)
// and returns back the memory offset to the response body.
//
//go:wasmimport extism:host/env http_request
func extism_http_request(request, body extismPointer) extismPointer

// `extism_http_status_code` returns the status code for the last-sent `extism_http_request` call.
//
//go:wasmimport extism:host/env http_status_code
func extism_http_status_code() int32

// `extism_log_info` logs an "info" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_info
func extism_log_info(offset extismPointer)

// `extism_log_debug` logs a "debug" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_debug
func extism_log_debug(offset extismPointer)

// `extism_log_warn` logs a "warning" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_warn
func extism_log_warn(offset extismPointer)

// `extism_log_error` logs an "error" string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log_error
func extism_log_error(offset extismPointer)
