# GoAWS
参考文档
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
https://pkg.go.dev/github.com/go-redis/redis/v8#pkg-functions
https://github.com/elastic/go-elasticsearch

https://www.cnblogs.com/Me1onRind/p/11534544.html

## 架构
1. 错误处理方式：函数内的错误不返回调用层，直接在函数内输出。


## 部署
cp config.template.yaml config.yaml