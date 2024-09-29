package server

import (
	audioPackageConfig "KeyTone/audioPackage/config"
	"KeyTone/config"
	"KeyTone/logger"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ServerRun() {
	// 启动gin
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	mainRouters(r)

	r.GET("/stream", func(c *gin.Context) {

		logger.Logger.Debug("新生成了一个线程............................")
		logger.Debug("新生成了一个线程............................")

		clientStoresChan := make(chan *config.Store)

		serverStoresChan := make(chan bool, 1)

		config.Clients_sse_stores.Store(clientStoresChan, serverStoresChan)

		defer func() {
			config.Clients_sse_stores.Delete(clientStoresChan)

			logger.Logger.Debug("一个线程退出了............................")
			logger.Debug("一个线程退出了............................")
		}()

		clientGone := c.Request.Context().Done()

		for {
			re := c.Stream(func(w io.Writer) bool {
				select {
				case <-clientGone:
					serverStoresChan <- false

					return false

				case message, ok := <-clientStoresChan:
					if !ok {
						logger.Error("通道clientStoresChan非正常关闭")
						return true
					}
					c.SSEvent("message", message)
					return true
				}
			})

			if !re {
				return
			}
		}

	})

	keytonePkgRouters(r)

	// 运行gin
	r.Run("0.0.0.0:38888")
}

func mainRouters(r *gin.Engine) {
	settingStoreRouters := r.Group("/store")

	// 给到'客户端'或'前端'使用, 供它们获取持久化的设置。
	settingStoreRouters.GET("/get", func(ctx *gin.Context) {

		// key := ctx.Query("key")
		key := ctx.DefaultQuery("key", "unknown")

		if key == "unknown" || key == "" {
			ctx.JSON(200, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:",
			})
			return
		}

		value := config.GetValue(key)

		fmt.Println("查询到的value= ", value)

		ctx.JSON(200, gin.H{
			"message": "ok",
			"key":     key,
			// 这里的value, 会自动转换为JSON字符串
			"value": value,
		})
	})

	// 给到'前端'使用, 供其ui界面实现应用的设置功能
	settingStoreRouters.POST("/set", func(ctx *gin.Context) {
		type SettingStore struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		}

		var store_setting SettingStore
		err := ctx.ShouldBind(&store_setting)
		if err != nil || store_setting.Key == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		config.SetValue(store_setting.Key, store_setting.Value)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

}

func keytonePkgRouters(r *gin.Engine) {

	keytonePkgRouters := r.Group("/keytone_pkg")

	// 接收前端上传的音频文件, 并存入本地路径
	keytonePkgRouters.POST("/add_new_sound_file", func(ctx *gin.Context) {
		audioPkgUUID := ctx.PostForm("audioPkgUUID")
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件添加失败, 传输问题:" + err.Error(),
			})
			return
		}

		// 保存文件
		err = ctx.SaveUploadedFile(file, filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", file.Filename))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件添加失败, 后端保存过程中发生错误:" + err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "ok",
		})

	})
}
