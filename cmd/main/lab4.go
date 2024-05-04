package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gost-lab/internal/control"
	"gost-lab/internal/lab4_crisp"
	"os"
	"time"
)

const (
	BlockSize = 16 // bytes
	KeySize   = 32 // bytes
)

func main() {
	ec := control.NewExecuteControl()
	defer ec.Wait()

	// input:
	// - plaintext
	// - key

	// output:
	// - message = Message(N)

	// Message(N):
	// return [
	//     ExternalKeyIdFlagWithVersion +  // 2 byte   [0:2]
	//     CS +							   // 1 byte   [2:3]
	//     KeyId +                         // 1 byte   [3:4]
	//     SeqNum(N) +                     // 4 bytes  [4:8]
	//     Payload(N) +                    // 16 bytes [8:24]
	//     Mac(N)                          // 32 bytes [24:56]
	// ]                                   // As sum: 56 bytes

	// SeqNum(N):
	// return (uint64) N

	// Payload(N):
	// return cipher(text[N], Key(N)).ciphertext

	// Mac(N):
	// return cipher(text[N], Key(N)).mac

	// Key(N):
	// return KDF(key, SeqNum(N), Random())

	// Random():
	// return Xorshift(time.toNano())

	// plainText
	plainText := [BlockSize * 4]byte{
		// part 0
		0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0, 0xA0, 0xB0,
		0xC0, 0xD0, 0xE0, 0xF0, 0xA0, 0xB0, 0xC0, 0xD0,

		// part 1
		0xA1, 0xB1, 0xC1, 0xD1, 0xE1, 0xF1, 0xA1, 0xB1,
		0xC1, 0xD1, 0xE1, 0xF1, 0xA1, 0xB1, 0xC1, 0xD1,

		// part 2
		0xA2, 0xB2, 0xC2, 0xD2, 0xE2, 0xF2, 0xA2, 0xB2,
		0xC2, 0xD2, 0xE2, 0xF2, 0xA2, 0xB2, 0xC2, 0xD2,

		// part 3
		0xA3, 0xB3, 0xC3, 0xD3, 0xE3, 0xF3, 0xA3, 0xB3,
		0xC3, 0xD3, 0xC3, 0xF3, 0xA3, 0xB3, 0xC3, 0xD3,
	}
	fmt.Printf("Plain: %s\n", hex.EncodeToString(plainText[:]))

	// key
	key := [KeySize]byte{
		0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
		0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
		0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
		0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF,
	}

	// encode
	var randomSeed [16]byte
	binary.BigEndian.PutUint16(randomSeed[:], uint16(time.Now().Nanosecond()))
	crisp := lab4_crisp.New(key[:], randomSeed)
	defer crisp.Close()

	b, err := os.ReadFile("xorshift_1mb.bin")
	if err != nil {
		panic("failed to open file: " + err.Error())
	}

	start := time.Now().UnixMilli()
	for i := 0; i < len(b); i += BlockSize {
		message := crisp.EncodeNextBlock(b[i : i+BlockSize])
		decoded := crisp.DecodeNextBlock(message.Digits)

		if !bytes.Equal(b[i:i+BlockSize], decoded) {
			panic("incorrect decode")
		}

		if i%(4096*4) == 0 {
			fmt.Printf("Complete %d/%d [%.2f%%]\n", i, len(b), float64(i)/float64(len(b))*100)
		}
	}
	fmt.Printf("Elapsed: %d msec.\n", time.Now().UnixMilli()-start)
}
