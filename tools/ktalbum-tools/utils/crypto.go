package utils

import (
	"crypto/sha256"
)

func XorCrypt(data []byte, key string) []byte {
	keyBytes := []byte(key)
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}
	return result
}

func CalculateChecksum(data []byte) [32]byte {
	return sha256.Sum256(data)
}