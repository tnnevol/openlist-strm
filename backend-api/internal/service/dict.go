package service

import (
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateDict(db *gorm.DB, d *model.Dict) error {
	logger.Info("[Service] CreateDict called", zap.String("type", d.Type), zap.String("key", d.Key))
	return model.CreateDict(db, d)
}

func UpdateDict(db *gorm.DB, d *model.Dict) error {
	logger.Info("[Service] UpdateDict called", zap.Int("id", d.ID))
	return model.UpdateDict(db, d)
}

func DeleteDict(db *gorm.DB, id int) error {
	logger.Info("[Service] DeleteDict called", zap.Int("id", id))
	return model.DeleteDict(db, id)
}

func GetDictByID(db *gorm.DB, id int) (*model.Dict, error) {
	logger.Info("[Service] GetDictByID called", zap.Int("id", id))
	return model.GetDictByID(db, id)
}

func ListDicts(db *gorm.DB, dictType string) ([]*model.Dict, error) {
	logger.Info("[Service] ListDicts called", zap.String("type", dictType))
	return model.ListDicts(db, dictType)
} 
