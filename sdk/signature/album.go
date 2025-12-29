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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// KeyToneAlbumSignatureEncryptionKey 专辑签名字段专用加密密钥
//
// 用途：加密专辑配置中的signature字段内容
// 安全等级：标准（固定密钥，源码可见，用于防止随手窥视）
// 长度：32字节（AES-256要求）
//
// 说明：
//   - 此密钥独立于签名管理的KeyA/KeyB，职责分离
//   - 专辑配置本身已有基于albumUUID的派生密钥保护（外层加密）
//   - 此密钥用于signature字段的内层加密，双重保护
const KeyToneAlbumSignatureEncryptionKey = "KeyTone2024Album_Signature_Field_EncryptionKey_32Bytes"

// GetAlbumSignatureKey 获取专辑签名加密密钥（32字节）
//
// 返回值：
//   - []byte: 32字节的AES-256密钥
//
// 说明：
//   - 确保密钥长度符合AES-256要求
//   - 如密钥字符串不足32字节，自动填充0
func GetAlbumSignatureKey() []byte {
	key := []byte(KeyToneAlbumSignatureEncryptionKey)
	// 确保密钥长度为32字节
	if len(key) < 32 {
		for len(key) < 32 {
			key = append(key, 0)
		}
	}
	return key[:32]
}

// EncryptAlbumSignatureField 加密专辑配置中的签名字段
//
// 用途：
//   - 将整个signature对象（JSON序列化后）加密为十六进制字符串
//   - 存储在专辑配置的signature字段中
//
// 参数：
//   - signatureJSON: 签名对象的JSON字符串（包含所有签名资格码及其数据）
//
// 返回值：
//   - string: 十六进制编码的加密密文
//   - error: 加密过程中的错误
//
// 加密算法：AES-256-GCM（复用现有的encryptData函数）
//
// 使用示例：
//
//	signatureJSON := `{"abc123...": {"name": "张三", "intro": "..."}}`
//	encrypted, err := EncryptAlbumSignatureField(signatureJSON)
//	// encrypted 可直接存入专辑配置的 signature 字段
func EncryptAlbumSignatureField(signatureJSON string) (string, error) {
	key := GetAlbumSignatureKey()
	encryptedHex, err := encryptData(signatureJSON, key)
	if err != nil {
		return "", fmt.Errorf("加密专辑签名字段失败: %w", err)
	}
	return encryptedHex, nil
}

// DecryptAlbumSignatureField 解密专辑配置中的签名字段
//
// 用途：
//   - 从专辑配置的signature字段读取加密密文
//   - 解密为JSON字符串，可解析为签名对象
//
// 参数：
//   - encryptedHex: 十六进制编码的加密密文
//
// 返回值：
//   - string: 解密后的JSON字符串
//   - error: 解密过程中的错误
//
// 说明：
//   - 与EncryptAlbumSignatureField配对使用
//   - 解密失败可能原因：密文损坏、密钥不匹配、数据被篡改
//
// 使用示例：
//
//	encryptedHex := config.GetValue("signature").(string)
//	signatureJSON, err := DecryptAlbumSignatureField(encryptedHex)
//	// 解析 signatureJSON 为 map[string]interface{} 使用
func DecryptAlbumSignatureField(encryptedHex string) (string, error) {
	key := GetAlbumSignatureKey()
	decryptedJSON, err := decryptData(encryptedHex, key)
	if err != nil {
		return "", fmt.Errorf("解密专辑签名字段失败: %w", err)
	}
	return decryptedJSON, nil
}

// GenerateQualificationCode 根据原始签名ID生成资格码
//
// 用途：
//   - 保护原始签名ID隐私，防止通过专辑配置反推签名创建者
//   - 提供确定性标识，相同签名ID总是生成相同资格码
//   - 用作专辑配置signature对象的key
//
// 参数：
//   - signatureID: 原始签名ID（未加密的UUID字符串）
//
// 返回值：
//   - string: SHA256哈希值（64字符十六进制字符串）
//
// 原理：
//   - 使用SHA256单向哈希算法
//   - 无法从资格码反推原始签名ID
//   - 满足确定性要求：相同输入总是得到相同输出
//
// 安全性：
//   - SHA256哈希不可逆，保护原始ID
//   - 抗碰撞性：不同签名ID极难生成相同资格码
//
// 使用示例：
//
//	originalID := "550e8400-e29b-41d4-a716-446655440000" // 未加密的UUID
//	qualCode := GenerateQualificationCode(originalID)
//	// qualCode: "a1b2c3d4..." (64字符)
//	// 用作专辑配置中的key: signature[qualCode] = {...}
func GenerateQualificationCode(signatureID string) string {
	// 计算SHA256哈希
	hash := sha256.Sum256([]byte(signatureID))
	// 转换为十六进制字符串（64字符）
	return hex.EncodeToString(hash[:])
}

// GenerateQualificationFingerprint 根据资格码生成资格码指纹
//
// TIPS: 资格码指纹用于在保护原始资格码不泄漏的前提下，保证签名的可追溯性。
// 计算方式：将资格码去除第2位（索引1）和第11位（索引10）字符后，计算SHA256哈希。
//
// 用途：
//   - 在前端展示时替代原始资格码，防止资格码泄漏
//   - 保持签名的可追溯性，相同资格码总是生成相同指纹
//
// 参数：
//   - qualificationCode: 资格码（64字符十六进制字符串）
//
// 返回值：
//   - string: 资格码指纹（64字符十六进制字符串）
//
// 安全性：
//   - 去除特定位置字符后再哈希，增加逆向难度
//   - 前端只接触指纹，不接触原始资格码
func GenerateQualificationFingerprint(qualificationCode string) string {
	if len(qualificationCode) < 12 {
		return qualificationCode // 无效输入时返回原值
	}

	// TIPS: 去除第2位（索引1）和第11位（索引10）字符
	// 索引:  0  1  2  3  4  5  6  7  8  9 10 11 ...
	// 去除:     ^                       ^
	modified := qualificationCode[0:1] + qualificationCode[2:10] + qualificationCode[11:]

	// 计算SHA256哈希
	hash := sha256.Sum256([]byte(modified))
	return hex.EncodeToString(hash[:])
}
