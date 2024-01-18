//go:build !std
// +build !std

package main

import (
	"os"

	"github.com/extism/go-pdk"
	_ "github.com/extism/go-pdk/pdk_reactor"
)

//export read_file
func read_file() {
	name := pdk.InputString()

	content, err := os.ReadFile(name)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return
	}

	pdk.Output(content)
}

func main() {}
