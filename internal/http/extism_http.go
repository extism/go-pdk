package http

import "github.com/extism/go-pdk/internal/memory"

// extismHTTPRequest sends the HTTP `request` to the Extism host with the provided `body` (0 means no body)
// and returns back the memory offset to the response body.
//
//go:wasmimport extism:host/env http_request
func ExtismHTTPRequest(request, body memory.ExtismPointer) memory.ExtismPointer

// extismHTTPStatusCode returns the status code for the last-sent `extism_http_request` call.
//
//go:wasmimport extism:host/env http_status_code
func ExtismHTTPStatusCode() int32

// extismHTTPHeaders returns the response headers for the last-sent `extism_http_request` call.
//
//go:wasmimport extism:host/env http_headers
func ExtismHTTPHeaders() memory.ExtismPointer
