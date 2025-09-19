package model

import (
	"time"
)

// User 定义用户表结构，对应数据库的 users 表
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                                // 默认主键
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"` // 邮箱，唯一索引
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`           // 用户名
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`                 // 密文密码，不输出到 JSON
	CreatedAt time.Time `json:"created_at"`                                          // 创建时间
}
