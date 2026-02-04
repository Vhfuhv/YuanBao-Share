package main

import (
	"yuanbao/config"
	"yuanbao/controllers"
	"yuanbao/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移数据库表
	config.DB.AutoMigrate(&models.Command{})

	// 创建 Gin 路由
	r := gin.Default()

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API 路由
	api := r.Group("/api/commands")
	{
		api.POST("", controllers.UploadCommand)
		api.GET("/random", controllers.GetRandomCommand)
		api.GET("/count", controllers.GetCount)
	}

	// 启动服务器
	r.Run(":18080")
}
