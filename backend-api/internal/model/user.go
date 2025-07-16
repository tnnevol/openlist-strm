package model

import (
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Username         string     `json:"username" gorm:"type:varchar(64);uniqueIndex"`
	Email            string     `json:"email" gorm:"type:varchar(128);uniqueIndex"`
	PasswordHash     string     `json:"passwordHash"`
	IsActive         bool       `json:"isActive"`
	Code             *string    `json:"code" gorm:"type:varchar(32)"`
	CodeExpireAt     *time.Time `json:"codeExpireAt"`
	FailedLoginCount int        `json:"failedLoginCount"`
	LockedUntil      *time.Time `json:"lockedUntil"`
	CreatedAt        time.Time  `json:"createdAt"`
	TokenInvalidBefore *time.Time `json:"tokenInvalidBefore"`
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var u User
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var u User
	if err := db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

func ActivateUser(db *gorm.DB, id int) error {
	return db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": true,
		"code":      "",
		"code_expire_at": nil,
	}).Error
}

// ActivateUserWithPassword 激活用户并设置密码
func ActivateUserWithPassword(db *gorm.DB, id int, passwordHash string) error {
	logger.Info("[DB] ActivateUserWithPassword", zap.Int("user_id", id), zap.String("hash", passwordHash))
	err := db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": true,
		"password_hash": passwordHash,
		"code":      "",
		"code_expire_at": nil,
		"token_invalid_before": nil,
	}).Error
	if err != nil {
		logger.Error("[DB] ActivateUserWithPassword update error", zap.Error(err))
		return err
	}
	return nil
}

// ActivateUserWithPasswordAndUsername 激活用户并设置密码和用户名
func ActivateUserWithPasswordAndUsername(db *gorm.DB, id int, passwordHash string, username string) error {
	logger.Info("[DB] ActivateUserWithPasswordAndUsername", zap.Int("user_id", id), zap.String("hash", passwordHash), zap.String("username", username))
	err := db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": true,
		"password_hash": passwordHash,
		"username": username,
		"code":      "",
		"code_expire_at": nil,
		"token_invalid_before": nil,
	}).Error
	if err != nil {
		logger.Error("[DB] ActivateUserWithPasswordAndUsername update error", zap.Error(err))
		return err
	}
	return nil
}

func UpdateLoginFail(db *gorm.DB, id int, failCount int, lockedUntil *time.Time) error {
	return db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"failed_login_count": failCount,
		"locked_until":       lockedUntil,
	}).Error
}

func ResetLoginFail(db *gorm.DB, id int) error {
	return db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"failed_login_count": 0,
		"locked_until":       nil,
	}).Error
}

// GetUserByCode 通过验证码获取用户
func GetUserByCode(db *gorm.DB, code string) (*User, error) {
	var u User
	if err := db.Where("code = ? AND is_active = 0", code).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateCode 更新验证码
func UpdateCode(db *gorm.DB, email, code string, expireAt time.Time) error {
	return db.Model(&User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"code": code,
		"code_expire_at": expireAt,
	}).Error
}

// GetActiveUserByCode 通过验证码查找已激活用户
func GetActiveUserByCode(db *gorm.DB, code string) (*User, error) {
	var u User
	if err := db.Where("code = ? AND is_active = 1", code).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateTokenInvalidBefore 更新token_invalid_before字段
func UpdateTokenInvalidBefore(db *gorm.DB, userID int, t time.Time) error {
	return db.Model(&User{}).Where("id = ?", userID).Update("token_invalid_before", t).Error
}

// GetUserByID 通过用户ID获取用户信息
func GetUserByID(db *gorm.DB, id int) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
