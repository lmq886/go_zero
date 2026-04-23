package middleware

import (
	"net/http"
	"strings"

	"d:\code\work\go_zero\api\internal\config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// PermissionMiddleware 权限验证中间件
// 用于验证用户是否有访问特定接口的权限
// 遵循 RBAC（基于角色的访问控制）模型
type PermissionMiddleware struct {
	// 配置信息
	config config.SystemConfig
	// 需要验证的权限编码
	requiredPermission string
}

// NewPermissionMiddleware 创建权限验证中间件实例
// 参数 config: 系统配置
// 参数 requiredPermission: 需要验证的权限编码
// 返回值: 权限验证中间件实例
func NewPermissionMiddleware(config config.SystemConfig, requiredPermission string) *PermissionMiddleware {
	return &PermissionMiddleware{
		config:             config,
		requiredPermission: requiredPermission,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文中获取用户角色
		roles, ok := GetRolesFromContext(r.Context())
		if !ok {
			logx.Errorf("Failed to get roles from context")
			httpx.Error(w, http.StatusForbidden, "禁止访问：无法获取用户角色信息")
			return
		}

		// 检查是否是超级管理员
		if m.isSuperAdmin(roles) {
			// 超级管理员拥有所有权限，直接放行
			next(w, r)
			return
		}

		// 如果不需要特定权限，直接放行
		if m.requiredPermission == "" {
			next(w, r)
			return
		}

		// 从上下文中获取用户权限
		permissions, ok := GetPermissionsFromContext(r.Context())
		if !ok {
			logx.Errorf("Failed to get permissions from context")
			httpx.Error(w, http.StatusForbidden, "禁止访问：无法获取用户权限信息")
			return
		}

		// 检查用户是否拥有所需权限
		if !m.hasPermission(permissions, m.requiredPermission) {
			logx.Errorf("User does not have required permission: %s", m.requiredPermission)
			httpx.Error(w, http.StatusForbidden, "禁止访问：没有足够的权限")
			return
		}

		// 权限验证通过，继续处理请求
		next(w, r)
	}
}

// isSuperAdmin 检查用户是否是超级管理员
// 参数 roles: 用户角色列表
// 返回值: 是否是超级管理员
func (m *PermissionMiddleware) isSuperAdmin(roles []string) bool {
	for _, role := range roles {
		if role == m.config.SuperAdminRoleCode {
			return true
		}
	}
	return false
}

// hasPermission 检查用户是否拥有指定权限
// 参数 permissions: 用户权限列表
// 参数 requiredPermission: 需要的权限
// 返回值: 是否拥有权限
func (m *PermissionMiddleware) hasPermission(permissions []string, requiredPermission string) bool {
	for _, permission := range permissions {
		// 支持通配符匹配
		// 例如：system:user:* 可以匹配 system:user:list, system:user:create 等
		if m.matchPermission(permission, requiredPermission) {
			return true
		}
	}
	return false
}

// matchPermission 权限匹配
// 支持通配符匹配
// 参数 userPermission: 用户拥有的权限
// 参数 requiredPermission: 需要的权限
// 返回值: 是否匹配
func (m *PermissionMiddleware) matchPermission(userPermission, requiredPermission string) bool {
	// 完全匹配
	if userPermission == requiredPermission {
		return true
	}

	// 通配符匹配
	// 例如：system:user:* 匹配 system:user:list
	if strings.HasSuffix(userPermission, ":*") {
		prefix := strings.TrimSuffix(userPermission, "*")
		if strings.HasPrefix(requiredPermission, prefix) {
			return true
		}
	}

	// 例如：system:* 匹配 system:user:list
	if strings.HasSuffix(userPermission, "*") && !strings.HasSuffix(userPermission, ":*") {
		prefix := strings.TrimSuffix(userPermission, "*")
		if strings.HasPrefix(requiredPermission, prefix) {
			return true
		}
	}

	return false
}

// RoleMiddleware 角色验证中间件
// 用于验证用户是否拥有特定角色
type RoleMiddleware struct {
	// 配置信息
	config config.SystemConfig
	// 需要验证的角色列表
	requiredRoles []string
	// 是否需要所有角色（true: 需要所有角色，false: 只需要其中一个）
	requireAll bool
}

// NewRoleMiddleware 创建角色验证中间件实例
// 参数 config: 系统配置
// 参数 requiredRoles: 需要验证的角色列表
// 参数 requireAll: 是否需要所有角色
// 返回值: 角色验证中间件实例
func NewRoleMiddleware(config config.SystemConfig, requiredRoles []string, requireAll bool) *RoleMiddleware {
	return &RoleMiddleware{
		config:        config,
		requiredRoles: requiredRoles,
		requireAll:    requireAll,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *RoleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文中获取用户角色
		roles, ok := GetRolesFromContext(r.Context())
		if !ok {
			logx.Errorf("Failed to get roles from context")
			httpx.Error(w, http.StatusForbidden, "禁止访问：无法获取用户角色信息")
			return
		}

		// 检查是否是超级管理员
		if m.isSuperAdmin(roles) {
			// 超级管理员拥有所有角色权限，直接放行
			next(w, r)
			return
		}

		// 检查角色
		if m.requireAll {
			// 需要所有角色
			if !m.hasAllRoles(roles, m.requiredRoles) {
				logx.Errorf("User does not have all required roles: %v", m.requiredRoles)
				httpx.Error(w, http.StatusForbidden, "禁止访问：没有足够的角色权限")
				return
			}
		} else {
			// 只需要其中一个角色
			if !m.hasAnyRole(roles, m.requiredRoles) {
				logx.Errorf("User does not have any of the required roles: %v", m.requiredRoles)
				httpx.Error(w, http.StatusForbidden, "禁止访问：没有足够的角色权限")
				return
			}
		}

		// 角色验证通过，继续处理请求
		next(w, r)
	}
}

// isSuperAdmin 检查用户是否是超级管理员
// 参数 roles: 用户角色列表
// 返回值: 是否是超级管理员
func (m *RoleMiddleware) isSuperAdmin(roles []string) bool {
	for _, role := range roles {
		if role == m.config.SuperAdminRoleCode {
			return true
		}
	}
	return false
}

// hasAllRoles 检查用户是否拥有所有指定角色
// 参数 userRoles: 用户角色列表
// 参数 requiredRoles: 需要的角色列表
// 返回值: 是否拥有所有角色
func (m *RoleMiddleware) hasAllRoles(userRoles, requiredRoles []string) bool {
	userRoleMap := make(map[string]bool)
	for _, role := range userRoles {
		userRoleMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if !userRoleMap[requiredRole] {
			return false
		}
	}

	return true
}

// hasAnyRole 检查用户是否拥有任意一个指定角色
// 参数 userRoles: 用户角色列表
// 参数 requiredRoles: 需要的角色列表
// 返回值: 是否拥有任意一个角色
func (m *RoleMiddleware) hasAnyRole(userRoles, requiredRoles []string) bool {
	userRoleMap := make(map[string]bool)
	for _, role := range userRoles {
		userRoleMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if userRoleMap[requiredRole] {
			return true
		}
	}

	return false
}
