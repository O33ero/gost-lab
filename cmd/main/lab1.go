package main

import (
	"bytes"
	"fmt"
	"gost-lab/internal/lab1_gost34122015"
	"os"
	"sync"
	"time"
)

func main() {
	key := [32]byte{
		0x80, 0x94, 0xA8, 0xBC, 0xC0, 0xD4, 0xE8, 0xFC,
		0x81, 0x95, 0xA9, 0xBD, 0xC1, 0xD5, 0xE9, 0xFD,
		0x82, 0x96, 0xAA, 0xBE, 0xC2, 0xD6, 0xEA, 0xFE,
		0x83, 0x97, 0xAB, 0xBF, 0xC3, 0xD7, 0xEB, 0xFF,
	}

	plainText := [16]byte{
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0xAA, 0xBB, 0xCC, 0xDD,
	}
	fmt.Printf("Plain: %v\n", plainText)

	b, err := os.ReadFile("xorshift_1000mb.bin")
	if err != nil {
		panic("failed to read file: " + err.Error())
	}

	cipher := lab1_gost34122015.NewCipher(key[:])
	defer cipher.Close()

	var wg sync.WaitGroup
	start := time.Now().UnixMilli()
	for i := 0; i < len(b); i += 16 {
		wg.Add(1)
		go func(part int) {
			defer wg.Done()
			encrypted := cipher.Encrypt(b[part : part+16])
			decrypt := cipher.Decrypt(encrypted[:])

			if !bytes.Equal(b[part:part+16], decrypt[:]) {
				panic("incorrect decrypt")
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("Complete in %d msec.\n", time.Now().UnixMilli()-start)
}
