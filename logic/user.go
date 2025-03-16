package logic

import (
	"bulebell/dao/mysql"
	"bulebell/models"
	"bulebell/pkg/jwt"
	"bulebell/pkg/snowflake"
	"errors"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	exist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户已存在")
	}
	// 2. 生成 UID
	userId := snowflake.GenID()
	user := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(&user); err != nil {
		return "", err
	}
	// 生成 JWT 
	return jwt.GenToken(user.UserID, user.Username)
}
