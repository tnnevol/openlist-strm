package model

import (
	"strings"

	"gorm.io/gorm"
)

// 统一管理所有表的迁移（仅在表不存在时创建表，不做结构变更）
func MigrateIfNotExists(db *gorm.DB) error {
	models := []interface{}{
		&User{},
		&OpenListService{},
		&StrmConfig{},
		&StrmTask{},
		&LogRecord{},
		&Dict{},
	}
	for _, m := range models {
		if !db.Migrator().HasTable(m) {
			err := db.Migrator().CreateTable(m)
			if err != nil && strings.Contains(err.Error(), "already exists") {
					// 忽略已存在错误
					continue
			} else if err != nil {
					return err
			}
	}
	}
	return nil
} 
