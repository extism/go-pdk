package pdk

import (
	"encoding/binary"
	"encoding/json"
	"strings"
)

type Memory struct {
	offset extismPointer
	length uint64
}

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogDebug
	LogWarn
	LogError
)

func load(offset extismPointer, buf []byte) {
	length := len(buf)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := extism_load_u64(offset + extismPointer(i))
			binary.LittleEndian.PutUint64(buf[i:i+8], x)
			i += 7
			continue
		}
		buf[i] = extism_load_u8(offset + extismPointer(i))
	}
}

func loadInput() []byte {
	length := int(extism_input_length())
	buf := make([]byte, length)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := extism_input_load_u64(extismPointer(i))
			binary.LittleEndian.PutUint64(buf[i:i+8], x)
			i += 7
			continue
		}
		buf[i] = extism_input_load_u8(extismPointer(i))
	}

	return buf
}

func store(offset extismPointer, buf []byte) {
	length := len(buf)

	for i := 0; i < length; i++ {
		if length-i >= 8 {
			x := binary.LittleEndian.Uint64(buf[i : i+8])
			extism_store_u64(offset+extismPointer(i), x)
			i += 7
			continue
		}

		extism_store_u8(offset+extismPointer(i), buf[i])
	}
}

func Input() []byte {
	return loadInput()
}

func Allocate(length int) Memory {
	clength := uint64(length)
	offset := extism_alloc(clength)

	return Memory{
		offset: offset,
		length: clength,
	}
}

func AllocateBytes(data []byte) Memory {
	clength := uint64(len(data))
	offset := extism_alloc(clength)

	store(offset, data)

	return Memory{
		offset: offset,
		length: clength,
	}

}

func AllocateString(data string) Memory {
	return AllocateBytes([]byte(data))
}

func InputString() string {
	return string(Input())
}

func OutputMemory(mem Memory) {
	extism_output_set(mem.offset, mem.length)
}

func Output(data []byte) {
	clength := uint64(len(data))
	offset := extism_alloc(clength)

	store(offset, data)
	extism_output_set(offset, clength)
}

func OutputString(s string) {
	Output([]byte(s))
}

func GetConfig(key string) (string, bool) {
	mem := AllocateBytes([]byte(key))
	defer mem.Free()

	offset := extism_config_get(mem.offset)
	clength := extism_length(offset)
	if offset == 0 || clength == 0 {
		return "", false
	}

	value := make([]byte, clength)
	load(offset, value)

	return string(value), true
}

func LogMemory(level LogLevel, memory Memory) {
	switch level {
	case LogInfo:
		extism_log_info(memory.offset)
	case LogDebug:
		extism_log_debug(memory.offset)
	case LogWarn:
		extism_log_warn(memory.offset)
	case LogError:
		extism_log_error(memory.offset)
	}
}

func Log(level LogLevel, s string) {
	mem := AllocateString(s)
	defer mem.Free()

	LogMemory(level, mem)
}

func GetVar(key string) []byte {
	mem := AllocateBytes([]byte(key))

	offset := extism_var_get(mem.offset)
	clength := extism_length(offset)
	if offset == 0 || clength == 0 {
		return nil
	}

	value := make([]byte, clength)
	load(offset, value)

	return value
}

func SetVar(key string, value []byte) {
	keyMem := AllocateBytes([]byte(key))
	defer keyMem.Free()

	valMem := AllocateBytes(value)
	defer valMem.Free()

	extism_var_set(keyMem.offset, valMem.offset)
}

func RemoveVar(key string) {
	mem := AllocateBytes([]byte(key))
	extism_var_set(mem.offset, 0)
}

type HTTPRequest struct {
	url     string
	headers map[string]string
	method  string
	body    []byte
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
	return &HTTPRequest{url: url, headers: nil, method: strings.ToUpper(method), body: nil}
}

func (r *HTTPRequest) SetHeader(key string, value string) *HTTPRequest {
	if r.headers == nil {
		r.headers = map[string]string{}
	}
	r.headers[key] = value
	return r
}

func (r *HTTPRequest) SetBody(body []byte) *HTTPRequest {
	r.body = body
	return r
}

type HTTPRequestMeta struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"header"`
}

func (r *HTTPRequest) Send() HTTPResponse {
	meta := HTTPRequestMeta{
		Url:     r.url,
		Method:  r.method,
		Headers: r.header,
	}

	enc, _ := json.Marshal(meta)

	req := AllocateBytes(enc)
	defer req.Free()
	data := AllocateBytes(r.body)
	defer data.Free()

	offset := extism_http_request(req.offset, data.offset)
	length := extism_length(offset)
	status := uint16(extism_http_status_code())

	memory := Memory{offset, length}

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
	extism_free(m.offset)
}

func (m *Memory) Length() uint64 {
	return m.length
}

func (m *Memory) Offset() uint64 {
	return uint64(m.offset)
}

func FindMemory(offset uint64) Memory {
	length := extism_length(extismPointer(offset))
	return Memory{extismPointer(offset), length}
}
