# GoAWS
- 参考文档  
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html  
https://pkg.go.dev/github.com/go-redis/redis/v8#pkg-functions  
https://github.com/elastic/go-elasticsearch  

https://www.cnblogs.com/Me1onRind/p/11534544.html

## 架构
1. 错误处理方式：函数内的错误不返回调用层，直接在函数内输出。

## 功能
1. 将aws的redis数据库慢日志拉取下来放进ES中

## 部署
1. 编译
go build

2. 填写配置文件
cp config.template.yaml .config.yaml  
vim .config.yaml  

3. 运行二进制文件