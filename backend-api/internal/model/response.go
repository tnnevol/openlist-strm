package model

// Response 统一响应结构
// swagger:model
type Response struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
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
