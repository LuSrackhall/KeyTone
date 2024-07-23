package server

import (
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

	// 运行gin
	r.Run("0.0.0.0:38888")
}

func mainRouters(r *gin.Engine) {
	mainSettingRouters := r.Group("/main-setting")

	// 给到'客户端'使用, 供其判断如何设置自启动功能
	mainSettingRouters.GET("/get-is-auto-run", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// 给到'前端'使用, 供其ui界面实现应用自启动的设置功能
	mainSettingRouters.GET("/setting-is-auto-run", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

}
