package utils

import (
	"encoding/hex"
	"time"
)

const (
	KeytoneFileSignature = "KTALBUM"
)

// xorMask 用于混淆密钥的掩码，必须与 tools/key-obfuscator/main.go 中的一致
var xorMask = []byte{0x55, 0xAA, 0x33, 0xCC, 0x99, 0x66, 0x11, 0xEE, 0x77, 0xBB, 0x22, 0xDD, 0x88, 0x44, 0xFF, 0x00}

// 默认开源密钥
const (
	DefaultEncryptKeyV1 = "KeyTone2024SecretKey"
	DefaultEncryptKeyV2 = "KeyTone2025AlbumSecureEncryptionKeyV2"
)

// 可通过 -ldflags -X 注入 XOR 混淆后的 hex
var (
	KeytoneEncryptKeyV1      = DefaultEncryptKeyV1
	KeytoneEncryptKeyV2      = DefaultEncryptKeyV2
	KeytoneEncryptKeyCurrent = KeytoneEncryptKeyV2

	// 兼容旧代码：默认等同 v1
	KeytoneEncryptKey = KeytoneEncryptKeyV1
)

func deobfuscateString(obfuscatedHex string) string {
	obfuscated, err := hex.DecodeString(obfuscatedHex)
	if err != nil {
		return obfuscatedHex
	}

	realKey := make([]byte, len(obfuscated))
	for i, b := range obfuscated {
		realKey[i] = b ^ xorMask[i%len(xorMask)]
	}

	return string(realKey)
}

func GetEncryptKeyV1() string {
	if KeytoneEncryptKeyV1 == DefaultEncryptKeyV1 {
		return DefaultEncryptKeyV1
	}
	return deobfuscateString(KeytoneEncryptKeyV1)
}

func GetEncryptKeyV2() string {
	if KeytoneEncryptKeyV2 == DefaultEncryptKeyV2 {
		return DefaultEncryptKeyV2
	}
	return deobfuscateString(KeytoneEncryptKeyV2)
}

func GetEncryptKeyByVersion(version uint8) string {
	switch version {
	case 1:
		return GetEncryptKeyV1()
	case 2:
		return GetEncryptKeyV2()
	default:
		return GetEncryptKeyV2()
	}
}

// GetDecryptKeyCandidates 返回按顺序尝试的解密 key 列表：
// 1) 版本对应的“注入 key（若已注入）”
// 2) 版本对应的默认 key
// 3) 若版本 != 1，则追加 v1 的注入/默认 key 作为回退
func GetDecryptKeyCandidates(version uint8) []string {
	appendUnique := func(list []string, key string) []string {
		for _, v := range list {
			if v == key {
				return list
			}
		}
		return append(list, key)
	}

	addVersionCandidates := func(list []string, injectedVal string, defaultVal string, getter func() string) []string {
		if injectedVal != defaultVal {
			list = appendUnique(list, getter())
		}
		list = appendUnique(list, defaultVal)
		return list
	}

	keys := []string{}
	switch version {
	case 1:
		keys = addVersionCandidates(keys, KeytoneEncryptKeyV1, DefaultEncryptKeyV1, GetEncryptKeyV1)
	case 2:
		keys = addVersionCandidates(keys, KeytoneEncryptKeyV2, DefaultEncryptKeyV2, GetEncryptKeyV2)
	default:
		keys = addVersionCandidates(keys, KeytoneEncryptKeyV2, DefaultEncryptKeyV2, GetEncryptKeyV2)
	}

	if version != 1 {
		keys = addVersionCandidates(keys, KeytoneEncryptKeyV1, DefaultEncryptKeyV1, GetEncryptKeyV1)
	}

	return keys
}

type KeytoneFileHeader struct {
	Signature [7]byte
	Version   uint8
	DataSize  uint64
	Checksum  [32]byte
}

type KeytoneAlbumMeta struct {
	MagicNumber string    `json:"magicNumber"`
	Version     string    `json:"version"`
	ExportTime  time.Time `json:"exportTime"`
	AlbumUUID   string    `json:"albumUUID"`
	AlbumName   string    `json:"albumName"`
}
