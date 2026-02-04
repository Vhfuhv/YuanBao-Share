package services

import (
	"yuanbao/config"
	"yuanbao/models"
	"yuanbao/repositories"

	"gorm.io/gorm"
)

// SaveCommand 保存口令
func SaveCommand(content string) (*models.Command, error) {
	return repositories.SaveCommand(content)
}

// GetRandomCommand 获取随机口令（带悲观锁和事务）
func GetRandomCommand() (*models.Command, error) {
	var command *models.Command
	var err error

	// 开启事务
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中查询并锁定
		command, err = repositories.FindRandomCommandWithLock()
		if err != nil {
			return err
		}

		if command != nil {
			// 更新展示次数
			command.DisplayCount++
			err = repositories.UpdateCommand(command)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return command, err
}

// GetCount 获取可用口令数量
func GetCount() (int64, error) {
	return repositories.CountAvailableCommands()
}
