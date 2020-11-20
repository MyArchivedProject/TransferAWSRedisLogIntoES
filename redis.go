package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// var rdb *redis.Client

// RedisSlowLog Redis慢日志
type RedisSlowLog struct {
	RedisID      string `type:"string"`
	RedisCluster string `type:"string"`
	RedisAddress string `type:"string"`
	RedisPort    int64  `type:"int"`
	ID           int64
	Time         time.Time
	Duration     time.Duration // 毫秒
	Args         []string
}

type redisConf struct {
	Addr     string
	Password string
	DB       int
}

// 功能：连接redis  失败输出：返回nil
func getRedisClient(redisConf redisConf) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,     // use default Addr
		Password: redisConf.Password, // no password set
		DB:       redisConf.DB,       // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		printLog("RDB连接失败 Addr=" + redisConf.Addr)
		errorTolerate(err)
		return nil
	} else {
		printLog("RDB连接成功 Addr=" + redisConf.Addr)
		return rdb
	}
}

// 功能：获取Redis慢日志  输入：(rdb, slowlog的行数)
func getSlowLog(rdb *redis.Client, num int64) []redis.SlowLog {
	// res, err := rdb.Do(ctx, "slowlog", "get", num).Result() // 返回接口
	res, err := rdb.SlowLogGet(ctx, num).Result() // 返回被格式化过的数组

	if err != nil {
		log.Println("Error getSlowLog():\n" + error.Error(err))
	}
	return res
}

func getRedisSlowLog() {

}

// GetMultiRedisSlowLog 输入：每个redis节点信息  成功输出：所有的慢日志  失败输出：nil
func GetMultiRedisSlowLog(redisNodeInfoArr []RedisNodeInfo) (redisSlowLogArr []RedisSlowLog) {
	nodeNum := len(redisNodeInfoArr)
	slowlogNum := viper.GetInt64("redis.slowlog_num")
	if nodeNum == 0 {
		printLog("输入的Redis节点个数为0")
		return nil
	}
	redisSlowLogArr = make([]RedisSlowLog, 0)

	// 遍历所有节点
	for i := 0; i < nodeNum; i++ {
		address := redisNodeInfoArr[i].RedisAddress
		port := redisNodeInfoArr[i].RedisPort
		clustName := redisNodeInfoArr[i].RedisCluster
		redisID := redisNodeInfoArr[i].RedisID
		redisConf := redisConf{
			Addr: address + ":" + strconv.FormatInt(port, 10),
		}

		// 连接Redis
		rdb := getRedisClient(redisConf)
		if rdb != nil {
			defer rdb.Close()
			// 遍历节点上所有的慢日志
			slowLogArr := getSlowLog(rdb, slowlogNum)
			for j := 0; j < len(slowLogArr); j++ {
				slowLog := &RedisSlowLog{
					RedisCluster: clustName,
					RedisID:      redisID,
					RedisAddress: address,
					RedisPort:    port,
					ID:           slowLogArr[j].ID,
					Time:         slowLogArr[j].Time,
					Duration:     slowLogArr[j].Duration / 1000,
					Args:         slowLogArr[j].Args,
				}
				redisSlowLogArr = append(redisSlowLogArr, *slowLog)
			}
		}
	}
	return redisSlowLogArr
}
