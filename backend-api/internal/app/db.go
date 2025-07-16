package app

import (
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := model.InitDB()
	if err != nil {
		logger.Error("InitDB failed", zap.Error(err))
		return nil, err
	}
	if err := model.MigrateIfNotExists(db); err != nil {
		logger.Error("MigrateIfNotExists failed", zap.Error(err))
		return nil, err
	}
	logger.Info("AutoMigrate executed")
	return db, nil
} 
