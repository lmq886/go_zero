/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @LastEditors: 羡鱼 lmqqq1435456124@163.com && 羡鱼
 * @LastEditTime: 2026-04-23 17:27:00
 * @FilePath: \go_zero\api\internal\middleware\jwt_auth.go
 * @Description: JWT认证中间件，用于验证用户身份
 */
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/config"
)

// JwtAuthMiddleware JWT认证中间件结构体
type JwtAuthMiddleware struct {
	Config config.Config // 系统配置
}

// NewJwtAuthMiddleware 创建JWT认证中间件实例
func NewJwtAuthMiddleware(c config.Config) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{Config: c}
}

// Handle 中间件处理函数
// 参数: next - 下一个处理函数
// 返回: http.HandlerFunc - 包装后的处理函数
func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 从HTTP请求头中获取Authorization Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logx.Error("Authorization header is empty")
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未登录或Token已过期",
			})
			return
		}

		// 2. 验证Token格式是否为Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logx.Error("Invalid authorization header format")
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token格式错误",
			})
			return
		}

		token := parts[1]

		// 3. 解析并验证Token有效性
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			// 使用配置中的密钥验证Token签名
			return []byte(m.Config.JwtAuth.AccessSecret), nil
		})

		if err != nil {
			logx.Error("Token parse error:", err)
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token无效或已过期",
			})
			return
		}

		// 4. 从Token中提取用户ID并存入Context
		userId := int64(claims["uid"].(float64))
		r = r.WithContext(setUserId(r.Context(), userId))

		// 5. 调用下一个处理函数
		next(w, r)
	}
}

// userKey 用于在Context中存储用户ID的键类型
type userKey struct{}

// setUserId 将用户ID存入Context
// 参数: ctx - 上下文对象
// 参数: userId - 用户ID
// 返回: interface{} - 新的上下文对象
func setUserId(ctx interface{}, userId int64) interface{} {
	return context.WithValue(ctx, userKey{}, userId)
}

// GetUserId 从Context中获取用户ID
// 参数: ctx - 上下文对象
// 返回: int64 - 用户ID，未找到则返回0
func GetUserId(ctx interface{}) int64 {
	if userId, ok := ctx.(interface{ Value(interface{}) interface{} }).Value(userKey{}).(int64); ok {
		return userId
	}
	return 0
}
