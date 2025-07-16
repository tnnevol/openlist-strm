package service

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"gorm.io/gorm"
)

func ListStrmConfigs(db *gorm.DB, userID, serviceID, page, pageSize int) ([]*model.StrmConfig, int64, error) {
	if serviceID > 0 {
		return model.GetStrmConfigsByServiceID(db, serviceID, userID, page, pageSize)
	}
	return model.GetAllStrmConfigs(db, userID, page, pageSize)
}

func CopyStrmConfigs(db *gorm.DB, ids []int) error {
	for _, id := range ids {
		cfg, err := model.GetStrmConfigByID(db, id)
		if err != nil || cfg == nil {
			return err
		}
		cfg.ID = 0
		cfg.Name = cfg.Name + "-复制"
		cfg.CreatedAt = time.Now()
		cfg.UpdatedAt = time.Now()
		err = model.CreateStrmConfig(db, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateStrmConfig(db *gorm.DB, config *model.StrmConfig) error {
	return model.CreateStrmConfig(db, config)
}

func UpdateStrmConfig(db *gorm.DB, config *model.StrmConfig) error {
	return model.UpdateStrmConfig(db, config)
}

func DeleteStrmConfig(db *gorm.DB, id int) error {
	return model.DeleteStrmConfig(db, id)
}

func GetStrmConfigByID(db *gorm.DB, id int) (*model.StrmConfig, error) {
	return model.GetStrmConfigByID(db, id)
}
