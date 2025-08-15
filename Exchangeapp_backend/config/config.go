package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct{
	App struct{
		Name string
		Port string
	}
	Database struct{
		Dsn string
	}
}

var AppConfig *Config

//配置config。通过InitConfig方法的viper
func InitConfig()  {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {//读取配置文件是否报错
		log.Fatalf("Error reading config file: %v",err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {//"Unmarshal()"方法是反序列化，将config.yml配置文件转为AppConfig结构体
		log.Fatalf("Unable to decode into struct: %v",err)
	}
	initDB()
}