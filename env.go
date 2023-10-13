package pdk

type extismPointer uint64

//go:wasmimport extism:env input_length
func extism_input_length() uint64

//go:wasmimport extism:env length
func extism_length(extismPointer) uint64

//go:wasmimport extism:env alloc
func extism_alloc(uint64) extismPointer

//go:wasmimport extism:env free
func extism_free(extismPointer)

//go:wasmimport extism:env input_load_u8
func extism_input_load_u8_(extismPointer) uint32

func extism_input_load_u8(p extismPointer) uint8 {
	return uint8(extism_input_load_u8_(p))
}

//go:wasmimport extism:env input_load_u64
func extism_input_load_u64(extismPointer) uint64

//go:wasmimport extism:env output_set
func extism_output_set(extismPointer, uint64)

//go:wasmimport extism:env error_set
func extism_error_set(extismPointer)

//go:wasmimport extism:env config_get
func extism_config_get(extismPointer) extismPointer

//go:wasmimport extism:env var_get
func extism_var_get(extismPointer) extismPointer

//go:wasmimport extism:env var_set
func extism_var_set(extismPointer, extismPointer)

//go:wasmimport extism:env store_u8
func extism_store_u8_(extismPointer, uint32)
func extism_store_u8(p extismPointer, v uint8) {
	extism_store_u8_(p, uint32(v))
}

//go:wasmimport extism:env load_u8
func extism_load_u8_(extismPointer) uint32
func extism_load_u8(p extismPointer) uint8 {
	return uint8(extism_load_u8_(p))
}

//go:wasmimport extism:env store_u64
func extism_store_u64(extismPointer, uint64)

//go:wasmimport extism:env load_u64
func extism_load_u64(extismPointer) uint64

//go:wasmimport extism:env http_request
func extism_http_request(extismPointer, extismPointer) extismPointer

//go:wasmimport extism:env http_status_code
func extism_http_status_code() int32

//go:wasmimport extism:env log_info
func extism_log_info(extismPointer)

//go:wasmimport extism:env log_debug
func extism_log_debug(extismPointer)

//go:wasmimport extism:env log_warn
func extism_log_warn(extismPointer)

//go:wasmimport extism:env log_error
func extism_log_error(extismPointer)
