package lab1_gost34122015

import (
	"fmt"
	"runtime"
)

const (
	BlockSize       = 16 // bytes
	KeySize         = 32 // bytes
	Rounds          = 9
	KeyConstantSize = 32
	KeySetSize      = 10
)

var (
	// L-function vector
	lVec = [BlockSize]byte{
		0x94, 0x20, 0x85, 0x10, 0xC2, 0xC0, 0x01, 0xFB,
		0x01, 0xC0, 0xC2, 0x10, 0x85, 0x20, 0x94, 0x01,
	}

	// S-substitution
	sBox = [256]byte{
		0xFC, 0xEE, 0xDD, 0x11, 0xCF, 0x6E, 0x31, 0x16,
		0xFB, 0xC4, 0xFA, 0xDA, 0x23, 0xC5, 0x04, 0x4D,
		0xE9, 0x77, 0xF0, 0xDB, 0x93, 0x2E, 0x99, 0xBA,
		0x17, 0x36, 0xF1, 0xBB, 0x14, 0xCD, 0x5F, 0xC1,
		0xF9, 0x18, 0x65, 0x5A, 0xE2, 0x5C, 0xEF, 0x21,
		0x81, 0x1C, 0x3C, 0x42, 0x8B, 0x01, 0x8E, 0x4F,
		0x05, 0x84, 0x02, 0xAE, 0xE3, 0x6A, 0x8F, 0xA0,
		0x06, 0x0B, 0xED, 0x98, 0x7F, 0xD4, 0xD3, 0x1F,
		0xEB, 0x34, 0x2C, 0x51, 0xEA, 0xC8, 0x48, 0xAB,
		0xF2, 0x2A, 0x68, 0xA2, 0xFD, 0x3A, 0xCE, 0xCC,
		0xB5, 0x70, 0x0E, 0x56, 0x08, 0x0C, 0x76, 0x12,
		0xBF, 0x72, 0x13, 0x47, 0x9C, 0xB7, 0x5D, 0x87,
		0x15, 0xA1, 0x96, 0x29, 0x10, 0x7B, 0x9A, 0xC7,
		0xF3, 0x91, 0x78, 0x6F, 0x9D, 0x9E, 0xB2, 0xB1,
		0x32, 0x75, 0x19, 0x3D, 0xFF, 0x35, 0x8A, 0x7E,
		0x6D, 0x54, 0xC6, 0x80, 0xC3, 0xBD, 0x0D, 0x57,
		0xDF, 0xF5, 0x24, 0xA9, 0x3E, 0xA8, 0x43, 0xC9,
		0xD7, 0x79, 0xD6, 0xF6, 0x7C, 0x22, 0xB9, 0x03,
		0xE0, 0x0F, 0xEC, 0xDE, 0x7A, 0x94, 0xB0, 0xBC,
		0xDC, 0xE8, 0x28, 0x50, 0x4E, 0x33, 0x0A, 0x4A,
		0xA7, 0x97, 0x60, 0x73, 0x1E, 0x00, 0x62, 0x44,
		0x1A, 0xB8, 0x38, 0x82, 0x64, 0x9F, 0x26, 0x41,
		0xAD, 0x45, 0x46, 0x92, 0x27, 0x5E, 0x55, 0x2F,
		0x8C, 0xA3, 0xA5, 0x7D, 0x69, 0xD5, 0x95, 0x3B,
		0x07, 0x58, 0xB3, 0x40, 0x86, 0xAC, 0x1D, 0xF7,
		0x30, 0x37, 0x6B, 0xE4, 0x88, 0xD9, 0xE7, 0x89,
		0xE1, 0x1B, 0x83, 0x49, 0x4C, 0x3F, 0xF8, 0xFE,
		0x8D, 0x53, 0xAA, 0x90, 0xCA, 0xD8, 0x85, 0x61,
		0x20, 0x71, 0x67, 0xA4, 0x2D, 0x2B, 0x09, 0x5B,
		0xCB, 0x9B, 0x25, 0xD0, 0xBE, 0xE5, 0x6C, 0x52,
		0x59, 0xA6, 0x74, 0xD2, 0xE6, 0xF4, 0xB4, 0xC0,
		0xD1, 0x66, 0xAF, 0xC2, 0x39, 0x4B, 0x63, 0xB6,
	}

	// S-substitution (inverse)
	sBoxInv = [256]byte{
		0xA5, 0x2D, 0x32, 0x8F, 0x0E, 0x30, 0x38, 0xC0,
		0x54, 0xE6, 0x9E, 0x39, 0x55, 0x7E, 0x52, 0x91,
		0x64, 0x03, 0x57, 0x5A, 0x1C, 0x60, 0x07, 0x18,
		0x21, 0x72, 0xA8, 0xD1, 0x29, 0xC6, 0xA4, 0x3F,
		0xE0, 0x27, 0x8D, 0x0C, 0x82, 0xEA, 0xAE, 0xB4,
		0x9A, 0x63, 0x49, 0xE5, 0x42, 0xE4, 0x15, 0xB7,
		0xC8, 0x06, 0x70, 0x9D, 0x41, 0x75, 0x19, 0xC9,
		0xAA, 0xFC, 0x4D, 0xBF, 0x2A, 0x73, 0x84, 0xD5,
		0xC3, 0xAF, 0x2B, 0x86, 0xA7, 0xB1, 0xB2, 0x5B,
		0x46, 0xD3, 0x9F, 0xFD, 0xD4, 0x0F, 0x9C, 0x2F,
		0x9B, 0x43, 0xEF, 0xD9, 0x79, 0xB6, 0x53, 0x7F,
		0xC1, 0xF0, 0x23, 0xE7, 0x25, 0x5E, 0xB5, 0x1E,
		0xA2, 0xDF, 0xA6, 0xFE, 0xAC, 0x22, 0xF9, 0xE2,
		0x4A, 0xBC, 0x35, 0xCA, 0xEE, 0x78, 0x05, 0x6B,
		0x51, 0xE1, 0x59, 0xA3, 0xF2, 0x71, 0x56, 0x11,
		0x6A, 0x89, 0x94, 0x65, 0x8C, 0xBB, 0x77, 0x3C,
		0x7B, 0x28, 0xAB, 0xD2, 0x31, 0xDE, 0xC4, 0x5F,
		0xCC, 0xCF, 0x76, 0x2C, 0xB8, 0xD8, 0x2E, 0x36,
		0xDB, 0x69, 0xB3, 0x14, 0x95, 0xBE, 0x62, 0xA1,
		0x3B, 0x16, 0x66, 0xE9, 0x5C, 0x6C, 0x6D, 0xAD,
		0x37, 0x61, 0x4B, 0xB9, 0xE3, 0xBA, 0xF1, 0xA0,
		0x85, 0x83, 0xDA, 0x47, 0xC5, 0xB0, 0x33, 0xFA,
		0x96, 0x6F, 0x6E, 0xC2, 0xF6, 0x50, 0xFF, 0x5D,
		0xA9, 0x8E, 0x17, 0x1B, 0x97, 0x7D, 0xEC, 0x58,
		0xF7, 0x1F, 0xFB, 0x7C, 0x09, 0x0D, 0x7A, 0x67,
		0x45, 0x87, 0xDC, 0xE8, 0x4F, 0x1D, 0x4E, 0x04,
		0xEB, 0xF8, 0xF3, 0x3E, 0x3D, 0xBD, 0x8A, 0x88,
		0xDD, 0xCD, 0x0B, 0x13, 0x98, 0x02, 0x93, 0x80,
		0x90, 0xD0, 0x24, 0x34, 0xCB, 0xED, 0xF4, 0xCE,
		0x99, 0x10, 0x44, 0x40, 0x92, 0x3A, 0x01, 0x26,
		0x12, 0x1A, 0x48, 0x68, 0xF5, 0x81, 0x8B, 0xC7,
		0xD6, 0x20, 0x0A, 0x08, 0x00, 0x4C, 0xD7, 0x74,
	}
	// pre-calVeculated generative key constants (init in initKeyConstant() function)
	keyConstant [32]*[BlockSize]byte
	// pre-calVeculated multiplication of each to each byte (init in initGfCache() function)
	gfCache [256][256]byte
)

type Cipher struct {
	keySet [KeySetSize][BlockSize]byte
}

func (c *Cipher) Close() {
	for i := 0; i < len(c.keySet); i++ {
		for j := 0; j < len(c.keySet[0]); j++ {
			c.keySet[i][j] = 0x00
		}
	}
	runtime.GC()
	fmt.Printf("Clear mem [Cipher]: %p\n", &c)
}

func init() {
	initGfCache()
	initKeyConstant()
}

// pre-calVeculate generative key constants
func initKeyConstant() {
	for i := 0; i < KeyConstantSize; i++ {
		keyConstant[i] = new([BlockSize]byte)
		keyConstant[i][15] = byte(i) + 1
		l(keyConstant[i])
	}
}

// pre-calVeculate multiply of each to each bytes in GF(2^8)
func initGfCache() {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			gfCache[a][b] = gf(byte(a), byte(b))
		}
	}
}

// Multiplication in GF(2^8) with with P(x)=x^8+x^7+x^6+x+1.
// Used in L-function
func gf(a, b byte) (c byte) {
	for b > 0 {
		if b&1 > 0 {
			c ^= a
		}
		if a&0x80 > 0 {
			a = (a << 1) ^ 0xC3
		} else {
			a <<= 1
		}
		b >>= 1
	}
	return
}

// L-function
func l(block *[BlockSize]byte) {
	for n := 0; n < BlockSize; n++ {
		var t byte
		for i := 0; i < BlockSize; i++ {
			t ^= gfCache[block[i]][lVec[i]]
		}
		for i := BlockSize - 1; i > 0; i-- {
			block[i] = block[i-1]
		}
		block[0] = t
	}
}

// L-function (inverse)
func lInverse(block *[BlockSize]byte) {
	var t byte
	for n := 0; n < BlockSize; n++ {
		t = block[0]
		copy(block[:], block[1:])
		for i := 0; i < BlockSize-1; i++ {
			t ^= gfCache[block[i]][lVec[i]]
		}
		block[15] = t
	}
}

// S-substitute
func s(block *[BlockSize]byte) {
	for i := 0; i < BlockSize; i++ { // substitute byte by S-Box
		block[i] = sBox[int(block[i])]
	}
}

// S-substitute (inverse)
func sInverse(block *[BlockSize]byte) {
	for i := 0; i < BlockSize; i++ { // substitute byte by inverse S-Box
		block[i] = sBoxInv[int(block[i])]
	}
}

func NewCipher(key []byte) *Cipher {
	if len(key) != KeySize {
		panic("invalid key size, expected key with 32 bytes")
	}
	var (
		keySet       [KeySetSize][BlockSize]byte
		keyRoundEven [BlockSize]byte
		keyRoundOdd  [BlockSize]byte
		keyRoundN    [BlockSize]byte
	)
	copy(keyRoundEven[:], key[:BlockSize])
	copy(keyRoundOdd[:], key[BlockSize:])
	copy(keySet[0][:], keyRoundEven[:]) // K0
	copy(keySet[1][:], keyRoundOdd[:])  // K1

	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			xor(keyRoundN[:], keyRoundEven[:], keyConstant[8*i+j][:])
			s(&keyRoundN)
			l(&keyRoundN)
			xor(keyRoundN[:], keyRoundN[:], keyRoundOdd[:])
			copy(keyRoundOdd[:], keyRoundEven[:])
			copy(keyRoundEven[:], keyRoundN[:])
		}
		copy(keySet[2+2*i][:], keyRoundEven[:])  // [K2, K4, K8, K10] - all even
		copy(keySet[2+2*i+1][:], keyRoundOdd[:]) // [K3, K5, K7, K9] - all odd
	}

	return &Cipher{keySet}
}

func (c *Cipher) Encrypt(src []byte) *[BlockSize]byte {
	if len(src) != BlockSize {
		panic("invalid block size, expected block with 16 bytes")
	}

	result := new([BlockSize]byte)
	block := new([BlockSize]byte)
	copy(block[:], src)
	for i := 0; i < Rounds; i++ {
		xor(block[:], block[:], c.keySet[i][:])
		s(block)
		l(block)
	}
	xor(block[:], block[:], c.keySet[9][:])
	copy(result[:], block[:])

	return result
}

func (c *Cipher) Decrypt(src []byte) *[BlockSize]byte {
	if len(src) != BlockSize {
		panic("invalid block size, expected block with 16 bytes")
	}

	result := new([BlockSize]byte)
	block := new([BlockSize]byte)
	copy(block[:], src)
	for i := Rounds; i > 0; i-- {
		xor(block[:], block[:], c.keySet[i][:])
		lInverse(block)
		sInverse(block)
	}
	xor(result[:], block[:], c.keySet[0][:])

	return result
}

func xor(dst, src1, src2 []byte) {
	if len(dst) != BlockSize {
		panic("dst is not 16 bytes")
	}
	if len(src1) != BlockSize {
		panic("src1 is not 16 bytes")
	}
	if len(src2) != BlockSize {
		panic("src2 is not 16 bytes")
	}
	for i := 0; i < BlockSize; i++ {
		dst[i] = src1[i] ^ src2[i]
	}
}
