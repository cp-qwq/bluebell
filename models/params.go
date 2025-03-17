package models

// 定义请求的结构体
// binding 可以对传入的参数进行校验
const (
	OrderTime = "time"
	OrderScore = "score"
)
// 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`                // 帖子 id
	Direction *int8  `json:"direction" binding:"required,oneof=1 0 -1"` // 赞成(1)还是反对(-1)取消投票(0)
}

// ParamPostList 获取帖子列表 query 参数
type ParamPostList struct {
	CommunityID int64 `json:"community_id" form:"community_id"`
	Page  int64  `form:"page"`
	Size  int64  `form:"size"`
	Order string `form:"order"`
}

