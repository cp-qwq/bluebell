package logic

import (
	"bulebell/dao/redis"
	"bulebell/models"
	"go.uber.org/zap"
	"strconv"
)

// 投票功能

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userId",userID),
		zap.String("postId", p.PostID),
		zap.Int8p("Direction",p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(*p.Direction))
}
