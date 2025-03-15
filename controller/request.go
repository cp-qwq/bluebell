package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)
var ErrorUserNotLogin = errors.New("用户未登录")
const CtxUserIdKey = "userID"
// GetCurrentUser 获取当前登录用户 ID
func GetCurrentUser(c *gin.Context) (int64, error) {
	uid, ok := c.Get(CtxUserIdKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	return userID, nil
}