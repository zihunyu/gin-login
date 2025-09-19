package config

import (
	"github.com/spf13/viper"
)

type MYSQL struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type App struct {
	Port string `mapstructure:"port"`
}

// Config 定义配置结构体
type Config struct {
	App   App   `mapstructure:"app"`
	MYSQL MYSQL `mapstructure:"mysql"`
}

// LoadConfig 读取config/config.yaml 并解析到Config结构体
func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil

}
