package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zihunyu/gin-login/controller"
	"github.com/zihunyu/gin-login/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authSvc *service.AuthService) *gin.Engine {
	r := gin.Default()

	// 注册接口（你原本的）
	r.POST("/register", controller.RegisterHandler(db))

	// 新增登录接口
	r.POST("/login", controller.LoginHandler(authSvc))

	return r
}
