package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Dict struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Type        string    `json:"type" gorm:"type:varchar(64);index"`
	Key         string    `json:"key" gorm:"type:varchar(64);index"`
	Value       string    `json:"value" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	ParentID    int       `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func CreateDict(db *gorm.DB, d *Dict) error {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	logger.Info("[DB] CreateDict", zap.String("type", d.Type), zap.String("key", d.Key))
	if err := db.Create(d).Error; err != nil {
		logger.Error("[DB] CreateDict error", zap.Error(err))
		return err
	}
	return nil
}

func UpdateDict(db *gorm.DB, d *Dict) error {
	d.UpdatedAt = time.Now()
	logger.Info("[DB] UpdateDict", zap.Int("id", d.ID))
	if err := db.Model(&Dict{}).Where("id = ?", d.ID).Updates(d).Error; err != nil {
		logger.Error("[DB] UpdateDict error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteDict(db *gorm.DB, id int) error {
	logger.Info("[DB] DeleteDict", zap.Int("id", id))
	if err := db.Delete(&Dict{}, id).Error; err != nil {
		logger.Error("[DB] DeleteDict error", zap.Error(err))
		return err
	}
	return nil
}

func GetDictByID(db *gorm.DB, id int) (*Dict, error) {
	var d Dict
	if err := db.First(&d, id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func ListDicts(db *gorm.DB, dictType string) ([]*Dict, error) {
	var dicts []*Dict
	var err error
	if dictType != "" {
		err = db.Where("type = ?", dictType).Order("id DESC").Find(&dicts).Error
	} else {
		err = db.Order("id DESC").Find(&dicts).Error
	}
	if err != nil {
		return nil, err
	}
	return dicts, nil
} 
