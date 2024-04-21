package lab1_gost34122015

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
)

type CtrAcpkm struct {
	InitialVector [BlockSize]byte
}

func NewCtrAcpkm() *CtrAcpkm {
	var iv [BlockSize]byte

	_, err := rand.Read(iv[:])
	if err != nil {
		panic(err)
	}

	return &CtrAcpkm{iv}
}

func (c *CtrAcpkm) Encrypt(plaintext, key []byte) ([BlockSize]byte, []byte) {
	var gamma []byte
	var ciphertext [BlockSize]byte
	var mac []byte

	gamma = initGamma(c.InitialVector[:], key)

	xor(ciphertext[:], plaintext, gamma)
	mac = createVerificationCode(ciphertext[:], key)

	return ciphertext, mac
}

func (c *CtrAcpkm) Decrypt(ciphertext, key, mac []byte) [BlockSize]byte {
	var gamma []byte
	var plaintext [BlockSize]byte
	var expectedMac []byte

	expectedMac = createVerificationCode(ciphertext[:], key)
	if !verifyVerificationCode(expectedMac, mac) {
		panic("Expected MAC isn't equal to received MAC")
	}

	gamma = initGamma(c.InitialVector[:], key)

	xor(plaintext[:], ciphertext, gamma)

	return plaintext
}

func initGamma(initialVector, key []byte) []byte {
	var gamma []byte

	cipher := NewCipher(key[:])
	encoded := cipher.Encrypt(initialVector[:])

	gamma = append(gamma, encoded[:]...)

	return gamma
}

func createVerificationCode(ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(ciphertext)
	ciphertextMac := mac.Sum(nil)

	return ciphertextMac
}

func verifyVerificationCode(expectedMac, mac []byte) bool {
	if string(expectedMac) != string(mac) {
		panic("Expected MAC isn't equal to received MAC")
	}

	return true
}
