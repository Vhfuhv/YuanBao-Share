package controllers

import (
	"net/http"
	"yuanbao/services"

	"github.com/gin-gonic/gin"
)

// getClientIP 获取并标准化客户端IP
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// 将IPv6本地地址转换为IPv4格式
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

// UploadCommandRequest 上传口令请求
type UploadCommandRequest struct {
	Content string `json:"content" binding:"required"`
}

// UploadCommand 上传口令
func UploadCommand(c *gin.Context) {
	var req UploadCommandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "口令内容不能为空",
		})
		return
	}

	// 获取客户端IP（标准化处理）
	clientIP := getClientIP(c)

	command, err := services.SaveCommand(req.Content, clientIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(), // 返回具体的错误信息
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "口令上传成功",
		"id":      command.ID,
	})
}

// GetRandomCommand 随机获取口令
func GetRandomCommand(c *gin.Context) {
	// 获取客户端IP（标准化处理）
	clientIP := getClientIP(c)

	command, err := services.GetRandomCommand(clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取失败",
		})
		return
	}

	if command == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "暂无可用口令",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"content":   command.Content,
		"createdAt": command.CreatedAt,
	})
}

// GetCount 获取可用口令数量
func GetCount(c *gin.Context) {
	count, err := services.GetCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// ReportInvalid 报告无效口令
func ReportInvalid(c *gin.Context) {
	var req UploadCommandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	err := services.MarkAsInvalid(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已标记为无效",
	})
}
