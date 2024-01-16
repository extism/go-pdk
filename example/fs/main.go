package main

import (
	"os"

	"github.com/extism/go-pdk"
)

//export __wasm_call_ctors
func __wasm_call_ctors()

//export _initialize
func _initialize() {
	__wasm_call_ctors()
}

//export write_file
func writeFile() int32 {
	input := pdk.Input()

	err := pluginMain("/mnt/wasm.txt", input)

	if err != nil {
		pdk.Log(pdk.LogTrace, err.Error())
		return 1
	}

	return 0
}

func pluginMain(filename string, data []byte) error {
	pdk.Log(pdk.LogInfo, "Writing following data to disk: "+string(data))

	// Write to the file, will be created if it doesn't exist
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {}
