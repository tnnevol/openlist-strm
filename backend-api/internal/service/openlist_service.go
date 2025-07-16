package service

import (
	"database/sql"

	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
)

func CreateOpenListService(db *sql.DB, service *model.OpenListService) error {
	return model.CreateOpenListService(db, service)
}

func GetOpenListServicesByUserID(db *sql.DB, userID int) ([]*model.OpenListService, error) {
	return model.GetOpenListServicesByUserID(db, userID)
}

func GetOpenListServiceByID(db *sql.DB, id int) (*model.OpenListService, error) {
	return model.GetOpenListServiceByID(db, id)
}

func UpdateOpenListService(db *sql.DB, service *model.OpenListService) error {
	return model.UpdateOpenListService(db, service)
}

func DeleteOpenListService(db *sql.DB, id int) error {
	return model.DeleteOpenListService(db, id)
} 
