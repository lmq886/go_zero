package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/svc"
)

type PermissionAuthMiddleware struct {
	SvcCtx *svc.ServiceContext
}

func NewPermissionAuthMiddleware(svcCtx *svc.ServiceContext) *PermissionAuthMiddleware {
	return &PermissionAuthMiddleware{SvcCtx: svcCtx}
}

func (m *PermissionAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Context获取用户ID
		userId := GetUserId(r.Context())
		if userId == 0 {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未登录",
			})
			return
		}

		// 获取用户权限
		permissions, err := m.SvcCtx.PermissionModel.FindByUserId(r.Context(), userId)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取权限失败",
			})
			return
		}

		// 获取当前请求路径
		path := r.URL.Path
		method := r.Method

		// 检查是否有超级管理员角色
		roles, err := m.SvcCtx.RoleModel.FindByUserId(r.Context(), userId)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取角色失败",
			})
			return
		}

		// 检查是否为超级管理员
		isSuperAdmin := false
		for _, role := range roles {
			if role.Code == "super_admin" {
				isSuperAdmin = true
				break
			}
		}

		if isSuperAdmin {
			next(w, r)
			return
		}

		// 检查权限
		hasPermission := false
		for _, permission := range permissions {
			// 简单的权限匹配逻辑
			if permission.Type == "api" {
				// 可以根据实际需求实现更复杂的权限匹配
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "没有权限访问此资源",
			})
			return
		}

		next(w, r)
	}
}

func contextWithValue(ctx context.Context, key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
