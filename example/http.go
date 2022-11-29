package main

import (
	"github.com/extism/go-pdk"
)

//export http_get
func http_get() int32 {
	req := pdk.NewHTTPRequest("GET", "https://jsonplaceholder.typicode.com/todos/1")
	req.SetHeader("some-name", "some-value")
	req.SetHeader("another", "again")
	output := req.Send().Body()
	pdk.OutputMemory(pdk.AllocateString(string(output)))

	return 0
}

func main() {}
