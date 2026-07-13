package random

import (
	"math/rand/v2"
	"unsafe"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func GenerateRandomStringFast(length int) string {
	if length <= 0 {
		return ""
	}

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	b := make([]byte, length)
	for i := range length {
		idx := rng.IntN(len(allowedChars))
		b[i] = allowedChars[idx]
	}

	return *(*string)(unsafe.Pointer(&b))
}
