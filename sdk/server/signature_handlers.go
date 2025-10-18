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
	"KeyTone/logger"
	"KeyTone/signature"
	"io"
	"net/http"

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
		encryptedID, err := signature.CreateSignature(id, signatureData, imageData, imageExt, encryptionKey)
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
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 删除签名
	signatureRouter.POST("/signature/delete", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 导出签名
	signatureRouter.POST("/signature/export", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 导入签名
	signatureRouter.POST("/signature/import", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

}
