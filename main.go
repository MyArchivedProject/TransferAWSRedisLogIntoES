package main

import "encoding/json"

// "github.com/aws/aws-sdk-go/service/s3"

func main() {
	InitConfig("")
	// 获取aws的redis节点连接信息
	redisNodeInfoArr := GetAwsRedisClusterInfo()

	// test
	// testNode := redisNodeInfoArr[:1]
	// testNode[0].RedisAddress = "54.83.160.248"
	// testNode[0].RedisPort = 6379
	// testNode[0].RedisID = "vova-multi-test-3-vova1"
	// testNode[0].RedisCluster = "vova-multi-test-3-vova1"

	// 获取redis慢日志
	allSlowLogArr := GetMultiRedisSlowLog(redisNodeInfoArr)

	// 插入进ES
	byteData, _ := json.Marshal(allSlowLogArr)
	var dataArr []map[string]interface{}
	_ = json.Unmarshal(byteData, &dataArr)
	PushDataToES(dataArr)

	// data, _ := json.Marshal(allSlowLogArr)
	// fmt.Printf(string(data))
	// var out bytes.Buffer
	// json.Indent(&out, data, "", "\t")
	// fmt.Printf("array=%v\n", out.String())
}
