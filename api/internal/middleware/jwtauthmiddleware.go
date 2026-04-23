package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go_zero/api/internal/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// 定义上下文键类型
// 使用自定义类型避免键冲突
type contextKey string

const (
	// UserIdKey 用户ID上下文键
	UserIdKey contextKey = "userId"
	// UsernameKey 用户名上下文键
	UsernameKey contextKey = "username"
	// RolesKey 用户角色上下文键
	RolesKey contextKey = "roles"
	// PermissionsKey 用户权限上下文键
	PermissionsKey contextKey = "permissions"
)

// 定义 JWT 相关错误
var (
	// ErrTokenInvalid 令牌无效错误
	ErrTokenInvalid = errors.New("token is invalid")
	// ErrTokenExpired 令牌过期错误
	ErrTokenExpired = errors.New("token is expired")
)

// JwtAuthMiddleware JWT 认证中间件
// 用于验证用户身份，提取用户信息并注入到上下文中
// 遵循 GoZero 中间件最佳实践，实现 rest.Middleware 接口
type JwtAuthMiddleware struct {
	// 配置信息
	config config.AuthConfig
}

// NewJwtAuthMiddleware 创建 JWT 认证中间件实例
// 参数 config: JWT 配置
// 返回值: JWT 认证中间件实例
func NewJwtAuthMiddleware(config config.AuthConfig) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		config: config,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 没有 Authorization 头，返回 401 未授权
			logx.Errorf("Missing Authorization header")
			WriteError(w, http.StatusUnauthorized, "未授权：缺少认证令牌")
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// Authorization 头格式错误
			logx.Errorf("Invalid Authorization header format")
			WriteError(w, http.StatusUnauthorized, "未授权：认证令牌格式错误")
			return
		}

		// 提取令牌
		tokenString := parts[1]

		// 解析和验证令牌
		claims, err := m.parseToken(tokenString)
		if err != nil {
			logx.Errorf("Token validation failed: %v", err)
			WriteError(w, http.StatusUnauthorized, "未授权：令牌无效或已过期")
			return
		}

		// 将用户信息注入到上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIdKey, claims.UserId)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, RolesKey, claims.Roles)
		ctx = context.WithValue(ctx, PermissionsKey, claims.Permissions)

		// 继续处理请求
		next(w, r.WithContext(ctx))
	}
}

// CustomClaims 自定义 JWT 声明
// 包含用户ID、用户名、角色、权限等信息
type CustomClaims struct {
	// 用户ID
	UserId int64 `json:"userId"`
	// 用户名
	Username string `json:"username"`
	// 角色列表
	Roles []string `json:"roles"`
	// 权限列表
	Permissions []string `json:"permissions"`
	// 标准 JWT 声明
	jwt.RegisteredClaims
}

// parseToken 解析和验证 JWT 令牌
// 参数 tokenString: JWT 令牌字符串
// 返回值: 自定义声明和错误信息
func (m *JwtAuthMiddleware) parseToken(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		// 返回签名密钥
		return []byte(m.config.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 检查令牌是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, ErrTokenExpired
		}
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GenerateToken 生成 JWT 令牌
// 参数 userId: 用户ID
// 参数 username: 用户名
// 参数 roles: 角色列表
// 参数 permissions: 权限列表
// 返回值: 访问令牌、刷新令牌、过期时间和错误信息
func GenerateToken(
	userId int64,
	username string,
	roles []string,
	permissions []string,
	config config.AuthConfig,
) (string, string, int64, error) {
	// 计算过期时间
	now := time.Now()
	accessExpire := now.Add(time.Duration(config.AccessExpire) * time.Second)
	refreshExpire := now.Add(time.Duration(config.RefreshExpire) * time.Second)

	// 创建访问令牌声明
	accessClaims := CustomClaims{
		UserId:      userId,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpire),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Issuer,
			Audience:  jwt.ClaimStrings{config.Audience},
		},
	}

	// 创建刷新令牌声明
	refreshClaims := CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpire),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Issuer,
		},
	}

	// 生成访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", "", 0, err
	}

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", "", 0, err
	}

	return accessTokenString, refreshTokenString, accessExpire.Unix(), nil
}

// ParseRefreshToken 解析刷新令牌
// 参数 tokenString: 刷新令牌字符串
// 参数 config: JWT 配置
// 返回值: 用户ID、用户名和错误信息
func ParseRefreshToken(tokenString string, config config.AuthConfig) (int64, string, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		// 返回签名密钥
		return []byte(config.RefreshSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 检查令牌是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return 0, "", ErrTokenExpired
		}
		return claims.UserId, claims.Username, nil
	}

	return 0, "", ErrTokenInvalid
}

// GetUserIdFromContext 从上下文中获取用户ID
// 参数 ctx: 上下文
// 返回值: 用户ID和是否存在
func GetUserIdFromContext(ctx context.Context) (int64, bool) {
	userId, ok := ctx.Value(UserIdKey).(int64)
	return userId, ok
}

// GetUsernameFromContext 从上下文中获取用户名
// 参数 ctx: 上下文
// 返回值: 用户名和是否存在
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// GetRolesFromContext 从上下文中获取角色列表
// 参数 ctx: 上下文
// 返回值: 角色列表和是否存在
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey).([]string)
	return roles, ok
}

// GetPermissionsFromContext 从上下文中获取权限列表
// 参数 ctx: 上下文
// 返回值: 权限列表和是否存在
func GetPermissionsFromContext(ctx context.Context) ([]string, bool) {
	permissions, ok := ctx.Value(PermissionsKey).([]string)
	return permissions, ok
}
