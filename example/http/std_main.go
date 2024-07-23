//go:build std
// +build std

package main

import (
	"github.com/extism/go-pdk"
)

// Currently, the standard Go compiler cannot export custom functions and is limited to exporting
// `_start` via WASI. So, `main` functions should contain the plugin behavior, that the host will
// invoke by explicitly calling `_start`.
func main() {
	httpGet()
}

func httpGet() int32 {
	// create an HTTP Request (withuot relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodGet, "https://jsonplaceholder.typicode.com/todos/1")
	req.SetHeader("some-name", "some-value")
	req.SetHeader("another", "again")
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()

	pdk.Output(res.Body())

	return 0
}
