package router

import (
	"bulebell/controller"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	return r
}