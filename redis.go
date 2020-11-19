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

// RedisSlowLog 所有的慢日志
type RedisSlowLog struct {
	RedisCluster  string `type:"string"`
	RedisAddress string `type:"string"`
	RedisPort    int64  `type:"int"`
	ID           int64
	Time         time.Time
	Duration     time.Duration
	Args         []string
}

type redisConf struct {
	Addr     string
	Password string
	DB       int
}

func getRedisClient(redisConf redisConf) (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,     // use default Addr
		Password: redisConf.Password, // no password set
		DB:       redisConf.DB,       // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("RDB连接失败 Addr=" + redisConf.Addr)
		panic(err)
	} else {
		log.Println("RDB连接成功 Addr=" + redisConf.Addr)
	}
	return
}

func getSlowLog(rdb *redis.Client, num int64) []redis.SlowLog {
	// res, err := rdb.Do(ctx, "slowlog", "get", "128").Result() //返回接口
	res, err := rdb.SlowLogGet(ctx, num).Result() //返回被格式化过的数组

	if err != nil {
		log.Println("Error getSlowLog():\n" + error.Error(err))
	}
	return res
}

// GetMultiRedisSlowLog 输入数组：包含每个redis节点信息；返回数组：所有的慢日志
func GetMultiRedisSlowLog(redisNodeInfoArr []RedisNodeInfo) (redisSlowLogArr []RedisSlowLog) {
	nodeNum := len(redisNodeInfoArr)
	slowlogNum := viper.GetInt64("redis.slowlog_num")
	if nodeNum == 0 {
		log.Println("connectMultiRedis() redisNodeInfoArr==[]")
		return nil
	}
	redisSlowLogArr = make([]RedisSlowLog, nodeNum*int(slowlogNum)/2)

	// 遍历所有节点
	for i := 0; i < nodeNum; i++ {
		address := redisNodeInfoArr[i].RedisAddress
		port := redisNodeInfoArr[i].RedisPort
		clustName := redisNodeInfoArr[i].RedisCluster

		redisConf := redisConf{
			Addr: address + ":" + strconv.FormatInt(port, 10),
		}
		rdb := getRedisClient(redisConf)
		// defer rdb.Close() // 选择手动close
		// 遍历节点上所有的慢日志
		slowLogArr := getSlowLog(rdb, slowlogNum)
		for j :=0; j< len(slowLogArr);j++{
			slowLog := &RedisSlowLog{
				RedisCluster:  clustName,
				RedisAddress: address,
				RedisPort:    port,
				ID:           slowLogArr[j].ID,
				Time:         slowLogArr[j].Time,
				Duration:     slowLogArr[j].Duration,
				Args:         slowLogArr[j].Args,
			}
			redisSlowLogArr = append(redisSlowLogArr,*slowLog)
		}
		rdb.Close()
	}
	return redisSlowLogArr
}
