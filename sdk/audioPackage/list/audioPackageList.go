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

package audioPackageList

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	apconfig "KeyTone/audioPackage/config"
	apenc "KeyTone/audioPackage/enc"
	"KeyTone/logger"
	"KeyTone/signature"

	"github.com/spf13/viper"
)

// ============================================================================
// AlbumSignatureSummary 专辑签名摘要信息
// 用于在专辑列表展示时提供签名信息预览，避免加载完整签名数据
//
// 包含原始作者和直接导出作者的基本信息，用于悬停卡片展示
// 当原始作者与直接导出作者相同时，前端只展示一个作者区块
// ============================================================================
type AlbumSignatureSummary struct {
	// HasSignature 是否有签名
	HasSignature bool `json:"hasSignature"`

	// ============================================================================
	// 原始作者信息
	// ============================================================================
	// OriginalAuthorName 原始作者名称
	OriginalAuthorName string `json:"originalAuthorName"`
	// OriginalAuthorImage 原始作者图片路径（相对于专辑目录）
	OriginalAuthorImage string `json:"originalAuthorImage"`
	// OriginalAuthorIntro 原始作者介绍
	OriginalAuthorIntro string `json:"originalAuthorIntro"`

	// ============================================================================
	// 直接导出作者信息
	// ============================================================================
	// DirectExportAuthorName 直接导出作者名称
	DirectExportAuthorName string `json:"directExportAuthorName"`
	// DirectExportAuthorImage 直接导出作者图片路径（相对于专辑目录）
	DirectExportAuthorImage string `json:"directExportAuthorImage"`
	// DirectExportAuthorIntro 直接导出作者介绍
	DirectExportAuthorIntro string `json:"directExportAuthorIntro"`

	// ============================================================================
	// 是否为同一作者标记
	// ============================================================================
	// IsSameAuthor 原始作者与直接导出作者是否为同一人
	// 前端可根据此标记决定只展示一个作者区块
	IsSameAuthor bool `json:"isSameAuthor"`
}

// migrateConfigFile 自动检测并迁移旧的 config.json 为 package.json
// 用于支持从旧版本到新版本的平滑升级
func migrateConfigFile(configPath string) error {
	oldConfigPath := filepath.Join(configPath, "config.json")
	newConfigPath := filepath.Join(configPath, "package.json")

	// 检查新的配置文件是否已存在
	if _, err := os.Stat(newConfigPath); err == nil {
		// 新文件已存在，检查旧文件是否也存在
		if _, err := os.Stat(oldConfigPath); err == nil {
			// 新旧文件都存在，这是一个异常情况，记录警告但不处理
			logger.Warn("检测到 config.json 和 package.json 同时存在，将保留 package.json", "path", configPath)
			return nil
		}
		// 新文件存在，旧文件不存在，无需迁移
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误（如权限问题）
		logger.Error("检查新配置文件时出错", "path", newConfigPath, "error", err.Error())
		return err
	}

	// 检查旧的配置文件是否存在
	if _, err := os.Stat(oldConfigPath); err == nil {
		// 旧文件存在，执行迁移
		if err := os.Rename(oldConfigPath, newConfigPath); err != nil {
			logger.Error("迁移配置文件失败", "from", oldConfigPath, "to", newConfigPath, "error", err.Error())
			return err
		}
		logger.Info("配置文件成功迁移", "from", "config.json", "to", "package.json", "path", configPath)
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误
		logger.Error("检查旧配置文件时出错", "path", oldConfigPath, "error", err.Error())
		return err
	}

	// 新旧文件都不存在，无需迁移
	return nil
}

// 用于获取键音包根目录中, 所有键音包的完整路径的列表。
func GetAudioPackageList(rootDir string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(rootDir)

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(rootDir, entry.Name()))
		}
	}

	return directories, err
}

// 用于获取具体路径下的键音包的名称。
var viperRWMutex sync.RWMutex

func GetAudioPackageName(configPath string) any {
	// 加载配置文件
	Viper := viper.New()
	defer func() {
		Viper = nil
	}()

	// 尝试迁移旧的 config.json 文件为 package.json
	if err := migrateConfigFile(configPath); err != nil {
		logger.Warn("配置文件迁移检查失败，继续使用现有配置", "error", err.Error())
	}

	Viper.SetConfigType("json")

	var plainJSON string
	stubInfo, pkgRaw, err := apconfig.ReadCoreStubInfo(configPath)
	if err != nil {
		logger.Error("读取指示 JSON 失败", "err", err.Error())
		return nil
	}

	albumUUID := filepath.Base(configPath)

	if stubInfo != nil {
		corePath := filepath.Join(configPath, stubInfo.Core)
		cipherBytes, err := os.ReadFile(corePath)
		if err != nil {
			logger.Error("读取 core 文件失败", "err", err.Error())
			return nil
		}
		plain, err := apenc.DecryptConfigBytes(cipherBytes, albumUUID)
		if err != nil {
			logger.Error("解密 core 文件失败", "err", err.Error())
			return nil
		}
		plainJSON = plain
	} else {
		if pkgRaw == nil {
			pkgPath := filepath.Join(configPath, "package.json")
			pkgRaw, err = os.ReadFile(pkgPath)
			if err != nil {
				logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
				return nil
			}
		}

		if apenc.IsLikelyHexCipher(pkgRaw) {
			plain, decErr := apenc.DecryptConfigHex(strings.TrimSpace(string(pkgRaw)), albumUUID)
			if decErr != nil {
				logger.Error("旧版密文解密失败", "err", decErr.Error())
				return nil
			}
			plainJSON = plain
		} else {
			plainJSON = string(pkgRaw)
		}
	}

	if err := Viper.ReadConfig(strings.NewReader(plainJSON)); err != nil {
		logger.Error("解析音频包配置失败", "err", err.Error())
		return nil
	}

	// 从加载的配置文件中, 读取 package_name 字段
	viperRWMutex.RLock()
	defer viperRWMutex.RUnlock()
	return Viper.Get("package_name")
}

// UpdateAlbumUUID 更新专辑配置文件中的 UUID，必要时处理加密 core 文件。
// originalUUID 用于解密旧密文；若为空则回退为当前目录名。
func UpdateAlbumUUID(albumPath string, newUUID string, originalUUID string) error {
	stubInfo, _, err := apconfig.ReadCoreStubInfo(albumPath)
	if err != nil {
		return err
	}

	if originalUUID == "" {
		originalUUID = filepath.Base(albumPath)
	}

	albumUUID := filepath.Base(albumPath)

	if stubInfo != nil {
		corePath := filepath.Join(albumPath, stubInfo.Core)
		cipherBytes, err := os.ReadFile(corePath)
		if err != nil {
			return err
		}
		plain, err := apenc.DecryptConfigBytes(cipherBytes, originalUUID)
		if err != nil {
			// 兼容旧行为：尝试使用目录名再次解密
			if fallback, fallbackErr := apenc.DecryptConfigBytes(cipherBytes, albumUUID); fallbackErr == nil {
				plain = fallback
			} else {
				return err
			}
		}
		var data map[string]any
		if err := json.Unmarshal([]byte(plain), &data); err != nil {
			return err
		}
		data["audio_pkg_uuid"] = newUUID
		updated, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		cipherBytes, err = apenc.EncryptConfigBytes(string(updated), newUUID)
		if err != nil {
			return err
		}
		tmpCore := corePath + ".tmp"
		if err := os.WriteFile(tmpCore, cipherBytes, 0644); err != nil {
			return err
		}
		if err := os.Rename(tmpCore, corePath); err != nil {
			_ = os.Remove(tmpCore)
			return err
		}
		return apconfig.WriteCoreStubFile(albumPath)
	}

	// 加载配置文件
	v := viper.New()
	v.SetConfigName("package")
	v.SetConfigType("json")
	v.AddConfigPath(albumPath)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 更新 UUID
	v.Set("audio_pkg_uuid", newUUID)

	// 保存更改
	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// ============================================================================
// GetAlbumSignatureSummary 获取专辑签名摘要信息
// 用于在专辑列表展示时提供签名信息预览
//
// 功能说明：
//   - 读取专辑配置的 signature 字段
//   - 解密并解析签名数据
//   - 提取原始作者和直接导出作者的名称、图片路径、介绍
//   - 返回轻量级的签名摘要供前端列表展示和悬停卡片使用
//
// 参数：
//   - albumPath: 专辑目录的绝对路径
//
// 返回值：
//   - *AlbumSignatureSummary: 签名摘要信息（包含原始作者和直接导出作者）
//   - error: 错误信息
//
// 设计考虑：
//   - 此函数为轻量级实现，仅提取展示所需的数据
//   - 独立创建 Viper 实例并在函数结束时释放，避免内存泄漏
//   - 用于遍历专辑列表时批量获取签名信息
//   - 当原始作者与直接导出作者相同时，设置 IsSameAuthor=true
// ============================================================================
func GetAlbumSignatureSummary(albumPath string) (*AlbumSignatureSummary, error) {
	// 创建独立的 Viper 实例（避免影响全局配置）
	v := viper.New()
	defer func() {
		v = nil // 辅助 GC 回收
	}()

	// 尝试迁移旧的 config.json 文件为 package.json
	if err := migrateConfigFile(albumPath); err != nil {
		logger.Warn("配置文件迁移检查失败，继续使用现有配置", "error", err.Error())
	}

	v.SetConfigType("json")

	var plainJSON string
	stubInfo, pkgRaw, err := apconfig.ReadCoreStubInfo(albumPath)
	if err != nil {
		logger.Error("读取指示 JSON 失败", "err", err.Error())
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	albumUUID := filepath.Base(albumPath)

	// 解密配置文件内容
	if stubInfo != nil {
		corePath := filepath.Join(albumPath, stubInfo.Core)
		cipherBytes, err := os.ReadFile(corePath)
		if err != nil {
			logger.Error("读取 core 文件失败", "err", err.Error())
			return &AlbumSignatureSummary{HasSignature: false}, nil
		}
		plain, err := apenc.DecryptConfigBytes(cipherBytes, albumUUID)
		if err != nil {
			logger.Error("解密 core 文件失败", "err", err.Error())
			return &AlbumSignatureSummary{HasSignature: false}, nil
		}
		plainJSON = plain
	} else {
		if pkgRaw == nil {
			pkgPath := filepath.Join(albumPath, "package.json")
			pkgRaw, err = os.ReadFile(pkgPath)
			if err != nil {
				logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
				return &AlbumSignatureSummary{HasSignature: false}, nil
			}
		}

		if apenc.IsLikelyHexCipher(pkgRaw) {
			plain, decErr := apenc.DecryptConfigHex(strings.TrimSpace(string(pkgRaw)), albumUUID)
			if decErr != nil {
				logger.Error("旧版密文解密失败", "err", decErr.Error())
				return &AlbumSignatureSummary{HasSignature: false}, nil
			}
			plainJSON = plain
		} else {
			plainJSON = string(pkgRaw)
		}
	}

	if err := v.ReadConfig(strings.NewReader(plainJSON)); err != nil {
		logger.Error("解析音频包配置失败", "err", err.Error())
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	// 读取 signature 字段
	signatureValue := v.Get("signature")
	if signatureValue == nil {
		// 专辑不包含签名
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	// 解密 signature 字段
	encryptedSigStr, ok := signatureValue.(string)
	if !ok {
		logger.Warn("signature 字段格式错误", "albumPath", albumPath)
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	decryptedSigJSON, err := signature.DecryptAlbumSignatureField(encryptedSigStr)
	if err != nil {
		logger.Warn("解密 signature 字段失败", "albumPath", albumPath, "err", err.Error())
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	// ============================================================================
	// 解析签名数据 - 提取原始作者和直接导出作者的完整信息
	// 包含 name, cardImagePath, intro 字段
	// ============================================================================
	var albumSignatureMap map[string]struct {
		Name          string `json:"name"`
		Intro         string `json:"intro"`
		CardImagePath string `json:"cardImagePath"`
		Authorization *struct {
			DirectExportAuthor string `json:"directExportAuthor"`
		} `json:"authorization,omitempty"`
	}

	if err := json.Unmarshal([]byte(decryptedSigJSON), &albumSignatureMap); err != nil {
		logger.Warn("解析 signature JSON 失败", "albumPath", albumPath, "err", err.Error())
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	// ============================================================================
	// 查找原始作者和直接导出作者
	// 原始作者：具有 authorization 字段的签名
	// 直接导出作者：由原始作者的 authorization.directExportAuthor 指向
	// ============================================================================
	var originalAuthorQualCode string
	var directExportAuthorQualCode string

	// 第一步：找到原始作者（具有 authorization 字段的签名）
	for qualCode, entry := range albumSignatureMap {
		if entry.Authorization != nil {
			originalAuthorQualCode = qualCode
			directExportAuthorQualCode = entry.Authorization.DirectExportAuthor
			break
		}
	}

	// 如果没有找到原始作者，返回无签名
	if originalAuthorQualCode == "" {
		return &AlbumSignatureSummary{HasSignature: false}, nil
	}

	// 获取原始作者信息
	originalAuthor := albumSignatureMap[originalAuthorQualCode]

	// 构建签名摘要
	summary := &AlbumSignatureSummary{
		HasSignature:        true,
		OriginalAuthorName:  originalAuthor.Name,
		OriginalAuthorImage: originalAuthor.CardImagePath,
		OriginalAuthorIntro: originalAuthor.Intro,
	}

	// 获取直接导出作者信息
	if directExportAuthorQualCode != "" {
		if directAuthor, exists := albumSignatureMap[directExportAuthorQualCode]; exists {
			summary.DirectExportAuthorName = directAuthor.Name
			summary.DirectExportAuthorImage = directAuthor.CardImagePath
			summary.DirectExportAuthorIntro = directAuthor.Intro
			// 判断是否为同一作者
			summary.IsSameAuthor = (originalAuthorQualCode == directExportAuthorQualCode)
		}
	}

	return summary, nil
}
