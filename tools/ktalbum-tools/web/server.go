package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ktalbum-tools/commands"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed all:frontend/dist
var distFS embed.FS

func StartServer(port int) error {
	r := gin.Default()
	r.Use(cors.Default())

	// 准备前端文件系统
	webFS, err := fs.Sub(distFS, "frontend/dist")
	if err != nil {
		return fmt.Errorf("准备前端文件系统失败: %v", err)
	}

	// 提供静态文件服务（前端页面）
	r.StaticFS("/", http.FS(webFS))

	// API 路由
	api := r.Group("/api")
	{
		// 解包 ktalbum 文件
		api.POST("/extract", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("文件上传失败: %v", err),
				})
				return
			}

			// 创建临时目录
			tempDir := filepath.Join(os.TempDir(), "ktalbum-tools")
			os.MkdirAll(tempDir, 0755)

			// 保存上传的文件
			inputPath := filepath.Join(tempDir, file.Filename)
			if err := c.SaveUploadedFile(file, inputPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("保存文件失败: %v", err),
				})
				return
			}

			// 解包文件
			outputPath := filepath.Join(tempDir, strings.TrimSuffix(file.Filename, ".ktalbum")+".zip")
			if err := commands.Extract(inputPath, outputPath, true); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("解包失败: %v", err),
				})
				return
			}

			// 返回 zip 文件
			c.FileAttachment(outputPath, filepath.Base(outputPath))

			// 清理临时文件
			os.Remove(inputPath)
			os.Remove(outputPath)
		})

		// 获取文件信息
		api.POST("/info", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("文件上传失败: %v", err),
				})
				return
			}

			// 创建临时目录
			tempDir := filepath.Join(os.TempDir(), "ktalbum-tools")
			os.MkdirAll(tempDir, 0755)

			// 保存上传的文件
			inputPath := filepath.Join(tempDir, file.Filename)
			if err := c.SaveUploadedFile(file, inputPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("保存文件失败: %v", err),
				})
				return
			}
			defer os.Remove(inputPath)

			// 读取文件信息
			info, err := commands.GetFileInfo(inputPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("读取文件信息失败: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, info)
		})
	}

	return r.Run(fmt.Sprintf(":%d", port))
} 