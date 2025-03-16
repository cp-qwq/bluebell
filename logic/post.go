package logic

import (
	"bulebell/dao/mysql"
	"bulebell/models"
	"bulebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	// 1. 生成 post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	return mysql.CreatePost(p)
}

func GetPostById(pid int64) (*models.ApiPostDetail, error) {
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
				zap.Int64("pid", pid),
				zap.Error(err))
		return nil, err
	}
	// 根据作者id 查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
		return nil, err
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return nil, err
	}
	rsp := &models.ApiPostDetail{
		AuthorName: user.Username, 
		Post : post,
		CommunityDetail : community,
	}
	return rsp, nil
}

func GetPostList(page, size int64) ([]*models.ApiPostDetail, error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data := make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return data, nil
}