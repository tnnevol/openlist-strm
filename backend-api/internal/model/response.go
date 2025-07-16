package model

// Response 统一响应结构
// swagger:model
type Response[T any] struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	// data 可为任意类型，如 PageResult、StrmConfig、OpenListServiceResponse 等
	Data    T `json:"data,omitempty"` // 泛型，具体类型见各接口注释
}


// PageResult 通用分页响应结构
// swagger:model
// 用于所有分页接口的响应
// list: 当前页数据，total: 总条数，page: 当前页码，pageSize: 每页条数
type PageResult[T any] struct {
	List     []T `json:"list"` // 当前页数据
	Total    int         `json:"total" example:"100"` // 总条数
	Page     int         `json:"page" example:"1"` // 当前页码
	PageSize int         `json:"pageSize" example:"10"` // 每页条数
}

// UserInfoResponse 用户信息接口响应结构
// swagger:model
type UserInfoResponse struct {
	Roles    []string `json:"roles" example:"[\"admin\"]"`
	RealName string   `json:"realName" example:"Super"`
	ID       int      `json:"id" example:"1"`
	Username string   `json:"username" example:"zhangsan"`
	Email    string   `json:"email" example:"user@example.com"`
}

type DownloadEnabled string

const (
	DownloadEnabledTrue DownloadEnabled = "1"
	DownloadEnabledFalse DownloadEnabled = "0"
)

type IsUseBackupUrl string

const (
	IsUseBackupUrlTrue IsUseBackupUrl = "1"
	IsUseBackupUrlFalse IsUseBackupUrl = "0"
)

// 新增用于接口返回的小驼峰结构体
// swagger:model
//
type StrmConfigResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	AlistBasePath    string    `json:"alistBasePath"`
	StrmOutputPath   string    `json:"strmOutputPath"`
	DownloadEnabled  DownloadEnabled      `json:"downloadEnabled"`
	DownloadInterval int       `json:"downloadInterval"`
	UpdateMode       string    `json:"updateMode"`
	ServiceID        int       `json:"serviceId"`
	IsUseBackupUrl   IsUseBackupUrl      `json:"isUseBackupUrl"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
} 


type Enabled string

const (
	EnabledTrue Enabled = "1"
	EnabledFalse Enabled = "0"
)

// OpenListServiceResponse 响应结构体，使用小驼峰格式，去除 CreatedAt 和 UserID 字段
// swagger:model
//
type OpenListServiceResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Account   string `json:"account"`
	Token     string `json:"token"`
	ServiceUrl string `json:"serviceUrl"`
	BackupUrl  string `json:"backupUrl"`
	Enabled    Enabled   `json:"enabled"`
	UpdatedAt  string `json:"updatedAt"`
}
