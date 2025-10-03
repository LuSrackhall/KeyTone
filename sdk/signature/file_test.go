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
	"testing"
)

func TestEncodeDecodeSignatureFile(t *testing.T) {
	// Create a test signature payload
	payload := &SignatureFilePayload{
		Key: "encrypted_protect_code_123",
		Value: map[string]interface{}{
			"name":          "Test Signature",
			"intro":         "A test signature",
			"cardImagePath": "test.png",
			"createdAt":     "2025-10-01T10:00:00Z",
		},
		Assets: &SignatureAssets{
			CardImage: "base64_encoded_image_data",
		},
	}

	// Encode
	encoded, err := EncodeSignatureFile(payload)
	if err != nil {
		t.Fatalf("Failed to encode signature file: %v", err)
	}

	if encoded == "" {
		t.Fatal("Encoded data is empty")
	}

	// Decode
	decoded, err := DecodeSignatureFile(encoded)
	if err != nil {
		t.Fatalf("Failed to decode signature file: %v", err)
	}

	// Verify
	if decoded.Key != payload.Key {
		t.Errorf("Key mismatch: expected %s, got %s", payload.Key, decoded.Key)
	}

	name, ok := decoded.Value["name"].(string)
	if !ok || name != "Test Signature" {
		t.Errorf("Name mismatch: expected Test Signature, got %v", name)
	}

	if decoded.Assets == nil {
		t.Error("Assets are nil")
	} else if decoded.Assets.CardImage != "base64_encoded_image_data" {
		t.Errorf("CardImage mismatch: expected base64_encoded_image_data, got %s", decoded.Assets.CardImage)
	}
}

func TestDecodeInvalidBase64(t *testing.T) {
	_, err := DecodeSignatureFile("invalid!@#$base64")
	if err == nil {
		t.Error("Expected error for invalid base64, got nil")
	}
}

func TestDecodeInvalidJSON(t *testing.T) {
	// Create invalid encrypted data
	invalidData := XOREncryptDecrypt([]byte("not json"), KeytoneSignatureKey)
	encoded := string(invalidData)
	
	_, err := DecodeSignatureFile(encoded)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestDecodeMissingName(t *testing.T) {
	payload := &SignatureFilePayload{
		Key: "test_key",
		Value: map[string]interface{}{
			"intro": "Missing name field",
		},
	}

	encoded, _ := EncodeSignatureFile(payload)
	_, err := DecodeSignatureFile(encoded)
	
	if err == nil {
		t.Error("Expected error for missing name field, got nil")
	}
}
