/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\middleware\permission_auth.go
 * @Description: 权限验证中间件，用于检查用户是否有权限访问资源
 */
package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/svc"
)

// PermissionAuthMiddleware 权限验证中间件结构体
type PermissionAuthMiddleware struct {
	SvcCtx *svc.ServiceContext // 服务上下文，用于数据库操作
}

// NewPermissionAuthMiddleware 创建权限验证中间件实例
func NewPermissionAuthMiddleware(svcCtx *svc.ServiceContext) *PermissionAuthMiddleware {
	return &PermissionAuthMiddleware{SvcCtx: svcCtx}
}

// Handle 中间件处理函数
// 参数: next - 下一个处理函数
// 返回: http.HandlerFunc - 包装后的处理函数
func (m *PermissionAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 从Context中获取用户ID
		userId := GetUserId(r.Context())
		if userId == 0 {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未登录",
			})
			return
		}

		// 2. 获取用户拥有的所有权限
		permissions, err := m.SvcCtx.PermissionModel.FindByUserId(r.Context(), userId)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取权限失败",
			})
			return
		}

		// 3. 获取当前请求的路径和方法
		path := r.URL.Path
		method := r.Method

		// 4. 检查用户是否拥有超级管理员角色
		roles, err := m.SvcCtx.RoleModel.FindByUserId(r.Context(), userId)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取角色失败",
			})
			return
		}

		// 5. 验证是否为超级管理员
		isSuperAdmin := false
		for _, role := range roles {
			if role.Code == "super_admin" {
				isSuperAdmin = true
				break
			}
		}

		// 超级管理员拥有所有权限，直接放行
		if isSuperAdmin {
			next(w, r)
			return
		}

		// 6. 检查用户是否拥有访问当前资源的权限
		hasPermission := false
		for _, permission := range permissions {
			// 简单的权限匹配逻辑
			// 实际项目中可以实现更复杂的权限匹配规则
			if permission.Type == "api" {
				// 可以根据permission.Code与当前请求路径和方法进行匹配
				hasPermission = true
				break
			}
		}

		// 没有权限则返回403
		if !hasPermission {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "没有权限访问此资源",
			})
			return
		}

		// 权限验证通过，继续执行下一个处理函数
		next(w, r)
	}
}

// contextWithValue 将值存入Context
// 参数: ctx - 上下文对象
// 参数: key - 键
// 参数: val - 值
// 返回: context.Context - 新的上下文对象
func contextWithValue(ctx context.Context, key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
