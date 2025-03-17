package mysql

import (
	"bulebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	sqlStr := `insert into post
	(post_id, title, content, author_id, community_id)
	values(?, ?, ?, ?, ?)
	`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// GetPostById 根据 id 查询单个帖子的数据
func GetPostById(pid int64) (*models.Post, error) {
	post := new(models.Post)

	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	err := db.Get(post, sqlStr, pid)
	return post, err
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) ([]*models.Post, error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post limit ?, ?`

	posts := make([]*models.Post, 0, 2)
	err := db.Select(&posts, sqlStr, (page-1)*size, size)
	return posts, err
}

// GetPostListByIDs 根据给定的 id 列表查询帖子的数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}