package main

import (
	"context"
	"log"
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

func operateES(es *elasticsearch.Client, redisSlowLogArr []RedisSlowLog) {
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
	
	req := esapi.IndexRequest{
		Index: "redis-slowlog-",                         // Index name
		Body:  strings.NewReader(`{"field1" : "Test"}`), // Document body
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
	operateES(es,redisSlowLogArr)

}
