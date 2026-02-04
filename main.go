package main

import (
	"time"
	"yuanbao/config"
	"yuanbao/controllers"
	"yuanbao/middleware"
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

	// 创建限流器
	uploadLimiter := middleware.NewRateLimiter(5, 1*time.Minute)   // 每分钟最多上传5次
	getLimiter := middleware.NewRateLimiter(10, 1*time.Minute)     // 每分钟最多获取10次

	// API 路由
	api := r.Group("/api/commands")
	{
		api.POST("", uploadLimiter.Middleware("upload"), controllers.UploadCommand)
		api.GET("/random", getLimiter.Middleware("get"), controllers.GetRandomCommand)
		api.GET("/count", controllers.GetCount) // 统计接口不限流
		api.POST("/report", controllers.ReportInvalid) // 报告无效口令
	}

	// 启动服务器
	r.Run(":18080")
}
