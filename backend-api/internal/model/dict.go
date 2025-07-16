package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

type Dict struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	ParentID    int       `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func CreateDict(db *sql.DB, d *Dict) error {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	logger.Info("[DB] CreateDict", zap.String("type", d.Type), zap.String("key", d.Key))
	_, err := db.Exec(`INSERT INTO dict (type, key, value, description, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Type, d.Key, d.Value, d.Description, d.ParentID, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		logger.Error("[DB] CreateDict error", zap.Error(err))
		return err
	}
	return nil
}

func UpdateDict(db *sql.DB, d *Dict) error {
	d.UpdatedAt = time.Now()
	logger.Info("[DB] UpdateDict", zap.Int("id", d.ID))
	_, err := db.Exec(`UPDATE dict SET type=?, key=?, value=?, description=?, parent_id=?, updated_at=? WHERE id=?`,
		d.Type, d.Key, d.Value, d.Description, d.ParentID, d.UpdatedAt, d.ID)
	if err != nil {
		logger.Error("[DB] UpdateDict error", zap.Error(err))
		return err
	}
	return nil
}

func DeleteDict(db *sql.DB, id int) error {
	logger.Info("[DB] DeleteDict", zap.Int("id", id))
	_, err := db.Exec(`DELETE FROM dict WHERE id=?`, id)
	if err != nil {
		logger.Error("[DB] DeleteDict error", zap.Error(err))
		return err
	}
	return nil
}

func GetDictByID(db *sql.DB, id int) (*Dict, error) {
	var d Dict
	err := db.QueryRow(`SELECT id, type, key, value, description, parent_id, created_at, updated_at FROM dict WHERE id=?`, id).
		Scan(&d.ID, &d.Type, &d.Key, &d.Value, &d.Description, &d.ParentID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func ListDicts(db *sql.DB, dictType string) ([]*Dict, error) {
	var rows *sql.Rows
	var err error
	if dictType != "" {
		rows, err = db.Query(`SELECT id, type, key, value, description, parent_id, created_at, updated_at FROM dict WHERE type=? ORDER BY id DESC`, dictType)
	} else {
		rows, err = db.Query(`SELECT id, type, key, value, description, parent_id, created_at, updated_at FROM dict ORDER BY id DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dicts []*Dict
	for rows.Next() {
		var d Dict
		err := rows.Scan(&d.ID, &d.Type, &d.Key, &d.Value, &d.Description, &d.ParentID, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dicts = append(dicts, &d)
	}
	return dicts, nil
} 
