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
	"io"
	"net/http"
	"os"

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

		// 定义加密密钥（在实际应用中应该从安全的地方获取）
		// 这里使用的是一个示例密钥，长度为32字节（256位AES密钥）
		encryptionKey := []byte("KeyTone2024SignatureEncryptionKey"[:32]) // 截取前32字节

		// 调用signature包创建签名
		encryptedID, err := signature.CreateSignature(id, signatureData, imageData, imageExt, fileName, encryptionKey)
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
		// TODO: 实现完整的签名更新逻辑
		// 需要：
		// 1. 从请求中获取签名 ID、name、intro、cardImage（Base64）
		// 2. 从配置中读取现有的签名存储条目
		// 3. 保留原有的 sort.time 时间戳（不更改排序时间）
		// 4. 处理图片：如果有新图片，保存并更新路径；如果没有新图片，保留原路径或清除
		// 5. 加密新的签名数据
		// 6. 更新配置文件中的 value，但保留 sort.time 不变
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
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
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
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
		}

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的请求参数",
			})
			return
		}

		// 定义加密密钥
		encryptionKey := []byte("KeyTone2024SignatureEncryptionKey"[:32])

		// 调用解密函数
		decryptedValue, err := signature.DecryptData(req.EncryptedValue, encryptionKey)
		if err != nil {
			logger.Error("解密签名数据失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "解密失败: " + err.Error(),
			})
			return
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
		// TODO: 实现完整的签名导入逻辑
		// 需要：
		// 1. 从请求中获取 .ktsign 文件内容
		// 2. 解析文件并验证格式和校验和
		// 3. 检查签名 ID 是否已存在
		// 4. 如果存在，询问用户是否覆盖
		// 5. 保存图片文件（Base64 解码）
		// 6. 加密签名数据
		// 7. 生成 sort.time（当前时间戳），用于排序
		// 8. 存储到配置文件
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

}
