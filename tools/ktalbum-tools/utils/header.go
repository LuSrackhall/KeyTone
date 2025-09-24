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
	MagicNumber      string    `json:"magicNumber"`
	Version          string    `json:"version"`
	ExportTime       time.Time `json:"exportTime"`
	AlbumUUID        string    `json:"albumUUID"`
	AlbumName        string    `json:"albumName"`
	// 新增作者相关字段
	AuthorName       string    `json:"authorName,omitempty"`       // 作者名称
	AuthorContact    string    `json:"authorContact,omitempty"`    // 联系方式文本
	AuthorContactImg string    `json:"authorContactImg,omitempty"` // 联系方式图片(MD5文件名)
	HistoryAuthors   []string  `json:"historyAuthors,omitempty"`   // 历史创作者列表
	AllowReExport    bool      `json:"allowReExport"`              // 是否允许二次导出，默认true
	ExportPassword   string    `json:"exportPassword,omitempty"`   // 导出密码(SHA512)
}