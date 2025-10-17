package signature

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	gonanoid "github.com/jaevor/go-nanoid"
)

const encryptionKey = "KeyTone2024SecretKey_SignatureProtection!!!!" // 32 bytes for AES-256

// GenerateProtectCode generates a 21-character protect code using nanoid algorithm
func GenerateProtectCode() (string, error) {
	canonicGenerator, err := gonanoid.Standard(21)
	if err != nil {
		return "", err
	}
	return canonicGenerator(), nil
}

// EncryptSignature encrypts signature data using AES-256-GCM
func EncryptSignature(data string) (string, error) {
	if data == "" {
		return "", errors.New("data cannot be empty")
	}

	// Create AES cipher
	block, err := aes.NewCipher([]byte(encryptionKey)[:32])
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSignature decrypts signature data using AES-256-GCM
func DecryptSignature(encryptedData string) (string, error) {
	if encryptedData == "" {
		return "", errors.New("encrypted data cannot be empty")
	}

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// Create AES cipher
	block, err := aes.NewCipher([]byte(encryptionKey)[:32])
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get nonce size
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt and open
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
