package handler

import (
	"net/http"

	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthHandler 健康检查处理器
// 用于检查服务是否正常运行
// 参数 svcCtx: 服务上下文
// 返回值: HTTP 处理函数
func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 构建健康检查响应
		resp := types.HealthResp{
			Status:  "ok",
			Version: "v1.0.0",
			Uptime:  "service is running",
		}

		// 返回响应
		httpx.OkJson(w, resp)
	}
}
