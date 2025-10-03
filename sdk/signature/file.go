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
	"encoding/base64"
	"encoding/json"
	"errors"
)

const (
	// KeytoneSignatureKey is used for simple symmetric encryption
	KeytoneSignatureKey = "KeyTone2024SecretKey"
)

// SignatureFilePayload represents the structure of a .ktsign file
type SignatureFilePayload struct {
	Key   string                 `json:"key"`   // encrypt(protectCode)
	Value map[string]interface{} `json:"value"` // Plain JSON with name, intro, cardImagePath, createdAt
	Assets *SignatureAssets      `json:"assets,omitempty"`
}

// SignatureAssets contains optional assets like card image
type SignatureAssets struct {
	CardImage string `json:"cardImage,omitempty"` // base64 encoded image
}

// XOREncryptDecrypt performs simple XOR encryption/decryption
func XOREncryptDecrypt(data []byte, key string) []byte {
	keyBytes := []byte(key)
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}
	return result
}

// EncodeSignatureFile encodes a signature payload to .ktsign format
// Returns base64 encoded binary data
func EncodeSignatureFile(payload *SignatureFilePayload) (string, error) {
	if payload == nil {
		return "", errors.New("payload cannot be nil")
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Encrypt the JSON data
	encrypted := XOREncryptDecrypt(jsonData, KeytoneSignatureKey)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return encoded, nil
}

// DecodeSignatureFile decodes a .ktsign file from base64 binary data
func DecodeSignatureFile(base64Data string) (*SignatureFilePayload, error) {
	if base64Data == "" {
		return nil, errors.New("base64Data cannot be empty")
	}

	// Decode from base64
	encrypted, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, errors.New("invalid_format: failed to decode base64")
	}

	// Decrypt the data
	decrypted := XOREncryptDecrypt(encrypted, KeytoneSignatureKey)

	// Unmarshal JSON
	var payload SignatureFilePayload
	if err := json.Unmarshal(decrypted, &payload); err != nil {
		return nil, errors.New("invalid_format: failed to parse JSON")
	}

	// Validate required fields
	if payload.Key == "" {
		return nil, errors.New("invalid_format: missing key field")
	}
	if payload.Value == nil {
		return nil, errors.New("invalid_format: missing value field")
	}

	name, ok := payload.Value["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("invalid_format: missing or invalid name field")
	}

	return &payload, nil
}
