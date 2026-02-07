package services

import (
	"errors"
	"strings"
	"yuanbao/config"
	"yuanbao/models"
	"yuanbao/repositories"

	"gorm.io/gorm"
)

// SaveCommand 保存口令（用户上传，带验证）
func SaveCommand(content string, uploaderIP string) (*models.Command, error) {
	// 1. 去除首尾空格
	content = strings.TrimSpace(content)

	// 2. 长度验证
	if len(content) < 10 {
		return nil, errors.New("口令长度不能少于10个字符")
	}
	if len(content) > 500 {
		return nil, errors.New("口令长度不能超过500个字符")
	}

	// 3. 基本内容验证
	if strings.Contains(content, "http://") || strings.Contains(content, "https://") {
		return nil, errors.New("口令不能包含链接")
	}

	// 4. 保存到数据库（数据库会自动检查重复）
	command, err := repositories.SaveCommand(content, uploaderIP)
	if err != nil {
		// 检查是否是重复错误
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
			return nil, errors.New("该口令已存在，请勿重复提交")
		}
		return nil, err
	}

	return command, nil
}

// SaveCrawlerCommand 保存爬虫口令（无需IP）
func SaveCrawlerCommand(content string) (*models.Command, error) {
	// 1. 去除首尾空格
	content = strings.TrimSpace(content)

	// 2. 长度验证
	if len(content) < 10 {
		return nil, errors.New("口令长度不能少于10个字符")
	}
	if len(content) > 500 {
		return nil, errors.New("口令长度不能超过500个字符")
	}

	// 3. 基本内容验证
	if strings.Contains(content, "http://") || strings.Contains(content, "https://") {
		return nil, errors.New("口令不能包含链接")
	}

	// 4. 保存到数据库
	command, err := repositories.SaveCrawlerCommand(content)
	if err != nil {
		// 检查是否是重复错误
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
			return nil, errors.New("该口令已存在")
		}
		return nil, err
	}

	return command, nil
}

// GetRandomCommand 获取随机口令（排除同IP上传的，带悲观锁和事务）
func GetRandomCommand(clientIP string) (*models.Command, error) {
	var command *models.Command
	var err error

	// 开启事务
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中查询并锁定
		command, err = repositories.FindRandomCommandWithLock(clientIP)
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

			// 暂时注释掉自动删除功能，先观察数据量
			// 如果达到3次，立即删除
			// if command.DisplayCount >= 3 {
			// 	err = repositories.DeleteCommand(command.ID)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		}

		return nil
	})

	return command, err
}

// GetCount 获取可用口令数量
func GetCount() (int64, error) {
	return repositories.CountAvailableCommands()
}

// MarkAsInvalid 标记口令为无效（直接删除）
func MarkAsInvalid(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("口令内容不能为空")
	}

	err := repositories.MarkCommandAsInvalid(content)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("口令不存在或已被删除")
		}
		return err
	}

	return nil
}
