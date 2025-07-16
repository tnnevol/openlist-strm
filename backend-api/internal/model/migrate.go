package model

import "database/sql"

// 统一管理所有表的迁移
func AutoMigrateAll(db *sql.DB) error {
	if err := AutoMigrateUser(db); err != nil {
		return err
	}
	if err := AutoMigrateOpenListService(db); err != nil {
		return err
	}
	if err := AutoMigrateStrmConfig(db); err != nil {
		return err
	}
	if err := AutoMigrateStrmTask(db); err != nil {
		return err
	}
	if err := AutoMigrateLogRecord(db); err != nil {
		return err
	}
	if err := AutoMigrateDict(db); err != nil {
		return err
	}
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
		created_at DATETIME,
		token_invalid_before DATETIME
	)`)
	if err != nil {
		return err
	}
	
	// 检查username字段是否存在，如果不存在则添加
	_, err = db.Exec(`ALTER TABLE user ADD COLUMN username TEXT UNIQUE`)
	if err != nil {
		// 字段已存在，忽略错误
	}
	
	// 检查token_invalid_before字段是否存在，如果不存在则添加
	_, err = db.Exec(`ALTER TABLE user ADD COLUMN token_invalid_before DATETIME`)
	if err != nil {
		// 字段已存在，忽略错误
	}
	
	return nil
}

// AutoMigrateOpenListService 迁移OpenList服务表
func AutoMigrateOpenListService(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS openlist_service (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		account TEXT NOT NULL,
		token TEXT NOT NULL,
		service_url TEXT NOT NULL,
		backup_url TEXT,
		enabled INTEGER DEFAULT 1,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	)`)
	return err
}

// AutoMigrateStrmConfig 迁移Strm配置表
func AutoMigrateStrmConfig(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS strm_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		alist_base_path TEXT NOT NULL,
		strm_output_path TEXT NOT NULL,
		download_enabled INTEGER DEFAULT 0,
		download_interval INTEGER DEFAULT 3600,
		update_mode TEXT DEFAULT 'incremental' CHECK (update_mode IN ('incremental', 'full')),
		service_id INTEGER NOT NULL,
		is_use_backup_url INTEGER DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (service_id) REFERENCES openlist_service(id) ON DELETE CASCADE
	)`)
	return err
}

// AutoMigrateStrmTask 迁移Strm任务表
func AutoMigrateStrmTask(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS strm_task (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		scheduled_time DATETIME NOT NULL,
		task_mode TEXT NOT NULL CHECK (task_mode IN ('create', 'check')),
		enabled INTEGER DEFAULT 1,
		service_id INTEGER NOT NULL,
		config_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (service_id) REFERENCES openlist_service(id) ON DELETE CASCADE,
		FOREIGN KEY (config_id) REFERENCES strm_config(id) ON DELETE CASCADE
	)`)
	return err
}

// AutoMigrateLogRecord 迁移日志记录表
func AutoMigrateLogRecord(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS log_record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL CHECK (name IN ('create', 'check')),
		log_path TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		task_status TEXT NOT NULL CHECK (task_status IN ('running', 'error', 'completed')),
		task_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (task_id) REFERENCES strm_task(id) ON DELETE CASCADE
	)`)
	return err
} 

// DropStrmConfigTable 删除 strm_config 表
func DropStrmConfigTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS strm_config`)
	return err
}

// DropStrmTaskTable 删除 strm_task 表
func DropStrmTaskTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS strm_task`)
	return err
}

// DropLogRecordTable 删除 log_record 表
func DropLogRecordTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS log_record`)
	return err
} 

func AutoMigrateDict(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS dict (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		description TEXT,
		parent_id INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		return err
	}
	// 检查 parent_id 字段是否存在，不存在则补充
	_, err = db.Exec(`ALTER TABLE dict ADD COLUMN parent_id INTEGER DEFAULT 0`)
	if err != nil {
		// 字段已存在可忽略
	}
	return nil
} 
