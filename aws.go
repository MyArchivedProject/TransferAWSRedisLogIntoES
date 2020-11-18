package main

// 错误直接在函数内输出

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type redisNodeInfo struct {
	ClusterName string `type:"string"`
	NodeAddress string `type:"string"`
	NodePort    int64  `type:"int"`
}

func connectAWS() (sess *session.Session) {
	// const AWS_ACCESS_KEY_ID string = "AKIA2RIJKRTSUE7434KR"
	// const AWS_SECRET_ACCESS_KEY string = ""
	// const AWS_SESSION_TOKEN string = "TOKEN"
	// const AWS_Region string = "us-east-1"

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(AWS_Region),
	// 	Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	// })

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	if err != nil {
		fmt.Println("Error creating session ", err)
		panic("Panic----")
	}
	return
}

func connectElasticache() (svc *elasticache.ElastiCache) {
	session := connectAWS()
	// svc := elasticache.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	svc = elasticache.New(session)
	fmt.Printf("%+v\n", svc)
	if svc == nil {
		fmt.Println("Err connectElasticache()---------------")
	}
	return
}

func getAllClusterNodeInfo(svc *elasticache.ElastiCache) []*elasticache.CacheCluster {
	input := &elasticache.DescribeCacheClustersInput{
		// CacheClusterId:    aws.String("vv-andes-test-0003-001"),
		CacheClusterId:    aws.String(""),
		ShowCacheNodeInfo: aws.Bool(true),
	}
	result, err := svc.DescribeCacheClusters(input)
	if err != nil {
		fmt.Println("Error getAllClusterNodeInfo()----------------")

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

// 获取node节点的集群名称,节点地址,节点端口
func customInfo() []redisNodeInfo {
	cacheClusters := getAllClusterNodeInfo(connectElasticache())
	length := len(cacheClusters)
	var redisNodeInfoArr []redisNodeInfo = make([]redisNodeInfo, length, length)
	clusterName := ""
	address := ""
	port := int64(0)
	// for _, v := range cacheClusters { //ok
	for i := 0; i < length; i++ {
		cacheNodes := cacheClusters[i].CacheNodes
		if cacheClusters[i].ReplicationGroupId != nil {
			clusterName = *cacheClusters[i].ReplicationGroupId
		}
		if cacheNodes[0].Endpoint.Address != nil {
			address = *cacheNodes[0].Endpoint.Address
		}
		if cacheNodes[0].Endpoint.Port != nil {
			port = *cacheNodes[0].Endpoint.Port
		}
		fmt.Println(cacheNodes[0].Endpoint.Address)
		fmt.Println(cacheNodes[0].Endpoint.Port)

		// redisNodeInfoArr = append(redisNodeInfoArr, *&redisNodeInfo{})  //ok
		redisNodeInfoArr[i] = *&redisNodeInfo{
			ClusterName: clusterName,
			NodeAddress: address,
			NodePort:    port,
		}

	}
	return redisNodeInfoArr
}
