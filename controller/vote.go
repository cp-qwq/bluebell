package controller

import (
	"bulebell/logic"
	"bulebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func VoteHandler(c *gin.Context) {
	
	// 参数校验，给哪个文章投什么票
	var p models.ParamVoteData
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("BindJSON failed", zap.String("postId", p.PostID), zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return 
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return 
	}
	// 获取当前请求的用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	if err := logic.VoteForPost(userID, &p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
