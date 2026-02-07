package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRecord IP访问记录
type IPRecord struct {
	Count      int
	ResetTime  time.Time
	mu         sync.Mutex
}

// RateLimiter 频率限制器
type RateLimiter struct {
	records sync.Map // map[string]*IPRecord
	limit   int
	window  time.Duration
}

// NewRateLimiter 创建频率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		limit:  limit,
		window: window,
	}

	// 启动清理协程，每分钟清理过期记录
	go limiter.cleanup()

	return limiter
}

// Allow 检查是否允许访问
func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()

	// 获取或创建记录
	value, _ := rl.records.LoadOrStore(ip, &IPRecord{
		Count:     0,
		ResetTime: now.Add(rl.window),
	})

	record := value.(*IPRecord)
	record.mu.Lock()
	defer record.mu.Unlock()

	// 如果已过期，重置计数
	if now.After(record.ResetTime) {
		record.Count = 0
		record.ResetTime = now.Add(rl.window)
	}

	// 检查是否超限
	if record.Count >= rl.limit {
		return false
	}

	// 增加计数
	record.Count++
	return true
}

// cleanup 定期清理过期记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		rl.records.Range(func(key, value interface{}) bool {
			record := value.(*IPRecord)
			record.mu.Lock()
			// 如果记录已过期超过5分钟，删除
			if now.After(record.ResetTime.Add(5 * time.Minute)) {
				rl.records.Delete(key)
			}
			record.mu.Unlock()
			return true
		})
	}
}

// Middleware 创建限流中间件
func (rl *RateLimiter) Middleware(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(ip) {
			var message string
			if action == "upload" {
				message = "同一IP每分钟最多上传5次，请稍后再试"
			} else if action == "get" {
				message = "您的获取次数已达上限（每分钟20次），请稍后再试"
			} else {
				message = "操作过于频繁，请稍后再试"
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": message,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
