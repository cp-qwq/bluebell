package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("投票时间已过")
)

// 本项目用简化版的投票分数
// 投一票就加432分 一天是86400秒，86400/200 -> 432 表示200张赞成票可以给你的帖子续一天

/*
v = 1时，有两种情况：
	1. 之前没有投过票，现在投赞成票  + 432
	2. 之前投反对票，现在改投赞成票  + 432 * 2
v = 0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票 - 432
	2. 之前投过反对票，现在要取消投票 + 432
v = -1时，有两种情况：
	1. 之前没有投过票，现在投反对票  - 432
	2. 之前投赞成票，现在投反对票    - 432 * 2

限制：
每个帖子自发表之日起，一个星期内允许用户投票，超过一个星期就不允许投票。
	1. 到期之后将redis中保存的赞成票和反对票存储到mysql中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子发布时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})
	
	// 帖子分数（以时间戳为初始分数，投一票续432）
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})
	
	// 把帖子 id 加到社区的 set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) (err error) {
	// 1. 判断投票限制
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	
	// 2. 更新帖子的分数
	oldValue := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == oldValue {
		return ErrVoteRepested
	}
	// 只要现在的value > oldValue, 就一定是正值
	var op float64
	if value > oldValue {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(oldValue - value)
	pipline := client.TxPipeline() // 开启事务
	pipline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3. 记录用户为该帖子投票数据
	if value == 0 {
		pipline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
	} else {
		pipline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err = pipline.Exec()
	return err
}
