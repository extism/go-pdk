//go:build std
// +build std

// extism call --wasi --allow-host "jsonplaceholder.typicode.com" std_main _start

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	pdk "github.com/extism/go-pdk"
)

// Currently, the standard Go compiler cannot export custom functions and is limited to exporting
// `_start` via WASI. So, `main` functions should contain the plugin behavior, that the host will
// invoke by explicitly calling `_start`.
func main() {
	body, err := httpGet()
	if err != nil {
		pdk.SetError(err)
		os.Exit(1)
	}

	pdk.OutputString(string(body))
}

func httpGet() ([]byte, error) {
	// Set the default transport to use Extism PDK HTTPTransport
	//
	// Alternativly, if using http.Client, specify the transport:
	//   client := http.Client{
	//   	Transport: &pdk.HTTPTransport{},
	//   }
	http.DefaultTransport = &pdk.HTTPTransport{}

	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %q", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %q", err)
	}

	return body, nil
}
