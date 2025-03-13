package main

import (
	"bulebell/dao/mysql"
	"bulebell/logger"
	"bulebell/router"
	"bulebell/settings"
	"fmt"
)
func main() {
	// 1. 加载配置
	settings.Init();
	// 2. 初始化日志
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// 3. 初始化 Mysql 连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化 Redis 连接
	//if err := redis.Init(); err != nil {
	//	fmt.Printf("init redis failed, err:%v\n", err)
	//	return
	//}
	//defer redis.Close()
	// 5. 初始化 Router
	r := router.SetUpRouter()
	r.Run("127.0.0.1:8081")
}