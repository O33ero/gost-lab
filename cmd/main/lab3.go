package main

import (
	"encoding/binary"
	"fmt"
	"gost-lab/internal/lab3_xorshiftplus"
	"time"
)

func main() {
	var seed [16]byte
	nanosecond := time.Now().Nanosecond()
	binary.LittleEndian.PutUint64(seed[:], uint64(nanosecond))

	generator := lab3_xorshiftplus.New(seed)
	fmt.Printf("Next [0]: %d\n", generator.Next())
	fmt.Printf("Next [1]: %d\n", generator.Next())
	fmt.Printf("Next [2]: %d\n", generator.Next())
	fmt.Printf("Next [3]: %d\n", generator.Next())
	fmt.Printf("Next [4]: %d\n", generator.Next())
	fmt.Printf("Next [5]: %d\n", generator.Next())
	fmt.Printf("Next [6]: %d\n", generator.Next())
	fmt.Printf("Next [7]: %d\n", generator.Next())

	start := time.Now().UnixMilli()
	lab3_xorshiftplus.Create1MbFile(generator) // ~ 300 msec
	fmt.Printf("Elapsed [1mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3_xorshiftplus.Create100MbFile(generator) // ~ 28000 msec (28 sec)
	fmt.Printf("Elapsed [100mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3_xorshiftplus.Create1000MbFile(generator) // ~ 300000 msec (300 sec)
	fmt.Printf("Elapsed [1000mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3_xorshiftplus.Create1000Values(generator) // ~ 5 msec
	fmt.Printf("Elapsed [1000val]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3_xorshiftplus.Create10000Values(generator) // ~ 25 msec
	fmt.Printf("Elapsed [10000val]: %d msec.\n", time.Now().UnixMilli()-start)
}
