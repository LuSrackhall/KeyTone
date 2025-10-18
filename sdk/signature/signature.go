/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package signature

import (
	"KeyTone/config"
	"KeyTone/logger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SignatureData 用于存储签名信息的结构体
type SignatureData struct {
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	CardImage string `json:"cardImage"` // 完整的图片文件路径
	// CardImageExt string `json:"cardImageExt"` // 图片文件扩展名（用于前端显示）
}

// CreateSignature 创建新签名的处理函数
// id: 签名唯一标识（未加密）
// signatureData: 签名数据
// imageData: 图片文件的二进制数据
// imageExt: 图片文件的扩展名（如 .jpg, .png，可以为空, 不包含点号）
// encryptionKey: 对称加密密钥
func CreateSignature(id string, signatureData SignatureData, imageData []byte, imageExt string, encryptionKey []byte) (string, error) {
	// 1. 对ID进行对称加密
	encryptedID, err := encryptData(id, encryptionKey)
	if err != nil {
		logger.Error("ID加密失败", "error", err.Error())
		return "", err
	}

	// 2. 处理图片：使用uuid命名策略, 存入文件系统
	// 使用"id+名称+原始图片名称+unix时间戳"经过SHA-1哈希后作为文件名
	if len(imageData) > 0 {
		// 构建用于生成文件名的字符串
		// 格式: id + "|" + name
		fileNameSeed := fmt.Sprintf("%s|%s",
			id,
			signatureData.Name,
		)

		logger.Debug("图片文件名种子生成",
			"seed", fileNameSeed,
		)

		// 对文件名种子的字符串计算SHA-1哈希
		sha1Hash := sha1.Sum([]byte(fileNameSeed))
		imageFileName := hex.EncodeToString(sha1Hash[:])

		// 添加扩展名
		if imageExt != "" {
			// 确保扩展名以点开头
			if !strings.HasPrefix(imageExt, ".") {
				imageExt = "." + imageExt
			}
			imageFileName = imageFileName + imageExt
		}

		logger.Debug("图片文件名生成完成",
			"种子文件名", fileNameSeed,
			"文件名", imageFileName,
		)

		// 创建 signature 目录
		signatureDir := filepath.Join(config.ConfigPath, "signature")
		if err := os.MkdirAll(signatureDir, os.ModePerm); err != nil {
			logger.Error("创建signature目录失败", "error", err.Error())
			return "", err
		}

		// 保存图片文件
		imagePath := filepath.Join(signatureDir, imageFileName)
		if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
			logger.Error("保存图片文件失败", "error", err.Error())
			return "", err
		}

		// 更新签名数据中的图片路径和扩展名
		signatureData.CardImage = imagePath
		// signatureData.CardImageExt = imageExt // 存储扩展名供前端使用
	}

	// 3. 对签名数据结构进行JSON序列化
	jsonData, err := json.Marshal(signatureData)
	if err != nil {
		logger.Error("签名数据JSON序列化失败", "error", err.Error())
		return "", err
	}

	// 4. 对JSON字符串进行对称加密
	encryptedValue, err := encryptData(string(jsonData), encryptionKey)
	if err != nil {
		logger.Error("签名数据加密失败", "error", err.Error())
		return "", err
	}

	// 5. 以调试级别打印加密前后的数据
	logger.Debug("签名加密处理完成",
		"原始ID", id,
		"加密后ID", encryptedID,
		"原始Value", string(jsonData),
		"加密后Value", encryptedValue,
	)

	// 6. 调用config包的SetValue函数，存储加密后的key:value键值对
	// 获取现有的signature配置值（如果存在）
	existingValue := config.GetValue("signature")
	var signatureMap map[string]string

	if existingValue != nil {
		// 如果已存在，将其转换为map
		if m, ok := existingValue.(map[string]interface{}); ok {
			signatureMap = make(map[string]string)
			for k, v := range m {
				if str, ok := v.(string); ok {
					signatureMap[k] = str
				}
			}
		}
	} else {
		signatureMap = make(map[string]string)
	}

	// 添加新的签名数据
	signatureMap[encryptedID] = encryptedValue

	// 存储回配置文件
	config.SetValue("signature", signatureMap)

	return encryptedID, nil
}

// encryptData 使用AES-GCM对数据进行对称加密
func encryptData(data string, key []byte) (string, error) {
	// 确保密钥长度正确（16, 24, 或 32 字节）
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("加密密钥长度错误: 应为 16, 24 或 32 字节，实际为 %d 字节", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// 返回16进制编码的密文
	return hex.EncodeToString(ciphertext), nil
}

// decryptData 使用AES-GCM对数据进行对称解密
func decryptData(encryptedData string, key []byte) (string, error) {
	// 确保密钥长度正确
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("解密密钥长度错误: 应为 16, 24 或 32 字节，实际为 %d 字节", len(key))
	}

	// 解码16进制
	ciphertext, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 提取nonce和密文
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("密文长度过短")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
