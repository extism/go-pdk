package reactor

//export __wasm_call_ctors
func wasmCallCtors()

//export _initialize
func initialize() {
	wasmCallCtors()
}
