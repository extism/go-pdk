package pdk

import (
	"encoding/binary"
	"encoding/json"
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
	LogTrace
)

func load(offset extismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		binary.LittleEndian.PutUint64(buf[i:i+8], extism_load_u64(offset+extismPointer(i)))
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		buf[index] = extism_load_u8(offset + extismPointer(index))
	}
}

func loadInput() []byte {
	length := int(extism_input_length())
	buf := make([]byte, length)

	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		binary.LittleEndian.PutUint64(buf[i:i+8], extism_input_load_u64(extismPointer(i)))
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		buf[index] = extism_input_load_u8(extismPointer(index))
	}

	return buf
}

func store(offset extismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		x := binary.LittleEndian.Uint64(buf[i : i+8])
		extism_store_u64(offset+extismPointer(i), x)
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		extism_store_u8(offset+extismPointer(index), buf[index])
	}
}

func Input() []byte {
	return loadInput()
}

func JSONFrom(offset uint64, v any) error {
	mem := FindMemory(offset)
	return json.Unmarshal(mem.ReadBytes(), v)
}

func InputJSON(v any) error {
	return json.Unmarshal(Input(), v)
}

func OutputJSON(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	OutputMemory(AllocateBytes(b))
	return nil
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

func SetError(err error) {
	SetErrorString(err.Error())
}

func SetErrorString(err string) {
	mem := AllocateString(err)
	extism_error_set(mem.offset)
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

func GetVarInt(key string) int {
	mem := AllocateBytes([]byte(key))

	offset := extism_var_get(mem.offset)
	clength := extism_length(offset)
	if offset == 0 || clength == 0 {
		return 0
	}

	value := make([]byte, clength)
	load(offset, value)

	return int(binary.LittleEndian.Uint64(value))
}

func SetVarInt(key string, value int) {
	keyMem := AllocateBytes([]byte(key))
	defer keyMem.Free()

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(value))

	valMem := AllocateBytes(bytes)
	defer valMem.Free()

	extism_var_set(keyMem.offset, valMem.offset)
}

func RemoveVar(key string) {
	mem := AllocateBytes([]byte(key))
	extism_var_set(mem.offset, 0)
}

type HTTPRequestMeta struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type HTTPRequest struct {
	meta HTTPRequestMeta
	body []byte
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

func NewHTTPRequest(method HTTPMethod, url string) *HTTPRequest {
	return &HTTPRequest{
		meta: HTTPRequestMeta{
			Url:     url,
			Headers: nil,
			Method:  method.String(),
		},
		body: nil,
	}
}

func (r *HTTPRequest) SetHeader(key string, value string) *HTTPRequest {
	if r.meta.Headers == nil {
		r.meta.Headers = make(map[string]string)
	}
	r.meta.Headers[key] = value
	return r
}

func (r *HTTPRequest) SetBody(body []byte) *HTTPRequest {
	r.body = body
	return r
}

func (r *HTTPRequest) Send() HTTPResponse {
	enc, _ := json.Marshal(r.meta)

	req := AllocateBytes(enc)
	defer req.Free()
	data := AllocateBytes(r.body)
	defer data.Free()

	offset := extism_http_request(req.offset, data.offset)
	length := extism_length_unsafe(offset)
	status := uint16(extism_http_status_code())

	memory := Memory{offset, length}

	return HTTPResponse{
		memory,
		status,
	}
}

func (m *Memory) ReadBytes() []byte {
	buff := make([]byte, m.length)
	m.Load(buff)
	return buff
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
	if length == 0 {
		return Memory{0, 0}
	}
	return Memory{extismPointer(offset), length}
}
