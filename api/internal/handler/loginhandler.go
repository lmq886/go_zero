package handler

import (
	"net/http"

	"go_zero/api/internal/logic"
	"go_zero/api/internal/middleware"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LoginHandler 登录处理器
// 处理用户登录请求
// 参数 svcCtx: 服务上下文
// 返回值: HTTP 处理函数
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// 创建业务逻辑实例
		l := logic.NewLoginLogic(r.Context(), svcCtx)

		// 执行登录逻辑
		resp, err := l.Login(&req)
		if err != nil {
			middleware.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// 返回成功响应
		httpx.OkJson(w, resp)
	}
}

// RegisterHandler 注册处理器
// 处理用户注册请求
// 参数 svcCtx: 服务上下文
// 返回值: HTTP 处理函数
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// 创建业务逻辑实例
		l := logic.NewRegisterLogic(r.Context(), svcCtx)

		// 执行注册逻辑
		resp, err := l.Register(&req)
		if err != nil {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// 返回成功响应
		httpx.OkJson(w, resp)
	}
}

// LogoutHandler 登出处理器
// 处理用户登出请求
// 参数 svcCtx: 服务上下文
// 返回值: HTTP 处理函数
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.LogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// 创建业务逻辑实例
		l := logic.NewLogoutLogic(r.Context(), svcCtx)

		// 执行登出逻辑
		resp, err := l.Logout(&req)
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// 返回成功响应
		httpx.OkJson(w, resp)
	}
}

// RefreshTokenHandler 刷新令牌处理器
// 处理令牌刷新请求
// 参数 svcCtx: 服务上下文
// 返回值: HTTP 处理函数
func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.RefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// 创建业务逻辑实例
		l := logic.NewRefreshTokenLogic(r.Context(), svcCtx)

		// 执行令牌刷新逻辑
		resp, err := l.RefreshToken(&req)
		if err != nil {
			middleware.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// 返回成功响应
		httpx.OkJson(w, resp)
	}
}
