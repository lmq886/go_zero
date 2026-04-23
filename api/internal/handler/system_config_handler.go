package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type CreateSystemConfigHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreateSystemConfigHandler(svcCtx *svc.ServiceContext) *CreateSystemConfigHandler {
	return &CreateSystemConfigHandler{svcCtx: svcCtx}
}

func (h *CreateSystemConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.CreateSystemConfigRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	config := &model.SystemConfig{
		Key:    req.Key,
		Value:  req.Value,
		Name:   req.Name,
		Remark: sql.NullString{String: req.Remark, Valid: req.Remark != ""},
	}

	_, err := h.svcCtx.SystemConfigModel.Insert(r.Context(), config)
	if err != nil {
		logx.Error("Failed to insert config:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建配置失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "创建成功",
	})
}

type UpdateSystemConfigHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateSystemConfigHandler(svcCtx *svc.ServiceContext) *UpdateSystemConfigHandler {
	return &UpdateSystemConfigHandler{svcCtx: svcCtx}
}

func (h *UpdateSystemConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的配置ID",
		})
		return
	}

	var req types.UpdateSystemConfigRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	config, err := h.svcCtx.SystemConfigModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "配置不存在",
		})
		return
	}

	if req.Value != "" {
		config.Value = req.Value
	}
	if req.Name != "" {
		config.Name = req.Name
	}
	if req.Remark != "" {
		config.Remark = sql.NullString{String: req.Remark, Valid: true}
	}

	err = h.svcCtx.SystemConfigModel.Update(r.Context(), config)
	if err != nil {
		logx.Error("Failed to update config:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新配置失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	})
}

type DeleteSystemConfigHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteSystemConfigHandler(svcCtx *svc.ServiceContext) *DeleteSystemConfigHandler {
	return &DeleteSystemConfigHandler{svcCtx: svcCtx}
}

func (h *DeleteSystemConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的配置ID",
		})
		return
	}

	_, err = h.svcCtx.SystemConfigModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "配置不存在",
		})
		return
	}

	err = h.svcCtx.SystemConfigModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete config:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除配置失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

type GetSystemConfigHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetSystemConfigHandler(svcCtx *svc.ServiceContext) *GetSystemConfigHandler {
	return &GetSystemConfigHandler{svcCtx: svcCtx}
}

func (h *GetSystemConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的配置ID",
		})
		return
	}

	config, err := h.svcCtx.SystemConfigModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "配置不存在",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":         config.Id,
			"key":        config.Key,
			"value":      config.Value,
			"name":       config.Name,
			"remark":     config.Remark.String,
			"created_at": config.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": config.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

type ListSystemConfigsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListSystemConfigsHandler(svcCtx *svc.ServiceContext) *ListSystemConfigsHandler {
	return &ListSystemConfigsHandler{svcCtx: svcCtx}
}

func (h *ListSystemConfigsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.SystemConfigListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	configs, total, err := h.svcCtx.SystemConfigModel.FindPage(r.Context(), page, pageSize, req.Key, req.Name)
	if err != nil {
		logx.Error("Failed to list configs:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取配置列表失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, config := range configs {
		list = append(list, map[string]interface{}{
			"id":         config.Id,
			"key":        config.Key,
			"value":      config.Value,
			"name":       config.Name,
			"remark":     config.Remark.String,
			"created_at": config.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": config.UpdatedAt.Format("2006-01-02 15:04:05"),
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
