package utils

import (
	"encoding/hex"
	"time"
)

const (
	KeytoneFileSignature = "KTALBUM"
)

// ==============================
// 专辑导出文件对称密钥（版本化）
// 注意：这些变量不再是 const，而是 var，以便在编译时通过 -ldflags 进行注入
// 注入的值应为经过 XOR 混淆后的 Hex 字符串（与 SDK 授权流一致）
// ==============================

// xorMask 用于混淆密钥的掩码，必须与 SDK 授权流一致
var xorMask = []byte{0x55, 0xAA, 0x33, 0xCC, 0x99, 0x66, 0x11, 0xEE, 0x77, 0xBB, 0x22, 0xDD, 0x88, 0x44, 0xFF, 0x00}

// 默认开源密钥常量（明文）
const (
	DefaultKeytoneEncryptKeyV1 = "KeyTone2024SecretKey"                  // v1
	DefaultKeytoneEncryptKeyV2 = "KeyTone2025AlbumSecureEncryptionKeyV2" // v2
)

// 版本化加密密钥（可注入）
var (
	KeytoneEncryptKeyV1 = DefaultKeytoneEncryptKeyV1
	KeytoneEncryptKeyV2 = DefaultKeytoneEncryptKeyV2
	// 向后兼容：旧变量名仍保留（等价于 v1）
	KeytoneEncryptKey = KeytoneEncryptKeyV1
)

func deobfuscateString(obfuscatedHex string) string {
	obfuscated, err := hex.DecodeString(obfuscatedHex)
	if err != nil {
		return obfuscatedHex
	}
	realBytes := make([]byte, len(obfuscated))
	for i, b := range obfuscated {
		realBytes[i] = b ^ xorMask[i%len(xorMask)]
	}
	return string(realBytes)
}

func getPlainEncryptKey(value string, defaultValue string) string {
	if value == defaultValue {
		return defaultValue
	}
	return deobfuscateString(value)
}

// GetEncryptKeyByVersion 根据文件头版本号返回对应的明文密钥。
func GetEncryptKeyByVersion(version uint8) string {
	switch version {
	case 1:
		return getPlainEncryptKey(KeytoneEncryptKeyV1, DefaultKeytoneEncryptKeyV1)
	case 2:
		return getPlainEncryptKey(KeytoneEncryptKeyV2, DefaultKeytoneEncryptKeyV2)
	default:
		// 未知版本：保守回退到 v2（与 SDK 保持一致的“当前版本”语义）
		return getPlainEncryptKey(KeytoneEncryptKeyV2, DefaultKeytoneEncryptKeyV2)
	}
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