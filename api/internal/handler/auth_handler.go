package handler

import (
	"database/sql"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"

	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type RegisterHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRegisterHandler(svcCtx *svc.ServiceContext) *RegisterHandler {
	return &RegisterHandler{svcCtx: svcCtx}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
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
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Status:    1,
	}

	_, err = h.svcCtx.UserModel.Insert(r.Context(), user)
	if err != nil {
		logx.Error("Failed to insert user:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "注册失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "注册成功",
	})
}

type LogoutHandler struct {
	svcCtx *svc.ServiceContext
}

func NewLogoutHandler(svcCtx *svc.ServiceContext) *LogoutHandler {
	return &LogoutHandler{svcCtx: svcCtx}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 简单实现，实际项目中可能需要处理Token黑名单等逻辑
	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "退出成功",
	})
}

type GetUserInfoHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoHandler(svcCtx *svc.ServiceContext) *GetUserInfoHandler {
	return &GetUserInfoHandler{svcCtx: svcCtx}
}

func (h *GetUserInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 从Context获取用户ID
	userId := GetUserIdFromCtx(r.Context())
	if userId == 0 {
		httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	user, err := h.svcCtx.UserModel.FindOne(r.Context(), userId)
	if err != nil {
		logx.Error("Failed to get user:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	// 获取用户角色
	roles, err := h.svcCtx.RoleModel.FindByUserId(r.Context(), userId)
	if err != nil {
		logx.Error("Failed to get roles:", err)
	}

	roleList := make([]map[string]interface{}, 0)
	for _, role := range roles {
		roleList = append(roleList, map[string]interface{}{
			"id":   role.Id,
			"name": role.Name,
			"code": role.Code,
		})
	}

	// 获取用户权限
	permissions, err := h.svcCtx.PermissionModel.FindByUserId(r.Context(), userId)
	if err != nil {
		logx.Error("Failed to get permissions:", err)
	}

	permissionList := make([]map[string]interface{}, 0)
	for _, permission := range permissions {
		permissionList = append(permissionList, map[string]interface{}{
			"id":       permission.Id,
			"name":     permission.Name,
			"code":     permission.Code,
			"type":     permission.Type,
			"parent_id": permission.ParentId,
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":          user.Id,
			"username":    user.Username,
			"nickname":    user.Nickname.String,
			"avatar":      user.Avatar.String,
			"email":       user.Email.String,
			"phone":       user.Phone.String,
			"status":      user.Status,
			"created_at":  user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  user.UpdatedAt.Format("2006-01-02 15:04:05"),
			"roles":       roleList,
			"permissions": permissionList,
		},
	})
}

// GetUserIdFromCtx 从Context获取用户ID
func GetUserIdFromCtx(ctx interface{}) int64 {
	type contextKey string
	key := contextKey("userId")
	if val, ok := ctx.(interface{ Value(interface{}) interface{} }).Value(key).(int64); ok {
		return val
	}
	return 0
}
