package model

import (
	"database/sql"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

type User struct {
	ID               int
	Username         sql.NullString
	Email            sql.NullString
	PasswordHash     string
	IsActive         bool
	Code             sql.NullString
	CodeExpireAt     sql.NullTime
	FailedLoginCount int
	LockedUntil      time.Time
	CreatedAt        time.Time
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var u User
	var lockedUntil sql.NullTime
	err := db.QueryRow("SELECT id, email, password_hash, is_active, code, code_expire_at, failed_login_count, locked_until, created_at, username FROM user WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.Code, &u.CodeExpireAt, &u.FailedLoginCount, &lockedUntil, &u.CreatedAt, &u.Username)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		u.LockedUntil = lockedUntil.Time
	} else {
		u.LockedUntil = time.Time{}
	}
	return &u, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var u User
	var lockedUntil sql.NullTime
	err := db.QueryRow("SELECT id, email, password_hash, is_active, code, code_expire_at, failed_login_count, locked_until, created_at, username FROM user WHERE username = ?", username).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.Code, &u.CodeExpireAt, &u.FailedLoginCount, &lockedUntil, &u.CreatedAt, &u.Username)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		u.LockedUntil = lockedUntil.Time
	} else {
		u.LockedUntil = time.Time{}
	}
	return &u, nil
}

func CreateUser(db *sql.DB, u *User) error {
	_, err := db.Exec("INSERT INTO user(username, email, password_hash, is_active, code, code_expire_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.Username, u.Email, u.PasswordHash, u.IsActive, u.Code, u.CodeExpireAt, u.CreatedAt)
	return err
}

func ActivateUser(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE user SET is_active = 1, code = '', code_expire_at = NULL WHERE id = ?", id)
	return err
}

// ActivateUserWithPassword 激活用户并设置密码
func ActivateUserWithPassword(db *sql.DB, id int, passwordHash string) error {
	logger.Info("[DB] ActivateUserWithPassword", zap.Int("user_id", id), zap.String("hash", passwordHash))
	res, err := db.Exec("UPDATE user SET is_active = 1, password_hash = ?, code = '', code_expire_at = NULL WHERE id = ?", passwordHash, id)
	if err != nil {
		logger.Error("[DB] ActivateUserWithPassword update error", zap.Error(err))
		return err
	}
	n, _ := res.RowsAffected()
	logger.Info("[DB] ActivateUserWithPassword rows affected", zap.Int64("rows", n))
	return nil
}

// ActivateUserWithPasswordAndUsername 激活用户并设置密码和用户名
func ActivateUserWithPasswordAndUsername(db *sql.DB, id int, passwordHash string, username string) error {
	logger.Info("[DB] ActivateUserWithPasswordAndUsername", zap.Int("user_id", id), zap.String("hash", passwordHash), zap.String("username", username))
	res, err := db.Exec("UPDATE user SET is_active = 1, password_hash = ?, username = ?, code = '', code_expire_at = NULL WHERE id = ?", passwordHash, username, id)
	if err != nil {
		logger.Error("[DB] ActivateUserWithPasswordAndUsername update error", zap.Error(err))
		return err
	}
	n, _ := res.RowsAffected()
	logger.Info("[DB] ActivateUserWithPasswordAndUsername rows affected", zap.Int64("rows", n))
	return nil
}

func UpdateLoginFail(db *sql.DB, id int, failCount int, lockedUntil time.Time) error {
	_, err := db.Exec("UPDATE user SET failed_login_count = ?, locked_until = ? WHERE id = ?", failCount, lockedUntil, id)
	return err
}

func ResetLoginFail(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE user SET failed_login_count = 0, locked_until = NULL WHERE id = ?", id)
	return err
}

// GetUserByCode 通过验证码获取用户
func GetUserByCode(db *sql.DB, code string) (*User, error) {
	var u User
	var lockedUntil sql.NullTime
	err := db.QueryRow("SELECT id, email, password_hash, is_active, code, code_expire_at, failed_login_count, locked_until, created_at, username FROM user WHERE code = ? AND is_active = 0", code).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.Code, &u.CodeExpireAt, &u.FailedLoginCount, &lockedUntil, &u.CreatedAt, &u.Username)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		u.LockedUntil = lockedUntil.Time
	} else {
		u.LockedUntil = time.Time{}
	}
	return &u, nil
}

// UpdateCode 更新验证码
func UpdateCode(db *sql.DB, email, code string, expireAt time.Time) error {
	_, err := db.Exec("UPDATE user SET code = ?, code_expire_at = ? WHERE email = ?", code, expireAt, email)
	return err
}

// GetActiveUserByCode 通过验证码查找已激活用户
func GetActiveUserByCode(db *sql.DB, code string) (*User, error) {
	var u User
	var lockedUntil sql.NullTime
	err := db.QueryRow("SELECT id, email, password_hash, is_active, code, code_expire_at, failed_login_count, locked_until, created_at, username FROM user WHERE code = ? AND is_active = 1", code).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.Code, &u.CodeExpireAt, &u.FailedLoginCount, &lockedUntil, &u.CreatedAt, &u.Username)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		u.LockedUntil = lockedUntil.Time
	} else {
		u.LockedUntil = time.Time{}
	}
	return &u, nil
}
