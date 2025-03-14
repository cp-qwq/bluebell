package main

import (
	"bulebell/controller"
	"bulebell/dao/mysql"
	"bulebell/logger"
	"bulebell/pkg/snowflake"
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
	
	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return 
	}
	
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return 
	}
	
	// 5. 初始化 Router
	
	r := router.SetUpRouter()
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}