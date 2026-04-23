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

type CreateRoleHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreateRoleHandler(svcCtx *svc.ServiceContext) *CreateRoleHandler {
	return &CreateRoleHandler{svcCtx: svcCtx}
}

func (h *CreateRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.CreateRoleRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Status:      int64(req.Status),
		Sort:        int64(req.Sort),
	}

	result, err := h.svcCtx.RoleModel.Insert(r.Context(), role)
	if err != nil {
		logx.Error("Failed to insert role:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建角色失败",
		})
		return
	}

	roleId, _ := result.LastInsertId()

	// 关联权限
	if len(req.PermissionIDs) > 0 {
		err = h.svcCtx.RolePermissionModel.BatchInsert(r.Context(), roleId, req.PermissionIDs)
		if err != nil {
			logx.Error("Failed to bind permissions:", err)
		}
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "创建成功",
	})
}

type UpdateRoleHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleHandler(svcCtx *svc.ServiceContext) *UpdateRoleHandler {
	return &UpdateRoleHandler{svcCtx: svcCtx}
}

func (h *UpdateRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	var req types.UpdateRoleRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	role, err := h.svcCtx.RoleModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "角色不存在",
		})
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.Status > 0 {
		role.Status = int64(req.Status)
	}
	if req.Sort >= 0 {
		role.Sort = int64(req.Sort)
	}

	err = h.svcCtx.RoleModel.Update(r.Context(), role)
	if err != nil {
		logx.Error("Failed to update role:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新角色失败",
		})
		return
	}

	// 更新权限关联
	if len(req.PermissionIDs) > 0 {
		err = h.svcCtx.RolePermissionModel.DeleteByRoleId(r.Context(), id)
		if err != nil {
			logx.Error("Failed to delete old permissions:", err)
		}
		err = h.svcCtx.RolePermissionModel.BatchInsert(r.Context(), id, req.PermissionIDs)
		if err != nil {
			logx.Error("Failed to bind permissions:", err)
		}
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	})
}

type DeleteRoleHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleHandler(svcCtx *svc.ServiceContext) *DeleteRoleHandler {
	return &DeleteRoleHandler{svcCtx: svcCtx}
}

func (h *DeleteRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	role, err := h.svcCtx.RoleModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "角色不存在",
		})
		return
	}

	// 不能删除超级管理员角色
	if role.Code == "super_admin" {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "不能删除超级管理员角色",
		})
		return
	}

	err = h.svcCtx.RoleModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete role:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除角色失败",
		})
		return
	}

	// 删除角色权限关联和用户角色关联
	h.svcCtx.RolePermissionModel.DeleteByRoleId(r.Context(), id)
	h.svcCtx.UserRoleModel.DeleteByRoleId(r.Context(), id)

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

type GetRoleHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetRoleHandler(svcCtx *svc.ServiceContext) *GetRoleHandler {
	return &GetRoleHandler{svcCtx: svcCtx}
}

func (h *GetRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	role, err := h.svcCtx.RoleModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "角色不存在",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":          role.Id,
			"name":        role.Name,
			"code":        role.Code,
			"description": role.Description.String,
			"status":      role.Status,
			"sort":        role.Sort,
			"created_at":  role.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  role.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

type ListRolesHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListRolesHandler(svcCtx *svc.ServiceContext) *ListRolesHandler {
	return &ListRolesHandler{svcCtx: svcCtx}
}

func (h *ListRolesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.RoleListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	roles, total, err := h.svcCtx.RoleModel.FindPage(r.Context(), page, pageSize, req.Name, int64(req.Status))
	if err != nil {
		logx.Error("Failed to list roles:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取角色列表失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, role := range roles {
		list = append(list, map[string]interface{}{
			"id":          role.Id,
			"name":        role.Name,
			"code":        role.Code,
			"description": role.Description.String,
			"status":      role.Status,
			"sort":        role.Sort,
			"created_at":  role.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  role.UpdatedAt.Format("2006-01-02 15:04:05"),
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
