package logic

import (
	"bulebell/dao/mysql"
	"bulebell/models"
	"errors"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if mysql.CheckUserExist(p.Username) {
		return errors.New("用户已存在")
	}
	// 2. 生成 UID
	// 3. 保存进数据库
}