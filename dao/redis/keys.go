package redis

// redis key

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // ZSet; 帖子以及发帖时间
	KeyPostScoreZSet   = "post:score"  // ZSet; 帖子以及投票的分数
	KeyPostVotedZSetPF = "post:voted:" // ZSet; 记录用户和投票的类型; 参数是 post_id
	KeyCommunitySetPF = "community:" // set; 保存分区下帖子的所有 id
)

func getRedisKey(key string) string {
	return Prefix + key
}
