package service

import (
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"gorm.io/gorm"
)

func CreateOpenListService(db *gorm.DB, service *model.OpenListService) error {
	return model.CreateOpenListService(db, service)
}

func GetOpenListServicesByUserID(db *gorm.DB, userID, page, pageSize int) ([]*model.OpenListService, int64, error) {
	return model.GetOpenListServicesByUserID(db, userID, page, pageSize)
}

func GetOpenListServiceByID(db *gorm.DB, id int) (*model.OpenListService, error) {
	return model.GetOpenListServiceByID(db, id)
}

func UpdateOpenListService(db *gorm.DB, service *model.OpenListService) error {
	return model.UpdateOpenListService(db, service)
}

func DeleteOpenListService(db *gorm.DB, id int) error {
	return model.DeleteOpenListService(db, id)
} 
