package src

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
	esapi6 "github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/spf13/viper"
)

func connectES6() *elasticsearch6.Client {
	if viper.GetString("es.address") == "" {
		PrintErrorExit(fmt.Errorf("%s", "Can not get es address from config"))
	}

	var addressArr []string = make([]string, 0)
	addressArr = append(addressArr, viper.GetString("es.address"))
	esConf := esConf{
		AddressArr: addressArr,
		Username:   "",
		Password:   "",
	}
	PrintLog("connecting es: address=" + addressArr[0] + "; elasticsearchSDKVersion=" + elasticsearch6.Version)

	cfg := elasticsearch6.Config{
		Addresses: esConf.AddressArr,
		Username:  "",
		Password:  "",
	}
	es, err := elasticsearch6.NewClient(cfg)

	PrintErrorExit(err)

	res, err := es.Info()
	PrintErrorExit(err)

	defer res.Body.Close()
	if res.IsError() {
		PrintErrorExit(fmt.Errorf("%s", res.String()))
	}
	PrintLog(res)
	return es
}

func insertBatchES6(es *elasticsearch6.Client, dataArr []map[string]interface{}, index string) {
	slowlogNum := len(dataArr)
	PrintLog("将向ES批量插入 " + strconv.Itoa(slowlogNum) + " 条数据")
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
				"_type":  index,
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

	req := esapi6.BulkRequest{
		Body: &bodyBuf,
	}
	res, err := req.Do(context.Background(), es)
	defer res.Body.Close()
	PrintErrorTolerate(err)

	PrintLog(res.String())
}

// PushDataToES6 批量推数据进es
func PushDataToES6(dataArr []map[string]interface{}) {
	es := connectES6()
	index := viper.GetString("es.index")
	insertBatchES6(es, dataArr, index)
	PrintLog("批量插入数据进ES结束")
}
