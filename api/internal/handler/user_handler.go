package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"

	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type CreateUserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreateUserHandler(svcCtx *svc.ServiceContext) *CreateUserHandler {
	return &CreateUserHandler{svcCtx: svcCtx}
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	// 检查用户名是否已存在
	_, err := h.svcCtx.UserModel.FindOneByUsername(r.Context(), req.Username)
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logx.Error("Failed to hash password:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	user := &model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Nickname:  sql.NullString{String: req.Nickname, Valid: req.Nickname != ""},
		Avatar:    sql.NullString{String: req.Avatar, Valid: req.Avatar != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Status:    int64(req.Status),
	}

	result, err := h.svcCtx.UserModel.Insert(r.Context(), user)
	if err != nil {
		logx.Error("Failed to insert user:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建用户失败",
		})
		return
	}

	// 获取插入的用户ID
	userId, _ := result.LastInsertId()

	// 如果有角色ID，关联角色
	if len(req.RoleIDs) > 0 {
		err = h.svcCtx.UserRoleModel.BatchInsert(r.Context(), userId, req.RoleIDs)
		if err != nil {
			logx.Error("Failed to bind roles:", err)
			httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "关联角色失败",
			})
			return
		}
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "创建成功",
	})
}

type UpdateUserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateUserHandler(svcCtx *svc.ServiceContext) *UpdateUserHandler {
	return &UpdateUserHandler{svcCtx: svcCtx}
}

func (h *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateUserRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	user, err := h.svcCtx.UserModel.FindOne(r.Context(), req.ID)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户不存在",
		})
		return
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = sql.NullString{String: req.Nickname, Valid: true}
	}
	if req.Email != "" {
		user.Email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.Phone != "" {
		user.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	if req.Avatar != "" {
		user.Avatar = sql.NullString{String: req.Avatar, Valid: true}
	}
	if req.Status > 0 {
		user.Status = int64(req.Status)
	}

	err = h.svcCtx.UserModel.Update(r.Context(), user)
	if err != nil {
		logx.Error("Failed to update user:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新用户失败",
		})
		return
	}

	// 更新角色关联
	if len(req.RoleIDs) > 0 {
		err = h.svcCtx.UserRoleModel.DeleteByUserId(r.Context(), req.ID)
		if err != nil {
			logx.Error("Failed to delete old roles:", err)
		}
		err = h.svcCtx.UserRoleModel.BatchInsert(r.Context(), req.ID, req.RoleIDs)
		if err != nil {
			logx.Error("Failed to bind roles:", err)
		}
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	})
}

type DeleteUserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteUserHandler(svcCtx *svc.ServiceContext) *DeleteUserHandler {
	return &DeleteUserHandler{svcCtx: svcCtx}
}

func (h *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 检查用户是否存在
	_, err = h.svcCtx.UserModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户不存在",
		})
		return
	}

	err = h.svcCtx.UserModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete user:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除用户失败",
		})
		return
	}

	// 删除用户角色关联
	err = h.svcCtx.UserRoleModel.DeleteByUserId(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete user roles:", err)
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

type GetUserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetUserHandler(svcCtx *svc.ServiceContext) *GetUserHandler {
	return &GetUserHandler{svcCtx: svcCtx}
}

func (h *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	user, err := h.svcCtx.UserModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户不存在",
		})
		return
	}

	// 获取用户角色
	roles, _ := h.svcCtx.RoleModel.FindByUserId(r.Context(), id)
	roleList := make([]map[string]interface{}, 0)
	for _, role := range roles {
		roleList = append(roleList, map[string]interface{}{
			"id":   role.Id,
			"name": role.Name,
			"code": role.Code,
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":         user.Id,
			"username":   user.Username,
			"nickname":   user.Nickname.String,
			"avatar":     user.Avatar.String,
			"email":      user.Email.String,
			"phone":      user.Phone.String,
			"status":     user.Status,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
			"roles":      roleList,
		},
	})
}

type ListUsersHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListUsersHandler(svcCtx *svc.ServiceContext) *ListUsersHandler {
	return &ListUsersHandler{svcCtx: svcCtx}
}

func (h *ListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.UserListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	users, total, err := h.svcCtx.UserModel.FindPage(r.Context(), page, pageSize, req.Username, req.Nickname, int64(req.Status))
	if err != nil {
		logx.Error("Failed to list users:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户列表失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, user := range users {
		list = append(list, map[string]interface{}{
			"id":         user.Id,
			"username":   user.Username,
			"nickname":   user.Nickname.String,
			"avatar":     user.Avatar.String,
			"email":      user.Email.String,
			"phone":      user.Phone.String,
			"status":     user.Status,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"total":    total,
			"page":     page,
			"page_size": pageSize,
			"list":     list,
		},
	})
}

type UpdatePasswordHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordHandler(svcCtx *svc.ServiceContext) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{svcCtx: svcCtx}
}

func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var req types.UpdatePasswordRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	user, err := h.svcCtx.UserModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户不存在",
		})
		return
	}

	// 验证旧密码
	if req.OldPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
				"code":    400,
				"message": "旧密码不正确",
			})
			return
		}
	}

	// 新密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logx.Error("Failed to hash password:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	user.Password = string(hashedPassword)
	err = h.svcCtx.UserModel.Update(r.Context(), user)
	if err != nil {
		logx.Error("Failed to update password:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新密码失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "密码更新成功",
	})
}
