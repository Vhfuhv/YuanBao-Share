package repositories

import (
	"yuanbao/config"
	"yuanbao/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SaveCommand 保存口令
func SaveCommand(content string) (*models.Command, error) {
	command := &models.Command{
		Content:      content,
		DisplayCount: 0,
	}

	result := config.DB.Create(command)
	return command, result.Error
}

// FindRandomCommandWithLock 使用悲观锁查询随机口令
func FindRandomCommandWithLock() (*models.Command, error) {
	var command models.Command

	// 使用悲观锁 (SELECT ... FOR UPDATE)
	// 并发安全：同一时刻只有一个事务能锁定该行
	err := config.DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("display_count < ?", 3).
		Order("RAND()").
		First(&command).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &command, err
}

// UpdateCommand 更新口令
func UpdateCommand(command *models.Command) error {
	return config.DB.Save(command).Error
}

// CountAvailableCommands 统计可用口令数量
func CountAvailableCommands() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Command{}).
		Where("display_count < ?", 3).
		Count(&count).Error
	return count, err
}
