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
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func signatureRouters(r *gin.Engine) {

	signatureRouter := r.Group("/")

	// 创建签名
	signatureRouter.POST("/signature/create", func(ctx *gin.Context) {
		// 解析表单数据
		id := ctx.PostForm("id")
		name := ctx.PostForm("name")
		intro := ctx.PostForm("intro")

		// 验证必填字段
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少必填字段: id",
			})
			return
		}

		if name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少必填字段: name",
			})
			return
		}

		// 获取上传的图片文件（可选）
		var imageData []byte
		var imageExt string
		var fileName string
		file, err := ctx.FormFile("cardImage")
		if err == nil && file != nil {
			// 获取文件名信息(包括文件名+扩展名)
			fileName = file.Filename
			// 从文件名中提取扩展名
			extIndex := -1
			for i := len(fileName) - 1; i >= 0; i-- {
				if fileName[i] == '.' {
					extIndex = i
					break
				}
			}
			if extIndex != -1 {
				imageExt = fileName[extIndex:]
			}

			// 读取文件内容
			fileContent, err := file.Open()
			if err != nil {
				logger.Error("打开上传的图片文件失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "处理图片文件失败",
				})
				return
			}
			defer fileContent.Close()

			imageData, err = io.ReadAll(fileContent)
			if err != nil {
				logger.Error("读取上传的图片文件内容失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "读取图片文件失败",
				})
				return
			}
		}

		// 构建签名数据结构
		signatureData := signature.SignatureData{
			Name:  name,
			Intro: intro,
		}

		// 使用密钥A（用于ID加密）
		// 动态密钥会在signature.CreateSignature中自动生成
		keyA := signature.GetKeyA()

		// 调用signature包创建签名
		encryptedID, err := signature.CreateSignature(id, signatureData, imageData, imageExt, fileName, keyA)
		if err != nil {
			logger.Error("创建签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建签名失败: " + err.Error(),
			})
			return
		}

		// 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "签名创建成功",
			"data": gin.H{
				"id": encryptedID,
			},
		})
	})

	// 更新签名
	signatureRouter.POST("/signature/update", func(ctx *gin.Context) {
		// 解析表单数据
		encryptedID := ctx.PostForm("encryptedId")
		name := ctx.PostForm("name")
		intro := ctx.PostForm("intro")
		removeImage := ctx.PostForm("removeImage")   // 获取删除图片标记
		imageChanged := ctx.PostForm("imageChanged") // 获取图片是否发生变更的标记

		// 验证必填字段
		if encryptedID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少必填字段: encryptedId",
			})
			return
		}

		// 获取上传的图片文件（可选）
		var imageData []byte
		var imageExt string
		var fileName string
		file, err := ctx.FormFile("cardImage")
		if err == nil && file != nil {
			// 获取文件名信息(包括文件名+扩展名)
			fileName = file.Filename
			// 从文件名中提取扩展名
			extIndex := -1
			for i := len(fileName) - 1; i >= 0; i-- {
				if fileName[i] == '.' {
					extIndex = i
					break
				}
			}
			if extIndex != -1 {
				imageExt = fileName[extIndex:]
			}

			// 读取文件内容
			fileContent, err := file.Open()
			if err != nil {
				logger.Error("打开上传的图片文件失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "处理图片文件失败",
				})
				return
			}
			defer fileContent.Close()

			imageData, err = io.ReadAll(fileContent)
			if err != nil {
				logger.Error("读取上传的图片文件内容失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "读取图片文件失败",
				})
				return
			}
		}

		// 构建签名数据结构
		signatureData := signature.SignatureData{
			Name:  name,
			Intro: intro,
		}

		// 使用密钥A（用于ID加密）
		keyA := signature.GetKeyA()

		// 调用signature包更新签名，传递 removeImage 和 imageChanged 标记
		shouldRemoveImage := removeImage == "true"
		hasImageChanged := imageChanged != "false" // 默认为 true（向后兼容）

		logger.Debug("更新签名图片状态",
			"encryptedId", encryptedID,
			"removeImage", shouldRemoveImage,
			"imageChanged", hasImageChanged,
		)

		err = signature.UpdateSignature(encryptedID, signatureData, imageData, imageExt, fileName, keyA, shouldRemoveImage, hasImageChanged)
		if err != nil {
			logger.Error("更新签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新签名失败: " + err.Error(),
			})
			return
		}

		// 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "签名更新成功",
		})
	})

	// TODO: 实现更新签名排序时间戳的 API 端点
	// 需要支持：
	// 1. 当用户执行拖动排序操作时，更新指定签名的 sort.time 值
	// 2. 后端需要读取现有的签名存储数据，更新排序时间戳，然后保存回配置文件
	// 3. 前端拖动排序完成后调用此 API 提交新的排序顺序
	signatureRouter.POST("/signature/update-sort", func(ctx *gin.Context) {
		var req struct {
			SortOrder []struct {
				ID       string `json:"id" binding:"required"`
				SortTime int64  `json:"sortTime" binding:"required"`
			} `json:"sortOrder" binding:"required"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			logger.Error("绑定请求参数失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的请求参数",
			})
			return
		}

		// 从配置中获取现有的签名存储数据
		signatureMapValue := config.GetValue("signature")
		if signatureMapValue == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "签名不存在",
			})
			return
		}

		// 类型转换
		if m, ok := signatureMapValue.(map[string]interface{}); ok {
			// 为每个签名更新排序时间戳
			for _, item := range req.SortOrder {
				if entry, ok := m[item.ID].(map[string]interface{}); ok {
					// 更新 sort.time
					if sort, ok := entry["sort"].(map[string]interface{}); ok {
						sort["time"] = item.SortTime
					} else {
						// 创建新的 sort 对象
						entry["sort"] = map[string]interface{}{
							"time": item.SortTime,
						}
					}
				}
			}

			// 保存回配置文件
			config.SetValue("signature", m)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "排序更新成功",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "配置数据格式错误",
			})
		}
	})

	// 删除签名
	signatureRouter.POST("/signature/delete", func(ctx *gin.Context) {
		var req struct {
			ID string `json:"id" binding:"required"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的请求参数",
			})
			return
		}

		// 从配置中获取签名
		signatureMap := config.GetValue("signature")
		if signatureMap == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "签名不存在",
			})
			return
		}

		// 类型转换并删除
		if m, ok := signatureMap.(map[string]interface{}); ok {
			if _, exists := m[req.ID]; !exists {
				ctx.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "签名不存在",
				})
				return
			}

			// 删除签名
			delete(m, req.ID)

			// 更新配置
			config.SetValue("signature", m)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "签名删除成功",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "配置数据格式错误",
			})
		}
	})

	// 导出签名
	signatureRouter.POST("/signature/export", func(ctx *gin.Context) {
		// 获取请求中的签名 ID（已加密）
		var req struct {
			EncryptedID string `json:"encryptedId" binding:"required"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			logger.Error("绑定请求参数失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少必填字段: encryptedId",
			})
			return
		}

		// 使用密钥A（用于解密ID）
		keyA := signature.GetKeyA()

		// 1. 调用导出函数获取导出数据
		exportData, err := signature.ExportSignature(req.EncryptedID, keyA)
		if err != nil {
			logger.Error("导出签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出签名失败: " + err.Error(),
			})
			return
		}

		// 2. 将导出数据转换为 JSON
		jsonData, err := json.Marshal(exportData)
		if err != nil {
			logger.Error("签名导出数据JSON序列化失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "序列化导出数据失败",
			})
			return
		}

		// 3. 对 JSON 字符串进行加密（使用密钥B）
		keyB := signature.GetKeyB()
		encryptedJSON, err := signature.EncryptData(string(jsonData), keyB)
		if err != nil {
			logger.Error("签名导出数据加密失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "加密导出数据失败",
			})
			return
		}

		// 4. 将加密后的数据转为二进制
		binaryData := []byte(encryptedJSON)

		// 5. 生成文件名（使用签名名称）
		fileName := exportData.Name + ".ktsign"
		// 清理文件名中的非法字符
		fileName = strings.Map(func(r rune) rune {
			if r == '<' || r == '>' || r == ':' || r == '"' || r == '/' || r == '\\' || r == '|' || r == '?' || r == '*' {
				return '_'
			}
			return r
		}, fileName)

		// 6. 设置响应头，返回文件
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
		ctx.Header("Content-Length", strconv.Itoa(len(binaryData)))

		logger.Info("签名导出完成",
			"encryptedID", req.EncryptedID,
			"文件名", fileName,
			"数据大小", len(binaryData),
		)

		ctx.Data(http.StatusOK, "application/octet-stream", binaryData)
	})

	// 获取签名列表（加密的key-value对）
	signatureRouter.GET("/signature/list", func(ctx *gin.Context) {
		// 从配置中获取所有签名数据
		signatureMap := config.GetValue("signature")
		if signatureMap == nil {
			// 没有签名时返回空对象
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    make(map[string]interface{}),
			})
			return
		}

		// 类型转换为 map[string]interface{}
		// 兼容新格式 map[string]SignatureStorageEntry 和旧格式 map[string]string
		encryptedSignatures := make(map[string]interface{})
		if m, ok := signatureMap.(map[string]interface{}); ok {
			for k, v := range m {
				if entry, ok := v.(map[string]interface{}); ok {
					// 新格式：SignatureStorageEntry
					encryptedSignatures[k] = entry
				} else if str, ok := v.(string); ok {
					// 旧格式：直接是加密字符串，需要升级为新格式
					logger.Warn("检测到旧格式的签名数据，正在进行格式升级", "key", k)
					encryptedSignatures[k] = map[string]interface{}{
						"value": str,
						"sort": map[string]interface{}{
							"time": 0, // TODO: 应该获取文件创建时间或使用其他策略
						},
					}
				}
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    encryptedSignatures,
		})
	})

	// 解密单个签名数据
	signatureRouter.POST("/signature/decrypt", func(ctx *gin.Context) {
		var req struct {
			EncryptedValue string `json:"encryptedValue" binding:"required"`
			EncryptedID    string `json:"encryptedId"` // 可选：用于动态密钥解密
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的请求参数",
			})
			return
		}

		var decryptedValue string
		var err error

		// 如果提供了encryptedID，使用动态密钥解密；否则使用旧方式
		if req.EncryptedID != "" {
			// 新方案：使用动态密钥解密
			decryptedValue, err = signature.DecryptValueWithDynamicKey(req.EncryptedValue, req.EncryptedID)
			if err != nil {
				logger.Error("使用动态密钥解密签名数据失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "解密失败: " + err.Error(),
				})
				return
			}
		} else {
			// 旧方案：使用KeyA解密（兼容旧数据）
			keyA := signature.GetKeyA()
			decryptedValue, err = signature.DecryptData(req.EncryptedValue, keyA)
			if err != nil {
				logger.Error("解密签名数据失败", "error", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "解密失败: " + err.Error(),
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    decryptedValue,
		})
	})

	// 获取签名图片
	signatureRouter.POST("/signature/get-image", func(ctx *gin.Context) {
		var req struct {
			ImagePath string `json:"imagePath" binding:"required"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的请求参数",
			})
			return
		}

		// 读取图片文件
		imageData, err := os.ReadFile(req.ImagePath)
		if err != nil {
			logger.Error("读取图片文件失败", "path", req.ImagePath, "error", err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "图片文件不存在",
			})
			return
		}

		// 设置响应头为图片类型
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "inline")
		ctx.Data(http.StatusOK, "application/octet-stream", imageData)
	})

	// 导入签名
	signatureRouter.POST("/signature/import", func(ctx *gin.Context) {
		// 获取上传的 .ktsign 文件
		file, err := ctx.FormFile("file")
		if err != nil {
			logger.Error("获取上传文件失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少文件或文件格式错误",
			})
			return
		}

		// 读取文件内容
		fileContent, err := file.Open()
		if err != nil {
			logger.Error("打开上传的文件失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "无法读取文件",
			})
			return
		}
		defer fileContent.Close()

		// 读取文件数据
		fileData, err := io.ReadAll(fileContent)
		if err != nil {
			logger.Error("读取文件内容失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "读取文件失败",
			})
			return
		}

		// 定义加密密钥
		// 导入时使用密钥B（与导出时相同）
		keyB := signature.GetKeyB()

		// 1. 解密文件数据（文件内容本身是加密的字符串，需要先转成字符串）
		encryptedJSON := string(fileData)

		// 2. 使用密钥B解密
		decryptedJSON, err := signature.DecryptData(encryptedJSON, keyB)
		if err != nil {
			logger.Error("解密导入文件失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "文件格式错误或密钥不匹配，无法解密",
			})
			return
		}

		// 3. 解析 JSON 为 SignatureExportData
		var exportData signature.SignatureExportData
		if err := json.Unmarshal([]byte(decryptedJSON), &exportData); err != nil {
			logger.Error("解析导入数据失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "导入文件格式错误",
			})
			return
		}

		// 4. 检查签名是否已存在（先不覆盖，返回冲突状态）
		// 使用密钥A（用于ID加密）
		keyA := signature.GetKeyA()
		encryptedID, conflict, err := signature.ImportSignature(&exportData, keyA)
		if err != nil {
			logger.Error("导入签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导入失败: " + err.Error(),
			})
			return
		}

		// 5. 如果存在冲突，返回冲突标记和现有签名的 ID，让前端决定是否覆盖
		if conflict {
			logger.Info("导入的签名已存在，等待用户确认是否覆盖", "encryptedID", encryptedID)
			ctx.JSON(http.StatusConflict, gin.H{
				"success":  false,
				"conflict": true,
				"message":  "签名已存在",
				"data": gin.H{
					"encryptedId": encryptedID,
					"name":        exportData.Name,
				},
			})
			return
		}

		// 6. 导入成功
		logger.Info("签名导入成功", "encryptedID", encryptedID, "名称", exportData.Name)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "签名导入成功",
			"data": gin.H{
				"encryptedId": encryptedID,
				"name":        exportData.Name,
			},
		})
	})

	// 导入签名（处理冲突 - 覆盖或保留）
	signatureRouter.POST("/signature/import-confirm", func(ctx *gin.Context) {
		var req struct {
			File      string `json:"file" binding:"required"`      // Base64 编码的文件内容或文件路径
			Overwrite bool   `json:"overwrite" binding:"required"` // 是否覆盖现有签名
		}

		if err := ctx.BindJSON(&req); err != nil {
			logger.Error("绑定请求参数失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少必填字段",
			})
			return
		}

		// 使用密钥B解密（与导出时相同）
		keyB := signature.GetKeyB()

		// 解密文件数据
		decryptedJSON, err := signature.DecryptData(req.File, keyB)
		if err != nil {
			logger.Error("解密导入文件失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "文件格式错误或密钥不匹配",
			})
			return
		}

		// 解析 JSON
		var exportData signature.SignatureExportData
		if err := json.Unmarshal([]byte(decryptedJSON), &exportData); err != nil {
			logger.Error("解析导入数据失败", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "导入文件格式错误",
			})
			return
		}

		// 调用带覆盖选项的导入函数
		// 使用密钥A（用于ID加密）
		keyA := signature.GetKeyA()
		encryptedID, err := signature.ImportSignatureWithOverwrite(&exportData, req.Overwrite, keyA)
		if err != nil {
			logger.Error("导入签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导入失败: " + err.Error(),
			})
			return
		}

		logger.Info("签名导入完成", "encryptedID", encryptedID, "覆盖", req.Overwrite)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "签名导入成功",
			"data": gin.H{
				"encryptedId": encryptedID,
				"name":        exportData.Name,
			},
		})
	})

}
