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
	"KeyTone/logger"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ==============================
// 授权相关密钥变量定义
// 注意：这些变量不再是 const，而是 var，以便在编译时通过 -ldflags 进行注入
// 注入的值应为经过 XOR 混淆后的 Hex 字符串
// ==============================

// xorMask 用于混淆密钥的掩码，必须与 tools/key-obfuscator/main.go 中的一致
var xorMask = []byte{0x55, 0xAA, 0x33, 0xCC, 0x99, 0x66, 0x11, 0xEE, 0x77, 0xBB, 0x22, 0xDD, 0x88, 0x44, 0xFF, 0x00}

// 定义默认的开源密钥常量，用于运行时比对
const (
	DefaultKeyF = "PLACEHOLDER_KEY_F_REPLACE_ME_32B"
	DefaultKeyK = "PLACEHOLDER_KEY_K_REPLACE_ME_32B"
	DefaultKeyY = "PLACEHOLDER_KEY_Y_REPLACE_ME_32B"
	DefaultKeyN = "PLACEHOLDER_KEY_N_REPLACE_ME_32B"
)

// deobfuscateKey 解混淆密钥
// 注意：此函数仅处理注入的混淆值（Hex -> XOR -> Plaintext）
// 如果注入的是非 Hex 的明文（异常情况），则直接返回补齐后的明文
func deobfuscateKey(obfuscatedHex string) []byte {
	// 1. 尝试解析 Hex 字符串
	obfuscated, err := hex.DecodeString(obfuscatedHex)
	if err != nil {
		// 解析失败，说明注入的不是 Hex 字符串（可能是用户错误注入了明文）
		// 直接作为明文处理，补齐或截断到 32 字节
		key := []byte(obfuscatedHex)
		if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
		}
		return key[:32]
	}

	// 2. 执行 XOR 解混淆
	realKey := make([]byte, len(obfuscated))
	for i, b := range obfuscated {
		realKey[i] = b ^ xorMask[i%len(xorMask)]
	}

	// 3. 长度补齐或截断
	if len(realKey) < 32 {
		for len(realKey) < 32 {
			realKey = append(realKey, 0)
		}
	}
	return realKey[:32]
}

// KeyToneAuthRequestEncryptionKeyF 密钥F：用于授权申请文件中字段级加密
// 用途：
//   - 加密请求方签名原始UUID的后10位
//   - 加密原始作者签名的SHA256值
//
// 安全等级：标准
// 建议长度: 32字节 (不足自动补齐)
var KeyToneAuthRequestEncryptionKeyF = DefaultKeyF // 默认开源密钥（明文）

// KeyToneAuthRequestEncryptionKeyK 密钥K：用于授权申请文件的整体加密
// 用途：
//   - 对授权申请文件内容进行整体加密
//
// 安全等级：高
// 建议长度: 32字节 (不足自动补齐)
var KeyToneAuthRequestEncryptionKeyK = DefaultKeyK // 默认开源密钥（明文）

// KeyToneAuthGrantEncryptionKeyY 密钥Y：用于授权文件中原始作者签名资格码前缀加密
// 用途：
//   - 加密原始作者签名资格码的前15位
//
// 安全等级：高
// 建议长度: 32字节 (不足自动补齐)
var KeyToneAuthGrantEncryptionKeyY = DefaultKeyY // 默认开源密钥（明文）

// KeyToneAuthGrantEncryptionKeyN 密钥N：用于授权文件的最终加密
// 用途：
//   - 对授权文件的最终组合哈希进行加密
//
// 安全等级：高
// 建议长度: 32字节 (不足自动补齐)
var KeyToneAuthGrantEncryptionKeyN = DefaultKeyN // 默认开源密钥（明文）

// ==============================
// 密钥获取函数
// ==============================

// GetKeyF 获取密钥F (32字节)
// 用途：授权申请文件中的字段级加密
func GetKeyF() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneAuthRequestEncryptionKeyF == DefaultKeyF {
		return []byte(DefaultKeyF)
	}
	// 2. 否则说明已被注入，执行解混淆逻辑
	return deobfuscateKey(KeyToneAuthRequestEncryptionKeyF)
}

// GetKeyK 获取密钥K (32字节)
// 用途：授权申请文件的整体加密
func GetKeyK() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneAuthRequestEncryptionKeyK == DefaultKeyK {
		return []byte(DefaultKeyK)
	}
	// 2. 否则说明已被注入，执行解混淆逻辑
	return deobfuscateKey(KeyToneAuthRequestEncryptionKeyK)
}

// GetKeyY 获取密钥Y (32字节)
// 用途：授权文件中原始作者签名资格码前缀加密
func GetKeyY() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneAuthGrantEncryptionKeyY == DefaultKeyY {
		return []byte(DefaultKeyY)
	}
	// 2. 否则说明已被注入，执行解混淆逻辑
	return deobfuscateKey(KeyToneAuthGrantEncryptionKeyY)
}

// GetKeyN 获取密钥N (32字节)
// 用途：授权文件的最终加密
func GetKeyN() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneAuthGrantEncryptionKeyN == DefaultKeyN {
		return []byte(DefaultKeyN)
	}
	// 2. 否则说明已被注入，执行解混淆逻辑
	return deobfuscateKey(KeyToneAuthGrantEncryptionKeyN)
}

// ==============================
// 授权申请文件数据结构
// ==============================

// AuthRequestData 授权申请文件的内部数据结构（加密前）
type AuthRequestData struct {
	// AuthorizationUUIDHash 专辑authorizationUUID的SHA256值
	// 用于标识申请针对哪个专辑
	AuthorizationUUIDHash string `json:"authorizationUUIDHash"`

	// RequesterSignatureIDSuffix 请求方签名原始UUID后10位（密钥F加密）
	// 用于后续校验请求方身份
	RequesterSignatureIDSuffix string `json:"requesterSignatureIDSuffix"`

	// OriginalAuthorQualCodeHash 原始作者签名资格码的SHA256值（密钥F加密）
	// 用于标识请求的是哪个原始作者的授权
	OriginalAuthorQualCodeHash string `json:"originalAuthorQualCodeHash"`

	// RequesterSignatureName 请求方签名名称（明文）
	// 方便原始作者了解请求方身份
	RequesterSignatureName string `json:"requesterSignatureName"`

	// RequesterQualificationFingerprint 请求方签名的资格码指纹（明文）
	// 便于原始作者核实申请方身份
	RequesterQualificationFingerprint string `json:"requesterQualificationFingerprint"`
}

// AuthRequestFile 授权申请文件的最终结构（整体加密后）
type AuthRequestFile struct {
	// EncryptedData 使用密钥K整体加密的AuthRequestData JSON
	EncryptedData string `json:"encryptedData"`

	// Version 文件版本号
	Version string `json:"version"`
}

// ==============================
// 授权文件数据结构
// ==============================

// AuthGrantFile 签名授权文件结构
type AuthGrantFile struct {
	// EncryptedAuthToken 加密的授权令牌
	// 内容：SHA256(authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15) 的密钥N加密结果
	EncryptedAuthToken string `json:"encryptedAuthToken"`

	// Version 文件版本号
	Version string `json:"version"`
}

// ==============================
// 授权申请文件生成
// ==============================

// GenerateAuthRequest 生成授权申请文件内容
//
// 功能说明：
//  1. 根据专辑的authorizationUUID生成SHA256哈希
//  2. 提取请求方签名原始UUID的后10位并加密
//  3. 计算原始作者签名资格码的SHA256并加密
//  4. 组装并整体加密授权申请数据
//
// 参数：
//   - authorizationUUID: 专辑的授权标识UUID
//   - requesterEncryptedSignatureID: 请求方签名的加密ID
//   - originalAuthorQualificationCode: 原始作者签名的资格码
//   - requesterSignatureName: 请求方签名的名称
//
// 返回值：
//   - []byte: 授权申请文件的二进制内容
//   - error: 错误信息
func GenerateAuthRequest(
	authorizationUUID string,
	requesterEncryptedSignatureID string,
	originalAuthorQualificationCode string,
	requesterSignatureName string,
) ([]byte, error) {
	logger.Info("开始生成授权申请文件",
		"requesterSignatureName", requesterSignatureName,
	)

	// 步骤1：计算authorizationUUID的SHA256哈希
	authUUIDHash := sha256.Sum256([]byte(authorizationUUID))
	authUUIDHashHex := hex.EncodeToString(authUUIDHash[:])

	logger.Debug("[授权申请-数据链路] authorizationUUID -> SHA256",
		"原始值(authorizationUUID)", authorizationUUID,
		"变换方式", "SHA256哈希",
		"变换结果(authUUIDHashHex)", authUUIDHashHex,
	)

	// 步骤2：解密请求方签名ID并提取后10位，然后加密
	requesterOriginalID, err := DecryptWithKeyA(requesterEncryptedSignatureID)
	if err != nil {
		return nil, fmt.Errorf("解密请求方签名ID失败: %w", err)
	}

	if len(requesterOriginalID) < 10 {
		return nil, fmt.Errorf("请求方签名ID长度不足10位")
	}
	requesterIDSuffix10 := requesterOriginalID[len(requesterOriginalID)-10:]

	logger.Debug("[授权申请-数据链路] 请求方签名ID -> 后10位",
		"原始值(requesterOriginalID)", requesterOriginalID,
		"变换方式", "提取后10位",
		"变换结果(requesterIDSuffix10)", requesterIDSuffix10,
	)

	// 使用密钥F加密请求方签名ID后10位
	keyF := GetKeyF()
	encryptedRequesterIDSuffix, err := encryptData(requesterIDSuffix10, keyF)
	if err != nil {
		return nil, fmt.Errorf("加密请求方签名ID后10位失败: %w", err)
	}

	logger.Debug("[授权申请-数据链路] 请求方签名ID后10位 -> 密钥F加密",
		"原始值(requesterIDSuffix10)", requesterIDSuffix10,
		"变换方式", "AES-256-GCM对称加密",
		"使用密钥", "KeyF (KeyToneAuthRequestEncryptionKeyF)",
		"变换结果(encryptedRequesterIDSuffix)", encryptedRequesterIDSuffix,
	)

	// 步骤3：计算原始作者签名资格码的SHA256并加密
	originalQualCodeHash := sha256.Sum256([]byte(originalAuthorQualificationCode))
	originalQualCodeHashHex := hex.EncodeToString(originalQualCodeHash[:])

	logger.Debug("[授权申请-数据链路] 原始作者资格码 -> SHA256",
		"原始值(originalAuthorQualificationCode)", originalAuthorQualificationCode,
		"变换方式", "SHA256哈希",
		"变换结果(originalQualCodeHashHex)", originalQualCodeHashHex,
	)

	// 使用密钥F加密原始作者资格码SHA256
	encryptedOriginalQualCodeHash, err := encryptData(originalQualCodeHashHex, keyF)
	if err != nil {
		return nil, fmt.Errorf("加密原始作者资格码SHA256失败: %w", err)
	}

	logger.Debug("[授权申请-数据链路] 原始作者资格码SHA256 -> 密钥F加密",
		"原始值(originalQualCodeHashHex)", originalQualCodeHashHex,
		"变换方式", "AES-256-GCM对称加密",
		"使用密钥", "KeyF (KeyToneAuthRequestEncryptionKeyF)",
		"变换结果(encryptedOriginalQualCodeHash)", encryptedOriginalQualCodeHash,
	)

	// 步骤3.5：计算请求方签名的资格码指纹
	requesterQualCode := GenerateQualificationCode(requesterOriginalID)
	requesterQualFingerprint := GenerateQualificationFingerprint(requesterQualCode)

	logger.Debug("[授权申请-数据链路] 请求方签名ID -> 资格码 -> 指纹",
		"原始签名ID", requesterOriginalID,
		"资格码", requesterQualCode,
		"指纹", requesterQualFingerprint,
	)

	// 步骤4：组装授权申请数据
	authRequestData := AuthRequestData{
		AuthorizationUUIDHash:             authUUIDHashHex,
		RequesterSignatureIDSuffix:        encryptedRequesterIDSuffix,
		OriginalAuthorQualCodeHash:        encryptedOriginalQualCodeHash,
		RequesterSignatureName:            requesterSignatureName,
		RequesterQualificationFingerprint: requesterQualFingerprint,
	}

	// 序列化为JSON
	authRequestJSON, err := json.Marshal(authRequestData)
	if err != nil {
		return nil, fmt.Errorf("序列化授权申请数据失败: %w", err)
	}

	logger.Debug("[授权申请-数据链路] 授权申请数据JSON",
		"内容", string(authRequestJSON),
	)

	// 步骤5：使用密钥K整体加密
	keyK := GetKeyK()
	encryptedAuthRequestData, err := encryptData(string(authRequestJSON), keyK)
	if err != nil {
		return nil, fmt.Errorf("整体加密授权申请数据失败: %w", err)
	}

	logger.Debug("[授权申请-数据链路] 授权申请数据 -> 密钥K整体加密",
		"原始值", string(authRequestJSON),
		"变换方式", "AES-256-GCM对称加密",
		"使用密钥", "KeyK (KeyToneAuthRequestEncryptionKeyK)",
		"变换结果(encryptedAuthRequestData)", encryptedAuthRequestData,
	)

	// 步骤6：构建最终文件结构
	authRequestFile := AuthRequestFile{
		EncryptedData: encryptedAuthRequestData,
		Version:       "1.0",
	}

	// 序列化为JSON并转换为二进制
	fileContent, err := json.Marshal(authRequestFile)
	if err != nil {
		return nil, fmt.Errorf("序列化授权申请文件失败: %w", err)
	}

	logger.Info("授权申请文件生成成功",
		"文件大小(字节)", len(fileContent),
	)

	return fileContent, nil
}

// ==============================
// 授权申请文件解析
// ==============================

// ParsedAuthRequest 解析后的授权申请数据
type ParsedAuthRequest struct {
	// AuthorizationUUIDHash 专辑authorizationUUID的SHA256值
	AuthorizationUUIDHash string `json:"authorizationUUIDHash"`

	// RequesterSignatureIDSuffix 请求方签名原始UUID后10位（已解密）
	RequesterSignatureIDSuffix string `json:"requesterSignatureIDSuffix"`

	// OriginalAuthorQualCodeHash 原始作者签名资格码的SHA256值（已解密）
	OriginalAuthorQualCodeHash string `json:"originalAuthorQualCodeHash"`

	// RequesterSignatureName 请求方签名名称
	RequesterSignatureName string `json:"requesterSignatureName"`

	// RequesterQualificationFingerprint 请求方签名的资格码指纹
	// 便于原始作者核实申请方身份
	RequesterQualificationFingerprint string `json:"requesterQualificationFingerprint"`
}

// ParseAuthRequest 解析授权申请文件
//
// 功能说明：
//  1. 解析文件JSON结构
//  2. 使用密钥K解密整体数据
//  3. 使用密钥F解密各个加密字段
//  4. 返回解析后的数据
//
// 参数：
//   - fileContent: 授权申请文件的二进制内容
//
// 返回值：
//   - *ParsedAuthRequest: 解析后的授权申请数据
//   - error: 错误信息
func ParseAuthRequest(fileContent []byte) (*ParsedAuthRequest, error) {
	logger.Info("开始解析授权申请文件",
		"文件大小(字节)", len(fileContent),
	)

	// 步骤1：解析文件JSON结构
	var authRequestFile AuthRequestFile
	if err := json.Unmarshal(fileContent, &authRequestFile); err != nil {
		return nil, fmt.Errorf("解析授权申请文件JSON失败: %w", err)
	}

	logger.Debug("[授权申请解析-数据链路] 文件结构",
		"Version", authRequestFile.Version,
		"EncryptedData长度", len(authRequestFile.EncryptedData),
	)

	// 步骤2：使用密钥K解密整体数据
	keyK := GetKeyK()
	decryptedData, err := decryptData(authRequestFile.EncryptedData, keyK)
	if err != nil {
		return nil, fmt.Errorf("解密授权申请数据失败（可能文件损坏或密钥不匹配）: %w", err)
	}

	logger.Debug("[授权申请解析-数据链路] 密钥K解密",
		"变换方式", "AES-256-GCM对称解密",
		"使用密钥", "KeyK (KeyToneAuthRequestEncryptionKeyK)",
		"解密结果", decryptedData,
	)

	// 步骤3：解析内部数据结构
	var authRequestData AuthRequestData
	if err := json.Unmarshal([]byte(decryptedData), &authRequestData); err != nil {
		return nil, fmt.Errorf("解析授权申请内部数据失败: %w", err)
	}

	// 步骤4：使用密钥F解密各个加密字段
	keyF := GetKeyF()

	// 解密请求方签名ID后10位
	requesterIDSuffix, err := decryptData(authRequestData.RequesterSignatureIDSuffix, keyF)
	if err != nil {
		return nil, fmt.Errorf("解密请求方签名ID后10位失败: %w", err)
	}

	logger.Debug("[授权申请解析-数据链路] 解密请求方签名ID后10位",
		"加密值", authRequestData.RequesterSignatureIDSuffix,
		"变换方式", "AES-256-GCM对称解密",
		"使用密钥", "KeyF (KeyToneAuthRequestEncryptionKeyF)",
		"解密结果", requesterIDSuffix,
	)

	// 解密原始作者资格码SHA256
	originalQualCodeHash, err := decryptData(authRequestData.OriginalAuthorQualCodeHash, keyF)
	if err != nil {
		return nil, fmt.Errorf("解密原始作者资格码SHA256失败: %w", err)
	}

	logger.Debug("[授权申请解析-数据链路] 解密原始作者资格码SHA256",
		"加密值", authRequestData.OriginalAuthorQualCodeHash,
		"变换方式", "AES-256-GCM对称解密",
		"使用密钥", "KeyF (KeyToneAuthRequestEncryptionKeyF)",
		"解密结果", originalQualCodeHash,
	)

	// 构建解析结果
	parsed := &ParsedAuthRequest{
		AuthorizationUUIDHash:             authRequestData.AuthorizationUUIDHash,
		RequesterSignatureIDSuffix:        requesterIDSuffix,
		OriginalAuthorQualCodeHash:        originalQualCodeHash,
		RequesterSignatureName:            authRequestData.RequesterSignatureName,
		RequesterQualificationFingerprint: authRequestData.RequesterQualificationFingerprint,
	}

	logger.Info("授权申请文件解析成功",
		"请求方签名名称", parsed.RequesterSignatureName,
		"请求方资格码指纹", parsed.RequesterQualificationFingerprint,
		"authorizationUUIDHash", parsed.AuthorizationUUIDHash,
	)

	return parsed, nil
}

// ==============================
// 签名授权文件生成
// ==============================

// GenerateAuthGrant 生成签名授权文件
//
// 功能说明：
//  1. 提取authorizationUUID SHA256的后10位
//  2. 计算请求方签名资格码并提取前11位
//  3. 提取原始作者签名资格码前15位并用密钥Y加密
//  4. 组合上述三部分并计算SHA256
//  5. 使用密钥N加密最终结果
//
// 参数：
//   - authorizationUUIDHash: 专辑authorizationUUID的SHA256值（从申请文件获取）
//   - requesterSignatureIDSuffix: 请求方签名原始UUID后10位（从申请文件获取，已解密）
//   - originalAuthorEncryptedSignatureID: 原始作者签名的加密ID
//
// 返回值：
//   - []byte: 签名授权文件的二进制内容
//   - error: 错误信息
func GenerateAuthGrant(
	authorizationUUIDHash string,
	requesterSignatureIDSuffix string,
	originalAuthorEncryptedSignatureID string,
) ([]byte, error) {
	logger.Info("开始生成签名授权文件")

	// 步骤1：提取authorizationUUID SHA256的后10位
	if len(authorizationUUIDHash) < 10 {
		return nil, fmt.Errorf("authorizationUUIDHash长度不足10位")
	}
	authUUIDHashSuffix10 := authorizationUUIDHash[len(authorizationUUIDHash)-10:]

	logger.Debug("[授权文件-数据链路] authorizationUUIDHash -> 后10位",
		"原始值(authorizationUUIDHash)", authorizationUUIDHash,
		"变换方式", "提取后10位",
		"变换结果(authUUIDHashSuffix10)", authUUIDHashSuffix10,
	)

	// 步骤2：根据请求方签名ID后10位计算其资格码（需要完整ID，这里用后10位近似）
	// 注意：由于我们只有后10位，我们使用一种特殊方式来生成用于校验的前缀
	// 实际上，我们需要请求方在导入时提供完整的加密签名ID来进行校验
	// 这里我们直接使用后10位的SHA256的前11位作为代表
	requesterSuffixHash := sha256.Sum256([]byte(requesterSignatureIDSuffix))
	requesterQualCodePrefix11 := hex.EncodeToString(requesterSuffixHash[:])[0:11]

	logger.Debug("[授权文件-数据链路] 请求方签名ID后10位 -> SHA256 -> 前11位",
		"原始值(requesterSignatureIDSuffix)", requesterSignatureIDSuffix,
		"变换方式", "SHA256哈希后取前11位",
		"变换结果(requesterQualCodePrefix11)", requesterQualCodePrefix11,
	)

	// 步骤3：解密原始作者签名ID并计算资格码，取前15位后加密
	originalAuthorOriginalID, err := DecryptWithKeyA(originalAuthorEncryptedSignatureID)
	if err != nil {
		return nil, fmt.Errorf("解密原始作者签名ID失败: %w", err)
	}

	originalAuthorQualCode := GenerateQualificationCode(originalAuthorOriginalID)
	if len(originalAuthorQualCode) < 15 {
		return nil, fmt.Errorf("原始作者资格码长度不足15位")
	}
	originalQualCodePrefix15 := originalAuthorQualCode[0:15]

	logger.Debug("[授权文件-数据链路] 原始作者签名ID -> 资格码 -> 前15位",
		"原始值(originalAuthorOriginalID)", originalAuthorOriginalID,
		"变换方式1", "GenerateQualificationCode (SHA256)",
		"中间结果(originalAuthorQualCode)", originalAuthorQualCode,
		"变换方式2", "提取前15位",
		"变换结果(originalQualCodePrefix15)", originalQualCodePrefix15,
	)

	// 使用密钥Y加密原始作者资格码前15位
	// 注意：为了保证验证时的确定性，这里使用确定性加密（固定Nonce）
	keyY := GetKeyY()
	encryptedOriginalQualCodePrefix15, err := encryptDataDeterministic(originalQualCodePrefix15, keyY)
	if err != nil {
		return nil, fmt.Errorf("加密原始作者资格码前15位失败: %w", err)
	}

	logger.Debug("[授权文件-数据链路] 原始作者资格码前15位 -> 密钥Y加密(确定性)",
		"原始值(originalQualCodePrefix15)", originalQualCodePrefix15,
		"变换方式", "AES-256-GCM确定性加密(固定Nonce)",
		"使用密钥", "KeyY (KeyToneAuthGrantEncryptionKeyY)",
		"变换结果(encryptedOriginalQualCodePrefix15)", encryptedOriginalQualCodePrefix15,
	)

	// 步骤4：组合三部分并计算SHA256
	// 组合方式：authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15
	combinedString := authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15
	combinedHash := sha256.Sum256([]byte(combinedString))
	combinedHashHex := hex.EncodeToString(combinedHash[:])

	logger.Debug("[授权文件-数据链路] 组合字符串 -> SHA256",
		"组合方式", "authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15",
		"组合字符串(combinedString)", combinedString,
		"变换方式", "SHA256哈希",
		"变换结果(combinedHashHex)", combinedHashHex,
	)

	// 步骤5：使用密钥N加密最终结果
	keyN := GetKeyN()
	encryptedAuthToken, err := encryptData(combinedHashHex, keyN)
	if err != nil {
		return nil, fmt.Errorf("加密授权令牌失败: %w", err)
	}

	logger.Debug("[授权文件-数据链路] 组合哈希 -> 密钥N加密",
		"原始值(combinedHashHex)", combinedHashHex,
		"变换方式", "AES-256-GCM对称加密",
		"使用密钥", "KeyN (KeyToneAuthGrantEncryptionKeyN)",
		"变换结果(encryptedAuthToken)", encryptedAuthToken,
	)

	// 步骤6：构建授权文件结构
	authGrantFile := AuthGrantFile{
		EncryptedAuthToken: encryptedAuthToken,
		Version:            "1.0",
	}

	// 序列化为JSON并转换为二进制
	fileContent, err := json.Marshal(authGrantFile)
	if err != nil {
		return nil, fmt.Errorf("序列化授权文件失败: %w", err)
	}

	logger.Info("签名授权文件生成成功",
		"文件大小(字节)", len(fileContent),
	)

	return fileContent, nil
}

// ==============================
// 签名授权文件导入验证
// ==============================

// VerifyAndImportAuthGrant 验证并导入签名授权文件
//
// 功能说明：
//  1. 解析授权文件
//  2. 使用密钥N解密授权令牌
//  3. 根据当前专辑和签名信息重新计算组合哈希
//  4. 比较哈希值，验证授权有效性
//
// 参数：
//   - fileContent: 授权文件的二进制内容
//   - authorizationUUID: 当前专辑的authorizationUUID
//   - requesterEncryptedSignatureID: 请求方（当前用户）签名的加密ID
//   - originalAuthorQualificationCode: 原始作者签名的资格码
//
// 返回值：
//   - bool: 验证是否通过
//   - string: 请求方签名的资格码（仅验证通过时返回）
//   - error: 错误信息
func VerifyAndImportAuthGrant(
	fileContent []byte,
	authorizationUUID string,
	requesterEncryptedSignatureID string,
	originalAuthorQualificationCode string,
) (bool, string, error) {
	logger.Info("开始验证签名授权文件")

	// 步骤1：解析授权文件
	var authGrantFile AuthGrantFile
	if err := json.Unmarshal(fileContent, &authGrantFile); err != nil {
		return false, "", fmt.Errorf("解析授权文件JSON失败: %w", err)
	}

	logger.Debug("[授权验证-数据链路] 文件结构",
		"Version", authGrantFile.Version,
		"EncryptedAuthToken长度", len(authGrantFile.EncryptedAuthToken),
	)

	// 步骤2：使用密钥N解密授权令牌
	keyN := GetKeyN()
	decryptedCombinedHash, err := decryptData(authGrantFile.EncryptedAuthToken, keyN)
	if err != nil {
		return false, "", fmt.Errorf("解密授权令牌失败（可能文件损坏或密钥不匹配）: %w", err)
	}

	logger.Debug("[授权验证-数据链路] 密钥N解密授权令牌",
		"加密值", authGrantFile.EncryptedAuthToken,
		"变换方式", "AES-256-GCM对称解密",
		"使用密钥", "KeyN (KeyToneAuthGrantEncryptionKeyN)",
		"解密结果(decryptedCombinedHash)", decryptedCombinedHash,
	)

	// 步骤3：重新计算组合哈希以进行验证

	// 3.1 计算authorizationUUID的SHA256并取后10位
	authUUIDHash := sha256.Sum256([]byte(authorizationUUID))
	authUUIDHashHex := hex.EncodeToString(authUUIDHash[:])
	authUUIDHashSuffix10 := authUUIDHashHex[len(authUUIDHashHex)-10:]

	logger.Debug("[授权验证-数据链路] authorizationUUID -> SHA256 -> 后10位",
		"原始值(authorizationUUID)", authorizationUUID,
		"SHA256结果", authUUIDHashHex,
		"后10位(authUUIDHashSuffix10)", authUUIDHashSuffix10,
	)

	// 3.2 解密请求方签名ID并取后10位，计算SHA256的前11位
	requesterOriginalID, err := DecryptWithKeyA(requesterEncryptedSignatureID)
	if err != nil {
		return false, "", fmt.Errorf("解密请求方签名ID失败: %w", err)
	}
	if len(requesterOriginalID) < 10 {
		return false, "", fmt.Errorf("请求方签名ID长度不足10位")
	}
	requesterIDSuffix10 := requesterOriginalID[len(requesterOriginalID)-10:]
	requesterSuffixHash := sha256.Sum256([]byte(requesterIDSuffix10))
	requesterQualCodePrefix11 := hex.EncodeToString(requesterSuffixHash[:])[0:11]

	// 计算完整的请求方资格码（用于返回）
	requesterFullHash := sha256.Sum256([]byte(requesterOriginalID))
	requesterQualificationCode := hex.EncodeToString(requesterFullHash[:])

	logger.Debug("[授权验证-数据链路] 请求方签名ID -> 后10位 -> SHA256 -> 前11位",
		"原始ID", requesterOriginalID,
		"后10位", requesterIDSuffix10,
		"SHA256前11位(requesterQualCodePrefix11)", requesterQualCodePrefix11,
	)

	// 3.3 提取原始作者资格码前15位并用密钥Y加密
	if len(originalAuthorQualificationCode) < 15 {
		return false, "", fmt.Errorf("原始作者资格码长度不足15位")
	}
	originalQualCodePrefix15 := originalAuthorQualificationCode[0:15]
	keyY := GetKeyY()
	// 使用确定性加密以匹配生成时的逻辑
	encryptedOriginalQualCodePrefix15, err := encryptDataDeterministic(originalQualCodePrefix15, keyY)
	if err != nil {
		return false, "", fmt.Errorf("加密原始作者资格码前15位失败: %w", err)
	}

	logger.Debug("[授权验证-数据链路] 原始作者资格码 -> 前15位 -> 密钥Y加密(确定性)",
		"原始资格码", originalAuthorQualificationCode,
		"前15位", originalQualCodePrefix15,
		"密钥Y加密结果", encryptedOriginalQualCodePrefix15,
	)

	// 步骤4：组合并计算SHA256
	// 组合方式：authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15
	combinedString := authUUIDHashSuffix10 + requesterQualCodePrefix11 + encryptedOriginalQualCodePrefix15
	combinedHash := sha256.Sum256([]byte(combinedString))
	calculatedHashHex := hex.EncodeToString(combinedHash[:])

	logger.Debug("[授权验证-数据链路] 组合字符串 -> SHA256",
		"组合字符串", combinedString,
		"计算得到的哈希(calculatedHashHex)", calculatedHashHex,
		"文件中的哈希(decryptedCombinedHash)", decryptedCombinedHash,
	)

	// 步骤5：比较哈希值
	if calculatedHashHex != decryptedCombinedHash {
		logger.Warn("授权验证失败：哈希不匹配",
			"计算值", calculatedHashHex,
			"文件值", decryptedCombinedHash,
		)
		return false, "", fmt.Errorf("授权验证失败：授权文件与当前专辑/签名不匹配")
	}

	logger.Info("签名授权验证成功")
	return true, requesterQualificationCode, nil
}
