package models

import (
	"time"
)

// Command 口令实体
type Command struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Content      string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"content"` // 添加唯一索引防重复
	Source       string    `gorm:"type:varchar(20);not null;default:'user';index" json:"source"` // 来源：crawler(爬虫) 或 user(用户上传)
	UploaderIP   string    `gorm:"type:varchar(50);index" json:"uploader_ip,omitempty"` // 上传者IP（仅用户上传时有值）
	DisplayCount int       `gorm:"not null;default:0" json:"display_count"`
	CreatedAt    time.Time `gorm:"not null;index" json:"created_at"` // 添加索引用于定时清理
}

// TableName 指定表名
func (Command) TableName() string {
	return "commands"
}
