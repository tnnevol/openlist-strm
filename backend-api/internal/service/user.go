package service

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"strings"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
	"github.com/tnnevol/openlist-strm/backend-api/internal/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// CheckEmailExists 检查邮箱是否已注册
func CheckEmailExists(db *sql.DB, email string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT COUNT(1) FROM user WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// SendCode 发送验证码（用于未注册用户或未激活用户）
func SendCode(db *sql.DB, email string) (string, error) {
	logger.Info("[SendCode] called", zap.String("email", email))
	code := util.GenerateVerificationCode()
	expire := time.Now().Add(10 * time.Minute)
	exists, err := CheckEmailExists(db, email)
	if err != nil {
		logger.Error("[SendCode] CheckEmailExists error", zap.Error(err))
		return "", err
	}
	if exists {
		user, err := model.GetUserByEmail(db, email)
		if err != nil {
			logger.Error("[SendCode] GetUserByEmail error", zap.Error(err))
			return "", err
		}
		if user.IsActive {
			logger.Info("[SendCode] user already activated", zap.String("email", email))
			return "", errors.New("用户已注册")
		}
		logger.Info("[SendCode] email exists, update code", zap.String("email", email))
		err = model.UpdateCode(db, email, code, expire)
		if err != nil {
			logger.Error("[SendCode] UpdateCode error", zap.Error(err))
			return "", err
		}
		_ = util.SendMail(email, code)
		return code, nil
	}
	logger.Info("[SendCode] email not exists, create user", zap.String("email", email))
	u := &model.User{
		Email:        email,
		PasswordHash: "",
		IsActive:     false,
		Code:         code,
		CodeExpireAt: sql.NullTime{Time: expire, Valid: true},
		CreatedAt:    time.Now(),
	}
	err = model.CreateUser(db, u)
	if err != nil {
		logger.Error("[SendCode] CreateUser error", zap.Error(err))
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			logger.Info("[SendCode] UNIQUE constraint, update code", zap.String("email", email))
			err = model.UpdateCode(db, email, code, expire)
			if err != nil {
				logger.Error("[SendCode] UpdateCode error", zap.Error(err))
				return "", err
			}
			_ = util.SendMail(email, code)
			return code, nil
		}
		return "", err
	}
	_ = util.SendMail(email, code)
	return code, nil
}

// ActivateUserWithPassword 使用密码和验证码激活用户
func ActivateUserWithPassword(db *sql.DB, email, password, code string) error {
	logger.Info("[Register] called", zap.String("email", email), zap.String("code", code))
	u, err := model.GetUserByCode(db, code)
	if err != nil {
		logger.Error("[Register] GetUserByCode error", zap.Error(err))
		return err
	}
	if u.Email != email {
		logger.Error("[Register] code/email mismatch", zap.String("db_email", u.Email), zap.String("req_email", email))
		return sql.ErrNoRows
	}
	if !u.CodeExpireAt.Valid || time.Now().After(u.CodeExpireAt.Time) {
		logger.Error("[Register] code expired", zap.Any("expire_at", u.CodeExpireAt))
		return sql.ErrNoRows
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("[Register] bcrypt error", zap.Error(err))
		return err
	}
	logger.Info("[Register] activate user", zap.Int("user_id", u.ID))
	return model.ActivateUserWithPassword(db, u.ID, string(hash))
}

func LoginUser(db *sql.DB, email, password string) (string, error) {
	u, err := model.GetUserByEmail(db, email)
	if err != nil {
		logger.Error("[Login] GetUserByEmail error", zap.Error(err))
		return "", err
	}
	logger.Info("[Login] user info", zap.String("email", u.Email), zap.String("hash", u.PasswordHash), zap.Bool("is_active", u.IsActive), zap.Int("id", u.ID))
	if !u.IsActive {
		return "not_activated", nil
	}
	if u.LockedUntil.After(time.Now()) {
		return "locked", nil
	}
	if u.PasswordHash == "" {
		logger.Error("[Login] password hash is empty", zap.String("email", u.Email))
		return "", errors.New("邮箱或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		fail := u.FailedLoginCount + 1
		lock := u.LockedUntil
		if fail >= 5 {
			lock = time.Now().Add(10 * time.Minute)
		}
		model.UpdateLoginFail(db, u.ID, fail, lock)
		logger.Error("[Login] password not match", zap.Error(err))
		return "", err
	}
	model.ResetLoginFail(db, u.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		jwtKey = []byte("secret")
	}
	tokenString, _ := token.SignedString(jwtKey)
	logger.Info("[Login] success", zap.String("email", u.Email))
	return tokenString, nil
}

// ForgotPasswordSendCode 仅允许已激活用户获取验证码
func ForgotPasswordSendCode(db *sql.DB, email string) (string, error) {
	logger.Info("[ForgotPasswordSendCode] called", zap.String("email", email))
	exists, err := CheckEmailExists(db, email)
	if err != nil {
		logger.Error("[ForgotPasswordSendCode] CheckEmailExists error", zap.Error(err))
		return "", err
	}
	if !exists {
		return "", errors.New("用户未注册")
	}
	user, err := model.GetUserByEmail(db, email)
	if err != nil {
		logger.Error("[ForgotPasswordSendCode] GetUserByEmail error", zap.Error(err))
		return "", err
	}
	if !user.IsActive {
		return "", errors.New("账户未激活")
	}
	code := util.GenerateVerificationCode()
	expire := time.Now().Add(10 * time.Minute)
	err = model.UpdateCode(db, email, code, expire)
	if err != nil {
		logger.Error("[ForgotPasswordSendCode] UpdateCode error", zap.Error(err))
		return "", err
	}
	_ = util.SendMail(email, code)
	return code, nil
}

// ForgotPasswordReset 校验验证码并重置密码
func ForgotPasswordReset(db *sql.DB, email, code, newPassword string) error {
	logger.Info("[ForgotPasswordReset] called", zap.String("email", email), zap.String("code", code))
	user, err := model.GetActiveUserByCode(db, code)
	if err != nil {
		logger.Error("[ForgotPasswordReset] GetActiveUserByCode error", zap.Error(err))
		return errors.New("验证码无效或已过期")
	}
	if user.Email != email {
		logger.Error("[ForgotPasswordReset] code/email mismatch", zap.String("db_email", user.Email), zap.String("req_email", email))
		return errors.New("验证码无效或已过期")
	}
	if !user.CodeExpireAt.Valid || time.Now().After(user.CodeExpireAt.Time) {
		logger.Error("[ForgotPasswordReset] code expired", zap.Any("expire_at", user.CodeExpireAt))
		return errors.New("验证码无效或已过期")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("[ForgotPasswordReset] bcrypt error", zap.Error(err))
		return err
	}
	logger.Info("[ForgotPasswordReset] reset password", zap.Int("user_id", user.ID))
	return model.ActivateUserWithPassword(db, user.ID, string(hash))
}
