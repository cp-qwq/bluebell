package mysql

func CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {}
	return true;
}

func InsertUser() {
	// 与数据库交互
}