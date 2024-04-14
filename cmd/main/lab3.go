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
}
