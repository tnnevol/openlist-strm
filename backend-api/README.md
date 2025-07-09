# Backend API (Root Level)

基于 Gin 框架的 Go 语言后端服务

## 开发

```bash
# 启动开发服务器
go run main.go
```

---

如遇已激活用户 password_hash 为空导致无法登录，可用如下 SQL 修复：

UPDATE user SET is_active=0 WHERE is_active=1 AND (password_hash IS NULL OR password_hash='');
