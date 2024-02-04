package util

import (
    "github.com/spf13/viper"
)

// 应用程序的配置变量
// 映射app.env文件中的配置变量
type Config struct {
    DBDriver      string `mapstructure:"DB_DRIVER"`
    DBSource      string `mapstructure:"DB_SOURCE"`
    ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// 从环境变量中加载配置
// 返回一个配置对象或者错误
func LoadConfig(path string) (config Config, err error) {
    //  设置配置文件的搜索路径
    viper.AddConfigPath(path)
    // 设置配置文件的名称
    viper.SetConfigName("app")
    // 设置配置文件的类型
    viper.SetConfigType("env")

    // 从环境变量中读取配置
    viper.AutomaticEnv()

    // 读取配置文件
    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    // 将配置文件中的配置变量映射到配置对象
    err = viper.Unmarshal(&config)
    return
}
