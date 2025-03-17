package controller

import (
	"bulebell/logic"
	"bulebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 ctx 里拿到当前发请求的userID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID

	// 2. 创建帖子
	if err := logic.CreatePost(&p); err != nil {
		zap.S().Error("logic.CreatePost() failed:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandle 获取帖子详细
func GetPostDetailHandle(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的 id）
	strID := c.Param("id")
	pid, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据 id 取出帖子的数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {

	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版的帖子列表接口
// 根据前端传来的参数动态获取帖子列表（分数，创建时间）

// 1. 获取参数
// 2. 去 redis 查询 id列表
// 3. 根据 id 去数据库查询帖子详细信息

func GetPostListHandler2(c *gin.Context) {
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	
	// 获取数据
	data, err := logic.GetPostListNew(&p)
	if err != nil {
		zap.L().Error("logic.GetPostList2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}
