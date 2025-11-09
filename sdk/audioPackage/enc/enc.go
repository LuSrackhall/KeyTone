package enc

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"KeyTone/signature"
)

// FixedSecret is the fixed secret prefix for album config enc/dec per spec.
const FixedSecret = "LuSrackhall_KeyTone_2024_Signature_66688868686688"

// DeriveKey derives a 32-byte AES key using SHA256(FixedSecret + last6(sha1(albumUUID))).
// Assumption: albumUUID is the directory name of the album folder unless specified otherwise.
func DeriveKey(albumUUID string) []byte {
	// sha1 hex of albumUUID
	h := sha1.Sum([]byte(albumUUID))
	hexStr := hex.EncodeToString(h[:]) // 40 chars
	// take last 6 chars
	if len(hexStr) < 6 {
		// extremely unlikely; guard anyway
		pad := strings.Repeat("0", 6-len(hexStr))
		hexStr = pad + hexStr
	}
	suffix := hexStr[len(hexStr)-6:]
	seed := FixedSecret + suffix
	sum := sha256.Sum256([]byte(seed))
	key := make([]byte, 32)
	copy(key, sum[:])
	return key
}

// EncryptConfigJSON encrypts a plaintext JSON string and returns hex-encoded ciphertext.
func EncryptConfigJSON(jsonPlain string, albumUUID string) (string, error) {
	key := DeriveKey(albumUUID)
	return signature.EncryptData(jsonPlain, key)
}

// EncryptConfigBytes 加密明文 JSON 并返回原始二进制（nonce + ciphertext）。
func EncryptConfigBytes(jsonPlain string, albumUUID string) ([]byte, error) {
	key := DeriveKey(albumUUID)
	hexCipher, err := signature.EncryptData(jsonPlain, key)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(hexCipher)
}

// DecryptConfigHex decrypts a hex-encoded ciphertext string and returns plaintext JSON.
func DecryptConfigHex(hexCipher string, albumUUID string) (string, error) {
	key := DeriveKey(albumUUID)
	return signature.DecryptData(hexCipher, key)
}

// DecryptConfigBytes 解密二进制密文（nonce + ciphertext），返回明文 JSON。
func DecryptConfigBytes(cipherBytes []byte, albumUUID string) (string, error) {
	key := DeriveKey(albumUUID)
	return signature.DecryptData(hex.EncodeToString(cipherBytes), key)
}

// IsLikelyHexCipher tries to quickly determine if the given content is a hex string (not JSON).
// It checks: trimmed does not start with '{' and consists only of hex chars with even length.
func IsLikelyHexCipher(content []byte) bool {
	s := strings.TrimSpace(string(content))
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return false
	}
	if len(s)%2 != 0 {
		return false
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	// quick sanity: attempt hex decode minimal length
	if len(s) < 32 { // smaller than typical nonce+ciphertext
		return false
	}
	return true
}

// ValidateJSONFast performs a minimal validation that string looks like JSON object.
func ValidateJSONFast(s string) error {
	st := strings.TrimSpace(s)
	if strings.HasPrefix(st, "{") && strings.HasSuffix(st, "}") {
		return nil
	}
	if strings.HasPrefix(st, "[") && strings.HasSuffix(st, "]") {
		return nil
	}
	return errors.New("not a json-looking string")
}
