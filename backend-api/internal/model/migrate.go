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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password_hash TEXT,
		is_active INTEGER,
		code TEXT,
		code_expire_at DATETIME,
		failed_login_count INTEGER DEFAULT 0,
		locked_until DATETIME,
		created_at DATETIME
	)`)
	return err
} 
