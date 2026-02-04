package models

import (
	"time"
)

// Command 口令实体
type Command struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Content      string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"content"` // 添加唯一索引防重复
	DisplayCount int       `gorm:"not null;default:0" json:"display_count"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (Command) TableName() string {
	return "commands"
}
