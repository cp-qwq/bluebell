package controller

import (
	"bulebell/logic"
	"bulebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		// 如果参数有误，直接返回响应
		zap.S().Error("SignUpHandler ShouldBindJSON err", zap.Error(err))
		c.JSON(200, gin.H{
			"msg" : "请求参数有误",
		})
		return 
	}
	
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
		// 如果参数有误，直接返回响应
		c.JSON(200, gin.H{
			"msg" : "请求参数有误",
		})
		return 
	}
	
	// 2. 业务处理	
	logic.SignUp(&p)
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg" : "success",
	})
}