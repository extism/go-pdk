package pdk

type ExtismPointer uint64

//go:wasmimport env extism_input_length
func extism_input_length() uint64

//go:wasmimport env extism_length
func extism_length(ExtismPointer) uint64

//go:wasmimport env extism_alloc
func extism_alloc(uint64) ExtismPointer

//go:wasmimport env extism_free
func extism_free(ExtismPointer)

//go:wasmimport env extism_input_load_u8
func extism_input_load_u8_(ExtismPointer) uint32

func extism_input_load_u8(p ExtismPointer) uint8 {
	return uint8(extism_input_load_u8_(p))
}

//go:wasmimport env extism_input_load_u64
func extism_input_load_u64(ExtismPointer) uint64

//go:wasmimport env extism_output_set
func extism_output_set(ExtismPointer, uint64)

//go:wasmimport env extism_error_set
func extism_error_set(ExtismPointer)

//go:wasmimport env extism_config_get
func extism_config_get(ExtismPointer) ExtismPointer

//go:wasmimport env extism_var_get
func extism_var_get(ExtismPointer) ExtismPointer

//go:wasmimport env extism_var_set
func extism_var_set(ExtismPointer, ExtismPointer)

//go:wasmimport env extism_store_u8
func extism_store_u8_(ExtismPointer, uint32)
func extism_store_u8(p ExtismPointer, v uint8) {
	extism_store_u8_(p, uint32(v))
}

//go:wasmimport env extism_load_u8
func extism_load_u8_(ExtismPointer) uint32
func extism_load_u8(p ExtismPointer) uint8 {
	return uint8(extism_load_u8_(p))
}

//go:wasmimport env extism_store_u64
func extism_store_u64(ExtismPointer, uint64)

//go:wasmimport env extism_load_u64
func extism_load_u64(ExtismPointer) uint64

//go:wasmimport env extism_http_request
func extism_http_request(ExtismPointer, ExtismPointer) ExtismPointer

//go:wasmimport env extism_http_status_code
func extism_http_status_code() int32

//go:wasmimport env extism_log_info
func extism_log_info(ExtismPointer)

//go:wasmimport env extism_log_debug
func extism_log_debug(ExtismPointer)

//go:wasmimport env extism_log_warn
func extism_log_warn(ExtismPointer)

//go:wasmimport env extism_log_error
func extism_log_error(ExtismPointer)
