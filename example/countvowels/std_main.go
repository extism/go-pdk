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
func countVowelsTyped() {
	var input CountVowelsInput
	if err := pdk.InputJSON(&input); err != nil {
		pdk.Error(err)
	}

	pdk.OutputString(input.Input)
}

//export count_vowels_json_output
func countVowelsJSONOutput() {
	output := CountVowelsOuptut{Count: 42, Total: 2.1e7, Vowels: "aAeEiIoOuUyY"}
	err := pdk.OutputJSON(output)
	if err != nil {
		pdk.Error(err)
	}
}

func countVowels() {
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
	if _, ok := pdk.GetVar("a"); !ok {
		pdk.SetVar("a", []byte("this is var a"))
	}
	varA, _ := pdk.GetVar("a")
	thing, ok := pdk.GetConfig("thing")

	if !ok {
		thing = "<unset by host>"
	}

	output := `{"count": ` + strconv.Itoa(count) + `, "config": "` + thing + `", "a": "` + string(varA) + `"}`

	pdk.OutputString(output)
}
