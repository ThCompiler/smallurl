package simple

import (
	"hash"
	"math"
)

const (
	BaseOfResult  = 63
	countUsedBits = 6
)

type Hash63 struct {
	hash       hash.Hash
	bytesCount int
	state      []byte
}

func NewHash63(hsh hash.Hash, bytesCount int) *Hash63 {
	return &Hash63{hash: hsh, bytesCount: bytesCount, state: []byte{}}
}

func (h *Hash63) Write(p []byte) (int, error) {
	return h.hash.Write(p)
}

const A = 0.61803398874989484820

// Sum appends the current hash to b and returns the resulting slice.
func (h *Hash63) Sum(b []byte) []byte {
	hsh := h.hash.Sum(nil)

	res := make([]byte, len(b))
	copy(res, b)

	for i := 0; i < h.bytesCount; i++ {
		upd := float64(hsh[i]) * A

		res = append(res, byte((upd-math.Trunc(upd))*BaseOfResult))
	}

	return res
}

// Reset resets the Hash63 to its initial state.
func (h *Hash63) Reset() {
	h.hash.Reset()
	h.state = []byte{}
}

// Size returns the number of bytes Sum will return.
func (h *Hash63) Size() int {
	return h.bytesCount
}

// BlockSize returns the hash's underlying block size.
func (h *Hash63) BlockSize() int {
	return h.hash.BlockSize()
}
