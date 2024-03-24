package lab1_gost34122015

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
