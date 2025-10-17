package signature

import (
	"testing"
)

func TestGenerateProtectCode(t *testing.T) {
	// Test basic generation
	code1, err := GenerateProtectCode()
	if err != nil {
		t.Fatalf("GenerateProtectCode failed: %v", err)
	}

	if len(code1) != 21 {
		t.Errorf("Expected protect code length 21, got %d", len(code1))
	}

	// Test uniqueness
	code2, err := GenerateProtectCode()
	if err != nil {
		t.Fatalf("GenerateProtectCode failed: %v", err)
	}

	if code1 == code2 {
		t.Error("Expected different protect codes, got same")
	}
}

func TestEncryptDecryptSignature(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"normal text", `{"id":"test123","name":"Test User"}`},
		{"empty object", "{}"},
		{"special characters", `{"intro":"Hello! 你好 🎉"}`},
		{"long text", `{"intro":"` + string(make([]byte, 1000)) + `"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptSignature(tc.data)
			if err != nil {
				t.Fatalf("EncryptSignature failed: %v", err)
			}

			if encrypted == "" {
				t.Error("Expected non-empty encrypted data")
			}

			// Decrypt
			decrypted, err := DecryptSignature(encrypted)
			if err != nil {
				t.Fatalf("DecryptSignature failed: %v", err)
			}

			if decrypted != tc.data {
				t.Errorf("Expected %q, got %q", tc.data, decrypted)
			}
		})
	}
}

func TestEncryptSignature_EmptyString(t *testing.T) {
	_, err := EncryptSignature("")
	if err == nil {
		t.Error("Expected error for empty string, got nil")
	}
}

func TestDecryptSignature_EmptyString(t *testing.T) {
	_, err := DecryptSignature("")
	if err == nil {
		t.Error("Expected error for empty string, got nil")
	}
}

func TestDecryptSignature_InvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"invalid base64", "not-valid-base64!@#$"},
		{"too short", "YWJj"},
		{"corrupted data", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptSignature(tc.data)
			if err == nil {
				t.Error("Expected error for invalid data, got nil")
			}
		})
	}
}

func TestEncryptSignature_Uniqueness(t *testing.T) {
	data := `{"id":"test","name":"User"}`

	encrypted1, err := EncryptSignature(data)
	if err != nil {
		t.Fatalf("EncryptSignature failed: %v", err)
	}

	encrypted2, err := EncryptSignature(data)
	if err != nil {
		t.Fatalf("EncryptSignature failed: %v", err)
	}

	// Different encryptions should produce different ciphertext (due to random nonce)
	if encrypted1 == encrypted2 {
		t.Error("Expected different encrypted outputs for same input")
	}

	// But both should decrypt to same plaintext
	decrypted1, _ := DecryptSignature(encrypted1)
	decrypted2, _ := DecryptSignature(encrypted2)

	if decrypted1 != data || decrypted2 != data {
		t.Error("Decryption failed to recover original data")
	}
}
