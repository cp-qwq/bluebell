package controller

import (
	"bulebell/logic"
	"bulebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		// 如果参数有误，直接返回响应
		zap.S().Error("SignUpHandler ShouldBindJSON err", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(200, gin.H{
				"msg": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"msg": errs.Translate(trans),
		})
		return 
	}
	// 2. 业务处理	
	if err := logic.SignUp(&p); err != nil {
		c.JSON(200, gin.H{
			"msg": "注册失败",
		})
	}
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg" : "success", 
	})
}