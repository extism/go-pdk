package pdk

import (
	"github.com/extism/go-pdk/internal/memory"
)

// extismInputLength returns the number of bytes provided by the host via its input methods.
//
//go:wasmimport extism:host/env input_length
func extismInputLength() uint64

// extismInputLoadU8 returns the byte at location `offset` of the "input" data from the host.
//
//go:wasmimport extism:host/env input_load_u8
func extismInputLoadU8_(offset memory.ExtismPointer) uint32
func extismInputLoadU8(offset memory.ExtismPointer) uint8 {
	return uint8(extismInputLoadU8_(offset))
}

// extismInputLoadU64 returns the 64-bit unsigned integer of the "input" data from the host.
// Note that `offset` must lie on an 8-byte boundary.
//
//go:wasmimport extism:host/env input_load_u64
func extismInputLoadU64(offset memory.ExtismPointer) uint64

// extismOutputSet sets the "output" data from the plugin to the host to be the memory that
// has been written at `offset` with the given `length`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env output_set
func extismOutputSet(offset memory.ExtismPointer, length uint64)

// extismErrorSet sets the "error" data from the plugin to the host to be the memory that
// has been written at `offset`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env error_set
func extismErrorSet(offset memory.ExtismPointer)

// extismConfigGet returns the host memory block offset for the "config" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env config_get
func extismConfigGet(offset memory.ExtismPointer) memory.ExtismPointer

// extismVarGet returns the host memory block offset for the "var" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env var_get
func extismVarGet(offset memory.ExtismPointer) memory.ExtismPointer

// extismVarSet sets the host "var" memory keyed by the UTF-8 string located at `offset`
// to be the value which has been previously written at `valueOffset`.
//
// A `valueOffset` of 0 causes the old value associated with this key to be freed on the host
// and the association to be completely removed.
//
// The memory for the key can be immediately freed because the host has its own copy.
// The memory for the value, however, should not be freed, as that erases the value from the host.
//
//go:wasmimport extism:host/env var_set
func extismVarSet(offset, valueOffset memory.ExtismPointer)

// extismLogInfo logs an "info" string to the host from the previously-written UTF-8 string written to `offset`.
//
//go:wasmimport extism:host/env log_info
func extismLogInfo(offset memory.ExtismPointer)

// extismLogDebug logs a "debug" string to the host from the previously-written UTF-8 string written to `offset`.
//
//go:wasmimport extism:host/env log_debug
func extismLogDebug(offset memory.ExtismPointer)

// extismLogWarn logs a "warning" string to the host from the previously-written UTF-8 string written to `offset`.
//
//go:wasmimport extism:host/env log_warn
func extismLogWarn(offset memory.ExtismPointer)

// extismLogError logs an "error" string to the host from the previously-written UTF-8 string written to `offset`.
//
//go:wasmimport extism:host/env log_error
func extismLogError(offset memory.ExtismPointer)

// extismLogTrace logs an "error" string to the host from the previously-written UTF-8 string written to `offset`.
//
//go:wasmimport extism:host/env log_error
func extismLogTrace(offset memory.ExtismPointer)

// extismGetLogLevel returns the configured log level
//
//go:wasmimport extism:host/env get_log_level
func extismGetLogLevel() int32
