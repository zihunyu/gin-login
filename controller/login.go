package controller

import (
	"net/http"

	"github.com/zihunyu/gin-login/service"

	"github.com/gin-gonic/gin"
)

// LoginRequest 前端请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler 返回一个 gin.HandlerFunc，封装 authService
func LoginHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
			return
		}

		token, user, err := authSvc.LoginUser(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			// 登录失败
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 登录成功：返回 token 和用户信息（不要返回密码）
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				// 其他非敏感字段可以返回
			},
		})
	}
}
