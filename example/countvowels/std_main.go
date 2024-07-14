//go:build std
// +build std

package main

import (
	// "fmt"
	"strconv"

	"github.com/extism/go-pdk"
)

// Currently, the standard Go compiler cannot export custom functions and is limited to exporting
// `_start` via WASI. So, `main` functions should contain the plugin behavior, that the host will
// invoke by explicitly calling `_start`.
func main() {
	countVowels()
	countVowelsTyped()
	countVowelsJSONOutput()
	countVowelsJSONRoundtripMem()
}

// CountVowelsInput represents the JSON input provided by the host.
type CountVowelsInput struct {
	Input string `json:"input"`
}

// CountVowelsOutput represents the JSON output sent to the host.
type CountVowelsOuptut struct {
	Count  int    `json:"count"`
	Total  int    `json:"total"`
	Vowels string `json:"vowels"`
}

//export count_vowels_typed
func countVowelsTyped() int32 {
	var input CountVowelsInput
	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.OutputString(input.Input)
	return 0
}

//export count_vowels_json_output
func countVowelsJSONOutput() int32 {
	output := CountVowelsOuptut{Count: 42, Total: 2.1e7, Vowels: "aAeEiIoOuUyY"}
	err := pdk.OutputJSON(output)
	if err != nil {
		pdk.SetError(err)
		return -1
	}
	return 0
}

//export count_vowels_roundtrip_json_mem
func countVowelsJSONRoundtripMem() int32 {
	a := CountVowelsOuptut{Count: 42, Total: 2.1e7, Vowels: "aAeEiIoOuUyY"}
	mem, err := pdk.AllocateJSON(&a)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	// find the data in mem and ensure it's the same once decoded
	var b CountVowelsOuptut
	err = pdk.JSONFrom(mem.Offset(), &b)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	if a.Count != b.Count || a.Total != b.Total || a.Vowels != b.Vowels {
		pdk.SetErrorString("roundtrip JSON failed")
		return -1
	}

	pdk.OutputString("JSON roundtrip: a === b")
	return 0
}

func countVowels() int32 {
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
