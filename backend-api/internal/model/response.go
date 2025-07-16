package model

// Response 统一响应结构
// swagger:model
type Response struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	// data 可为任意类型，如 PageResult、StrmConfig、OpenListServiceResponse 等
	Data    interface{} `json:"data,omitempty"` // 泛型，具体类型见各接口注释
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

// PageResult 通用分页响应结构
// swagger:model
// 用于所有分页接口的响应
// list: 当前页数据，total: 总条数，page: 当前页码，pageSize: 每页条数
// 泛型 T 仅作注释，实际用 interface{}
type PageResult struct {
	List     interface{} `json:"list"` // 当前页数据
	Total    int         `json:"total" example:"100"` // 总条数
	Page     int         `json:"page" example:"1"` // 当前页码
	PageSize int         `json:"pageSize" example:"10"` // 每页条数
}

// OpenListServicePageResult OpenListService 分页响应
// swagger:model
// 用于 OpenListService 分页接口
// list: OpenListServiceResponse 数组
// total/page/pageSize 同 PageResult
//
type OpenListServicePageResult struct {
	List     []OpenListServiceResponse `json:"list"`
	Total    int                      `json:"total" example:"100"`
	Page     int                      `json:"page" example:"1"`
	PageSize int                      `json:"pageSize" example:"10"`
}

// StrmConfigPageResult StrmConfig 分页响应
// swagger:model
// 用于 StrmConfig 分页接口
// list: StrmConfig 数组
// total/page/pageSize 同 PageResult
//
type StrmConfigPageResult struct {
	List     []StrmConfigResponse `json:"list"`
	Total    int          `json:"total" example:"100"`
	Page     int          `json:"page" example:"1"`
	PageSize int          `json:"pageSize" example:"10"`
}

// OpenListServiceResponse 响应结构体，使用小驼峰格式，去除 CreatedAt 和 UserID 字段
// swagger:model
// 用于 OpenListServicePageResult 的 list 元素
//
type OpenListServiceResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Account   string `json:"account"`
	Token     string `json:"token"`
	ServiceUrl string `json:"serviceUrl"`
	BackupUrl  string `json:"backupUrl"`
	Enabled    bool   `json:"enabled"`
	UpdatedAt  string `json:"updatedAt"`
}
