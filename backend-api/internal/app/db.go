package app

import (
	"database/sql"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"go.uber.org/zap"
)

func InitDatabase() (*sql.DB, error) {
	db, err := model.InitDB()
	if err != nil {
		logger.Error("InitDB failed", zap.Error(err))
		return nil, err
	}
	if err := model.AutoMigrateAll(db); err != nil {
		logger.Error("AutoMigrateAll failed", zap.Error(err))
		return nil, err
	}
	logger.Info("AutoMigrate executed")
	return db, nil
} 
