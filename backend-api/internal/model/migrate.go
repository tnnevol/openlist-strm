package model

import (
	"gorm.io/gorm"
)

// 统一管理所有表的迁移（GORM实现）
func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&OpenListService{},
		&StrmConfig{},
		&StrmTask{},
		&LogRecord{},
		&Dict{},
	)
} 
