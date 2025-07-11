package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TestUser 测试用户结构
type TestUser struct {
	Name        string `yaml:"name"`
	Email       string `yaml:"email"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	IsActive    bool   `yaml:"is_active"`
	Description string `yaml:"description"`
}

// TestConfig 测试环境配置
type TestConfig struct {
	Database struct {
		Type string `yaml:"type"`
		DSN  string `yaml:"dsn"`
	} `yaml:"database"`
	JWT struct {
		Secret       string `yaml:"secret"`
		ExpiresHours int    `yaml:"expires_hours"`
	} `yaml:"jwt"`
	Email struct {
		Enabled   bool   `yaml:"enabled"`
		MockCode  string `yaml:"mock_code"`
	} `yaml:"email"`
}

// TestStep 测试步骤结构
type TestStep struct {
	Step     string                 `yaml:"step"`
	Endpoint string                 `yaml:"endpoint"`
	Method   string                 `yaml:"method"`
	Body     map[string]interface{} `yaml:"body,omitempty"`
	Headers  map[string]string      `yaml:"headers,omitempty"`
	Expected struct {
		Code    int    `yaml:"code"`
		Message string `yaml:"message,omitempty"`
		Contains string `yaml:"contains,omitempty"`
	} `yaml:"expected"`
}

// TestScenario 测试场景结构
type TestScenario struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Steps       []TestStep `yaml:"steps"`
}

// TestScenarios 测试场景配置
type TestScenarios struct {
	Registration []TestScenario `yaml:"registration"`
	Login        []TestScenario `yaml:"login"`
}

// UserGroups 用户组配置
type UserGroups struct {
	NormalUsers    []TestUser `yaml:"normal_users"`
	EdgeCaseUsers  []TestUser `yaml:"edge_case_users"`
	ErrorTestUsers []TestUser `yaml:"error_test_users"`
}

// CleanupConfig 清理配置
type CleanupConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Strategy  string   `yaml:"strategy"`
	DataTypes []string `yaml:"data_types"`
}

// TestConfigFile 完整的测试配置文件结构
type TestConfigFile struct {
	TestConfig    TestConfig     `yaml:"test_config"`
	UserGroups    UserGroups     `yaml:"user_groups"`
	TestScenarios TestScenarios  `yaml:"test_scenarios"`
	Cleanup       CleanupConfig  `yaml:"cleanup"`
}

// LoadTestConfig 加载测试配置文件
func LoadTestConfig(configPath string) (*TestConfigFile, error) {
	// 如果路径为空，使用默认路径
	if configPath == "" {
		configPath = "tests/fixtures/test_users.yml"
	}

	// 尝试多个可能的路径
	possiblePaths := []string{
		configPath,
		"../tests/fixtures/test_users.yml",
		"../../tests/fixtures/test_users.yml",
	}

	var data []byte
	var err error
	
	for _, path := range possiblePaths {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var config TestConfigFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %v", err)
	}

	return &config, nil
}

// GetTestUser 根据用户名获取测试用户
func (cfg *TestConfigFile) GetTestUser(username string) *TestUser {
	// 在正常用户组中查找
	for _, user := range cfg.UserGroups.NormalUsers {
		if user.Username == username {
			return &user
		}
	}

	// 在边界测试用户组中查找
	for _, user := range cfg.UserGroups.EdgeCaseUsers {
		if user.Username == username {
			return &user
		}
	}

	// 在错误测试用户组中查找
	for _, user := range cfg.UserGroups.ErrorTestUsers {
		if user.Username == username {
			return &user
		}
	}

	return nil
}

// GetTestUserByEmail 根据邮箱获取测试用户
func (cfg *TestConfigFile) GetTestUserByEmail(email string) *TestUser {
	// 在正常用户组中查找
	for _, user := range cfg.UserGroups.NormalUsers {
		if user.Email == email {
			return &user
		}
	}

	// 在边界测试用户组中查找
	for _, user := range cfg.UserGroups.EdgeCaseUsers {
		if user.Email == email {
			return &user
		}
	}

	// 在错误测试用户组中查找
	for _, user := range cfg.UserGroups.ErrorTestUsers {
		if user.Email == email {
			return &user
		}
	}

	return nil
}

// GetAllTestUsers 获取所有测试用户
func (cfg *TestConfigFile) GetAllTestUsers() []TestUser {
	var users []TestUser
	
	// 添加正常用户
	users = append(users, cfg.UserGroups.NormalUsers...)
	
	// 添加边界测试用户
	users = append(users, cfg.UserGroups.EdgeCaseUsers...)
	
	// 添加错误测试用户
	users = append(users, cfg.UserGroups.ErrorTestUsers...)
	
	return users
}

// GetActiveUsers 获取所有活跃用户
func (cfg *TestConfigFile) GetActiveUsers() []TestUser {
	var activeUsers []TestUser
	
	allUsers := cfg.GetAllTestUsers()
	for _, user := range allUsers {
		if user.IsActive {
			activeUsers = append(activeUsers, user)
		}
	}
	
	return activeUsers
}

// GetInactiveUsers 获取所有非活跃用户
func (cfg *TestConfigFile) GetInactiveUsers() []TestUser {
	var inactiveUsers []TestUser
	
	allUsers := cfg.GetAllTestUsers()
	for _, user := range allUsers {
		if !user.IsActive {
			inactiveUsers = append(inactiveUsers, user)
		}
	}
	
	return inactiveUsers
}

// GetTestScenario 根据场景名称获取测试场景
func (cfg *TestConfigFile) GetTestScenario(scenarioType, scenarioName string) *TestScenario {
	switch scenarioType {
	case "registration":
		for _, scenario := range cfg.TestScenarios.Registration {
			if scenario.Name == scenarioName {
				return &scenario
			}
		}
	case "login":
		for _, scenario := range cfg.TestScenarios.Login {
			if scenario.Name == scenarioName {
				return &scenario
			}
		}
	}
	
	return nil
}

// ValidateConfig 验证配置文件
func (cfg *TestConfigFile) ValidateConfig() error {
	// 只校验正常用户组
	for _, user := range cfg.UserGroups.NormalUsers {
		if user.Email == "" {
			return fmt.Errorf("正常用户 %s 缺少邮箱", user.Name)
		}
		if user.Username == "" {
			return fmt.Errorf("正常用户 %s 缺少用户名", user.Name)
		}
		if user.Password == "" {
			return fmt.Errorf("正常用户 %s 缺少密码", user.Name)
		}
	}

	// 边界和错误用户组只做警告
	for _, user := range cfg.UserGroups.EdgeCaseUsers {
		if user.Email == "" || user.Username == "" || user.Password == "" {
			fmt.Printf("[警告] 边界测试用户 %s 存在空字段，允许用于异常测试\n", user.Name)
		}
	}
	for _, user := range cfg.UserGroups.ErrorTestUsers {
		if user.Email == "" || user.Username == "" || user.Password == "" {
			fmt.Printf("[警告] 错误测试用户 %s 存在空字段，允许用于异常测试\n", user.Name)
		}
	}

	// 验证JWT配置
	if cfg.TestConfig.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// SaveConfig 保存配置到文件
func (cfg *TestConfigFile) SaveConfig(configPath string) error {
	// 如果路径为空，使用默认路径
	if configPath == "" {
		configPath = "tests/fixtures/test_users.yml"
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 序列化为YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化YAML失败: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
} 
