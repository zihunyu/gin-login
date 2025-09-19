package service

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/zihunyu/gin-login/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser 实现用户注册逻辑：校验、加密、入库
func RegisterUser(db *gorm.DB, email, username, password string) error {
	// 输入校验
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("无效的邮箱格式")
	}
	if strings.TrimSpace(username) == "" {
		return errors.New("用户名不能为空")
	}
	if len(password) < 8 {
		return errors.New("密码长度必须至少8位")
	}

	// 检查用户是否已存在（按邮箱查重）
	var existing model.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("用户已存在")
	}

	// 密码哈希加密（bcrypt）
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost 通常为 10
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户记录
	user := model.User{
		Email:    email,
		Username: username,
		Password: string(hash),
	}
	if err := db.Create(&user).Error; err != nil {
		return errors.New("注册失败")
	}
	return nil
}
