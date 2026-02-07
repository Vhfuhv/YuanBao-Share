package repositories

import (
	"time"
	"yuanbao/config"
	"yuanbao/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SaveCommand 保存口令（用户上传）
func SaveCommand(content string, uploaderIP string) (*models.Command, error) {
	command := &models.Command{
		Content:      content,
		Source:       "user",
		UploaderIP:   uploaderIP,
		DisplayCount: 0,
	}

	result := config.DB.Create(command)
	return command, result.Error
}

// SaveCrawlerCommand 保存爬虫口令
func SaveCrawlerCommand(content string) (*models.Command, error) {
	command := &models.Command{
		Content:      content,
		Source:       "crawler",
		DisplayCount: 0,
	}

	result := config.DB.Create(command)
	return command, result.Error
}

// FindRandomCommandWithLock 使用悲观锁查询随机口令（优先用户上传，排除同IP）
func FindRandomCommandWithLock(clientIP string) (*models.Command, error) {
	var command models.Command

	// 使用悲观锁 (SELECT ... FOR UPDATE)
	// 并发安全：同一时刻只有一个事务能锁定该行
	// SQLite 使用 RANDOM()，MySQL 使用 RAND()

	// 1. 优先查找用户上传的token（排除同IP）
	err := config.DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("display_count < ?", 3).
		Where("source = ?", "user").
		Where("uploader_ip != ? OR uploader_ip IS NULL OR uploader_ip = ''", clientIP).
		Order("RANDOM()").
		First(&command).Error

	// 如果找到用户上传的token，直接返回
	if err == nil {
		return &command, nil
	}

	// 2. 如果没有用户上传的token，查找爬虫token
	if err == gorm.ErrRecordNotFound {
		err = config.DB.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("display_count < ?", 3).
			Where("source = ?", "crawler").
			Order("RANDOM()").
			First(&command).Error

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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

// MarkCommandAsInvalid 标记口令为无效（直接删除）
func MarkCommandAsInvalid(content string) error {
	result := config.DB.Where("content = ?", content).Delete(&models.Command{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CleanOldCrawlerCommands 清理1小时前的爬虫口令
func CleanOldCrawlerCommands() (int64, error) {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	result := config.DB.
		Where("source = ?", "crawler").
		Where("created_at < ?", oneHourAgo).
		Delete(&models.Command{})

	return result.RowsAffected, result.Error
}

// CleanAllCommands 清空所有口令（每天0点执行）
func CleanAllCommands() (int64, error) {
	result := config.DB.Delete(&models.Command{}, "1=1")
	return result.RowsAffected, result.Error
}

