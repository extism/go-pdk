package pdk

import (
	"encoding/binary"
)

/*
#include "extism-pdk.h"
*/
import "C"

type Host struct {
	input []byte
}

type Memory struct {
	offset uint64
	length uint64
}

type Variables struct {
	host *Host
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

func NewHost() Host {
	input := loadInput()
	return Host{input}
}

func (h *Host) Allocate(length int) Memory {
	clength := C.uint64_t(length)
	offset := C.extism_alloc(clength)

	return Memory{
		offset: uint64(offset),
		length: uint64(clength),
	}
}

func (h *Host) AllocateBytes(data []byte) Memory {
	clength := C.uint64_t(len(data))
	offset := C.extism_alloc(clength)

	store(offset, data)

	return Memory{
		offset: uint64(offset),
		length: uint64(clength),
	}

}

func (h *Host) AllocateString(data string) Memory {
	return h.AllocateBytes([]byte(data))
}

func (h *Host) Input() []byte {
	return h.input
}

func (h *Host) InputString() string {
	return string(h.input)
}

func (h *Host) OutputMemory(mem Memory) {
	C.extism_output_set(mem.offset, mem.length)
}

func (h *Host) Output(data []byte) {
	clength := C.uint64_t(len(data))
	offset := C.extism_alloc(clength)

	store(offset, data)
	C.extism_output_set(offset, clength)
}

func (h *Host) Config(key string) string {
	mem := h.AllocateBytes([]byte(key))

	offset := C.extism_config_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return ""
	}

	value := make([]byte, uint64(clength))
	load(offset, value)

	return string(value)
}

func (h *Host) LogMemory(level LogLevel, memory Memory) {
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

func (h *Host) Log(level LogLevel, s string) {
	mem := h.AllocateString(s)
	h.LogMemory(level, mem)
}

func (h *Host) Vars() *Variables {
	return &Variables{host: h}
}

func (v *Variables) Get(key string) []byte {
	mem := v.host.AllocateBytes([]byte(key))

	offset := C.extism_var_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return nil
	}

	value := make([]byte, uint64(clength))
	load(offset, value)

	return value
}

func (v *Variables) Set(key string, value []byte) {
	keyMem := v.host.AllocateBytes([]byte(key))
	valMem := v.host.AllocateBytes(value)

	C.extism_var_set(
		C.uint64_t(keyMem.offset),
		C.uint64_t(valMem.offset),
	)
}

func (v *Variables) Remove(key string) {
	mem := v.host.AllocateBytes([]byte(key))
	C.extism_var_set(
		C.uint64_t(mem.offset),
		0,
	)
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
