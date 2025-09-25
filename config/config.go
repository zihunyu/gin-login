package config

import (
	"github.com/spf13/viper"
)

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	ExpireMinutes int    `mapstructure:"expire_minutes"`
}

type LoginLimitConfig struct {
	WindowMinutes int `mapstructure:"window_minutes"`
	MaxFailures   int `mapstructure:"max_failures"`
	LockMinutes   int `mapstructure:"lock_minutes"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Config 定义配置结构体
type Config struct {
	JWT            JWTConfig        `mapstructure:"jwt"`
	LoginLimit     LoginLimitConfig `mapstructure:"login_limit"`
	Redis          RedisConfig      `mapstructure:"redis"`
	RequestTimeout int              `mapstructure:"request_timeout"`
	App            struct {
		Port string
	}
	MySQL struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"mysql"`
}

// LoadConfig 读取 config/config.yaml 并解析到 Config 结构体
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
