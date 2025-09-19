package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zihunyu/gin-login/controller"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 注册用户路由
	r.POST("/register", controller.RegisterHandler(db))
	return r
}
