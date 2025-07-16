package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tnnevol/openlist-strm/backend-api/internal/app"
	"github.com/tnnevol/openlist-strm/backend-api/internal/model"
)

func main() {
	// 定义命令行参数
	var (
		action   = flag.String("action", "", "操作类型: getUser, addTestData")
		userID   = flag.Int("userid", 0, "用户ID")
		help     = flag.Bool("help", false, "显示帮助信息")
		count    = flag.Int("count", 10, "要创建的StrmConfig测试数据数量（仅createStrmConfigTestData用）")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 初始化数据库连接
	db, err := app.InitDatabase()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	switch *action {
	case "getUser":
		if *userID <= 0 {
			log.Fatal("请提供有效的用户ID")
		}
		getUserByID(db, *userID)
	case "addTestData":
		if *userID <= 0 {
			log.Fatal("请提供有效的用户ID (-userid)")
		}
		addTestDataToOpenListService(db, *userID)
	case "createStrmConfigTestData":
		if *userID <= 0 {
			log.Fatal("请提供有效的用户ID (-userid)")
		}
		createStrmConfigTestData(db, *userID, *count)
	default:
		log.Fatal("请指定有效的操作类型，使用 -help 查看帮助")
	}
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println(`
用户管理脚本使用说明:

用法:
  go run scripts/user_management.go -action <操作类型> [参数]

操作类型:
  getUser     - 通过用户ID获取用户信息
  addTestData - 为openlist_service表添加测试数据

参数:
  -userid <用户ID>  - 指定用户ID (getUser和addTestData操作都需要)

示例:
  # 获取用户ID为1的用户信息
  go run scripts/user_management.go -action getUser -userid 1

  # 为用户ID为1添加openlist_service测试数据
  go run scripts/user_management.go -action addTestData -userid 1

  # 显示帮助信息
  go run scripts/user_management.go -help
`)
}

// getUserByID 通过用户ID获取用户信息
func getUserByID(db *sql.DB, userID int) {
	user, err := model.GetUserByID(db, userID)
	if err != nil {
		log.Fatalf("获取用户信息失败: %v", err)
	}

	fmt.Printf("\n=== 用户信息 (ID: %d) ===\n", user.ID)
	fmt.Printf("用户名: %s\n", getStringValue(user.Username))
	fmt.Printf("邮箱: %s\n", getStringValue(user.Email))
	fmt.Printf("是否激活: %t\n", user.IsActive)
	fmt.Printf("验证码: %s\n", getStringValue(user.Code))
	fmt.Printf("验证码过期时间: %s\n", getTimeValue(user.CodeExpireAt))
	fmt.Printf("登录失败次数: %d\n", user.FailedLoginCount)
	fmt.Printf("锁定时间: %s\n", user.LockedUntil.Format("2006-01-02 15:04:05"))
	fmt.Printf("创建时间: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Token失效时间: %s\n", getTimeValue(user.TokenInvalidBefore))
	fmt.Println("========================\n")
}

// addTestDataToOpenListService 为openlist_service表添加测试数据
func addTestDataToOpenListService(db *sql.DB, userID int) {
	// 检查指定用户是否存在
	user, err := model.GetUserByID(db, userID)
	if err != nil || user == nil {
		log.Fatalf("未找到用户ID为%d的用户", userID)
	}

	// 批量生成50条测试数据
	testServices := make([]*model.OpenListService, 0, 50)
	for i := 1; i <= 50; i++ {
		service := &model.OpenListService{
			Name: fmt.Sprintf("测试服务%d", i),
			Account: fmt.Sprintf("test_account_%d", i),
			Token: fmt.Sprintf("test_token_%06d", i),
			ServiceUrl: fmt.Sprintf("https://api.example%d.com", i),
			BackupUrl: fmt.Sprintf("https://backup.example%d.com", i),
			Enabled: i%2 == 1, // 奇数启用，偶数禁用
			UserID: userID,
		}
		testServices = append(testServices, service)
	}

	fmt.Printf("开始为用户ID %d 批量添加50条测试数据到openlist_service表...\n", userID)

	for i, service := range testServices {
		err := model.CreateOpenListService(db, service)
		if err != nil {
			log.Printf("添加测试服务 %d 失败: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ 成功添加测试服务: %s (用户ID: %d)\n", service.Name, service.UserID)
	}

	fmt.Println("\n测试数据添加完成！")
	
	// 显示添加的服务
	services, err := model.GetOpenListServicesByUserID(db, userID)
	if err != nil {
		log.Printf("获取服务列表失败: %v", err)
		return
	}

	fmt.Printf("\n=== 用户 %d 的服务列表（前10条） ===\n", userID)
	for i, service := range services {
		if i >= 10 {
			fmt.Printf("... 共%d条，已省略 ...\n", len(services))
			break
		}
		fmt.Printf("ID: %d, 服务名: %s, 账户: %s, 启用状态: %t\n", 
			service.ID, service.Name, service.Account, service.Enabled)
	}
	fmt.Println("========================\n")
}

// createStrmConfigTestData 为指定用户创建StrmConfig测试数据
func createStrmConfigTestData(db *sql.DB, userID, count int) {
	// 检查指定用户是否存在
	user, err := model.GetUserByID(db, userID)
	if err != nil || user == nil {
		log.Fatalf("未找到用户ID为%d的用户", userID)
	}

	// 获取任意一个 serviceId
	var serviceID int
	err = db.QueryRow("SELECT id FROM openlist_service WHERE user_id = ? LIMIT 1", userID).Scan(&serviceID)
	if err != nil {
		log.Fatalf("未找到该用户的 openlist_service 记录，请先添加服务")
	}

	testConfigs := make([]*model.StrmConfig, 0, count)
	for i := 1; i <= count; i++ {
		cfg := &model.StrmConfig{
			UserID: userID,
			Name:   fmt.Sprintf("测试配置%d", i),
			AlistBasePath: fmt.Sprintf("/test/path/%d", i),
			StrmOutputPath: fmt.Sprintf("/output/strm/%d", i),
			DownloadEnabled: i%2 == 1,
			DownloadInterval: 3600 + i,
			UpdateMode: model.UpdateModeIncremental,
			ServiceID: serviceID,
			IsUseBackupUrl: true,
		}
		testConfigs = append(testConfigs, cfg)
	}

	fmt.Printf("开始为用户ID %d 批量添加%d条StrmConfig测试数据...\n", userID, count)
	for i, cfg := range testConfigs {
		err := model.CreateStrmConfig(db, cfg)
		if err != nil {
			log.Printf("添加测试配置 %d 失败: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ 成功添加测试配置: %s (用户ID: %d)\n", cfg.Name, cfg.UserID)
	}
	fmt.Println("\n测试数据添加完成！")

	// 显示前10条
	rows, err := db.Query("SELECT id, name, alist_base_path, strm_output_path, download_enabled, download_interval, update_mode, service_id, is_use_backup_url, created_at, updated_at FROM strm_config WHERE user_id = ? ORDER BY created_at DESC LIMIT 10", userID)
	if err != nil {
		log.Printf("获取配置列表失败: %v", err)
		return
	}
	defer rows.Close()
	fmt.Printf("\n=== 用户 %d 的StrmConfig列表（前10条） ===\n", userID)
	for rows.Next() {
		var c model.StrmConfig
		err := rows.Scan(&c.ID, &c.Name, &c.AlistBasePath, &c.StrmOutputPath, &c.DownloadEnabled, &c.DownloadInterval, &c.UpdateMode, &c.ServiceID, &c.IsUseBackupUrl, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			log.Printf("解析配置失败: %v", err)
			continue
		}
		fmt.Printf("ID: %d, 名称: %s, 路径: %s, 启用: %t\n", c.ID, c.Name, c.AlistBasePath, c.DownloadEnabled)
	}
	fmt.Println("========================\n")
}

// getAllUsers 获取所有用户
func getAllUsers(db *sql.DB) ([]*model.User, error) {
	rows, err := db.Query("SELECT id, email, password_hash, is_active, code, code_expire_at, failed_login_count, locked_until, created_at, username, token_invalid_before FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		var lockedUntil sql.NullTime
		err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.Code, &u.CodeExpireAt, &u.FailedLoginCount, &lockedUntil, &u.CreatedAt, &u.Username, &u.TokenInvalidBefore)
		if err != nil {
			return nil, err
		}
		if lockedUntil.Valid {
			u.LockedUntil = lockedUntil.Time
		} else {
			u.LockedUntil = time.Time{}
		}
		users = append(users, &u)
	}

	return users, nil
}

// getStringValue 安全获取字符串值
func getStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "未设置"
}

// getTimeValue 安全获取时间值
func getTimeValue(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02 15:04:05")
	}
	return "未设置"
} 
