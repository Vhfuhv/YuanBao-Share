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
	// SQLite 使用 RANDOM()，MySQL 使用 RAND()
	err := config.DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("display_count < ?", 3).
		Order("RANDOM()").
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

// DeleteCommand 删除口令
func DeleteCommand(id uint) error {
	return config.DB.Delete(&models.Command{}, id).Error
}

// CountAvailableCommands 统计可用口令数量
func CountAvailableCommands() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Command{}).
		Where("display_count < ?", 3).
		Count(&count).Error
	return count, err
}

// MarkCommandAsInvalid 标记口令为无效
func MarkCommandAsInvalid(content string) error {
	return config.DB.Model(&models.Command{}).
		Where("content = ?", content).
		Update("display_count", 3).Error
}
