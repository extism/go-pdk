package pdk

type extismHandle uint64

type extismStream int32

const (
	extismStreamInput extismStream = iota
	extismStreamOutput
)

// extismInputLength returns the number of bytes provided by the host via its input methods.
//
//go:wasmimport extism:host/env read
func extismRead(stream extismStream, buf extismHandle) int64

//go:wasmimport extism:host/env write
func extismWrite(stream extismStream, buf extismHandle) int64

//go:wasmimport extism:host/env close
func extismClose(stream extismStream)

//go:wasmimport extism:host/env stack_push
func extismStackPush()

//go:wasmimport extism:host/env stack_pop
func extismStackPop()

// extismErrorSet sets the "error" data from the plugin to the host to be the memory that
// has been written at `offset`.
// The memory can be immediately freed because the host makes a copy for its use.
//
//go:wasmimport extism:host/env error
func extismError(buf extismHandle)

// extismConfigGet returns the host memory block offset for the "config" data associated with
// the key which is represented by the UTF-8 string which as been previously written at `offset`.
// The memory for the key can be immediately freed because the host has its own copy.
//
//go:wasmimport extism:host/env config_read
func extismConfigRead(key extismHandle, buf extismHandle) int64

//go:wasmimport extism:host/env config_length
func extismConfigLength(key extismHandle) int64

// extismHTTPRequest sends the HTTP `request` to the Extism host with the provided `body` (0 means no body)
// and returns back the memory offset to the response body.
//
//go:wasmimport extism:host/env http_request
func extismHTTPRequest(request, body extismHandle) int64

//go:wasmimport extism:host/env http_body
func extismHTTPBody(buf extismHandle) int64

// extismHTTPStatusCode returns the status code for the last-sent `extism_http_request` call.
//
//go:wasmimport extism:host/env http_status_code
func extismHTTPStatusCode() int32

// extismLog logs a string to the host from the previously-written UTF-8 string written to `offset`.
// Note that the memory at `offset` can be immediately freed because it is immediately logged.
//
//go:wasmimport extism:host/env log
func extismLog(level LogLevel, offset extismHandle)
