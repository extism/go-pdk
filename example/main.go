package main

import (
	"strconv"

	"github.com/extism/go-pdk"
)

//export count_vowels
func count_vowels() int32 {
	host := pdk.NewHost()
	input := host.Input()

	count := 0
	for _, a := range input {
		switch a {
		case 'A', 'I', 'E', 'O', 'U', 'a', 'e', 'i', 'o', 'u':
			count++
		default:
		}
	}

	// test some extra pdk functionality
	vars := host.Vars()
	if vars.Get("a") == nil {
		vars.Set("a", []byte("this is var a"))
	}
	varA := vars.Get("a")
	thing := host.Config("thing")

	if thing == "" {
		thing = "<unset by host>"
	}

	output := `{"count": ` + strconv.Itoa(count) + `, "config": "` + thing + `", "a": "` + string(varA) + `"}`
	mem := host.AllocateString(output)

	// zero-copy output to host
	host.OutputMemory(mem)

	return 0
}

func main() {}
