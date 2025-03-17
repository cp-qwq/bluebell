package logic

import (
	"bulebell/dao/mysql"
	"bulebell/dao/redis"
	"bulebell/models"
	"bulebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	// 1. 生成 post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return err
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
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
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

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去redis里面查寻 id 列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, err
	}
	// 3. 根据ids 去数据库里面查询帖子的详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}

	// 4. 查询每篇帖子的投票数据
	voteData, err := redis.GetPostVoteDate(ids)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return data, err
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	
	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteDate(ids)
	if err != nil {
		return
	}


	// 将帖子的作者及分区信息查询出来填充到帖子中 
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return 
}


// GetPostListNew 将两个查询接口合二为一
func GetPostListNew(p * models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p)
	} else {
		// 根据社区 ID 查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return 
}