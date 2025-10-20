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
	"time"
)

// SignatureData 用于存储签名信息的结构体
type SignatureData struct {
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	CardImage string `json:"cardImage"` // 完整的图片文件路径
	// CardImageExt string `json:"cardImageExt"` // 图片文件扩展名（用于前端显示）
}

// SignatureSortMetadata 签名排序元数据（仅用于本地排序，不参与导出）
type SignatureSortMetadata struct {
	Time int64 `json:"time"` // Unix 时间戳，仅在创建或导入时生成；更新签名不会改变此值；但可通过拖动排序操作更改
}

// SignatureStorageEntry 签名在配置文件中的存储结构
type SignatureStorageEntry struct {
	Value string                `json:"value"` // 加密的签名数据
	Sort  SignatureSortMetadata `json:"sort"`  // 排序元数据
}

// CreateSignature 创建新签名的处理函数
// id: 签名唯一标识（未加密）
// signatureData: 签名数据
// imageData: 图片文件的二进制数据
// imageExt: 图片文件的扩展名（如 .jpg, .png，可以为空）
// originalImageName: 图片原始文件名
// encryptionKey: 对称加密密钥
//
// 说明：
// - 生成的签名存储条目包含 value（加密数据）和 sort.time（Unix 时间戳）
// - sort.time 仅在首次创建或导入时生成，代表排序顺序
// - 更新签名时不会改变 sort.time（需要在更新逻辑中保留原值）
// - 用户的拖动排序操作可以改变 sort.time 值（需要单独的 API 实现）
func CreateSignature(id string, signatureData SignatureData, imageData []byte, imageExt string, originalImageName string, encryptionKey []byte) (string, error) {
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
		// 格式: id + "|" + name + "|" + originalImageName + "|" + unix时间戳
		fileNameSeed := fmt.Sprintf("%s|%s|%s|%d",
			id,
			signatureData.Name,
			originalImageName,
			time.Now().Unix(),
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
			"原始文件名", originalImageName,
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
	var signatureMap map[string]SignatureStorageEntry

	if existingValue != nil {
		// 如果已存在，将其转换为map[string]SignatureStorageEntry
		if m, ok := existingValue.(map[string]interface{}); ok {
			signatureMap = make(map[string]SignatureStorageEntry)
			for k, v := range m {
				// 尝试解析为新格式的 SignatureStorageEntry
				if entry, ok := v.(map[string]interface{}); ok {
					// 新格式
					storageEntry := SignatureStorageEntry{}
					if value, ok := entry["value"].(string); ok {
						storageEntry.Value = value
					}
					if sort, ok := entry["sort"].(map[string]interface{}); ok {
						if time, ok := sort["time"].(float64); ok {
							storageEntry.Sort.Time = int64(time)
						}
					}
					signatureMap[k] = storageEntry
				} else if str, ok := v.(string); ok {
					// 兼容旧格式：如果是字符串，创建新的存储条目
					logger.Warn("检测到旧格式的签名数据，正在进行格式升级", "key", k)
					signatureMap[k] = SignatureStorageEntry{
						Value: str,
						Sort: SignatureSortMetadata{
							Time: time.Now().Unix(), // 使用当前时间作为迁移时间戳
						},
					}
				}
			}
		}
	} else {
		signatureMap = make(map[string]SignatureStorageEntry)
	}

	// 添加新的签名数据，生成排序时间戳
	signatureMap[encryptedID] = SignatureStorageEntry{
		Value: encryptedValue,
		Sort: SignatureSortMetadata{
			Time: time.Now().Unix(), // 仅在创建时生成
		},
	}

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

// DecryptData 使用AES-GCM对数据进行对称解密（导出函数，供其他包使用）
func DecryptData(encryptedData string, key []byte) (string, error) {
	return decryptData(encryptedData, key)
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

// CleanupOrphanCardImages 清理不在配置中的孤立签名图片
// 该函数会扫描 signature 目录下的所有文件，
// 与配置中的签名数据进行比对，删除不在配置中的图片文件
func CleanupOrphanCardImages(encryptionKey []byte) error {
	logger.Info("开始执行签名名片图片清理操作...")

	// 1. 获取 signature 目录路径
	signatureDir := filepath.Join(config.ConfigPath, "signature")

	// 检查目录是否存在
	if _, err := os.Stat(signatureDir); os.IsNotExist(err) {
		logger.Info("签名目录不存在，无需执行清理操作", "path", signatureDir)
		return nil
	}

	// 2. 从配置中获取所有签名数据，解析出有效的图片路径集合
	validImagePaths := make(map[string]bool)
	signatureMapValue := config.GetValue("signature")

	if signatureMapValue != nil {
		if signatureMap, ok := signatureMapValue.(map[string]interface{}); ok {
			// 遍历所有的签名配置
			for _, v := range signatureMap {
				var encryptedValueStr string

				// 兼容新格式 SignatureStorageEntry
				if entry, ok := v.(map[string]interface{}); ok {
					if value, ok := entry["value"].(string); ok {
						encryptedValueStr = value
					} else {
						logger.Warn("无法从 SignatureStorageEntry 中提取 value 字段")
						continue
					}
				} else if str, ok := v.(string); ok {
					// 兼容旧格式：直接是加密字符串
					encryptedValueStr = str
				} else {
					logger.Warn("无法识别签名数据格式")
					continue
				}

				// 解密签名数据
				decryptedData, err := decryptData(encryptedValueStr, encryptionKey)
				if err != nil {
					logger.Warn("签名数据解密失败，跳过此签名", "error", err.Error())
					continue
				}

				// 解析 JSON 数据
				var sigData SignatureData
				if err := json.Unmarshal([]byte(decryptedData), &sigData); err != nil {
					logger.Warn("签名数据 JSON 解析失败，跳过此签名", "error", err.Error())
					continue
				}

				// 如果有图片路径，添加到有效路径集合
				if sigData.CardImage != "" {
					validImagePaths[sigData.CardImage] = true
					logger.Debug("记录有效的签名图片路径", "path", sigData.CardImage)
				}
			}
		}
	}

	logger.Info("配置中的有效签名图片数量", "count", len(validImagePaths))

	// 3. 获取 signature 目录下的所有文件
	files, err := os.ReadDir(signatureDir)
	if err != nil {
		logger.Error("读取签名目录失败", "error", err.Error())
		return err
	}

	// 4. 遍历目录中的所有文件，删除不在有效路径中的文件
	deletedCount := 0
	for _, file := range files {
		// 只处理文件，跳过目录
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(signatureDir, file.Name())

		// 检查文件路径是否在有效路径集合中
		if _, exists := validImagePaths[filePath]; !exists {
			// 文件不在有效路径中，删除它
			if err := os.Remove(filePath); err != nil {
				logger.Warn("删除孤立图片文件失败", "path", filePath, "error", err.Error())
			} else {
				logger.Debug("成功删除孤立图片文件", "path", filePath)
				deletedCount++
			}
		}
	}

	logger.Info("签名名片图片清理操作完成", "deleted_count", deletedCount)
	return nil
}
