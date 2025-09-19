package main

import (
	"log"

	"github.com/zihunyu/gin-login/config"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("加载配置失败：", err)
	}
	// 打印配置
	log.Printf("mysql host: %s", cfg.MYSQL.Host)
	log.Printf("mysql port: %s", cfg.MYSQL.Port)
	log.Printf("mysql user: %s", cfg.MYSQL.User)
	log.Printf("mysql password: %s", cfg.MYSQL.Password)
	log.Printf("mysql dbname: %s", cfg.MYSQL.DBName)
}
