package lab4_crisp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gost-lab/internal/lab1_gost34122015"
	"gost-lab/internal/lab2_gost34112012_256"
	"gost-lab/internal/lab3_xorshiftplus"
	"runtime"
)

var (
	ExternalKeyIdFlagWithVersion = []byte{ // 1 bit + 15 bits
		0x00, 0x00,
	}
	CS = []byte{ // 8 bits
		0xf8,
	}
	KeyId = []byte{ // 8 bits
		0x80,
	}
)

const (
	BlockSize  = 16 // bytes
	KeySize    = 32 // bytes
	PacketSize = 56 // byte
)

type Crisp struct {
	Decoder    Decoder
	Encoder    Encoder
	randomSeed [16]byte
}

func (c *Crisp) Close() {
	for i := 0; i < len(c.randomSeed); i++ {
		c.randomSeed[i] = 0x00
	}
	c.Decoder.kdf.Close()
	c.Decoder.cipher.Close()
	c.Decoder.seqNum = 0
	c.Encoder.kdf.Close()
	c.Encoder.cipher.Close()
	c.Encoder.seqNum = 0
	runtime.GC()
	fmt.Printf("Clear mem [Crisp]: %p\n", &c)
}

type Decoder struct {
	random *lab3_xorshiftplus.XorShift128Plus
	kdf    *lab2_gost34112012_256.KDF
	cipher *lab1_gost34122015.CtrAcpkm
	seqNum uint32
}

type Encoder struct {
	random *lab3_xorshiftplus.XorShift128Plus
	kdf    *lab2_gost34112012_256.KDF
	cipher *lab1_gost34122015.CtrAcpkm
	seqNum uint32
}

type Message struct {
	ExternalKeyIdFlagWithVersion []byte
	CS                           []byte
	KeyId                        []byte
	SeqNum                       []byte
	Payload                      []byte
	ICV                          []byte
	Digits                       []byte
}

func New(key []byte, randomSeed [16]byte) *Crisp {
	if len(key) != KeySize {
		panic("Key size should be 32 bytes")
	}

	cipher := lab1_gost34122015.NewCtrAcpkm()
	kdf := lab2_gost34112012_256.NewKDF(key[:])
	return &Crisp{
		Decoder: Decoder{
			random: lab3_xorshiftplus.New(randomSeed),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		Encoder: Encoder{
			random: lab3_xorshiftplus.New(randomSeed),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		randomSeed: randomSeed,
	}
}

func (c *Crisp) Reset() {
	c.Encoder.seqNum = 0
	c.Encoder.random = lab3_xorshiftplus.New(c.randomSeed)
	c.Decoder.seqNum = 0
	c.Decoder.random = lab3_xorshiftplus.New(c.randomSeed)
}

func (c *Crisp) Encode(plainText []byte) []Message {
	var res []Message

	c.Reset()
	for i := 0; i < len(plainText); i += BlockSize {
		message := c.EncodeNextBlock(plainText[i : i+BlockSize])
		res = append(res, message)
	}

	return res
}

func (c *Crisp) EncodeNextBlock(plainText []byte) Message {
	if len(plainText) != BlockSize {
		panic("Block size should be 16 bytes")
	}
	e := c.Encoder

	var seqNum [4]byte
	var seed [8]byte
	binary.BigEndian.PutUint32(seqNum[:], e.seqNum)
	binary.BigEndian.PutUint64(seed[:], e.random.Next())

	// Key(N)
	key := e.kdf.Derive(seqNum[:], seed[:], 1)

	// text[N]
	block := plainText[:]

	// Payload(N), Mac(N)
	ciphertext, mac := e.cipher.Encrypt(block, key)

	var message []byte
	message = append(message, ExternalKeyIdFlagWithVersion...)
	message = append(message, CS...)
	message = append(message, KeyId...)
	message = append(message, seqNum[:]...)
	message = append(message, ciphertext[:]...)
	message = append(message, mac...)

	e.seqNum += 1 // complete current iteration and prepare next
	return Message{
		ExternalKeyIdFlagWithVersion: ExternalKeyIdFlagWithVersion,
		CS:                           CS,
		KeyId:                        KeyId,
		SeqNum:                       seqNum[:],
		Payload:                      ciphertext[:],
		ICV:                          mac[:],
		Digits:                       message,
	}
}

func (c *Crisp) Decode(cipherText [][]byte) [][]byte {
	for i, b := range cipherText {
		if len(b) != PacketSize {
			panic(fmt.Sprintf("Block size of block [%d] should be 56 bytes", i))
		}
	}

	var res [][]byte
	for _, b := range cipherText {
		decoded := c.DecodeNextBlock(b)
		res = append(res, decoded)
	}

	return res
}

func (c *Crisp) DecodeNextBlock(cipherText []byte) []byte {
	if len(cipherText) != PacketSize {
		panic("Block size should be equal 56 bytes")
	}
	d := c.Decoder

	var seqNum [4]byte
	var seed [8]byte
	binary.BigEndian.PutUint64(seed[:], d.random.Next())

	// parse
	seqNum = [4]byte(cipherText[4:8])
	payload := cipherText[8:24]
	mac := cipherText[24:56]

	key := d.kdf.Derive(seqNum[:], seed[:], 1)
	decrypt := d.cipher.Decrypt(payload, key, mac)

	return decrypt[:]
}

func (m *Message) String() string {
	format :=
		`Message:
    ExternalKeyIdFlagWithVersion: %s
    CS:                           %s
    KeyId:                        %s
    SeqNum:                       %s
    Payload:                      %s
    ICV:                          %s
    As block:                     %s`

	return fmt.Sprintf(format,
		hex.EncodeToString(m.ExternalKeyIdFlagWithVersion),
		hex.EncodeToString(m.CS),
		hex.EncodeToString(m.KeyId),
		hex.EncodeToString(m.SeqNum),
		hex.EncodeToString(m.Payload),
		hex.EncodeToString(m.ICV),
		hex.EncodeToString(m.Digits))
}
