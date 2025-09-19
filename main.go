package main

import (
	"fmt"
	"log"

	"github.com/zihunyu/gin-login/config"
	"github.com/zihunyu/gin-login/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("加载配置失败：", err)
	}

	// 构造MySQL DSN并连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MYSQL.User, cfg.MYSQL.Password,
		cfg.MYSQL.Host, cfg.MYSQL.Port,
		cfg.MYSQL.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移用户表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 打印配置
	log.Printf("mysql host: %s", cfg.MYSQL.Host)
	log.Printf("mysql port: %s", cfg.MYSQL.Port)
	log.Printf("mysql user: %s", cfg.MYSQL.User)
	log.Printf("mysql password: %s", cfg.MYSQL.Password)
	log.Printf("mysql dbname: %s", cfg.MYSQL.DBName)
}
