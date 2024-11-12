package rand

import (
	"crypto/rand"
)

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
