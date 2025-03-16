package logic

import (
	"bulebell/dao/mysql"
	"bulebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找到所有的community并且返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailById(id)
}