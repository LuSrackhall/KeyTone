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

package audioPackageConfig

import (
	"KeyTone/config"
	"KeyTone/logger"
	"KeyTone/signature"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AuthorizationMetadata 授权元数据结构
//
// 说明：
//   - 仅原始作者签名包含此对象
//   - 用于二次导出授权控制和授权记录管理
type AuthorizationMetadata struct {
	// RequireAuthorization 是否需要二次导出授权
	RequireAuthorization bool `json:"requireAuthorization"`

	// ContactEmail 联系邮箱（需要授权时必需）
	ContactEmail string `json:"contactEmail"`

	// ContactAdditional 补充联系信息（可选）
	ContactAdditional string `json:"contactAdditional,omitempty"`

	// AuthorizedList 已授权的资格码列表
	// 存储被授权者的资格码（签名ID的SHA256哈希）
	AuthorizedList []string `json:"authorizedList"`

	// DirectExportAuthor 直接导出作者的资格码
	// 记录每次导出时的签名者，用于前端展示
	// 每次导出时更新为当前导出者的资格码
	DirectExportAuthor string `json:"directExportAuthor"`
}

// AlbumSignatureEntry 专辑配置中的签名条目结构
//
// 说明：
//   - 存储在专辑配置的signature字段中
//   - 以资格码（SHA256哈希）为key索引
//   - 整个signature对象使用专用密钥加密
type AlbumSignatureEntry struct {
	// Name 签名名称
	Name string `json:"name"`

	// Intro 个人介绍
	Intro string `json:"intro"`

	// CardImagePath 名片图片相对路径（相对于专辑目录）
	// 格式示例："audioFiles/abc123.jpg"
	CardImagePath string `json:"cardImagePath"`

	// Authorization 授权元数据（仅原始作者签名包含）
	// 非原始作者签名不包含此字段
	Authorization *AuthorizationMetadata `json:"authorization,omitempty"`
}

// ApplySignatureToAlbum 将签名应用到专辑配置
//
// 功能说明：
//  1. 从签名管理配置读取并解密签名数据
//  2. 生成资格码（SHA256哈希）
//  3. 处理签名名片图片复制到专辑目录
//  4. 构建专辑签名对象（包含授权信息）
//  5. 加密签名数据并写入专辑配置
//  6. 输出调试日志
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//   - encryptedSignatureID: 加密的签名ID（从签名管理系统获取）
//   - requireAuthorization: 是否需要二次导出授权
//   - contactEmail: 联系邮箱（requireAuthorization=true时必需）
//   - contactAdditional: 补充联系信息（可选）
//
// 返回值：
//   - qualificationCode: 生成的资格码（SHA256哈希，64字符）
//   - error: 错误信息
//
// 使用示例：
//
//	qualCode, err := ApplySignatureToAlbum(
//	    "/path/to/album",
//	    "encryptedSigID",
//	    true,
//	    "author@example.com",
//	    "微信: author123",
//	)
func ApplySignatureToAlbum(
	albumPath string,
	encryptedSignatureID string,
	requireAuthorization bool,
	contactEmail string,
	contactAdditional string,
) (string, error) {
	logger.Info("开始应用签名到专辑配置",
		"albumPath", albumPath,
		"encryptedSignatureID", encryptedSignatureID,
		"requireAuthorization", requireAuthorization,
	)

	// 步骤1：参数验证
	if strings.TrimSpace(albumPath) == "" {
		return "", fmt.Errorf("专辑路径不能为空")
	}
	if strings.TrimSpace(encryptedSignatureID) == "" {
		return "", fmt.Errorf("签名ID不能为空")
	}
	if requireAuthorization && strings.TrimSpace(contactEmail) == "" {
		return "", fmt.Errorf("需要授权时必须提供联系邮箱")
	}

	// 验证专辑目录存在
	if info, err := os.Stat(albumPath); err != nil || !info.IsDir() {
		return "", fmt.Errorf("专辑目录不存在或无效: %v", err)
	}

	// 步骤2：从签名管理配置读取签名数据
	signatureMapValue := config.GetValue("signature")
	if signatureMapValue == nil {
		return "", fmt.Errorf("签名配置不存在，请先创建签名")
	}

	signatureMap, ok := signatureMapValue.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("签名配置数据格式错误")
	}

	// 检查签名是否存在
	entryData, exists := signatureMap[encryptedSignatureID]
	if !exists {
		return "", fmt.Errorf("签名不存在或已被删除: %s", encryptedSignatureID)
	}

	// 提取加密的签名Value
	var encryptedValueStr string
	if entry, ok := entryData.(map[string]interface{}); ok {
		if value, ok := entry["value"].(string); ok {
			encryptedValueStr = value
		} else {
			return "", fmt.Errorf("无法从签名存储条目中提取value字段")
		}
	} else if str, ok := entryData.(string); ok {
		// 兼容旧格式
		encryptedValueStr = str
	} else {
		return "", fmt.Errorf("无法识别签名数据格式")
	}

	// 步骤3：解密签名ID获取原始UUID
	originalSignatureID, err := signature.DecryptWithKeyA(encryptedSignatureID)
	if err != nil {
		return "", fmt.Errorf("解密签名ID失败: %w", err)
	}

	logger.Debug("签名ID解密成功", "originalID", originalSignatureID)

	// 步骤4：使用动态密钥解密签名Value
	decryptedValueJSON, err := signature.DecryptValueWithDynamicKey(encryptedValueStr, encryptedSignatureID)
	if err != nil {
		return "", fmt.Errorf("解密签名数据失败: %w", err)
	}

	// 解析签名数据
	var sigData signature.SignatureData
	if err := json.Unmarshal([]byte(decryptedValueJSON), &sigData); err != nil {
		return "", fmt.Errorf("解析签名数据JSON失败: %w", err)
	}

	logger.Debug("签名数据解密成功",
		"name", sigData.Name,
		"intro", sigData.Intro,
		"cardImage", sigData.CardImage,
	)

	// 步骤5：生成资格码
	qualificationCode := signature.GenerateQualificationCode(originalSignatureID)
	logger.Info("资格码生成成功", "qualificationCode", qualificationCode)

	// 步骤6：处理签名名片图片复制
	var cardImageRelPath string
	if sigData.CardImage != "" {
		relPath, err := copySignatureCardImageToAlbum(
			sigData.CardImage,
			albumPath,
			qualificationCode,
		)
		if err != nil {
			// 图片复制失败不中止流程，记录警告即可
			logger.Warn("签名图片复制失败，将跳过图片", "error", err.Error())
			cardImageRelPath = ""
		} else {
			cardImageRelPath = relPath
			logger.Info("签名图片复制成功", "relativePath", cardImageRelPath)
		}
	}

	// 步骤7：加载专辑配置，检查是否已有签名
	LoadConfig(albumPath, false)
	if Viper == nil {
		return "", fmt.Errorf("加载专辑配置失败")
	}

	// 读取现有的signature字段（可能已加密）
	var albumSignatureMap map[string]AlbumSignatureEntry
	var isFirstExport bool = true
	var originalAuthorEntry *AlbumSignatureEntry
	existingSignatureValue := GetValue("signature")

	if existingSignatureValue != nil {
		// 尝试解密现有的signature字段
		if encryptedSigStr, ok := existingSignatureValue.(string); ok {
			// 解密
			decryptedSigJSON, err := signature.DecryptAlbumSignatureField(encryptedSigStr)
			if err != nil {
				logger.Warn("解密现有signature字段失败，将创建新字段", "error", err.Error())
				albumSignatureMap = make(map[string]AlbumSignatureEntry)
			} else {
				// 解析JSON
				if err := json.Unmarshal([]byte(decryptedSigJSON), &albumSignatureMap); err != nil {
					logger.Warn("解析现有signature JSON失败，将创建新字段", "error", err.Error())
					albumSignatureMap = make(map[string]AlbumSignatureEntry)
				} else {
					isFirstExport = false
					// 找到原始作者签名（包含authorization字段的那个）
					for qualCode, entry := range albumSignatureMap {
						if entry.Authorization != nil {
							entryCopy := entry
							originalAuthorEntry = &entryCopy
							logger.Debug("找到原始作者签名", "qualCode", qualCode)
							break
						}
					}
				}
			}
		} else {
			// 非字符串格式，创建新字段
			albumSignatureMap = make(map[string]AlbumSignatureEntry)
		}
	} else {
		// 不存在signature字段，创建新map
		albumSignatureMap = make(map[string]AlbumSignatureEntry)
	}

	// 步骤8：构建专辑签名对象
	albumSigEntry := AlbumSignatureEntry{
		Name:          sigData.Name,
		Intro:         sigData.Intro,
		CardImagePath: cardImageRelPath,
	}

	// 步骤9：处理authorization字段
	if isFirstExport {
		// 首次导出：创建原始作者签名，包含authorization对象
		albumSigEntry.Authorization = &AuthorizationMetadata{
			RequireAuthorization: requireAuthorization,
			ContactEmail:         contactEmail,
			ContactAdditional:    contactAdditional,
			AuthorizedList:       []string{},        // 初始化为空数组
			DirectExportAuthor:   qualificationCode, // 设置为当前导出者的资格码
		}
		logger.Info("首次导出：创建原始作者签名",
			"qualificationCode", qualificationCode,
			"requireAuthorization", requireAuthorization,
		)
	} else {
		// 再次导出：需要更新原始作者签名的directExportAuthor
		if originalAuthorEntry != nil && originalAuthorEntry.Authorization != nil {
			// 更新原始作者签名的directExportAuthor
			for qualCode, entry := range albumSignatureMap {
				if entry.Authorization != nil {
					entry.Authorization.DirectExportAuthor = qualificationCode
					albumSignatureMap[qualCode] = entry
					logger.Info("更新原始作者签名的directExportAuthor",
						"originalAuthor", qualCode,
						"directExportAuthor", qualificationCode,
					)
					break
				}
			}
		}
		// 非原始作者签名不包含authorization字段
		logger.Info("再次导出：添加贡献者签名", "qualificationCode", qualificationCode)
	}

	// 添加或更新当前签名
	albumSignatureMap[qualificationCode] = albumSigEntry
	logger.Debug("签名已添加到专辑signature map", "count", len(albumSignatureMap))

	// 步骤9：序列化为JSON并加密
	albumSignatureJSON, err := json.MarshalIndent(albumSignatureMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化专辑签名数据失败: %w", err)
	}

	encryptedAlbumSignature, err := signature.EncryptAlbumSignatureField(string(albumSignatureJSON))
	if err != nil {
		return "", fmt.Errorf("加密专辑签名字段失败: %w", err)
	}

	// 步骤10：写入专辑配置
	SetValue("signature", encryptedAlbumSignature)
	logger.Info("签名已成功写入专辑配置", "qualificationCode", qualificationCode)

	// 步骤11：输出调试日志（包含未加密的签名内容）
	fmt.Printf("\n[专辑签名调试] 签名已成功应用到专辑配置 - 未加密内容：\n%s\n\n", string(albumSignatureJSON))

	return qualificationCode, nil
}

// copySignatureCardImageToAlbum 复制签名名片图片到专辑目录
//
// 功能说明：
//   - 从签名存储目录读取图片文件
//   - 生成新文件名：SHA1(资格码 + 原始文件名 + 时间戳) + 扩展名
//   - 复制到专辑的audioFiles目录
//   - 返回相对路径（相对于专辑目录）
//
// 参数：
//   - sourceImagePath: 签名图片源路径（绝对路径）
//   - albumPath: 专辑目录路径（绝对路径）
//   - qualificationCode: 资格码（用于生成文件名）
//
// 返回值：
//   - relativePath: 相对于专辑目录的图片路径，格式："audioFiles/{filename}.jpg"
//   - error: 错误信息
//
// 说明：
//   - 复用audioFiles目录，与音频文件共用
//   - 文件名使用SHA1哈希确保唯一性
//   - 图片文件不存在时返回错误，由调用方决定处理方式
func copySignatureCardImageToAlbum(
	sourceImagePath string,
	albumPath string,
	qualificationCode string,
) (string, error) {
	// 验证源图片文件存在
	if _, err := os.Stat(sourceImagePath); err != nil {
		return "", fmt.Errorf("源图片文件不存在: %w", err)
	}

	// 读取源图片文件
	imageData, err := os.ReadFile(sourceImagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %w", err)
	}

	// 获取文件扩展名
	ext := filepath.Ext(sourceImagePath)
	originalFileName := filepath.Base(sourceImagePath)

	// 生成新文件名：SHA1(资格码 + 原始文件名 + 时间戳)
	fileNameSeed := fmt.Sprintf("%s|%s|%d",
		qualificationCode,
		originalFileName,
		time.Now().Unix(),
	)
	sha1Hash := sha1.Sum([]byte(fileNameSeed))
	newFileName := hex.EncodeToString(sha1Hash[:]) + ext

	// 确保专辑的audioFiles目录存在
	audioFilesDir := filepath.Join(albumPath, "audioFiles")
	if err := os.MkdirAll(audioFilesDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建audioFiles目录失败: %w", err)
	}

	// 目标文件路径
	destPath := filepath.Join(audioFilesDir, newFileName)

	// 写入图片文件
	if err := os.WriteFile(destPath, imageData, 0644); err != nil {
		return "", fmt.Errorf("写入图片文件失败: %w", err)
	}

	// 返回相对路径
	relativePath := filepath.Join("audioFiles", newFileName)

	logger.Debug("签名图片复制完成",
		"source", sourceImagePath,
		"destination", destPath,
		"relativePath", relativePath,
	)

	return relativePath, nil
}
