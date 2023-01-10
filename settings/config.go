package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() error {
	//viper.SetConfigType("yaml")            // 如果配置文件的名称中没有扩展名，则需要配置此项
	//viper.AddConfigPath(".")               // 还可以在工作目录中查找配置
	viper.SetConfigFile("config_dev.yaml") // 指定配置文件
	err := viper.ReadInConfig()            // 查找并读取配置文件
	if err != nil {                        // 处理读取配置文件的错误
		fmt.Printf("viper.ReadConfig() failed, err:%v\n", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return err
	}
	// 监控配置文件变化
	// viper.WatchConfig()
	return nil
}
