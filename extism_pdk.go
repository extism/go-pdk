package pdk

import (
	"encoding/binary"
	"encoding/json"
	"github.com/extism/go-pdk/internal/models"
	"strings"
)

/*
#include "extism-pdk.h"
*/
import "C"

type Memory struct {
	offset uint64
	length uint64
}

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogDebug
	LogWarn
	LogError
)

func load(offset uint64, buf []byte) {
	length := len(buf)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := C.extism_load_u64(offset + uint64(i))
			binary.LittleEndian.PutUint64(buf[i:i+8], x)
			i += 7
			continue
		}
		buf[i] = byte(C.extism_load_u8(offset + uint64(i)))
	}
}

func loadInput() []byte {
	length := int(C.extism_input_length())
	buf := make([]byte, length)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := C.extism_input_load_u64(uint64(i))
			binary.LittleEndian.PutUint64(buf[i:i+8], x)
			i += 7
			continue
		}
		buf[i] = byte(C.extism_input_load_u8(uint64(i)))
	}

	return buf
}

func store(offset uint64, buf []byte) {
	length := len(buf)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := binary.LittleEndian.Uint64(buf[i : i+8])
			C.extism_store_u64(offset+uint64(i), C.uint64_t(x))
			i += 7
			continue
		}

		C.extism_store_u8(offset+uint64(i), C.uint8_t(buf[i]))
	}
}

func Input() []byte {
	return loadInput()
}

func Allocate(length int) Memory {
	clength := C.uint64_t(length)
	offset := C.extism_alloc(clength)

	return Memory{
		offset: uint64(offset),
		length: uint64(clength),
	}
}

func AllocateBytes(data []byte) Memory {
	clength := C.uint64_t(len(data))
	offset := C.extism_alloc(clength)

	store(offset, data)

	return Memory{
		offset: uint64(offset),
		length: uint64(clength),
	}

}

func AllocateString(data string) Memory {
	return AllocateBytes([]byte(data))
}

func InputString() string {
	return string(Input())
}

func OutputMemory(mem Memory) {
	C.extism_output_set(mem.offset, mem.length)
}

func Output(data []byte) {
	clength := C.uint64_t(len(data))
	offset := C.extism_alloc(clength)

	store(offset, data)
	C.extism_output_set(offset, clength)
}

func OutputString(s string) {
	Output([]byte(s))
}

func GetConfig(key string) (string, bool) {
	mem := AllocateBytes([]byte(key))

	offset := C.extism_config_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return "", false
	}

	value := make([]byte, uint64(clength))
	load(offset, value)

	return string(value), true
}

func LogMemory(level LogLevel, memory Memory) {
	switch level {
	case LogInfo:
		C.extism_log_info(memory.offset)
	case LogDebug:
		C.extism_log_debug(memory.offset)
	case LogWarn:
		C.extism_log_warn(memory.offset)
	case LogError:
		C.extism_log_error(memory.offset)
	}
}

func Log(level LogLevel, s string) {
	mem := AllocateString(s)
	LogMemory(level, mem)
}

func GetVar(key string) []byte {
	mem := AllocateBytes([]byte(key))

	offset := C.extism_var_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return nil
	}

	value := make([]byte, uint64(clength))
	load(offset, value)

	return value
}

func SetVar(key string, value []byte) {
	keyMem := AllocateBytes([]byte(key))
	valMem := AllocateBytes(value)

	C.extism_var_set(
		C.uint64_t(keyMem.offset),
		C.uint64_t(valMem.offset),
	)
}

func RemoveVar(key string) {
	mem := AllocateBytes([]byte(key))
	C.extism_var_set(
		C.uint64_t(mem.offset),
		0,
	)
}

type HTTPRequest struct {
	url    string
	header map[string]string
	method string
	body   []byte
}

type HTTPResponse struct {
	memory Memory
	status uint16
}

func (r HTTPResponse) Memory() Memory {
	return r.memory
}

func (r HTTPResponse) Body() []byte {
	buf := make([]byte, r.memory.length)
	r.memory.Load(buf)
	return buf
}

func (r HTTPResponse) Status() uint16 {
	return r.status
}

func NewHTTPRequest(method string, url string) *HTTPRequest {
	return &HTTPRequest{url: url, header: nil, method: strings.ToUpper(method), body: nil}
}

func (r *HTTPRequest) SetHeader(key string, value string) *HTTPRequest {
	if r.header == nil {
		r.header = map[string]string{}
	}
	r.header[key] = value
	return r
}

func (r *HTTPRequest) SetBody(body []byte) *HTTPRequest {
	r.body = body
	return r
}

func (r *HTTPRequest) Send() HTTPResponse {
	meta := models.MetaData{
		Url:     r.url,
		Method:  r.method,
		Headers: r.header,
	}

	enc, _ := json.Marshal(meta)

	req := AllocateBytes(enc)
	defer req.Free()
	data := AllocateBytes(r.body)
	defer data.Free()

	offset := C.extism_http_request(C.uint64_t(req.offset), data.offset)
	length := uint64(C.extism_length(offset))
	status := uint16(C.extism_http_status_code())

	memory := Memory{offset, length}
	defer memory.Free()

	return HTTPResponse{
		memory,
		status,
	}
}

func (m *Memory) Load(buffer []byte) {
	load(m.offset, buffer)
}

func (m *Memory) Store(data []byte) {
	store(m.offset, data)
}

func (m *Memory) Free() {
	C.extism_free(m.offset)
}

func (m *Memory) Length() uint64 {
	return m.length
}

func (m *Memory) Offset() uint64 {
	return m.offset
}

func FindMemory(offset uint64) Memory {
	length := uint64(C.extism_length(offset))
	return Memory{offset, length}
}
