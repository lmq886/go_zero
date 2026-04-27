package routes

import (
	"go_zero/api/internal/handler"
	"go_zero/api/internal/middleware"
	"go_zero/api/internal/svc"

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
