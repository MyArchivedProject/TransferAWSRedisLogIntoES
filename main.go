package main

// "github.com/aws/aws-sdk-go/service/s3"

func main() {
	InitConfig("")
	redisNodeInfoArr := GetAwsRedisClusterInfo()

	// test
	testNode := redisNodeInfoArr[:1]
	testNode[0].RedisAddress = "54.83.160.248"
	testNode[0].RedisPort = 6379
	allSlowLogArr := GetMultiRedisSlowLog(testNode)

	PushDataToES(allSlowLogArr)

	// data, _ := json.Marshal(allSlowLogArr)
	// var out bytes.Buffer
	// json.Indent(&out, data, "", "\t")
	// fmt.Printf("array=%v\n", out.String())
}
