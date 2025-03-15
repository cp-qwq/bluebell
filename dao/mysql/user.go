package mysql

import (
	"bulebell/models"
	"golang.org/x/crypto/bcrypt"
)

// CheckUserExist 检查
func CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertUser 像数据库插入一条新的用户记录
func InsertUser(user *models.User) error {
	// 对密码进行加密
	password, _ := GetPwd(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user (user_id, username,password) values(?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, password)
	return err
}

// GetPwd 给密码加密
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}

// ComparePwd 比对密码
// pwd1是已经加密的密码，pwd2是用户实际的明文密码
func ComparePwd(pwd1 string, pwd2 string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2))
	if err != nil {
		return false
	} else {
		return true
	}
}
