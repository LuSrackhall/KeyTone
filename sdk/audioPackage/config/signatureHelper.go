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
	"encoding/json"
	"fmt"
)

// AlbumSignatureInfo 专辑签名信息（前端使用）
//
// 说明：
//   - 用于前端读取和解析专辑配置中的签名信息
//   - 包含原始作者、历史贡献作者、直接导出作者等信息
type AlbumSignatureInfo struct {
	// HasSignature 专辑是否包含签名
	HasSignature bool `json:"hasSignature"`

	// OriginalAuthor 原始作者签名信息
	OriginalAuthor *SignatureAuthorInfo `json:"originalAuthor,omitempty"`

	// ContributorAuthors 历史贡献作者列表
	ContributorAuthors []SignatureAuthorInfo `json:"contributorAuthors"`

	// DirectExportAuthor 直接导出作者信息
	DirectExportAuthor *SignatureAuthorInfo `json:"directExportAuthor,omitempty"`

	// AllSignatures 所有签名条目（用于选择页面）
	AllSignatures map[string]AlbumSignatureEntry `json:"allSignatures"`
}

// SignatureAuthorInfo 签名作者信息
//
// 说明：
//   - 用于前端展示签名作者详情
type SignatureAuthorInfo struct {
	// QualificationCode 资格码（签名ID的SHA256哈希）
	QualificationCode string `json:"qualificationCode"`

	// Name 签名名称
	Name string `json:"name"`

	// Intro 个人介绍
	Intro string `json:"intro"`

	// CardImagePath 名片图片路径
	CardImagePath string `json:"cardImagePath"`

	// IsOriginalAuthor 是否为原始作者
	IsOriginalAuthor bool `json:"isOriginalAuthor"`

	// RequireAuthorization 是否需要授权（仅原始作者有效）
	RequireAuthorization bool `json:"requireAuthorization,omitempty"`

	// AuthorizedList 已授权的资格码列表（仅原始作者有效）
	AuthorizedList []string `json:"authorizedList,omitempty"`
}

// GetAlbumSignatureInfo 获取专辑签名信息
//
// 功能说明：
//   - 读取专辑配置的signature字段
//   - 解密并解析签名数据
//   - 识别原始作者、历史贡献作者、直接导出作者
//   - 返回结构化的签名信息供前端使用
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//
// 返回值：
//   - *AlbumSignatureInfo: 专辑签名信息
//   - error: 错误信息
//
// 使用场景：
//   - 前端需求2：再次导出时的签名识别
//   - 前端需求4：签名作者信息展示
func GetAlbumSignatureInfo(albumPath string) (*AlbumSignatureInfo, error) {
	logger.Info("开始获取专辑签名信息", "albumPath", albumPath)

	// 加载专辑配置
	LoadConfig(albumPath, false)
	if Viper == nil {
		return nil, fmt.Errorf("加载专辑配置失败")
	}

	// 读取signature字段
	existingSignatureValue := GetValue("signature")
	if existingSignatureValue == nil {
		// 专辑不包含签名
		return &AlbumSignatureInfo{
			HasSignature:       false,
			ContributorAuthors: []SignatureAuthorInfo{},
			AllSignatures:      make(map[string]AlbumSignatureEntry),
		}, nil
	}

	// 解密signature字段
	var albumSignatureMap map[string]AlbumSignatureEntry
	if encryptedSigStr, ok := existingSignatureValue.(string); ok {
		decryptedSigJSON, err := signature.DecryptAlbumSignatureField(encryptedSigStr)
		if err != nil {
			return nil, fmt.Errorf("解密signature字段失败: %w", err)
		}

		if err := json.Unmarshal([]byte(decryptedSigJSON), &albumSignatureMap); err != nil {
			return nil, fmt.Errorf("解析signature JSON失败: %w", err)
		}
	} else {
		return nil, fmt.Errorf("signature字段格式错误")
	}

	// 构建签名信息
	info := &AlbumSignatureInfo{
		HasSignature:       true,
		ContributorAuthors: []SignatureAuthorInfo{},
		AllSignatures:      albumSignatureMap,
	}

	// 识别原始作者和贡献者
	var directExportAuthorQualCode string
	for qualCode, entry := range albumSignatureMap {
		if entry.Authorization != nil {
			// 原始作者签名
			info.OriginalAuthor = &SignatureAuthorInfo{
				QualificationCode:    qualCode,
				Name:                 entry.Name,
				Intro:                entry.Intro,
				CardImagePath:        entry.CardImagePath,
				IsOriginalAuthor:     true,
				RequireAuthorization: entry.Authorization.RequireAuthorization,
				AuthorizedList:       entry.Authorization.AuthorizedList,
			}
			directExportAuthorQualCode = entry.Authorization.DirectExportAuthor
			logger.Debug("找到原始作者签名", "qualCode", qualCode, "directExportAuthor", directExportAuthorQualCode)
		} else {
			// 贡献者签名
			info.ContributorAuthors = append(info.ContributorAuthors, SignatureAuthorInfo{
				QualificationCode: qualCode,
				Name:              entry.Name,
				Intro:             entry.Intro,
				CardImagePath:     entry.CardImagePath,
				IsOriginalAuthor:  false,
			})
		}
	}

	// 识别直接导出作者
	if directExportAuthorQualCode != "" {
		if entry, exists := albumSignatureMap[directExportAuthorQualCode]; exists {
			info.DirectExportAuthor = &SignatureAuthorInfo{
				QualificationCode: directExportAuthorQualCode,
				Name:              entry.Name,
				Intro:             entry.Intro,
				CardImagePath:     entry.CardImagePath,
				IsOriginalAuthor:  (entry.Authorization != nil),
			}
			logger.Debug("找到直接导出作者", "qualCode", directExportAuthorQualCode)
		}
	}

	logger.Info("专辑签名信息获取成功",
		"hasOriginalAuthor", info.OriginalAuthor != nil,
		"contributorsCount", len(info.ContributorAuthors),
		"hasDirectExportAuthor", info.DirectExportAuthor != nil,
	)

	return info, nil
}

// CheckSignatureInAlbum 检查签名是否已在专辑中
//
// 功能说明：
//   - 比对当前用户签名的资格码与专辑中已有的签名
//   - 用于前端需求3：标记已在专辑中的签名
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//   - encryptedSignatureID: 加密的签名ID
//
// 返回值：
//   - isInAlbum: 签名是否已在专辑中
//   - qualificationCode: 签名的资格码
//   - error: 错误信息
func CheckSignatureInAlbum(albumPath string, encryptedSignatureID string) (bool, string, error) {
	// 解密签名ID获取原始UUID
	originalSignatureID, err := signature.DecryptWithKeyA(encryptedSignatureID)
	if err != nil {
		return false, "", fmt.Errorf("解密签名ID失败: %w", err)
	}

	// 生成资格码
	qualificationCode := signature.GenerateQualificationCode(originalSignatureID)

	// 获取专辑签名信息
	signatureInfo, err := GetAlbumSignatureInfo(albumPath)
	if err != nil {
		return false, qualificationCode, fmt.Errorf("获取专辑签名信息失败: %w", err)
	}

	// 检查资格码是否存在
	_, exists := signatureInfo.AllSignatures[qualificationCode]

	logger.Debug("检查签名是否在专辑中",
		"qualificationCode", qualificationCode,
		"exists", exists,
	)

	return exists, qualificationCode, nil
}

// CheckSignatureAuthorization 检查签名是否有导出授权
//
// 功能说明：
//   - 检查当前签名是否在原始作者的authorizedList中
//   - 用于前端需求3：根据authorizedList使能/失能签名选项
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//   - encryptedSignatureID: 加密的签名ID
//
// 返回值：
//   - isAuthorized: 是否有授权（或是原始作者本人）
//   - requireAuthorization: 是否需要授权验证
//   - qualificationCode: 签名的资格码
//   - error: 错误信息
func CheckSignatureAuthorization(albumPath string, encryptedSignatureID string) (bool, bool, string, error) {
	// 解密签名ID获取原始UUID
	originalSignatureID, err := signature.DecryptWithKeyA(encryptedSignatureID)
	if err != nil {
		return false, false, "", fmt.Errorf("解密签名ID失败: %w", err)
	}

	// 生成资格码
	qualificationCode := signature.GenerateQualificationCode(originalSignatureID)

	// 获取专辑签名信息
	signatureInfo, err := GetAlbumSignatureInfo(albumPath)
	if err != nil {
		return false, false, qualificationCode, fmt.Errorf("获取专辑签名信息失败: %w", err)
	}

	// 如果专辑没有签名，视为有授权（首次导出）
	if !signatureInfo.HasSignature {
		return true, false, qualificationCode, nil
	}

	// 如果没有原始作者签名，出现异常
	if signatureInfo.OriginalAuthor == nil {
		return false, false, qualificationCode, fmt.Errorf("专辑缺少原始作者签名")
	}

	// 如果不需要授权，所有签名都可以导出
	if !signatureInfo.OriginalAuthor.RequireAuthorization {
		logger.Debug("专辑不需要授权，允许导出", "qualificationCode", qualificationCode)
		return true, false, qualificationCode, nil
	}

	// 检查是否为原始作者本人
	if qualificationCode == signatureInfo.OriginalAuthor.QualificationCode {
		logger.Debug("是原始作者本人，允许导出", "qualificationCode", qualificationCode)
		return true, true, qualificationCode, nil
	}

	// 检查是否在授权列表中
	for _, authorizedQualCode := range signatureInfo.OriginalAuthor.AuthorizedList {
		if qualificationCode == authorizedQualCode {
			logger.Debug("在授权列表中，允许导出", "qualificationCode", qualificationCode)
			return true, true, qualificationCode, nil
		}
	}

	logger.Debug("不在授权列表中，需要导入授权文件", "qualificationCode", qualificationCode)
	return false, true, qualificationCode, nil
}

// GetAvailableSignaturesForExport 获取可用于导出的签名列表
//
// 功能说明：
//   - 从用户配置中读取所有签名
//   - 标记哪些签名已在专辑中
//   - 标记哪些签名有导出授权
//   - 用于前端需求3：签名选择页面增强
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//
// 返回值：
//   - []AvailableSignature: 可用签名列表
//   - error: 错误信息
type AvailableSignature struct {
	// EncryptedID 加密的签名ID
	EncryptedID string `json:"encryptedId"`

	// QualificationCode 资格码
	QualificationCode string `json:"qualificationCode"`

	// Name 签名名称
	Name string `json:"name"`

	// Intro 个人介绍
	Intro string `json:"intro"`

	// IsInAlbum 是否已在专辑中
	IsInAlbum bool `json:"isInAlbum"`

	// IsAuthorized 是否有导出授权
	IsAuthorized bool `json:"isAuthorized"`

	// IsOriginalAuthor 是否为原始作者
	IsOriginalAuthor bool `json:"isOriginalAuthor"`
}

func GetAvailableSignaturesForExport(albumPath string) ([]AvailableSignature, error) {
	logger.Info("开始获取可用签名列表", "albumPath", albumPath)

	// 从用户配置中获取所有签名
	signatureMapValue := config.GetValue("signature")
	if signatureMapValue == nil {
		return []AvailableSignature{}, nil
	}

	signatureMap, ok := signatureMapValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("签名配置数据格式错误")
	}

	// 获取专辑签名信息
	albumSigInfo, err := GetAlbumSignatureInfo(albumPath)
	if err != nil {
		return nil, fmt.Errorf("获取专辑签名信息失败: %w", err)
	}

	var availableSignatures []AvailableSignature

	// 遍历所有签名
	for encryptedID, entryData := range signatureMap {
		// 提取加密的签名Value
		var encryptedValueStr string
		if entry, ok := entryData.(map[string]interface{}); ok {
			if value, ok := entry["value"].(string); ok {
				encryptedValueStr = value
			} else {
				continue
			}
		} else if str, ok := entryData.(string); ok {
			encryptedValueStr = str
		} else {
			continue
		}

		// 解密签名ID
		originalSignatureID, err := signature.DecryptWithKeyA(encryptedID)
		if err != nil {
			logger.Warn("解密签名ID失败，跳过", "encryptedID", encryptedID, "error", err.Error())
			continue
		}

		// 生成资格码
		qualCode := signature.GenerateQualificationCode(originalSignatureID)

		// 解密签名Value获取名称和介绍
		decryptedValueJSON, err := signature.DecryptValueWithDynamicKey(encryptedValueStr, encryptedID)
		if err != nil {
			logger.Warn("解密签名数据失败，跳过", "encryptedID", encryptedID, "error", err.Error())
			continue
		}

		var sigData signature.SignatureData
		if err := json.Unmarshal([]byte(decryptedValueJSON), &sigData); err != nil {
			logger.Warn("解析签名数据失败，跳过", "encryptedID", encryptedID, "error", err.Error())
			continue
		}

		// 检查是否在专辑中
		_, isInAlbum := albumSigInfo.AllSignatures[qualCode]

		// 检查是否有授权
		isAuthorized := false
		isOriginalAuthor := false

		if albumSigInfo.HasSignature && albumSigInfo.OriginalAuthor != nil {
			// 检查是否为原始作者
			if qualCode == albumSigInfo.OriginalAuthor.QualificationCode {
				isAuthorized = true
				isOriginalAuthor = true
			} else if !albumSigInfo.OriginalAuthor.RequireAuthorization {
				// 不需要授权，所有签名都可用
				isAuthorized = true
			} else {
				// 检查是否在授权列表中
				for _, authQualCode := range albumSigInfo.OriginalAuthor.AuthorizedList {
					if qualCode == authQualCode {
						isAuthorized = true
						break
					}
				}
			}
		} else {
			// 首次导出，所有签名都可用
			isAuthorized = true
		}

		availableSignatures = append(availableSignatures, AvailableSignature{
			EncryptedID:       encryptedID,
			QualificationCode: qualCode,
			Name:              sigData.Name,
			Intro:             sigData.Intro,
			IsInAlbum:         isInAlbum,
			IsAuthorized:      isAuthorized,
			IsOriginalAuthor:  isOriginalAuthor,
		})
	}

	logger.Info("可用签名列表获取成功", "count", len(availableSignatures))
	return availableSignatures, nil
}
