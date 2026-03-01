//go:build goexperiment.simd

// This file is selected automatically by the Go toolchain when the caller
// builds with GOEXPERIMENT=simd (Go 1.26+). No action is required by library
// consumers — the scalar fallback in ascii_scalar.go is used on all other
// toolchains and configurations.

package utils

import (
	"simd/archsimd"
	"unsafe"
)

// isASCII reports whether s contains only ASCII bytes.
// SIMD-accelerated path: processes 16 bytes per SSE2 iteration using Int8x16.
// Bytes >= 0x80 are negative when interpreted as signed int8, so
// Less(zero).ToBits() != 0 indicates at least one non-ASCII byte in the chunk.
func isASCII(s string) bool {
	n := len(s)
	if n == 0 {
		return true
	}

	// Reinterpret the string's backing bytes as []int8 — same memory layout,
	// required by the archsimd load functions.
	ptr := (*int8)(unsafe.Pointer(unsafe.StringData(s)))
	data := unsafe.Slice(ptr, n)

	zero := archsimd.BroadcastInt8x16(0)
	i := 0

	// Main loop: 16 bytes per SSE2 iteration.
	for ; i+16 <= n; i += 16 {
		if archsimd.LoadInt8x16((*[16]int8)(data[i:])).Less(zero).ToBits() != 0 {
			return false
		}
	}

	// Tail: fewer than 16 bytes remaining — LoadInt8x16SlicePart zero-pads.
	// Zero padding is safe: int8(0) is not < 0, so it won't be flagged.
	if i < n {
		if archsimd.LoadInt8x16SlicePart(data[i:]).Less(zero).ToBits() != 0 {
			return false
		}
	}

	return true
}
