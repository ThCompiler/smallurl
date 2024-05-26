package sponge

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

func (h *Hash63) transform() {
	h.hash.Reset()
	_, _ = h.hash.Write(h.state)
	h.state = h.hash.Sum(nil)
}

const A = 0.61803398874989484820

func (h *Hash63) GetByte() byte {
	h.transform()

	upd := float64(h.state[0]) * A

	return byte((upd - math.Trunc(upd)) * BaseOfResult)
}

func (h *Hash63) Write(p []byte) (int, error) {
	h.hash.Reset()

	n, err := h.hash.Write(p)
	if err != nil {
		return 0, err
	}

	h.state = h.hash.Sum(nil)

	return n, err
}

// Sum appends the current hash to b and returns the resulting slice.
func (h *Hash63) Sum(b []byte) []byte {
	res := make([]byte, len(b))
	copy(res, b)

	for i := 0; i < h.bytesCount; i++ {
		res = append(res, h.GetByte())
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
