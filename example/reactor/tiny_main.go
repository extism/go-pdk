//go:build !std
// +build !std

package main

import (
	"os"

	"github.com/extism/go-pdk"
)

//go:wasmexport read_file
func readFile() {
	name := pdk.InputString()

	content, err := os.ReadFile(name)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return
	}

	pdk.Output(content)
}
