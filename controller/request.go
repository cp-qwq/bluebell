package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)
var ErrorUserNotLogin = errors.New("用户未登录")
const CtxUserIdKey = "userID"
// GetCurrentUser 获取当前登录用户 ID
func GetCurrentUserID(c *gin.Context) (int64, error) {
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

func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}