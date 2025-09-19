package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zihunyu/gin-login/service"
	"gorm.io/gorm"
)

// RegisterRequest 定义前端传入的注册参数结构
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler 返回一个 Gin 处理函数，用于用户注册
func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
			return
		}

		// 调用服务层注册用户
		if err := service.RegisterUser(db, req.Email, req.Username, req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		// 返回注册成功
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功"})
	}
}
