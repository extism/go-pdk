package pdk

import (
	"encoding/binary"
	"encoding/json"
)

// Memory represents memory allocated by (and shared with) the host.
type Memory struct {
	offset extismPointer
	length uint64
}

// LogLevel represents a logging level.
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogDebug
	LogWarn
	LogError
)

func load(offset extismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		binary.LittleEndian.PutUint64(buf[i:i+8], extismLoadU64(offset+extismPointer(i)))
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		buf[index] = extismLoadU8(offset + extismPointer(index))
	}
}

func loadInput() []byte {
	length := int(extismInputLength())
	buf := make([]byte, length)

	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		binary.LittleEndian.PutUint64(buf[i:i+8], extismInputLoadU64(extismPointer(i)))
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		buf[index] = extismInputLoadU8(extismPointer(index))
	}

	return buf
}

func store(offset extismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> 3

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << 3
		x := binary.LittleEndian.Uint64(buf[i : i+8])
		extismStoreU64(offset+extismPointer(i), x)
	}

	remainder := length & 7
	remainderOffset := chunkCount << 3
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		extismStoreU8(offset+extismPointer(index), buf[index])
	}
}

// Input returns a slice of bytes from the host.
func Input() []byte {
	return loadInput()
}

// JSONFrom unmarshals a `Memory` block located at `offset` from the host
// into the provided data `v`.
func JSONFrom(offset uint64, v any) error {
	mem := FindMemory(offset)
	return json.Unmarshal(mem.ReadBytes(), v)
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

	mem := AllocateBytes(b)
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer mem.Free()
	OutputMemory(mem)
	return nil
}

// Allocate allocates `length` uninitialized bytes on the host.
func Allocate(length int) Memory {
	clength := uint64(length)
	offset := extismAlloc(clength)

	return Memory{
		offset: offset,
		length: clength,
	}
}

// AllocateBytes allocates and saves the `data` into Memory on the host.
func AllocateBytes(data []byte) Memory {
	clength := uint64(len(data))
	offset := extismAlloc(clength)

	store(offset, data)

	return Memory{
		offset: offset,
		length: clength,
	}

}

// AllocateString allocates and saves the UTF-8 string `data` into Memory on the host.
func AllocateString(data string) Memory {
	return AllocateBytes([]byte(data))
}

// AllocateJSON allocates and saves the type `any` into Memory on the host.
func AllocateJSON(v any) (Memory, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return Memory{}, err
	}

	return AllocateBytes(b), nil
}

// InputString returns the input data from the host as a UTF-8 string.
func InputString() string {
	return string(Input())
}

// OutputMemory sends the `mem` Memory to the host output.
// Note that the `mem` is _NOT_ freed and is your responsibility to free when finished with it.
func OutputMemory(mem Memory) {
	extismOutputSet(mem.offset, mem.length)
}

// Output sends the `data` slice of bytes to the host output.
func Output(data []byte) {
	clength := uint64(len(data))
	offset := extismAlloc(clength)

	store(offset, data)
	extismOutputSet(offset, clength)
	extismFree(offset)
}

// OutputString sends the UTF-8 string `s` to the host output.
func OutputString(s string) {
	Output([]byte(s))
}

// SetError sets the host error string from `err`.
func SetError(err error) {
	SetErrorString(err.Error())
}

// SetErrorString sets the host error string from `err`.
func SetErrorString(err string) {
	mem := AllocateString(err)
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer mem.Free()
	extismErrorSet(mem.offset)
}

// GetConfig returns the config string associated with `key` (if any).
func GetConfig(key string) (string, bool) {
	mem := AllocateBytes([]byte(key))
	defer mem.Free()

	offset := extismConfigGet(mem.offset)
	clength := extismLength(offset)
	if offset == 0 || clength == 0 {
		return "", false
	}

	value := make([]byte, clength)
	load(offset, value)

	return string(value), true
}

// LogMemory logs the `memory` block on the host using the provided log `level`.
func LogMemory(level LogLevel, memory Memory) {
	switch level {
	case LogInfo:
		extismLogInfo(memory.offset)
	case LogDebug:
		extismLogDebug(memory.offset)
	case LogWarn:
		extismLogWarn(memory.offset)
	case LogError:
		extismLogError(memory.offset)
	}
}

// Log logs the provided UTF-8 string `s` on the host using the provided log `level`.
func Log(level LogLevel, s string) {
	mem := AllocateString(s)
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer mem.Free()

	LogMemory(level, mem)
}

// GetVar returns the byte slice (if any) associated with `key`.
func GetVar(key string) []byte {
	mem := AllocateBytes([]byte(key))
	defer mem.Free()

	offset := extismVarGet(mem.offset)
	clength := extismLength(offset)
	if offset == 0 || clength == 0 {
		return nil
	}

	value := make([]byte, clength)
	load(offset, value)

	return value
}

// SetVar sets the host variable associated with `key` to the `value` byte slice.
func SetVar(key string, value []byte) {
	keyMem := AllocateBytes([]byte(key))
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer keyMem.Free()

	valMem := AllocateBytes(value)
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer valMem.Free()

	extismVarSet(keyMem.offset, valMem.offset)
}

// GetVarInt returns the int associated with `key` (or 0 if none).
func GetVarInt(key string) int {
	mem := AllocateBytes([]byte(key))
	defer mem.Free()

	offset := extismVarGet(mem.offset)
	clength := extismLength(offset)
	if offset == 0 || clength == 0 {
		return 0
	}

	value := make([]byte, clength)
	load(offset, value)

	return int(binary.LittleEndian.Uint64(value))
}

// SetVarInt sets the host variable associated with `key` to the `value` int.
func SetVarInt(key string, value int) {
	keyMem := AllocateBytes([]byte(key))
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer keyMem.Free()

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(value))

	valMem := AllocateBytes(bytes)
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer valMem.Free()

	extismVarSet(keyMem.offset, valMem.offset)
}

// RemoveVar removes (and frees) the host variable associated with `key`.
func RemoveVar(key string) {
	mem := AllocateBytes([]byte(key))
	// TODO: coordinate replacement of call to free based on SDK alignment
	// defer mem.Free()
	extismVarSet(mem.offset, 0)
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
	memory Memory
	status uint16
}

// Memory returns the memory associated with the `HTTPResponse`.
func (r HTTPResponse) Memory() Memory {
	return r.memory
}

// Body returns the body byte slice (if any) from the `HTTPResponse`.
func (r HTTPResponse) Body() []byte {
	if r.memory.length == 0 {
		return nil
	}

	buf := make([]byte, r.memory.length)
	r.memory.Load(buf)
	return buf
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

	req := AllocateBytes(enc)
	defer req.Free()
	var dataOffset extismPointer
	if len(r.body) > 0 {
		data := AllocateBytes(r.body)
		defer data.Free()
		dataOffset = data.offset
	}

	offset := extismHTTPRequest(req.offset, dataOffset)
	length := extismLengthUnsafe(offset)
	status := uint16(extismHTTPStatusCode())

	memory := Memory{offset, length}

	return HTTPResponse{
		memory,
		status,
	}
}

// ReadBytes returns the host memory block as a slice of bytes.
func (m *Memory) ReadBytes() []byte {
	buff := make([]byte, m.length)
	m.Load(buff)
	return buff
}

// Load copies the host memory block to the provided `buffer` byte slice.
func (m *Memory) Load(buffer []byte) {
	load(m.offset, buffer)
}

// Store copies the `data` byte slice into host memory.
func (m *Memory) Store(data []byte) {
	store(m.offset, data)
}

// Free frees the host memory block.
func (m *Memory) Free() {
	extismFree(m.offset)
}

// Length returns the number of bytes in the host memory block.
func (m *Memory) Length() uint64 {
	return m.length
}

// Offset returns the offset of the host memory block.
func (m *Memory) Offset() uint64 {
	return uint64(m.offset)
}

// FindMemory finds the host memory block at the given `offset`.
func FindMemory(offset uint64) Memory {
	length := extismLength(extismPointer(offset))
	if length == 0 {
		return Memory{0, 0}
	}
	return Memory{extismPointer(offset), length}
}

// ParamBytes returns bytes from Extism host memory given an offset.
func ParamBytes(offset uint64) []byte {
	mem := FindMemory(offset)
	return mem.ReadBytes()
}

// ParamString returns UTF-8 string data from Extism host memory given an offset.
func ParamString(offset uint64) string {
	return string(ParamBytes(offset))
}

// ParamU32 returns a uint32 from Extism host memory given an offset.
func ParamU32(offset uint64) uint32 {
	return binary.LittleEndian.Uint32(ParamBytes(offset))
}

// ParamU64 returns a uint64 from Extism host memory given an offset.
func ParamU64(offset uint64) uint64 {
	return binary.LittleEndian.Uint64(ParamBytes(offset))
}

// ResultBytes allocates bytes and returns the offset in Extism host memory.
func ResultBytes(d []byte) uint64 {
	mem := AllocateBytes(d)
	return mem.Offset()
}

// ResultString allocates a UTF-8 string and returns the offset in Extism host memory.
func ResultString(s string) uint64 {
	mem := AllocateString(s)
	return mem.Offset()
}

// ResultU32 allocates a uint32 and returns the offset in Extism host memory.
func ResultU32(d uint32) uint64 {
	mem := AllocateBytes(binary.LittleEndian.AppendUint32([]byte{}, d))
	return mem.Offset()
}

// ResultU64 allocates a uint64 and returns the offset in Extism host memory.
func ResultU64(d uint64) uint64 {
	mem := AllocateBytes(binary.LittleEndian.AppendUint64([]byte{}, d))
	return mem.Offset()
}
