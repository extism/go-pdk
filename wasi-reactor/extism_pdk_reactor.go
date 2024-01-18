package reactor

//export __wasm_call_ctors
func __wasm_call_ctors()

//export _initialize
func _initialize() {
	__wasm_call_ctors()
}
