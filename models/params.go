package models

// 定义请求的结构体
// binding 可以对传入的参数进行校验
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}