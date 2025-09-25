package main

import (
	"fmt"
	"log"

	"github.com/zihunyu/gin-login/config"
	"github.com/zihunyu/gin-login/model"
	"github.com/zihunyu/gin-login/router"
	"github.com/zihunyu/gin-login/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 构造 MySQL DSN 并连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password,
		cfg.MySQL.Host, cfg.MySQL.Port,
		cfg.MySQL.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移：根据模型创建或更新表结构
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	// --- 新增：初始化 Redis 客户端（用于登录失败限流）
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// --- 新增：创建 AuthService，注入 db / redis / cfg
	authSvc := service.NewAuthService(db, redisClient, cfg)

	// 初始化路由
	r := router.SetupRouter(db, authSvc)

	// 启动服务
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("服务启动失败:", err)
	}

	// 后续：初始化路由并启动服务
	// ...
}
