package commands

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"ktalbum-tools/utils"
)

func Extract(inputFile, outputFile string, verbose bool) error {
	if verbose {
		fmt.Printf("正在解包: %s\n", inputFile)
	}

	// 读取输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer inFile.Close()

	// 读取并验证文件头
	var header utils.KeytoneFileHeader
	if err := binary.Read(inFile, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("读取文件头失败: %v", err)
	}

	if string(header.Signature[:]) != utils.KeytoneFileSignature {
		return fmt.Errorf("无效的文件格式：不是 KeyTone 专辑文件")
	}

	if verbose {
		fmt.Printf("文件版本: %d\n", header.Version)
		fmt.Printf("数据大小: %d bytes\n", header.DataSize)
	}

	// 读取加密数据
	encryptedData := make([]byte, header.DataSize)
	if _, err := io.ReadFull(inFile, encryptedData); err != nil {
		return fmt.Errorf("读取加密数据失败: %v", err)
	}

	// 解密数据（按版本选择候选密钥；私有密钥构建优先注入，回退默认，兼容开源产物）
	var zipData []byte
	decrypted := false

	for _, decryptKey := range utils.GetDecryptKeyCandidatesByVersion(header.Version) {
		candidate := utils.XorCrypt(encryptedData, decryptKey)
		checksum := sha256.Sum256(candidate)
		if checksum == header.Checksum {
			zipData = candidate
			decrypted = true
			break
		}
	}

	// 与 SDK 一致：若校验失败且版本不是 v1，尝试 v1 候选回退
	if !decrypted && header.Version != 1 {
		for _, decryptKey := range utils.GetDecryptKeyCandidatesByVersion(1) {
			candidate := utils.XorCrypt(encryptedData, decryptKey)
			checksum := sha256.Sum256(candidate)
			if checksum == header.Checksum {
				zipData = candidate
				decrypted = true
				break
			}
		}
	}

	if !decrypted {
		return fmt.Errorf("文件校验失败，文件可能已损坏或密钥不匹配")
	}

	// 写入解密后的zip数据
	if err := os.WriteFile(outputFile, zipData, 0644); err != nil {
		return fmt.Errorf("写入输出文件失败: %v", err)
	}

	if verbose {
		fmt.Printf("成功解包到: %s\n", outputFile)
	}

	return nil
} 