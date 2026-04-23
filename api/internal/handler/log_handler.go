package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type ListOperationLogsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListOperationLogsHandler(svcCtx *svc.ServiceContext) *ListOperationLogsHandler {
	return &ListOperationLogsHandler{svcCtx: svcCtx}
}

func (h *ListOperationLogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.OperationLogListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	logs, total, err := h.svcCtx.OperationLogModel.FindPage(r.Context(), page, pageSize,
		req.Username, req.Operation, req.Method, int64(req.Status), req.StartTime, req.EndTime)
	if err != nil {
		logx.Error("Failed to list operation logs:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取操作日志失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, log := range logs {
		list = append(list, map[string]interface{}{
			"id":             log.Id,
			"user_id":        log.UserId.Int64,
			"username":       log.Username.String,
			"operation":      log.Operation,
			"method":         log.Method,
			"request_uri":    log.RequestUri,
			"request_params": log.RequestParams.String,
			"response_data":  log.ResponseData.String,
			"ip":             log.Ip.String,
			"location":       log.Location.String,
			"browser":        log.Browser.String,
			"os":             log.Os.String,
			"status":         log.Status,
			"error_msg":      log.ErrorMsg.String,
			"duration":       log.Duration,
			"created_at":     log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"list":      list,
		},
	})
}

type DeleteOperationLogHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteOperationLogHandler(svcCtx *svc.ServiceContext) *DeleteOperationLogHandler {
	return &DeleteOperationLogHandler{svcCtx: svcCtx}
}

func (h *DeleteOperationLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的日志ID",
		})
		return
	}

	err = h.svcCtx.OperationLogModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete operation log:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除操作日志失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

type ListLoginLogsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListLoginLogsHandler(svcCtx *svc.ServiceContext) *ListLoginLogsHandler {
	return &ListLoginLogsHandler{svcCtx: svcCtx}
}

func (h *ListLoginLogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.LoginLogListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	logs, total, err := h.svcCtx.LoginLogModel.FindPage(r.Context(), page, pageSize,
		req.Username, int64(req.Status), req.StartTime, req.EndTime)
	if err != nil {
		logx.Error("Failed to list login logs:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取登录日志失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, log := range logs {
		list = append(list, map[string]interface{}{
			"id":         log.Id,
			"user_id":    log.UserId.Int64,
			"username":   log.Username.String,
			"ip":         log.Ip.String,
			"location":   log.Location.String,
			"browser":    log.Browser.String,
			"os":         log.Os.String,
			"status":     log.Status,
			"msg":        log.Msg.String,
			"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"list":      list,
		},
	})
}

type DeleteLoginLogHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteLoginLogHandler(svcCtx *svc.ServiceContext) *DeleteLoginLogHandler {
	return &DeleteLoginLogHandler{svcCtx: svcCtx}
}

func (h *DeleteLoginLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的日志ID",
		})
		return
	}

	err = h.svcCtx.LoginLogModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete login log:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除登录日志失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}
