package pdk

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"unsafe"
)

// LogLevel represents a logging level.
type LogLevel int32

const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

func makeHandle(data []byte) extismHandle {
	if data == nil {
		return 0
	}
	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))
	len := uint64(len(data))
	return extismHandle((ptr << 32) | (len & uint64(0xffffffff)))
}

func splitHandle(handle extismHandle) (uint32, uint32) {
	ptr := (handle >> 32) & 0xffffffff
	len := handle & 0xffffffff
	return uint32(ptr), uint32(len)
}

// Input returns a slice of bytes from the host.
func Input() []byte {
	buffer := make([]byte, 1024)
	data := bytes.NewBuffer([]byte{})
	for {
		n := extismRead(extismStreamInput, makeHandle(buffer))
		if n <= 0 {
			break
		}

		data.Write(buffer[:n])
	}
	return data.Bytes()
}

// InputJSON returns unmartialed JSON data from the host "input".
func InputJSON(v any) error {
	return json.Unmarshal(Input(), v)
}

// OutputJSON marshals the provided data `v` as output to the host.
func OutputJSON(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	Output(b)
	return nil
}

// InputString returns the input data from the host as a UTF-8 string.
func InputString() string {
	return string(Input())
}

// Output sends the `data` slice of bytes to the host output.
func Output(data []byte) {
	extismWrite(extismStreamOutput, makeHandle(data))
}

// OutputString sends the UTF-8 string `s` to the host output.
func OutputString(s string) {
	Output([]byte(s))
}

// Error sets the host error string from `err`.
func Error(err error) {
	ErrorString(err.Error())
}

// ErrorString sets the host error string from `err`.
func ErrorString(err string) {
	data := []byte(err)
	extismError(makeHandle(data))
}

// GetConfig returns the config string associated with `key` (if any).
func GetConfig(key string) (string, bool) {
	keyData := []byte(key)
	keyHandle := makeHandle(keyData)
	length := extismConfigLength(keyHandle)
	if length < 0 {
		return "", false
	}
	buf := make([]byte, length)
	extismConfigRead(keyHandle, makeHandle(buf))
	return string(buf), true
}

// LogMemory logs the `memory` block on the host using the provided log `level`.
func Log(level LogLevel, s string) {
	data := []byte(s)
	extismLog(level, makeHandle(data))
}

var vars = map[string][]byte{}

// GetVar returns the byte slice (if any) associated with `key`.
func GetVar(key string) ([]byte, bool) {
	x, ok := vars[key]
	return x, ok
}

// SetVar sets the host variable associated with `key` to the `value` byte slice.
func SetVar(key string, value []byte) {
	vars[key] = value
}

// GetVarInt returns the int associated with `key` (or 0 if none).
func GetVarInt(key string) (int, bool) {
	value, ok := GetVar(key)
	if !ok {
		return 0, false
	}
	return int(binary.LittleEndian.Uint64(value)), true
}

// SetVarInt sets the host variable associated with `key` to the `value` int.
func SetVarInt(key string, value int) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(value))
	SetVar(key, bytes)
}

// RemoveVar removes (and frees) the host variable associated with `key`.
func RemoveVar(key string) {
	delete(vars, key)
}

// HTTPRequestMeta represents the metadata associated with an HTTP request on the host.
type HTTPRequestMeta struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

// HTTPRequest represents an HTTP request sent by the host.
type HTTPRequest struct {
	meta HTTPRequestMeta
	body []byte
}

// HTTPResponse represents an HTTP response returned from the host.
type HTTPResponse struct {
	length int
	status uint16
}

// Memory returns the length of the response body
func (r HTTPResponse) Length() int {
	return r.length
}

// Body returns the body byte slice (if any) from the `HTTPResponse`.
func (r HTTPResponse) Body() []byte {
	if r.length == 0 {
		return nil
	}

	buf := make([]byte, 1024)
	out := bytes.NewBuffer([]byte{})
	for {
		n := extismHTTPBody(makeHandle(buf))
		if n < 0 {
			break
		}
		out.Write(buf[:n])
	}
	return out.Bytes()
}

// Status returns the status code from the `HTTPResponse`.
func (r HTTPResponse) Status() uint16 {
	return r.status
}

// HTTPMethod represents an HTTP method.
type HTTPMethod int32

const (
	MethodGet HTTPMethod = iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch // RFC 5789
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
)

func (m HTTPMethod) String() string {
	switch m {
	case MethodGet:
		return "GET"
	case MethodHead:
		return "HEAD"
	case MethodPost:
		return "POST"
	case MethodPut:
		return "PUT"
	case MethodPatch:
		return "PATCH"
	case MethodDelete:
		return "DELETE"
	case MethodConnect:
		return "CONNECT"
	case MethodOptions:
		return "OPTIONS"
	case MethodTrace:
		return "TRACE"
	default:
		return ""
	}
}

// NewHTTPRequest returns a new `HTTPRequest`.
func NewHTTPRequest(method HTTPMethod, url string) *HTTPRequest {
	return &HTTPRequest{
		meta: HTTPRequestMeta{
			URL:     url,
			Headers: nil,
			Method:  method.String(),
		},
		body: nil,
	}
}

// SetHeader sets an HTTP header `key` to `value`.
func (r *HTTPRequest) SetHeader(key string, value string) *HTTPRequest {
	if r.meta.Headers == nil {
		r.meta.Headers = make(map[string]string)
	}
	r.meta.Headers[key] = value
	return r
}

// SetBody sets an HTTP request body to the provided byte slice.
func (r *HTTPRequest) SetBody(body []byte) *HTTPRequest {
	r.body = body
	return r
}

// Send sends the `HTTPRequest` from the host and returns the `HTTPResponse`.
func (r *HTTPRequest) Send() HTTPResponse {
	enc, _ := json.Marshal(r.meta)

	length := int(extismHTTPRequest(makeHandle(enc), makeHandle(r.body)))
	status := uint16(extismHTTPStatusCode())

	return HTTPResponse{
		length,
		status,
	}
}

// // ParamBytes returns bytes from Extism host memory given an offset.
// func ParamBytes(offset uint64) []byte {
// 	mem := FindMemory(offset)
// 	return mem.ReadBytes()
// }

// // ParamString returns UTF-8 string data from Extism host memory given an offset.
// func ParamString(offset uint64) string {
// 	return string(ParamBytes(offset))
// }

// // ParamU32 returns a uint32 from Extism host memory given an offset.
// func ParamU32(offset uint64) uint32 {
// 	return binary.LittleEndian.Uint32(ParamBytes(offset))
// }

// // ParamU64 returns a uint64 from Extism host memory given an offset.
// func ParamU64(offset uint64) uint64 {
// 	return binary.LittleEndian.Uint64(ParamBytes(offset))
// }

// // ResultBytes allocates bytes and returns the offset in Extism host memory.
// func ResultBytes(d []byte) uint64 {
// 	mem := AllocateBytes(d)
// 	return mem.Offset()
// }

// // ResultString allocates a UTF-8 string and returns the offset in Extism host memory.
// func ResultString(s string) uint64 {
// 	mem := AllocateString(s)
// 	return mem.Offset()
// }

// // ResultU32 allocates a uint32 and returns the offset in Extism host memory.
// func ResultU32(d uint32) uint64 {
// 	mem := AllocateBytes(binary.LittleEndian.AppendUint32([]byte{}, d))
// 	return mem.Offset()
// }

// // ResultU64 allocates a uint64 and returns the offset in Extism host memory.
// func ResultU64(d uint64) uint64 {
// 	mem := AllocateBytes(binary.LittleEndian.AppendUint64([]byte{}, d))
// 	return mem.Offset()
// }
