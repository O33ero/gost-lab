package lab3_xorshiftplus

import "encoding/binary"

const (
	SeedSize = 16 // bytes
)

type XorShift128Plus struct {
	seed [SeedSize]byte
}

func New(seed [SeedSize]byte) *XorShift128Plus {
	return &XorShift128Plus{seed: seed}
}

func (x *XorShift128Plus) Next() uint64 {
	s1 := binary.BigEndian.Uint64(x.seed[0 : SeedSize/2]) // 8 bytes
	s0 := binary.BigEndian.Uint64(x.seed[SeedSize/2:])    // 8 bytes

	s1 ^= s1 << 23
	s1 = s1 ^ s0 ^ (s1 >> 18) ^ (s0 >> 5)

	// update the generator state
	binary.BigEndian.PutUint64(x.seed[:SeedSize/2], s0)
	binary.BigEndian.PutUint64(x.seed[SeedSize/2:], s1)

	return s1 + s0
}
