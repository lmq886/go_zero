/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @LastEditors: 羡鱼 lmqqq1435456124@163.com && 羡鱼
 * @LastEditTime: 2026-04-23 09:37:35
 * @FilePath: \go_zero\api\internal\middleware\jwt_auth.go
 * @Description: 
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

type JwtAuthMiddleware struct {
	Config config.Config
}

func NewJwtAuthMiddleware(c config.Config) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{Config: c}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Header中获取Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logx.Error("Authorization header is empty")
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未登录或Token已过期",
			})
			return
		}

		// Bearer token格式
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

		// 解析Token
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
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

		// 将用户信息存入Context
		userId := int64(claims["uid"].(float64))
		r = r.WithContext(setUserId(r.Context(), userId))

		next(w, r)
	}
}

// Context key
type userKey struct{}

func setUserId(ctx interface{}, userId int64) interface{} {
	return context.WithValue(ctx, userKey{}, userId)
}

func GetUserId(ctx interface{}) int64 {
	if userId, ok := ctx.(interface{ Value(interface{}) interface{} }).Value(userKey{}).(int64); ok {
		return userId
	}
	return 0
}
