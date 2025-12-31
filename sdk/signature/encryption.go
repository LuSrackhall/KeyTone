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
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// ==============================
// 签名相关对称密钥变量定义
// 注意：这些变量不再是 const，而是 var，以便在编译时通过 -ldflags 进行注入
// 注入的值应为经过 XOR 混淆后的 Hex 字符串（与授权流一致）
// ==============================

// 定义默认的开源密钥常量，用于运行时比对
const (
	DefaultKeyA = "KeyTone2024Signature_KeyA_SecureEncryptionKeyForIDEncryption"
	DefaultKeyB = "KeyTone2024Signature_KeyB_SuperSecureEncryptionKeyForExportImportOperation"
)

// KeyToneSignatureEncryptionKeyA 密钥A：用于加密签名ID和生成动态密钥
// 安全等级：标准
// 长度: 32字节
var KeyToneSignatureEncryptionKeyA = DefaultKeyA

// KeyToneSignatureEncryptionKeyB 密钥B：用于导出/导入加密，安全级别更高
// 安全等级：高
// 长度: 32字节
var KeyToneSignatureEncryptionKeyB = DefaultKeyB

// GetKeyA 获取密钥A (32字节)
// 用途：加密签名ID、生成动态密钥
func GetKeyA() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneSignatureEncryptionKeyA == DefaultKeyA {
		key := []byte(DefaultKeyA)
		if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
		}
		return key[:32]
	}
	// 2. 否则说明已被注入，执行解混淆逻辑（hex -> xor -> plaintext）
	return deobfuscateKey(KeyToneSignatureEncryptionKeyA)
}

// GetKeyB 获取密钥B (32字节)
// 用途：导出/导入签名文件加密
func GetKeyB() []byte {
	// 1. 如果变量值等于默认常量，说明未注入，直接使用默认明文密钥
	if KeyToneSignatureEncryptionKeyB == DefaultKeyB {
		key := []byte(DefaultKeyB)
		if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
		}
		return key[:32]
	}
	// 2. 否则说明已被注入，执行解混淆逻辑（hex -> xor -> plaintext）
	return deobfuscateKey(KeyToneSignatureEncryptionKeyB)
}

// GenerateDynamicKey 根据加密的签名ID生成动态密钥
//
// 工作原理：
// 1. 使用KeyA解密encryptedID得到原始UUID
// 2. 提取UUID的后7位字符
// 3. 使用PBKDF2将KeyA和后7位组合成32字节的对称密钥
// 4. 返回生成的密钥
//
// 参数：
//
//	encryptedID: 已加密的签名ID
//
// 返回值：
//
//	[]byte: 32字节的对称加密密钥
//	error: 生成过程中的错误
//
// 用途：
// - 加密/解密各个签名的数据Value
// - 每个签名都有唯一的动态密钥，增加安全性
func GenerateDynamicKey(encryptedID string) ([]byte, error) {
	keyA := GetKeyA()

	// 1. 使用KeyA解密以获取原始UUID
	unencryptedID, err := decryptData(encryptedID, keyA)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ID for dynamic key generation: %w", err)
	}

	// 2. 提取UUID的后7位
	if len(unencryptedID) < 7 {
		return nil, fmt.Errorf("ID too short for dynamic key generation (expected >= 7 chars, got %d)", len(unencryptedID))
	}
	suffix := unencryptedID[len(unencryptedID)-7:]

	// 3. 使用PBKDF2生成动态密钥
	// 参数：
	//   - password: KeyA (32字节)
	//   - salt: UUID的后7位
	//   - iterations: 10000 (OWASP推荐)
	//   - keylen: 32 (AES-256需要32字节)
	//   - hash: SHA256
	dynamicKey := pbkdf2.Key(keyA, []byte(suffix), 10000, 32, sha256.New)

	return dynamicKey, nil
}

// EncryptValueWithDynamicKey 使用动态密钥加密数据
//
// 这是一个便捷函数，结合了GenerateDynamicKey和encryptData两个步骤
//
// 参数：
//
//	data: 要加密的数据字符串（通常是JSON）
//	encryptedID: 加密的签名ID（用于生成对应的动态密钥）
//
// 返回值：
//
//	string: 加密后的数据（16进制编码）
//	error: 加密过程中的错误
//
// 用途：
// - 在创建/更新签名时加密Value
// - 确保使用与ID对应的唯一动态密钥
func EncryptValueWithDynamicKey(data string, encryptedID string) (string, error) {
	dynamicKey, err := GenerateDynamicKey(encryptedID)
	if err != nil {
		return "", fmt.Errorf("failed to generate dynamic key: %w", err)
	}

	encryptedData, err := encryptData(data, dynamicKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data with dynamic key: %w", err)
	}

	return encryptedData, nil
}

// DecryptValueWithDynamicKey 使用动态密钥解密数据
//
// 这是一个便捷函数，结合了GenerateDynamicKey和decryptData两个步骤
// 前端应优先使用后端的 /signature/decrypt API，本函数供后端内部使用
//
// 参数：
//
//	encryptedData: 加密的数据（16进制编码）
//	encryptedID: 加密的签名ID（用于生成对应的动态密钥）
//
// 返回值：
//
//	string: 解密后的JSON字符串
//	error: 解密过程中的错误
//
// 用途：
// - 在获取签名列表时解密Value
// - 在获取单个签名时解密Value
// - 在导出/导入时解密Value
// - 提供方便的API给前端通过HTTP调用
func DecryptValueWithDynamicKey(encryptedData string, encryptedID string) (string, error) {
	dynamicKey, err := GenerateDynamicKey(encryptedID)
	if err != nil {
		return "", fmt.Errorf("failed to generate dynamic key: %w", err)
	}

	decryptedData, err := decryptData(encryptedData, dynamicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data with dynamic key: %w", err)
	}

	return decryptedData, nil
}

// EncryptWithKeyA 使用密钥A加密数据
// 用途：加密签名ID
func EncryptWithKeyA(data string) (string, error) {
	keyA := GetKeyA()
	return encryptData(data, keyA)
}

// DecryptWithKeyA 使用密钥A解密数据
// 用途：解密签名ID
func DecryptWithKeyA(encryptedData string) (string, error) {
	keyA := GetKeyA()
	return decryptData(encryptedData, keyA)
}

// EncryptWithKeyB 使用密钥B加密数据
// 用途：加密导出的签名文件
func EncryptWithKeyB(data string) (string, error) {
	keyB := GetKeyB()
	return encryptData(data, keyB)
}

// DecryptWithKeyB 使用密钥B解密数据
// 用途：解密导入的签名文件
func DecryptWithKeyB(encryptedData string) (string, error) {
	keyB := GetKeyB()
	return decryptData(encryptedData, keyB)
}
