package router

import (
	"bulebell/controller"
	"bulebell/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	
	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	
	// 登录业务路由
	v1.POST("/login", controller.LoginHandler)
	
	v1.Use(middlewares.JWTAuthMiddleware())
	
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler) // 路径参数
		
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandle)
		v1.GET("/posts", controller.GetPostListHandler)
	}
	return r
}
