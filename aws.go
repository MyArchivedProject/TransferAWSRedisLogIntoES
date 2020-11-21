package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/spf13/viper"
)

// RedisNodeInfo 保存从aws上获取的redis节点信息
type RedisNodeInfo struct {
	RedisCluster string `type:"string"`
	RedisID      string `type:"string"`
	RedisAddress string `type:"string"`
	RedisPort    int64  `type:"int"`
}

type awsConf struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

// 未真实连接aws
func initAWS() (sess *session.Session) {
	// const AWS_ACCESS_KEY_ID string = "AKIA2RIJKRTSUE7434KR"
	// const AWS_SECRET_ACCESS_KEY string = ""
	// const AWS_SESSION_TOKEN string = "TOKEN"
	// const AWS_REGION string = "us-east-1"

	awsConfig := awsConf{
		AccessKeyID:     viper.GetString("aws.AWS_ACCESS_KEY_ID"),
		SecretAccessKey: viper.GetString("aws.AWS_SECRET_ACCESS_KEY"),
		Region:          viper.GetString("aws.AWS_REGION"),
	}

	// 方式一 连接aws
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfig.Region),
		Credentials: credentials.NewStaticCredentials(awsConfig.AccessKeyID, awsConfig.SecretAccessKey, ""),
	})

	// 方式二 连接aws
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-east-1"),
	// 	Credentials: credentials.NewSharedCredentials("", "default"), // 从~/.aws/credentials加载key和密钥
	// })

	errorExit(err)

	return
}

// 未真实连接aws
func connectElasticache(session *session.Session) (svc *elasticache.ElastiCache) {
	// svc := elasticache.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	svc = elasticache.New(session)

	// if svc == nil {
	// 	// errorExit(errors.New("Can not get connect to AWS Elasticache"))
	// 	errorExit(fmt.Errorf("%s", "Can not get connect to AWS Elasticache"))
	// }
	return
}

// 通过aws API拉取aws elasticache 的节点信息
func getAllRedisNodeInfo(svc *elasticache.ElastiCache) []*elasticache.CacheCluster {
	input := &elasticache.DescribeCacheClustersInput{
		// CacheClusterId:    aws.String("test-01"),
		CacheClusterId:    aws.String(""), // 为空则会拉取所有的redis节点信息
		ShowCacheNodeInfo: aws.Bool(true),
	}
	result, err := svc.DescribeCacheClusters(input)
	if err != nil {
		printLog("aws-sdk-go SDK 自己处理异常")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elasticache.ErrCodeCacheClusterNotFoundFault:
				fmt.Println(elasticache.ErrCodeCacheClusterNotFoundFault, aerr.Error())
			case elasticache.ErrCodeInvalidParameterValueException:
				fmt.Println(elasticache.ErrCodeInvalidParameterValueException, aerr.Error())
			case elasticache.ErrCodeInvalidParameterCombinationException:
				fmt.Println(elasticache.ErrCodeInvalidParameterCombinationException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result.CacheClusters
}

// 数据处理操作: 获取redis的集群名称, redis节点ID, redis节点地址, redis节点端口
func customInfo(cacheClusters []*elasticache.CacheCluster) []RedisNodeInfo {
	length := len(cacheClusters)
	var redisNodeInfoArr []RedisNodeInfo = make([]RedisNodeInfo, length, length)
	redisCluster := "--"
	redisID := "--"
	address := "--"
	port := int64(6379)

	for i := 0; i < length; i++ {
		cacheNodes := cacheClusters[i].CacheNodes
		if cacheClusters[i].CacheClusterId != nil {
			redisID = *cacheClusters[i].CacheClusterId
		}
		if cacheClusters[i].ReplicationGroupId != nil {
			redisCluster = *cacheClusters[i].ReplicationGroupId
		}

		if cacheNodes[0].Endpoint.Address != nil {
			address = *cacheNodes[0].Endpoint.Address
		}
		if cacheNodes[0].Endpoint.Port != nil {
			port = *cacheNodes[0].Endpoint.Port
		}

		redisNodeInfoArr[i] = *&RedisNodeInfo{
			RedisCluster: redisCluster,
			RedisAddress: address,
			RedisPort:    port,
			RedisID:      redisID,
		}

	}
	return redisNodeInfoArr
}

// GetAwsRedisClusterInfo 对外暴露的接口
func GetAwsRedisClusterInfo() []RedisNodeInfo {
	// 获取aws会话
	session := connectElasticache(initAWS())
	// 连接aws获取redis节点数据
	cacheClusters := getAllRedisNodeInfo(session)
	// 处理redis节点数据
	return customInfo(cacheClusters)
}
