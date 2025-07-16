package model

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitDB 负责数据库连接，数据库名为 OpenListStrm
func InitDB() (*gorm.DB, error) {
	// 确保 db 目录存在
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}
	db, err := gorm.Open(sqlite.Open("db/OpenListStrm.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 表名不要加 s
			SingularTable: true,
			// 单词之间不要加下划线
			NoLowerCase: true,
		},
	})
	return db, err
} 
