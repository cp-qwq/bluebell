package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 雪花算法生成全局唯一的 ID
var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	if startTime == "" {
		startTime = time.Now().Format("2006-01-02") // 如果没有提供 startTime, 使用当前日期
	}
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}