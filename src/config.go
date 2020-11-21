package src

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Name string
}

// InitConfig (配置文件路径)
func InitConfig(cfg string) {
	c := config{
		Name: cfg,
	}
	// 初始化配置文件
	c.readConfig()
	return
}

// 加载配置
func (c *config) readConfig() {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("./") // 路径
		viper.SetConfigName(".config")
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML


	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		log.Panic(err)
		return
	}
	log.Println("Init config file successfully")
	return
}
