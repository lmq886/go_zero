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

type CreatePermissionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionHandler(svcCtx *svc.ServiceContext) *CreatePermissionHandler {
	return &CreatePermissionHandler{svcCtx: svcCtx}
}

func (h *CreatePermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.CreatePermissionRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	permission := &model.Permission{
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		ParentId:  req.ParentID,
		Path:      sql.NullString{String: req.Path, Valid: req.Path != ""},
		Icon:      sql.NullString{String: req.Icon, Valid: req.Icon != ""},
		Component: sql.NullString{String: req.Component, Valid: req.Component != ""},
		Status:    int64(req.Status),
		Sort:      int64(req.Sort),
	}

	_, err := h.svcCtx.PermissionModel.Insert(r.Context(), permission)
	if err != nil {
		logx.Error("Failed to insert permission:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建权限失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "创建成功",
	})
}

type UpdatePermissionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionHandler(svcCtx *svc.ServiceContext) *UpdatePermissionHandler {
	return &UpdatePermissionHandler{svcCtx: svcCtx}
}

func (h *UpdatePermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的权限ID",
		})
		return
	}

	var req types.UpdatePermissionRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	permission, err := h.svcCtx.PermissionModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "权限不存在",
		})
		return
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.ParentID >= 0 {
		permission.ParentId = req.ParentID
	}
	if req.Path != "" {
		permission.Path = sql.NullString{String: req.Path, Valid: true}
	}
	if req.Icon != "" {
		permission.Icon = sql.NullString{String: req.Icon, Valid: true}
	}
	if req.Component != "" {
		permission.Component = sql.NullString{String: req.Component, Valid: true}
	}
	if req.Status > 0 {
		permission.Status = int64(req.Status)
	}
	if req.Sort >= 0 {
		permission.Sort = int64(req.Sort)
	}

	err = h.svcCtx.PermissionModel.Update(r.Context(), permission)
	if err != nil {
		logx.Error("Failed to update permission:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新权限失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	})
}

type DeletePermissionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionHandler(svcCtx *svc.ServiceContext) *DeletePermissionHandler {
	return &DeletePermissionHandler{svcCtx: svcCtx}
}

func (h *DeletePermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的权限ID",
		})
		return
	}

	_, err = h.svcCtx.PermissionModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "权限不存在",
		})
		return
	}

	err = h.svcCtx.PermissionModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete permission:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除权限失败",
		})
		return
	}

	// 删除角色权限关联
	h.svcCtx.RolePermissionModel.DeleteByPermissionId(r.Context(), id)

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

type GetPermissionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetPermissionHandler(svcCtx *svc.ServiceContext) *GetPermissionHandler {
	return &GetPermissionHandler{svcCtx: svcCtx}
}

func (h *GetPermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的权限ID",
		})
		return
	}

	permission, err := h.svcCtx.PermissionModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "权限不存在",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":         permission.Id,
			"name":       permission.Name,
			"code":       permission.Code,
			"type":       permission.Type,
			"parent_id":  permission.ParentId,
			"path":       permission.Path.String,
			"icon":       permission.Icon.String,
			"component":  permission.Component.String,
			"status":     permission.Status,
			"sort":       permission.Sort,
			"created_at": permission.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

type ListPermissionsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListPermissionsHandler(svcCtx *svc.ServiceContext) *ListPermissionsHandler {
	return &ListPermissionsHandler{svcCtx: svcCtx}
}

func (h *ListPermissionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.PermissionListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	permissions, total, err := h.svcCtx.PermissionModel.FindPage(r.Context(), page, pageSize, req.Name, req.Type, int64(req.Status))
	if err != nil {
		logx.Error("Failed to list permissions:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取权限列表失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, permission := range permissions {
		list = append(list, map[string]interface{}{
			"id":         permission.Id,
			"name":       permission.Name,
			"code":       permission.Code,
			"type":       permission.Type,
			"parent_id":  permission.ParentId,
			"path":       permission.Path.String,
			"icon":       permission.Icon.String,
			"component":  permission.Component.String,
			"status":     permission.Status,
			"sort":       permission.Sort,
			"created_at": permission.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": permission.UpdatedAt.Format("2006-01-02 15:04:05"),
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
