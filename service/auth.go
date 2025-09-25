package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zihunyu/gin-login/config"
	"github.com/zihunyu/gin-login/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
}

func NewAuthService(db *gorm.DB, r *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{
		DB:     db,
		Redis:  r,
		Config: cfg,
	}
}

// 自定义 JWT Claims
type MyClaims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// LoginUser 执行登录：检查锁定 -> 校验用户名密码 -> 失败计数/锁定 -> 成功生成 JWT
func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, *model.User, error) {
	// 1) 检查是否锁定
	lockKey := fmt.Sprintf("login:lock:%s", username)
	exists, err := s.Redis.Exists(ctx, lockKey).Result()
	if err != nil {
		// Redis 出问题时，为了可用性选择不过早阻断（或你可以选择返回错误）
		return "", nil, errors.New("登录服务暂时不可用")
	}
	if exists > 0 {
		return "", nil, errors.New("账号已被锁定，请稍后再试")
	}

	// 2) 查询用户（按 username）
	var user model.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 计为一次失败（避免用户名枚举）
			_ = s.incrementFailure(ctx, username)
			return "", nil, errors.New("用户名或密码错误")
		}
		return "", nil, errors.New("数据库查询错误")
	}

	// 3) 比对 bcrypt 密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// 密码错误 -> 记录失败并可能锁定
		_ = s.incrementFailure(ctx, username)
		return "", nil, errors.New("用户名或密码错误")
	}

	// 4) 登录成功 -> 清除失败计数（可选，但推荐）
	failKey := fmt.Sprintf("login:fail:%s", username)
	_ = s.Redis.Del(ctx, failKey).Err()

	// 5) 生成 JWT
	expMinutes := s.Config.JWT.ExpireMinutes
	expireAt := time.Now().Add(time.Duration(expMinutes) * time.Minute)
	claims := MyClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.Config.JWT.Secret))
	if err != nil {
		return "", nil, errors.New("生成 token 失败")
	}

	// 返回 token 与用户信息（包含密码哈希，但 controller 不会返回给前端）
	return tokenStr, &user, nil
}

// incrementFailure 内部函数：增加失败计数并在达到阈值时锁定账号
func (s *AuthService) incrementFailure(ctx context.Context, username string) error {
	failKey := fmt.Sprintf("login:fail:%s", username)

	// 增加失败计数
	failures, err := s.Redis.Incr(ctx, failKey).Result()
	if err != nil {
		return err
	}

	// 如果这是第一次失败，则设置窗口期过期时间
	if failures == 1 {
		_ = s.Redis.Expire(ctx, failKey, time.Duration(s.Config.LoginLimit.WindowMinutes)*time.Minute).Err()
	}

	// 达到阈值则设置锁定
	if failures >= int64(s.Config.LoginLimit.MaxFailures) {
		lockKey := fmt.Sprintf("login:lock:%s", username)
		_ = s.Redis.Set(ctx, lockKey, "1", time.Duration(s.Config.LoginLimit.LockMinutes)*time.Minute).Err()
	}

	return nil
}
