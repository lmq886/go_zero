package routes

import (
	"d:\code\work\go_zero\api\internal\handler"
	"d:\code\work\go_zero\api\internal\middleware"
	"d:\code\work\go_zero\api\internal\svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 注册所有路由处理器
// 参数 server: REST 服务器实例
// 参数 serverCtx: 服务上下文
// 遵循 GoZero 最佳实践，统一管理路由和中间件
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// ==================== 全局中间件 ====================
	// CORS 跨域中间件（应用于所有路由）
	server.Use(
		middleware.NewCORSMiddleware(serverCtx.Config.CORS).Handle,
	)

	// 错误处理中间件（应用于所有路由）
	server.Use(
		middleware.NewErrorMiddleware().Handle,
	)

	// 限流中间件（应用于所有路由）
	if serverCtx.Config.RateLimit.Enabled {
		server.Use(
			middleware.NewRateLimitMiddleware(serverCtx.Config.RateLimit).Handle,
		)
	}

	// ==================== 公开路由（不需要认证） ====================
	// 认证模块 - 公开路由
	server.AddRoutes(
		[]rest.Route{
			{
				// 登录
				Method:  "POST",
				Path:    "/api/v1/auth/login",
				Handler: handler.LoginHandler(serverCtx),
			},
			{
				// 注册
				Method:  "POST",
				Path:    "/api/v1/auth/register",
				Handler: handler.RegisterHandler(serverCtx),
			},
			{
				// 刷新令牌
				Method:  "POST",
				Path:    "/api/v1/auth/refresh",
				Handler: handler.RefreshTokenHandler(serverCtx),
			},
		},
	)

	// ==================== 需要认证的路由 ====================
	// 创建 JWT 认证中间件
	jwtAuthMiddleware := middleware.NewJwtAuthMiddleware(serverCtx.Config.Auth)

	// 创建操作日志中间件
	operationLogMiddleware := middleware.NewOperationLogMiddleware(serverCtx)

	// 添加需要认证的路由组
	server.AddRoutes(
		[]rest.Route{
			// ==================== 认证模块 - 需要认证的路由 ====================
			{
				// 登出
				Method:  "POST",
				Path:    "/api/v1/auth/logout",
				Handler: handler.LogoutHandler(serverCtx),
			},

			// ==================== 用户管理模块 ====================
			{
				// 获取用户列表
				Method:  "GET",
				Path:    "/api/v1/users",
				Handler: handler.GetUserListHandler(serverCtx),
			},
			{
				// 根据ID获取用户
				Method:  "GET",
				Path:    "/api/v1/users/:id",
				Handler: handler.GetUserByIdHandler(serverCtx),
			},
			{
				// 创建用户
				Method:  "POST",
				Path:    "/api/v1/users",
				Handler: handler.CreateUserHandler(serverCtx),
			},
			{
				// 更新用户
				Method:  "PUT",
				Path:    "/api/v1/users/:id",
				Handler: handler.UpdateUserHandler(serverCtx),
			},
			{
				// 删除用户
				Method:  "DELETE",
				Path:    "/api/v1/users/:id",
				Handler: handler.DeleteUserHandler(serverCtx),
			},
			{
				// 重置密码
				Method:  "POST",
				Path:    "/api/v1/users/:id/reset-password",
				Handler: handler.ResetPasswordHandler(serverCtx),
			},
			{
				// 获取个人资料
				Method:  "GET",
				Path:    "/api/v1/users/profile",
				Handler: handler.GetProfileHandler(serverCtx),
			},
			{
				// 更新个人资料
				Method:  "PUT",
				Path:    "/api/v1/users/profile",
				Handler: handler.UpdateProfileHandler(serverCtx),
			},

			// ==================== 角色管理模块 ====================
			{
				// 获取角色列表
				Method:  "GET",
				Path:    "/api/v1/roles",
				Handler: handler.GetRoleListHandler(serverCtx),
			},
			{
				// 根据ID获取角色
				Method:  "GET",
				Path:    "/api/v1/roles/:id",
				Handler: handler.GetRoleByIdHandler(serverCtx),
			},
			{
				// 创建角色
				Method:  "POST",
				Path:    "/api/v1/roles",
				Handler: handler.CreateRoleHandler(serverCtx),
			},
			{
				// 更新角色
				Method:  "PUT",
				Path:    "/api/v1/roles/:id",
				Handler: handler.UpdateRoleHandler(serverCtx),
			},
			{
				// 删除角色
				Method:  "DELETE",
				Path:    "/api/v1/roles/:id",
				Handler: handler.DeleteRoleHandler(serverCtx),
			},
			{
				// 分配权限
				Method:  "POST",
				Path:    "/api/v1/roles/:id/permissions",
				Handler: handler.AssignPermissionsHandler(serverCtx),
			},
			{
				// 获取角色权限
				Method:  "GET",
				Path:    "/api/v1/roles/:id/permissions",
				Handler: handler.GetRolePermissionsHandler(serverCtx),
			},

			// ==================== 权限管理模块 ====================
			{
				// 获取权限列表
				Method:  "GET",
				Path:    "/api/v1/permissions",
				Handler: handler.GetPermissionListHandler(serverCtx),
			},
			{
				// 根据ID获取权限
				Method:  "GET",
				Path:    "/api/v1/permissions/:id",
				Handler: handler.GetPermissionByIdHandler(serverCtx),
			},
			{
				// 创建权限
				Method:  "POST",
				Path:    "/api/v1/permissions",
				Handler: handler.CreatePermissionHandler(serverCtx),
			},
			{
				// 更新权限
				Method:  "PUT",
				Path:    "/api/v1/permissions/:id",
				Handler: handler.UpdatePermissionHandler(serverCtx),
			},
			{
				// 删除权限
				Method:  "DELETE",
				Path:    "/api/v1/permissions/:id",
				Handler: handler.DeletePermissionHandler(serverCtx),
			},

			// ==================== 菜单管理模块 ====================
			{
				// 获取菜单列表
				Method:  "GET",
				Path:    "/api/v1/menus",
				Handler: handler.GetMenuListHandler(serverCtx),
			},
			{
				// 根据ID获取菜单
				Method:  "GET",
				Path:    "/api/v1/menus/:id",
				Handler: handler.GetMenuByIdHandler(serverCtx),
			},
			{
				// 创建菜单
				Method:  "POST",
				Path:    "/api/v1/menus",
				Handler: handler.CreateMenuHandler(serverCtx),
			},
			{
				// 更新菜单
				Method:  "PUT",
				Path:    "/api/v1/menus/:id",
				Handler: handler.UpdateMenuHandler(serverCtx),
			},
			{
				// 删除菜单
				Method:  "DELETE",
				Path:    "/api/v1/menus/:id",
				Handler: handler.DeleteMenuHandler(serverCtx),
			},
			{
				// 获取用户菜单
				Method:  "GET",
				Path:    "/api/v1/menus/user",
				Handler: handler.GetUserMenusHandler(serverCtx),
			},

			// ==================== 日志管理模块 ====================
			{
				// 获取操作日志列表
				Method:  "GET",
				Path:    "/api/v1/logs/operation",
				Handler: handler.GetOperationLogListHandler(serverCtx),
			},
			{
				// 获取登录日志列表
				Method:  "GET",
				Path:    "/api/v1/logs/login",
				Handler: handler.GetLoginLogListHandler(serverCtx),
			},
			{
				// 删除操作日志
				Method:  "DELETE",
				Path:    "/api/v1/logs/operation/:id",
				Handler: handler.DeleteOperationLogHandler(serverCtx),
			},
			{
				// 删除登录日志
				Method:  "DELETE",
				Path:    "/api/v1/logs/login/:id",
				Handler: handler.DeleteLoginLogHandler(serverCtx),
			},
			{
				// 清空操作日志
				Method:  "DELETE",
				Path:    "/api/v1/logs/operation",
				Handler: handler.ClearOperationLogsHandler(serverCtx),
			},
			{
				// 清空登录日志
				Method:  "DELETE",
				Path:    "/api/v1/logs/login",
				Handler: handler.ClearLoginLogsHandler(serverCtx),
			},

			// ==================== 系统配置模块 ====================
			{
				// 获取配置列表
				Method:  "GET",
				Path:    "/api/v1/configs",
				Handler: handler.GetConfigListHandler(serverCtx),
			},
			{
				// 根据键获取配置
				Method:  "GET",
				Path:    "/api/v1/configs/:key",
				Handler: handler.GetConfigByKeyHandler(serverCtx),
			},
			{
				// 创建配置
				Method:  "POST",
				Path:    "/api/v1/configs",
				Handler: handler.CreateConfigHandler(serverCtx),
			},
			{
				// 更新配置
				Method:  "PUT",
				Path:    "/api/v1/configs/:key",
				Handler: handler.UpdateConfigHandler(serverCtx),
			},
			{
				// 删除配置
				Method:  "DELETE",
				Path:    "/api/v1/configs/:key",
				Handler: handler.DeleteConfigHandler(serverCtx),
			},
		},
		// 应用 JWT 认证中间件
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		// 应用 JWT 认证中间件（自定义）
		rest.WithMiddlewares(
			// JWT 认证中间件
			jwtAuthMiddleware.Handle,
			// 操作日志中间件
			operationLogMiddleware.Handle,
		),
	)

	// ==================== 健康检查路由（公开） ====================
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "GET",
				Path:    "/health",
				Handler: handler.HealthHandler(serverCtx),
			},
		},
	)
}
