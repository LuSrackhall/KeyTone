package utils

import "time"

const (
	KeytoneFileSignature = "KTALBUM"
	KeytoneEncryptKey    = "KeyTone2024SecretKey"
)

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