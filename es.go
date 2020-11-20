package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	// elasticsearch6 "github.com/elastic/go-elasticsearch/v6"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type esConf struct {
	AddressArr []string
	Username   string
	Password   string
}

func connectES() *elasticsearch.Client {
	if viper.GetString("es.address") == "" {
		errorExit(fmt.Errorf("%s", "Can not get es address from config"))
	}

	var addressArr []string = make([]string, 0)
	addressArr = append(addressArr, viper.GetString("es.address"))
	esConf := esConf{
		AddressArr: addressArr,
		Username:   "",
		Password:   "",
	}
	printLog("connecting es: address=" + addressArr[0] + "; elasticsearchSDKVersion=" + elasticsearch.Version)

	cfg := elasticsearch.Config{
		Addresses: esConf.AddressArr,
		Username:  "",
		Password:  "",
	}
	es, err := elasticsearch.NewClient(cfg)

	errorExit(err)

	res, err := es.Info()
	errorExit(err)

	defer res.Body.Close()
	if res.IsError() {
		errorExit(fmt.Errorf("%s", res.String()))
	}
	printLog(res)
	return es
}

func insertBatch(es *elasticsearch.Client, dataArr []map[string]interface{}, index string) {
	slowlogNum := len(dataArr)
	var bodyBuf bytes.Buffer

	// 遍历慢日志 生成Buffer
	for i := 0; i < slowlogNum; i++ {

		// 创建唯一ID,防止重复插入
		timeTemp, _ := time.Parse("2006-01-02T15:04:05+08:00", dataArr[i]["Time"].(string))
		timeStamp := timeTemp.Unix()
		// uniqID := dataArr[i]["RedisAddress"].(string) + strconv.FormatFloat(dataArr[i]["ID"].(float64), 'E', -1, 64) + strconv.FormatInt(timeStamp, 10)
		uniqID := fmt.Sprint(dataArr[i]["RedisAddress"], dataArr[i]["ID"], timeStamp)

		createLine := map[string]interface{}{
			"create": map[string]interface{}{
				"_index": index,
				"_id":    uniqID,
			},
		}

		// fmt.Println(dataArr[i]["Duration"])
		// fmt.Println(dataArr[i]["ID"])
		// fmt.Println(dataArr[i]["RedisAddress"])
		// fmt.Println(dataArr[i]["Time"])

		jsonStr, _ := json.Marshal(createLine)
		bodyBuf.Write(jsonStr)
		bodyBuf.WriteByte('\n')

		// body := map[string]interface{}{
		// 	"num": i % 3,
		// 	"v":   i,
		// 	"str": "test" + strconv.Itoa(i),
		// }
		body := dataArr[i]
		jsonStr, _ = json.Marshal(body)
		bodyBuf.Write(jsonStr)
		bodyBuf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Body: &bodyBuf,
	}
	res, err := req.Do(context.Background(), es)
	defer res.Body.Close()
	errorTolerate(err)

	printLog(res.String())
}

// 未使用到
func insertSiingle(es *elasticsearch.Client, redisSlowLogArr []RedisSlowLog) {
	// 方式一
	// res, err = es.Index(
	// 	"test",                                  // Index name
	// 	strings.NewReader(`{"title" : "Test"}`), // Document body
	// 	es.Index.WithDocumentID("1"),            // Document ID
	// 	es.Index.WithRefresh("true"),            // Refresh
	// )
	// if err != nil {
	// 	log.Fatalf("ERROR: %s", err)
	// }
	// defer res.Body.Close()
	// log.Println(res)

	// 单条插入 方式二
	data, _ := json.Marshal(redisSlowLogArr[0])
	req := esapi.IndexRequest{
		Index: "redis-slowlog-", // Index name
		// Body:  strings.NewReader(`{"field1" : "Test"}`), // Document body
		Body: strings.NewReader(string(data)), // Document body
		// DocumentID: "1",                                     // 指定 Document ID
		Refresh: "true", // Refresh
	}

	res, err := req.Do(context.Background(), es)

	errorTolerate(err)

	defer res.Body.Close()

	printLog(res)
}

// PushDataToES 批量推数据进es
func PushDataToES(dataArr []map[string]interface{}) {
	es := connectES()
	index := viper.GetString("es.index")
	insertBatch(es, dataArr, index)
	printLog("批量插入数据进ES结束")
}
