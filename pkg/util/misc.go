package util

import (
	"unsafe"
)

func LogDebugBytes(t []byte) {
	ptr := unsafe.Pointer(unsafe.SliceData(t))
	logDebug(ptr, uint32(len(t)))
}

//go:wasmimport misc log_debug
func logDebug(ptr unsafe.Pointer, size uint32)
