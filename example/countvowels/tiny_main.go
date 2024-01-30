//go:build !std
// +build !std

package main

import (
	"strconv"

	"github.com/extism/go-pdk"
)

type CountVowelsInput struct {
	Input string `json:"input"`
}

type CountVowelsOuptut struct {
	Count  int    `json:"count"`
	Total  int    `json:"total"`
	Vowels string `json:"vowels"`
}

//export count_vowels_typed
func count_vowels_typed() int32 {
	var input CountVowelsInput
	if err := pdk.InputJson(&input); err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.OutputString(input.Input)
	return 0
}

//export count_vowels_json_output
func count_vowels_json_output() int32 {
	output := CountVowelsOuptut{Count: 42, Total: 2.1e7, Vowels: "aAeEiIoOuUyY"}
	err := pdk.OutputJson(output)
	if err != nil {
		pdk.SetError(err)
		return -1
	}
	return 0
}

//export count_vowels
func count_vowels() int32 {
	input := pdk.Input()

	count := 0
	for _, a := range input {
		switch a {
		case 'A', 'I', 'E', 'O', 'U', 'a', 'e', 'i', 'o', 'u':
			count++
		default:
		}
	}

	// test some extra pdk functionality
	if pdk.GetVar("a") == nil {
		pdk.SetVar("a", []byte("this is var a"))
	}
	varA := pdk.GetVar("a")
	thing, ok := pdk.GetConfig("thing")

	if !ok {
		thing = "<unset by host>"
	}

	output := `{"count": ` + strconv.Itoa(count) + `, "config": "` + thing + `", "a": "` + string(varA) + `"}`
	mem := pdk.AllocateString(output)

	// zero-copy output to host
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
