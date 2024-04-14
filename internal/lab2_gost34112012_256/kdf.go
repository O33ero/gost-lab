package lab2_gost34112012_256

import (
	"crypto/hmac"
	"hash"
)

type KDF struct {
	h hash.Hash
}

func NewKDF(key []byte) *KDF {
	return &KDF{hmac.New(newHash, key)}
}

func (kdf *KDF) Derive(label, seed []byte, r int) (res []byte) {
	if r < 0 || r > 4 {
		panic("R should be between 1 and 4 inclusive")
	}

	for i := 1; i <= r; i++ {
		kdf.h.Write([]byte{byte(i)})
		kdf.h.Write(label)
		kdf.h.Write([]byte{0x00})
		kdf.h.Write(seed)
		kdf.h.Write([]byte{0x01})
		kdf.h.Write([]byte{0x00})
	}

	res = kdf.h.Sum(nil)
	kdf.h.Reset()
	return res
}

func newHash() hash.Hash {
	return NewHash()
}
