package controllers

import (
	"net/http"
	"yuanbao/services"

	"github.com/gin-gonic/gin"
)

// UploadCommandRequest 上传口令请求
type UploadCommandRequest struct {
	Content string `json:"content" binding:"required"`
}

// UploadCommand 上传口令
func UploadCommand(c *gin.Context) {
	var req UploadCommandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "口令内容不能为空",
		})
		return
	}

	command, err := services.SaveCommand(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存失败",
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
	command, err := services.GetRandomCommand()
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
