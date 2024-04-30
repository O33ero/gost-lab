package lab3_xorshiftplus

import (
	"encoding/binary"
	"os"
)

func Create1MbFile(rand *XorShift128Plus) {
	iters := 131072 // 8 byte * 131072 = 1mb
	createFile(rand, iters, "xorshift_1mb.bin")
}

func Create100MbFile(rand *XorShift128Plus) {
	iters := 13107200 // 8 byte * 13107200 = 100mb
	createFile(rand, iters, "xorshift_100mb.bin")
}

func Create1000MbFile(rand *XorShift128Plus) {
	iters := 131072000 // 8 byte * 131072000 = 1000mb
	createFile(rand, iters, "xorshift_1000mb.bin")
}

func Create1000Values(rand *XorShift128Plus) {
	iters := 1000 // 10^3
	createFile(rand, iters, "xorshift_1000values.bin")
}

func Create10000Values(rand *XorShift128Plus) {
	iters := 10000 // 10^4
	createFile(rand, iters, "xorshift_10000values.bin")
}

func createFile(rand *XorShift128Plus, iters int, filename string) {
	err := os.Truncate(filename, 0)
	if err != nil {
		panic("failed to clear existed file")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("failed to create file: " + err.Error())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for i := 0; i < iters; i++ {
		var nextBatch [8]byte
		value := rand.Next()
		binary.LittleEndian.PutUint64(nextBatch[:], value)

		_, err := file.Write(nextBatch[:])
		if err != nil {
			panic("failed to append to file: " + err.Error())
		}
	}
}
