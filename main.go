package main

import (
	"encoding/json"
	"fmt"
	"os"

	goaws "goaws/src"
)

// "github.com/aws/aws-sdk-go/service/s3"

func main() {
	defer goaws.TimeCost()()
	args()
}
func args() {
	fmt.Println(os.Args)
	if len(os.Args) == 1 {
		run()
	}else if os.Args[1] == "test"{
		test()
	}
}

func run() {
	goaws.InitConfig("")
	// 获取aws的redis节点连接信息
	redisNodeInfoArr := goaws.GetAwsRedisClusterInfo()

	// TODO 先去检查redis和es是否可连接,防止耗时获取redis慢日志后发现ES连接不上
	// 获取redis慢日志
	allSlowLogArr := goaws.GetMultiRedisSlowLog(redisNodeInfoArr)

	// 插入数据进ES
	byteData, _ := json.Marshal(allSlowLogArr)
	var dataArr []map[string]interface{}
	_ = json.Unmarshal(byteData, &dataArr)
	goaws.PushDataToES(dataArr)
}
func test(){
	goaws.InitConfig("")
	// 获取aws的redis节点连接信息
	redisNodeInfoArr := goaws.GetAwsRedisClusterInfo()

	// test
	testNode := redisNodeInfoArr[:1]
	testNode[0].RedisAddress = "54.83.160.248"
	testNode[0].RedisPort = 6379
	testNode[0].RedisID = "vova-multi-test-3-vova1"
	testNode[0].RedisCluster = "vova-multi-test-3-vova1"
	allSlowLogArr := goaws.GetMultiRedisSlowLog(testNode)


	// 插入数据进ES
	byteData, _ := json.Marshal(allSlowLogArr)
	var dataArr []map[string]interface{}
	_ = json.Unmarshal(byteData, &dataArr)
	goaws.PushDataToES(dataArr)

	// test
	// data, _ := json.Marshal(allSlowLogArr)
	// fmt.Printf(string(data))
	// var out bytes.Buffer
	// json.Indent(&out, data, "", "\t")
	// fmt.Printf("array=%v\n", out.String())
}