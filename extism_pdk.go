package pdk

/*
#include "extism-pdk.h"
*/
import "C"

import (
	"unsafe"
)

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

func NewHost() Host {
	inputOffset := C.extism_input_offset()
	inputLength := C.extism_length(inputOffset)
	input := make([]byte, int(inputLength))

	C.extism_load(
		C.uint64_t(inputOffset),
		(*C.uchar)(unsafe.Pointer(&input[0])),
		C.ulong(inputLength),
	)

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

	return Memory{
		offset: uint64(offset),
		length: uint64(clength),
	}

}

func (h *Host) Input() []byte {
	return h.input
}

func (h *Host) InputString() string {
	return string(h.input)
}

func (h *Host) Output(data []byte) {
	clength := C.uint64_t(len(data))
	offset := C.extism_alloc(clength)

	C.extism_store(
		offset,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(clength),
	)
}

func (h *Host) Config(key string) string {
	mem := h.AllocateBytes([]byte(key))

	offset := C.extism_config_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return ""
	}

	value := make([]byte, uint64(clength))
	C.extism_load(
		offset,
		(*C.uchar)(unsafe.Pointer(&value[0])),
		C.ulong(clength),
	)

	return string(value)
}

func (h *Host) Vars() *Variables {
	return &Variables{host: h}
}

func (v *Variables) Get(key string) []byte {
	mem := v.host.AllocateBytes([]byte(key))

	offset := C.extism_kv_get(C.uint64_t(mem.offset))
	clength := C.extism_length(offset)
	if offset == 0 || clength == 0 {
		return nil
	}

	value := make([]byte, uint64(clength))
	C.extism_load(
		offset,
		(*C.uchar)(unsafe.Pointer(&value[0])),
		C.ulong(clength),
	)

	return value
}

func (v *Variables) Set(key string, value []byte) {
	keyMem := v.host.AllocateBytes([]byte(key))
	valMem := v.host.AllocateBytes(value)

	C.extism_kv_set(
		C.uint64_t(keyMem.offset),
		C.uint64_t(valMem.offset),
	)
}

func (v *Variables) Remove(key string) {
	mem := v.host.AllocateBytes([]byte(key))
	C.extism_kv_set(
		C.uint64_t(mem.offset),
		0,
	)
}

func (m *Memory) Load(buffer []byte) {
	C.extism_load(
		C.uint64_t(m.offset),
		(*C.uchar)(unsafe.Pointer(&buffer[0])),
		C.ulong(m.length),
	)
}

func (m *Memory) Store(data []byte) {
	C.extism_store(
		C.uint64_t(m.offset),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(m.length),
	)
}