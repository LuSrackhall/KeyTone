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

// SignatureExportData 用于导出的签名数据结构，仅包含必要信息
// Key: 签名唯一标识（已加密）
// Name: 签名名称
// Intro: 签名介绍
// CardImage: 名片图片数据（Base64 编码）
// CardImageName: 名片图片文件名（包含扩展名，如 "avatar.jpg"）
type SignatureExportData struct {
	Key           string `json:"key"`           // 加密后的签名ID
	Name          string `json:"name"`          // 签名名称
	Intro         string `json:"intro"`         // 签名介绍
	CardImage     string `json:"cardImage"`     // Base64编码的图片数据
	CardImageName string `json:"cardImageName"` // 原始图片文件名（带扩展名）
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

// UpdateSignature 更新签名的处理函数
// encryptedID: 签名唯一标识（已加密）
// signatureData: 签名数据
// imageData: 图片文件的二进制数据（可为空表示不更改图片）
// imageExt: 图片文件的扩展名（如 .jpg, .png，可以为空）
// originalImageName: 图片原始文件名
// encryptionKey: 对称加密密钥
// removeImage: 是否需要删除图片（当用户主动删除时设为 true）
// imageChanged: 图片是否发生变更（前端报告的变更状态，用于区分"用户未修改"和"用户上传相同图片"）
//
// 说明：
// - 接收加密的ID，对其进行对称解密以验证和查找原有的签名存储条目
// - 保留原有的 sort.time（排序时间戳）不变，不更新此值
// - 根据 imageChanged 标记判断是否需要处理图片：
//   - 如果 imageChanged 为 false，表示用户未修改图片，直接使用原有的 cardImage URL
//   - 如果 imageChanged 为 true，则按照原有逻辑处理：删除或保存新图片
//
// - 图片删除由程序启动时的清理逻辑处理，不需要在此关心
// - 加密新的签名数据并替换原有的 value 字段
func UpdateSignature(encryptedID string, signatureData SignatureData, imageData []byte, imageExt string, originalImageName string, encryptionKey []byte, removeImage bool, imageChanged bool) error {
	// 1. 从配置中获取现有的签名存储数据
	existingValue := config.GetValue("signature")
	if existingValue == nil {
		logger.Error("签名配置不存在", "encryptedID", encryptedID)
		return fmt.Errorf("签名不存在")
	}

	// 类型转换
	signatureMap, ok := existingValue.(map[string]interface{})
	if !ok {
		logger.Error("签名配置数据格式错误")
		return fmt.Errorf("签名配置数据格式错误")
	}

	// 2. 检查要更新的签名是否存在
	entryData, exists := signatureMap[encryptedID]
	if !exists {
		logger.Error("要更新的签名不存在", "encryptedID", encryptedID)
		return fmt.Errorf("签名不存在")
	}

	// 3. 提取现有的排序时间戳
	var originalSortTime int64
	if entry, ok := entryData.(map[string]interface{}); ok {
		if sort, ok := entry["sort"].(map[string]interface{}); ok {
			if time, ok := sort["time"].(float64); ok {
				originalSortTime = int64(time)
			}
		}
	}

	// 如果没有找到有效的排序时间戳，使用当前时间（为了安全）
	if originalSortTime == 0 {
		logger.Warn("未找到原有的排序时间戳，使用当前时间", "encryptedID", encryptedID)
		originalSortTime = time.Now().Unix()
	}

	// 4. 处理图片逻辑
	// 首先检查 imageChanged 标记：
	// - 如果 imageChanged 为 false，表示图片未发生变更，直接使用原有的 CardImage URL
	// - 否则按照以下逻辑处理：
	//   情况1: 用户明确删除了图片（removeImage == true），则设置 CardImage 为空
	//   情况2: 用户上传了新图片（len(imageData) > 0），则使用新图片
	//   情况3: 既没有删除也没有上传（len(imageData) == 0 && removeImage == false），则保留原图片
	var originalCardImage string

	if !imageChanged {
		// 图片未发生变更，从现有数据中读取原始的图片路径
		logger.Debug("图片未发生变更，保留原有图片URL", "encryptedID", encryptedID)
		if entry, ok := entryData.(map[string]interface{}); ok {
			if value, ok := entry["value"].(string); ok {
				// 解密现有的签名数据
				decryptedData, err := decryptData(value, encryptionKey)
				if err == nil {
					// 反序列化以获取原始的 CardImage
					var originalSignatureData SignatureData
					if err := json.Unmarshal([]byte(decryptedData), &originalSignatureData); err == nil {
						originalCardImage = originalSignatureData.CardImage
						signatureData.CardImage = originalCardImage
					}
				}
			}
		}
	} else if removeImage {
		// 用户明确删除了图片，设置为空字符串
		signatureData.CardImage = ""
		logger.Debug("删除签名图片", "encryptedID", encryptedID)
	} else if len(imageData) == 0 {
		// 没有新图片，从现有数据中读取原始的图片路径
		// 尝试从现有的加密数据中解密并获取原始的 CardImage
		if entry, ok := entryData.(map[string]interface{}); ok {
			if value, ok := entry["value"].(string); ok {
				// 解密现有的签名数据
				decryptedData, err := decryptData(value, encryptionKey)
				if err == nil {
					// 反序列化以获取原始的 CardImage
					var originalSignatureData SignatureData
					if err := json.Unmarshal([]byte(decryptedData), &originalSignatureData); err == nil {
						originalCardImage = originalSignatureData.CardImage
					}
				}
			}
		}
		// 如果成功读取了原始的 CardImage，则使用它
		if originalCardImage != "" {
			signatureData.CardImage = originalCardImage
		}
	}

	// 处理图片：如果有新图片数据，按照创建逻辑处理
	// 解密加密的ID用于生成图片文件名
	unencryptedID, err := decryptData(encryptedID, encryptionKey)
	if err != nil {
		logger.Error("ID解密失败", "error", err.Error())
		return err
	}

	if len(imageData) > 0 {
		// 构建用于生成文件名的字符串
		fileNameSeed := fmt.Sprintf("%s|%s|%s|%d",
			unencryptedID,
			signatureData.Name,
			originalImageName,
			time.Now().Unix(),
		)

		logger.Debug("图片文件名种子生成（编辑模式）",
			"seed", fileNameSeed,
		)

		// 对文件名种子的字符串计算SHA-1哈希
		sha1Hash := sha1.Sum([]byte(fileNameSeed))
		imageFileName := hex.EncodeToString(sha1Hash[:])

		// 添加扩展名
		if imageExt != "" {
			if !strings.HasPrefix(imageExt, ".") {
				imageExt = "." + imageExt
			}
			imageFileName = imageFileName + imageExt
		}

		logger.Debug("图片文件名生成完成（编辑模式）",
			"原始文件名", originalImageName,
			"文件名", imageFileName,
		)

		// 创建 signature 目录
		signatureDir := filepath.Join(config.ConfigPath, "signature")
		if err := os.MkdirAll(signatureDir, os.ModePerm); err != nil {
			logger.Error("创建signature目录失败", "error", err.Error())
			return err
		}

		// 保存新的图片文件
		imagePath := filepath.Join(signatureDir, imageFileName)
		if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
			logger.Error("保存图片文件失败", "error", err.Error())
			return err
		}

		// 更新签名数据中的图片路径
		signatureData.CardImage = imagePath
	}

	// 5. 对签名数据结构进行JSON序列化
	jsonData, err := json.Marshal(signatureData)
	if err != nil {
		logger.Error("签名数据JSON序列化失败", "error", err.Error())
		return err
	}

	// 6. 对JSON字符串进行对称加密
	encryptedValue, err := encryptData(string(jsonData), encryptionKey)
	if err != nil {
		logger.Error("签名数据加密失败", "error", err.Error())
		return err
	}

	// 7. 以调试级别打印加密前后的数据
	logger.Debug("签名加密处理完成（编辑模式）",
		"加密ID", encryptedID,
		"原始Value", string(jsonData),
		"加密后Value", encryptedValue,
	)

	// 8. 更新签名存储条目（保留原有的排序时间戳）
	signatureMap[encryptedID] = SignatureStorageEntry{
		Value: encryptedValue,
		Sort: SignatureSortMetadata{
			Time: originalSortTime, // 保留原有的排序时间戳
		},
	}

	// 9. 存储回配置文件
	config.SetValue("signature", signatureMap)

	logger.Info("签名更新完成", "encryptedID", encryptedID)
	return nil
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

// EncryptData 使用AES-GCM对数据进行对称加密（导出函数，供其他包使用）
func EncryptData(data string, key []byte) (string, error) {
	return encryptData(data, key)
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

// ExportSignature 导出签名为结构化数据，用于生成 .ktsign 文件
// encryptedID: 签名唯一标识（已加密）
// encryptionKey: 对称加密密钥
//
// 返回值：
// - SignatureExportData 结构体，包含导出所需的所有信息
// - error 错误信息
//
// 说明：
// - 从配置中读取加密的签名数据
// - 解密签名数据获得原始的 SignatureData
// - 读取名片图片文件并转换为 Base64 编码
// - 构建 SignatureExportData 结构体，包含加密的 Key（encryptedID）和 Base64 图片
// - 调用方需要将此结构体转换为 JSON，再进行加密和二进制存储
func ExportSignature(encryptedID string, encryptionKey []byte) (*SignatureExportData, error) {
	// 1. 从配置中获取签名数据
	signatureMapValue := config.GetValue("signature")
	if signatureMapValue == nil {
		logger.Error("签名配置不存在", "encryptedID", encryptedID)
		return nil, fmt.Errorf("签名不存在")
	}

	// 类型转换
	signatureMap, ok := signatureMapValue.(map[string]interface{})
	if !ok {
		logger.Error("签名配置数据格式错误")
		return nil, fmt.Errorf("签名配置数据格式错误")
	}

	// 2. 检查要导出的签名是否存在
	entryData, exists := signatureMap[encryptedID]
	if !exists {
		logger.Error("要导出的签名不存在", "encryptedID", encryptedID)
		return nil, fmt.Errorf("签名不存在")
	}

	// 3. 提取加密的签名值
	var encryptedValueStr string
	if entry, ok := entryData.(map[string]interface{}); ok {
		if value, ok := entry["value"].(string); ok {
			encryptedValueStr = value
		} else {
			logger.Error("无法从 SignatureStorageEntry 中提取加密值")
			return nil, fmt.Errorf("签名数据格式错误")
		}
	} else if str, ok := entryData.(string); ok {
		// 兼容旧格式
		encryptedValueStr = str
	} else {
		logger.Error("无法识别签名数据格式")
		return nil, fmt.Errorf("签名数据格式错误")
	}

	// 4. 解密签名数据
	decryptedData, err := decryptData(encryptedValueStr, encryptionKey)
	if err != nil {
		logger.Error("签名数据解密失败", "error", err.Error())
		return nil, fmt.Errorf("签名数据解密失败: %w", err)
	}

	// 5. 解析 JSON 数据为 SignatureData
	var signatureData SignatureData
	if err := json.Unmarshal([]byte(decryptedData), &signatureData); err != nil {
		logger.Error("签名数据 JSON 解析失败", "error", err.Error())
		return nil, fmt.Errorf("签名数据解析失败: %w", err)
	}

	// 6. 构建导出数据结构
	exportData := &SignatureExportData{
		Key:   encryptedID,
		Name:  signatureData.Name,
		Intro: signatureData.Intro,
	}

	// 7. 处理图片：读取文件并转换为 Base64
	if signatureData.CardImage != "" {
		// 读取图片文件
		imageBytes, err := os.ReadFile(signatureData.CardImage)
		if err != nil {
			logger.Warn("读取签名图片文件失败，将跳过图片导出", "path", signatureData.CardImage, "error", err.Error())
			// 不中止导出流程，继续处理其他数据
			exportData.CardImage = ""
			exportData.CardImageName = ""
		} else {
			// 转换为 Base64
			exportData.CardImage = hex.EncodeToString(imageBytes)

			// 提取文件名（包含扩展名）
			exportData.CardImageName = filepath.Base(signatureData.CardImage)

			logger.Debug("签名图片转换为 Base64",
				"原始路径", signatureData.CardImage,
				"文件名", exportData.CardImageName,
				"Base64长度", len(exportData.CardImage),
			)
		}
	}

	logger.Info("签名导出数据构建完成",
		"encryptedID", encryptedID,
		"名称", exportData.Name,
		"包含图片", exportData.CardImage != "",
	)

	return exportData, nil
}

// ImportSignature 导入签名数据，检查冲突并保存
// exportData: 导出的签名数据结构体
// encryptionKey: 对称加密密钥
//
// 返回值：
// - encryptedID: 导入后的签名加密ID（若已存在则返回现有ID）
// - conflict: 是否存在冲突（签名已存在）
// - error: 错误信息
//
// 说明：
// - 检查 exportData.Key（加密的签名ID）是否已在配置中存在
// - 如果存在，返回 conflict=true，调用方应提示用户是否覆盖
// - 如果不存在或用户选择覆盖，则调用 CreateSignature 复用创建逻辑
// - 需要先将 Base64 图片数据解码为二进制，再调用 CreateSignature
// - 返回的 encryptedID 应该与 exportData.Key 保持一致
func ImportSignature(exportData *SignatureExportData, encryptionKey []byte) (string, bool, error) {
	if exportData == nil {
		logger.Error("导入的签名数据为空")
		return "", false, fmt.Errorf("签名数据为空")
	}

	// 1. 检查签名是否已存在
	signatureMapValue := config.GetValue("signature")
	var signatureMap map[string]interface{}

	if signatureMapValue != nil {
		if m, ok := signatureMapValue.(map[string]interface{}); ok {
			signatureMap = m
		} else {
			logger.Error("签名配置数据格式错误")
			return "", false, fmt.Errorf("签名配置数据格式错误")
		}
	} else {
		signatureMap = make(map[string]interface{})
	}

	// 检查 Key 是否存在
	if _, exists := signatureMap[exportData.Key]; exists {
		logger.Warn("尝试导入的签名已存在", "key", exportData.Key)
		return exportData.Key, true, nil // 返回冲突标记，由调用方决定是否覆盖
	}

	// 2. 解密导入的 Key 以获得原始的签名 ID
	unencryptedID, err := decryptData(exportData.Key, encryptionKey)
	if err != nil {
		logger.Error("签名 Key 解密失败", "error", err.Error())
		return "", false, fmt.Errorf("签名 Key 解密失败: %w", err)
	}

	// 3. 处理图片数据
	var imageData []byte
	var imageExt string
	var imageFileName string

	if exportData.CardImage != "" {
		// 从 Base64（十六进制编码）解码
		var err error
		imageData, err = hex.DecodeString(exportData.CardImage)
		if err != nil {
			logger.Warn("图片数据解码失败，将跳过图片导入", "error", err.Error())
			imageData = []byte{}
		} else {
			imageFileName = exportData.CardImageName
			// 提取文件扩展名
			lastDotIndex := strings.LastIndex(imageFileName, ".")
			if lastDotIndex != -1 && lastDotIndex < len(imageFileName)-1 {
				imageExt = imageFileName[lastDotIndex:]
			}

			logger.Debug("签名图片数据已解码",
				"文件名", imageFileName,
				"扩展名", imageExt,
				"数据大小", len(imageData),
			)
		}
	}

	// 4. 构建签名数据结构
	signatureData := SignatureData{
		Name:  exportData.Name,
		Intro: exportData.Intro,
	}

	// 5. 调用 CreateSignature 复用创建逻辑
	// 注意：CreateSignature 会生成新的 encryptedID，但我们需要使用导入的 Key
	// 为了保持 Key 的一致性，我们需要直接保存到配置，而不是调用 CreateSignature

	// 处理图片文件保存
	if len(imageData) > 0 {
		// 创建 signature 目录
		signatureDir := filepath.Join(config.ConfigPath, "signature")
		if err := os.MkdirAll(signatureDir, os.ModePerm); err != nil {
			logger.Error("创建signature目录失败", "error", err.Error())
			return "", false, fmt.Errorf("创建签名目录失败: %w", err)
		}

		// 生成文件名（使用与导出时相同的策略）
		fileNameSeed := fmt.Sprintf("%s|%s|%s|%d",
			unencryptedID,
			signatureData.Name,
			imageFileName,
			time.Now().Unix(),
		)

		sha1Hash := sha1.Sum([]byte(fileNameSeed))
		generatedFileName := hex.EncodeToString(sha1Hash[:])

		// 添加扩展名
		if imageExt != "" {
			if !strings.HasPrefix(imageExt, ".") {
				imageExt = "." + imageExt
			}
			generatedFileName = generatedFileName + imageExt
		}

		// 保存图片文件
		imagePath := filepath.Join(signatureDir, generatedFileName)
		if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
			logger.Error("保存导入的图片文件失败", "error", err.Error())
			return "", false, fmt.Errorf("保存图片失败: %w", err)
		}

		// 更新签名数据中的图片路径
		signatureData.CardImage = imagePath

		logger.Debug("导入的签名图片已保存",
			"目标路径", imagePath,
		)
	}

	// 6. 对签名数据进行 JSON 序列化
	jsonData, err := json.Marshal(signatureData)
	if err != nil {
		logger.Error("签名数据JSON序列化失败", "error", err.Error())
		return "", false, fmt.Errorf("签名数据序列化失败: %w", err)
	}

	// 7. 对 JSON 字符串进行加密
	encryptedValue, err := encryptData(string(jsonData), encryptionKey)
	if err != nil {
		logger.Error("签名数据加密失败", "error", err.Error())
		return "", false, fmt.Errorf("签名数据加密失败: %w", err)
	}

	// 8. 保存到配置
	signatureMap[exportData.Key] = SignatureStorageEntry{
		Value: encryptedValue,
		Sort: SignatureSortMetadata{
			Time: time.Now().Unix(), // 导入时生成新的排序时间戳
		},
	}

	config.SetValue("signature", signatureMap)

	logger.Info("签名导入完成",
		"encryptedID", exportData.Key,
		"名称", signatureData.Name,
	)

	return exportData.Key, false, nil
}

// ImportSignatureWithOverwrite 导入签名数据并强制覆盖现有签名
// exportData: 导出的签名数据结构体
// overwrite: 是否覆盖现有签名
// encryptionKey: 对称加密密钥
//
// 返回值：
// - encryptedID: 导入后的签名加密ID
// - error: 错误信息
func ImportSignatureWithOverwrite(exportData *SignatureExportData, overwrite bool, encryptionKey []byte) (string, error) {
	// 1. 先执行检查
	encryptedID, conflict, err := ImportSignature(exportData, encryptionKey)
	if err != nil {
		return "", err
	}

	// 2. 如果有冲突且不覆盖，返回错误
	if conflict && !overwrite {
		logger.Warn("签名已存在，未选择覆盖", "key", encryptedID)
		return "", fmt.Errorf("签名已存在，用户未选择覆盖")
	}

	// 3. 如果有冲突且选择覆盖，需要删除旧的签名并重新导入
	if conflict && overwrite {
		logger.Info("覆盖现有签名", "key", encryptedID)

		// 删除旧的签名配置
		signatureMapValue := config.GetValue("signature")
		if signatureMap, ok := signatureMapValue.(map[string]interface{}); ok {
			delete(signatureMap, encryptedID)
			config.SetValue("signature", signatureMap)
		}

		// 重新调用 ImportSignature 进行导入
		newEncryptedID, newConflict, err := ImportSignature(exportData, encryptionKey)
		if err != nil {
			return "", err
		}

		if newConflict {
			// 不应该再次出现冲突
			logger.Error("覆盖导入时意外出现冲突", "key", newEncryptedID)
			return "", fmt.Errorf("覆盖导入失败")
		}

		return newEncryptedID, nil
	}

	// 4. 没有冲突，直接返回
	return encryptedID, nil
}
