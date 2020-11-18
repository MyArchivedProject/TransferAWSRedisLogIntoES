package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	a := customInfo()
	data, _ := json.Marshal(a)
	var out bytes.Buffer
	json.Indent(&out, data, "", "\t")
	fmt.Printf("array=%v\n", out.String())

}
