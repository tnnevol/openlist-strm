package model

import "database/sql"

// 统一管理所有表的迁移
func AutoMigrateAll(db *sql.DB) error {
	if err := AutoMigrateUser(db); err != nil {
		return err
	}
	// 未来可添加更多表的迁移
	return nil
}

// 单独拆分每个表的迁移
func AutoMigrateUser(db *sql.DB) error {
	// 创建表（如果不存在）
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password_hash TEXT,
		is_active INTEGER,
		code TEXT,
		code_expire_at DATETIME,
		failed_login_count INTEGER DEFAULT 0,
		locked_until DATETIME,
		created_at DATETIME
	)`)
	if err != nil {
		return err
	}
	
	// 检查username字段是否存在，如果不存在则添加
	_, err = db.Exec(`ALTER TABLE user ADD COLUMN username TEXT UNIQUE`)
	if err != nil {
		// 字段已存在，忽略错误
	}
	
	return nil
} 
