package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	// elasticsearch6 "github.com/elastic/go-elasticsearch/v6"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func connectES() *elasticsearch.Client {
	address := viper.GetString("es.address")
	if address == "" {
		log.Fatalln("Fatal connectES(): Can not get es address")
	}
	log.Println("connecting es: address=" + address + "; elasticsearchSDKVersion=" + elasticsearch.Version)
	cfg := elasticsearch.Config{
		Addresses: []string{
			address,
		},
		Username: "",
		Password: "",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalln("Error connectES() creating the client:\n" + error.Error(err))
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalln("Error connectES() getting response:\n" + error.Error(err))
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatalf("Error connectES(): %s", res.String())
	}
	log.Println(res)
	return es
}

// func operateES1(es *elasticsearch.Client, redisSlowLogArr []RedisSlowLog) {
// 	// Create the BulkIndexer
// 	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
// 		Index:      "redis-slowlog-", // The default index name
// 		Client:     es,               // The Elasticsearch client
// 		NumWorkers: 2,                // The number of worker goroutines
// 		//FlushBytes:    int(flushBytes),  // The flush threshold in bytes
// 		//FlushInterval: 30 * time.Second, // The periodic flush interval
// 	})
// }
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func insertBatch(es *elasticsearch.Client, redisSlowLogArr []RedisSlowLog) {
	slowlogNum := len(redisSlowLogArr)
	var bodyBuf bytes.Buffer
	for i := 0; i < slowlogNum; i++ {
		createLine := map[string]interface{}{
			"create": map[string]interface{}{
				"_index": "redis-slowlog-",
				// "_id":    "test_" + strconv.Itoa(i),
				// "_type":  "test_type",
			},
		}
		jsonStr, _ := json.Marshal(createLine)
		bodyBuf.Write(jsonStr)
		bodyBuf.WriteByte('\n')

		// body := map[string]interface{}{
		// 	"num": i % 3,
		// 	"v":   i,
		// 	"str": "test" + strconv.Itoa(i),
		// }
		body := redisSlowLogArr[i]
		jsonStr, _ = json.Marshal(body)
		bodyBuf.Write(jsonStr)
		bodyBuf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Body: &bodyBuf,
	}
	res, err := req.Do(context.Background(), es)
	checkError(err)
	defer res.Body.Close()
	fmt.Println(res.String())
}

// v1 非批量插入
func operateESSingle(es *elasticsearch.Client, redisSlowLogArr []RedisSlowLog) {
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
	if err != nil {
		log.Fatalf("Error operateES() Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
}
func PushDataToES(redisSlowLogArr []RedisSlowLog) {
	es := connectES()
	insertBatch(es, redisSlowLogArr)

}
