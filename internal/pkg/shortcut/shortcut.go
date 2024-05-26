package shortcut

import (
	"crypto/md5" // nolint: gosec // Нет потребности в криптобезопастности хэша
	"hash"
	"math"
	"strings"

	"github.com/pkg/errors"
)

const (
	alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	BaseOfResult = len(alphabet)

	A = 0.61803398874989484820
)

var ErrorBadHashFunction = errors.New("set hash function return bytes less than alphabet power")

type HashShortcut struct {
	hash hash.Hash
}

func NewHashShortcut() *HashShortcut {
	return &HashShortcut{hash: md5.New()} // nolint: gosec // Нет потребности в криптобезопастности хэша
}

func (h *HashShortcut) GetShort(s string, n int) string {
	if h.hash.Size() < n {
		panic(ErrorBadHashFunction)
	}

	_, _ = h.hash.Write([]byte(s))

	sum := h.hash.Sum(nil)

	short := strings.Builder{}
	short.Grow(n)

	for i := 0; i < n; i++ {
		upd := float64(sum[i]) * A

		indx := byte((upd - math.Trunc(upd)) * float64(BaseOfResult))

		_ = short.WriteByte(alphabet[indx])
	}

	return short.String()
}
