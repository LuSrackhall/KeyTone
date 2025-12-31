package commands

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"ktalbum-tools/utils"
)

type FileInfo struct {
	Name       string    `json:"name"`
	Version    uint8     `json:"version"`
	ExportTime string    `json:"exportTime"`
	AlbumUUID  string    `json:"albumUUID"`
}

func GetFileInfo(filePath string) (*FileInfo, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件头
	var header utils.KeytoneFileHeader
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("读取文件头失败: %v", err)
	}

	// 验证文件签名
	if string(header.Signature[:]) != utils.KeytoneFileSignature {
		return nil, fmt.Errorf("无效的文件格式：不是 KeyTone 专辑文件")
	}

	// 读取加密数据
	encryptedData := make([]byte, header.DataSize)
	if _, err := io.ReadFull(file, encryptedData); err != nil {
		return nil, fmt.Errorf("读取加密数据失败: %v", err)
	}

	// 解密数据（按版本选择密钥）
	decryptKey := utils.GetEncryptKeyByVersion(header.Version)
	zipData := utils.XorCrypt(encryptedData, decryptKey)

	// 验证校验和
	checksum := utils.CalculateChecksum(zipData)
	if !bytes.Equal(checksum[:], header.Checksum[:]) {
		// 与 SDK 一致：若版本不是 v1，尝试使用 v1 密钥回退
		if header.Version != 1 {
			zipData = utils.XorCrypt(encryptedData, utils.GetEncryptKeyByVersion(1))
			checksum = utils.CalculateChecksum(zipData)
		}
		if !bytes.Equal(checksum[:], header.Checksum[:]) {
			return nil, fmt.Errorf("文件校验失败，文件可能已损坏或密钥不匹配")
		}
	}

	// 从 zip 数据中读取 .keytone-album 文件
	zipReader, err := utils.ReadZipData(zipData)
	if err != nil {
		return nil, fmt.Errorf("读取 zip 数据失败: %v", err)
	}

	var meta utils.KeytoneAlbumMeta
	for _, f := range zipReader.File {
		if f.Name == ".keytone-album" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开元数据文件失败: %v", err)
			}
			defer rc.Close()

			if err := json.NewDecoder(rc).Decode(&meta); err != nil {
				return nil, fmt.Errorf("解析元数据失败: %v", err)
			}
			break
		}
	}

	if meta.AlbumName == "" {
		return nil, fmt.Errorf("未找到专辑元数据")
	}

	return &FileInfo{
		Name:       meta.AlbumName,
		Version:    header.Version,
		ExportTime: meta.ExportTime.Format("2006-01-02 15:04:05"),
		AlbumUUID:  meta.AlbumUUID,
	}, nil
} 