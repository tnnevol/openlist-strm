package service

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/logger"
	"go.uber.org/zap"
)

// TokenBlacklist 管理已失效的token
type TokenBlacklist struct {
	blacklist map[string]time.Time
	mutex     sync.RWMutex
}

var (
	blacklist *TokenBlacklist
	once      sync.Once
)

// GetTokenBlacklist 获取token黑名单实例（单例模式）
func GetTokenBlacklist() *TokenBlacklist {
	once.Do(func() {
		blacklist = &TokenBlacklist{
			blacklist: make(map[string]time.Time),
		}
		// 启动清理过期token的goroutine
		go blacklist.cleanupExpiredTokens()
	})
	return blacklist
}

// AddToBlacklist 将token添加到黑名单
func (tb *TokenBlacklist) AddToBlacklist(token string, expireTime time.Time) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// 使用SHA256哈希token，避免存储明文token
	tokenHash := hashToken(token)
	tb.blacklist[tokenHash] = expireTime

	logger.Info("[TokenBlacklist] token已添加到黑名单", 
		zap.String("tokenHash", tokenHash[:8]+"..."),
		zap.Time("expireTime", expireTime))
}

// IsBlacklisted 检查token是否在黑名单中
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()

	tokenHash := hashToken(token)
	expireTime, exists := tb.blacklist[tokenHash]
	
	if !exists {
		return false
	}

	// 检查token是否已过期
	if time.Now().After(expireTime) {
		// token已过期，从黑名单中移除
		tb.mutex.RUnlock()
		tb.mutex.Lock()
		delete(tb.blacklist, tokenHash)
		tb.mutex.Unlock()
		tb.mutex.RLock()
		return false
	}

	return true
}

// cleanupExpiredTokens 定期清理过期的token
func (tb *TokenBlacklist) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时清理一次
	defer ticker.Stop()

	for range ticker.C {
		tb.mutex.Lock()
		now := time.Now()
		count := 0
		
		for tokenHash, expireTime := range tb.blacklist {
			if now.After(expireTime) {
				delete(tb.blacklist, tokenHash)
				count++
			}
		}
		
		tb.mutex.Unlock()
		
		if count > 0 {
			logger.Info("[TokenBlacklist] 清理过期token", zap.Int("count", count))
		}
	}
}

// GetBlacklistSize 获取黑名单大小（用于监控）
func (tb *TokenBlacklist) GetBlacklistSize() int {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()
	return len(tb.blacklist)
}

// hashToken 对token进行SHA256哈希
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
} 
