package service

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"github.com/tnnevol/openlist-strm/backend-api/internal/util"
	"gorm.io/gorm"
)

func ListStrmConfigs(db *gorm.DB, userID, serviceID, page, pageSize int) ([]*model.StrmConfig, int, error) {
	var configs []*model.StrmConfig
	var err error
	if serviceID > 0 {
		configs, err = model.GetStrmConfigsByServiceID(db, serviceID)
	} else {
		configs, err = model.GetAllStrmConfigs(db)
	}
	if err != nil {
		return nil, 0, err
	}
	// 只保留当前用户的数据
	var filtered []*model.StrmConfig
	for _, cfg := range configs {
		if cfg != nil && cfg.UserID == userID {
			filtered = append(filtered, cfg)
		}
	}
	paged, total := util.Paginate(filtered, page, pageSize)
	return paged, total, nil
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
