package model

import (
	"database/sql"
	"os"
)

// InitDB 负责数据库连接，数据库名为 OpenListStrm
func InitDB() (*sql.DB, error) {
	// 确保 db 目录存在
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}
	db, err := sql.Open("sqlite3", "file:db/OpenListStrm.db?cache=shared&mode=rwc")
	return db, err
} 
