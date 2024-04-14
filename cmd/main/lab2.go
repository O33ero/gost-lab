package main

import (
	"encoding/hex"
	"fmt"
	"gost-lab/internal/lab2_gost34112012_256"
)

func main() {
	// any plaintext
	plaintext := []byte{
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
		0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,
		0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,
		0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
		0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
	}
	fmt.Printf("Plain:\n%s\n", hex.EncodeToString(plaintext))

	// test hash
	h256 := lab2_gost34112012_256.NewHash()
	h256.Write(plaintext)
	hash := h256.Sum(nil)
	fmt.Printf("Hash:\n%s\n", hex.EncodeToString(hash))
	// expect:
	//		9d151eefd8590b89
	//		daa6ba6cb74af927
	//		5dd051026bb149a4
	//		52fd84e5e57b5500

	// test kdf
	key := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	label := []byte{
		0x26, 0xbd, 0xb8, 0x78,
	}
	seed := []byte{
		0xaf, 0x21, 0x43, 0x41, 0x45, 0x65, 0x63, 0x78,
	}
	kdf := lab2_gost34112012_256.NewKDF(key)
	res := kdf.Derive(label, seed, 1)
	fmt.Printf("Derive:\n%s\n", hex.EncodeToString(res))
	// expect:
	//		a1aa5f7de402d7b3
	//		d323f2991c8d4534
	//		013137010a83754f
	//		d0af6d7cd4922ed9
}