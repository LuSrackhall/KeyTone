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

package server

import (
	"KeyTone/config"
	"KeyTone/logger"
	"KeyTone/signature"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SignatureData represents the signature structure
type SignatureData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	CardImage string `json:"cardImage"`
}

// SignatureFileFormat represents the .ktsign file structure
type SignatureFileFormat struct {
	Version   string        `json:"version"`
	Signature SignatureData `json:"signature"`
	Checksum  string        `json:"checksum"`
}

// getSignatureImageDir returns the directory for signature card images
func getSignatureImageDir() string {
	configPath := config.GetValue("config_path")
	if configPath == nil {
		return filepath.Join(".", "signatures", "card_images")
	}
	return filepath.Join(configPath.(string), "signatures", "card_images")
}

// EnsureSignatureImageDir creates the signature image directory if it doesn't exist
// (exported for initialization)
func EnsureSignatureImageDir() error {
	dir := getSignatureImageDir()
	return os.MkdirAll(dir, os.ModePerm)
}

// ensureSignatureImageDir creates the signature image directory if it doesn't exist
func ensureSignatureImageDir() error {
	return EnsureSignatureImageDir()
}

// saveBase64Image saves a base64 encoded image and returns the filename
func saveBase64Image(base64Data string) (string, error) {
	// Remove data URL prefix if present
	if strings.Contains(base64Data, ",") {
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Calculate SHA-256 hash as filename
	hash := sha256.Sum256(imageData)
	filename := fmt.Sprintf("%x.png", hash)

	// Ensure directory exists
	if err := ensureSignatureImageDir(); err != nil {
		return "", fmt.Errorf("failed to create signature image directory: %v", err)
	}

	// Save file
	filePath := filepath.Join(getSignatureImageDir(), filename)
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image file: %v", err)
	}

	return filename, nil
}

// loadImageAsBase64 loads an image file and returns it as base64
func loadImageAsBase64(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	filePath := filepath.Join(getSignatureImageDir(), filename)
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

// signatureRouters sets up all signature-related routes
func signatureRouters(r *gin.Engine) {
	signatureGroup := r.Group("/signature")

	// Create signature
	signatureGroup.POST("/create", func(ctx *gin.Context) {
		var req struct {
			Name      string `json:"name" binding:"required,min=1,max=50"`
			Intro     string `json:"intro" binding:"max=500"`
			CardImage string `json:"cardImage"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request data: " + err.Error(),
			})
			return
		}

		// Generate ID and protect code
		id, err := signature.GenerateProtectCode()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to generate ID: " + err.Error(),
			})
			return
		}

		// Save image if provided
		var imagePath string
		if req.CardImage != "" {
			imagePath, err = saveBase64Image(req.CardImage)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to save image: " + err.Error(),
				})
				return
			}
		}

		// Create signature data
		sigData := SignatureData{
			ID:        id,
			Name:      req.Name,
			Intro:     req.Intro,
			CardImage: imagePath,
		}

		// Encrypt signature data
		sigJSON, err := json.Marshal(sigData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to marshal signature data: " + err.Error(),
			})
			return
		}

		encryptedData, err := signature.EncryptSignature(string(sigJSON))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt signature: " + err.Error(),
			})
			return
		}

		// Encrypt ID for key
		encryptedID, err := signature.EncryptSignature(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt ID: " + err.Error(),
			})
			return
		}

		// Store in config
		key := fmt.Sprintf("signature_manager.%s", encryptedID)
		config.SetValue(key, encryptedData)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Signature created successfully",
			"data":    sigData,
		})
	})

	// List all signatures
	signatureGroup.GET("/list", func(ctx *gin.Context) {
		signatureManager := config.GetValue("signature_manager")
		if signatureManager == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    make(map[string]SignatureData),
			})
			return
		}

		sigMap, ok := signatureManager.(map[string]interface{})
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    make(map[string]SignatureData),
			})
			return
		}

		result := make(map[string]SignatureData)

		for _, encryptedData := range sigMap {
			encryptedStr, ok := encryptedData.(string)
			if !ok {
				logger.Error("Invalid signature data type")
				continue
			}

			// Decrypt signature data
			decryptedJSON, err := signature.DecryptSignature(encryptedStr)
			if err != nil {
				logger.Error("Failed to decrypt signature", "error", err)
				continue
			}

			var sigData SignatureData
			if err := json.Unmarshal([]byte(decryptedJSON), &sigData); err != nil {
				logger.Error("Failed to unmarshal signature", "error", err)
				continue
			}

			result[sigData.ID] = sigData
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	// Update signature
	signatureGroup.PUT("/update", func(ctx *gin.Context) {
		var req struct {
			ID        string `json:"id" binding:"required"`
			Name      string `json:"name" binding:"required,min=1,max=50"`
			Intro     string `json:"intro" binding:"max=500"`
			CardImage string `json:"cardImage"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request data: " + err.Error(),
			})
			return
		}

		// Encrypt ID to find the signature
		encryptedID, err := signature.EncryptSignature(req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt ID: " + err.Error(),
			})
			return
		}

		key := fmt.Sprintf("signature_manager.%s", encryptedID)
		existingData := config.GetValue(key)
		if existingData == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Signature not found",
			})
			return
		}

		// Get existing signature to preserve image if not updated
		existingEncrypted, ok := existingData.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid signature data",
			})
			return
		}

		decryptedJSON, err := signature.DecryptSignature(existingEncrypted)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to decrypt existing signature: " + err.Error(),
			})
			return
		}

		var existingSig SignatureData
		if err := json.Unmarshal([]byte(decryptedJSON), &existingSig); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to parse existing signature: " + err.Error(),
			})
			return
		}

		// Handle image update
		imagePath := existingSig.CardImage
		if req.CardImage != "" {
			// New image provided, save it
			imagePath, err = saveBase64Image(req.CardImage)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to save image: " + err.Error(),
				})
				return
			}
		}

		// Create updated signature data
		sigData := SignatureData{
			ID:        req.ID,
			Name:      req.Name,
			Intro:     req.Intro,
			CardImage: imagePath,
		}

		// Encrypt and save
		sigJSON, err := json.Marshal(sigData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to marshal signature data: " + err.Error(),
			})
			return
		}

		encryptedData, err := signature.EncryptSignature(string(sigJSON))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt signature: " + err.Error(),
			})
			return
		}

		config.SetValue(key, encryptedData)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Signature updated successfully",
			"data":    sigData,
		})
	})

	// Delete signature
	signatureGroup.DELETE("/delete/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Signature ID is required",
			})
			return
		}

		// Encrypt ID to find the signature
		encryptedID, err := signature.EncryptSignature(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt ID: " + err.Error(),
			})
			return
		}

		key := fmt.Sprintf("signature_manager.%s", encryptedID)

		// Check if signature exists
		existingData := config.GetValue(key)
		if existingData == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Signature not found",
			})
			return
		}

		// Delete from config (set to nil)
		config.SetValue(key, nil)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Signature deleted successfully",
		})
	})

	// Export signature as .ktsign file
	signatureGroup.GET("/export/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Signature ID is required",
			})
			return
		}

		// Encrypt ID to find the signature
		encryptedID, err := signature.EncryptSignature(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt ID: " + err.Error(),
			})
			return
		}

		key := fmt.Sprintf("signature_manager.%s", encryptedID)
		existingData := config.GetValue(key)
		if existingData == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Signature not found",
			})
			return
		}

		// Decrypt signature
		encryptedStr, ok := existingData.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid signature data",
			})
			return
		}

		decryptedJSON, err := signature.DecryptSignature(encryptedStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to decrypt signature: " + err.Error(),
			})
			return
		}

		var sigData SignatureData
		if err := json.Unmarshal([]byte(decryptedJSON), &sigData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to parse signature: " + err.Error(),
			})
			return
		}

		// Load image as base64
		var cardImageBase64 string
		if sigData.CardImage != "" {
			cardImageBase64, err = loadImageAsBase64(sigData.CardImage)
			if err != nil {
				logger.Error("Failed to load signature image", "error", err)
				// Continue without image
			}
		}

		// Create file format
		fileFormat := SignatureFileFormat{
			Version: "1.0.0",
			Signature: SignatureData{
				ID:        sigData.ID,
				Name:      sigData.Name,
				Intro:     sigData.Intro,
				CardImage: cardImageBase64,
			},
			Checksum: "", // Will be calculated
		}

		// Calculate checksum
		fileJSON, err := json.Marshal(fileFormat.Signature)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to marshal signature for checksum: " + err.Error(),
			})
			return
		}

		hash := sha256.Sum256(fileJSON)
		fileFormat.Checksum = fmt.Sprintf("%x", hash)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    fileFormat,
		})
	})

	// Import signature from .ktsign file
	signatureGroup.POST("/import", func(ctx *gin.Context) {
		var fileFormat SignatureFileFormat

		if err := ctx.ShouldBindJSON(&fileFormat); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid file format: " + err.Error(),
			})
			return
		}

		// Verify checksum
		fileJSON, err := json.Marshal(fileFormat.Signature)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to marshal signature for checksum: " + err.Error(),
			})
			return
		}

		hash := sha256.Sum256(fileJSON)
		expectedChecksum := fmt.Sprintf("%x", hash)

		if fileFormat.Checksum != expectedChecksum {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Checksum mismatch - file may be corrupted",
			})
			return
		}

		// Save image if provided
		var imagePath string
		if fileFormat.Signature.CardImage != "" {
			imagePath, err = saveBase64Image(fileFormat.Signature.CardImage)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to save image: " + err.Error(),
				})
				return
			}
		}

		// Create signature data
		sigData := SignatureData{
			ID:        fileFormat.Signature.ID,
			Name:      fileFormat.Signature.Name,
			Intro:     fileFormat.Signature.Intro,
			CardImage: imagePath,
		}

		// Encrypt signature data
		sigJSON, err := json.Marshal(sigData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to marshal signature data: " + err.Error(),
			})
			return
		}

		encryptedData, err := signature.EncryptSignature(string(sigJSON))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt signature: " + err.Error(),
			})
			return
		}

		// Encrypt ID for key
		encryptedID, err := signature.EncryptSignature(fileFormat.Signature.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt ID: " + err.Error(),
			})
			return
		}

		// Store in config
		key := fmt.Sprintf("signature_manager.%s", encryptedID)
		config.SetValue(key, encryptedData)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Signature imported successfully",
			"data":    sigData,
		})
	})

	// Serve signature card images
	signatureGroup.GET("/image/:filename", func(ctx *gin.Context) {
		filename := ctx.Param("filename")
		if filename == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Filename is required",
			})
			return
		}

		// Validate filename to prevent path traversal
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid filename",
			})
			return
		}

		filePath := filepath.Join(getSignatureImageDir(), filename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Image not found",
			})
			return
		}

		// Serve the file
		ctx.Header("Content-Type", "image/png")
		ctx.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		ctx.File(filePath)
	})
}
