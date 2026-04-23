package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"

	"go_zero/api/internal/config"
	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type LoginHandler struct {
	svcCtx *svc.ServiceContext
}

func NewLoginHandler(svcCtx *svc.ServiceContext) *LoginHandler {
	return &LoginHandler{svcCtx: svcCtx}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	// 查询用户
	user, err := h.svcCtx.UserModel.FindOneByUsername(r.Context(), req.Username)
	if err != nil {
		logx.Error("User not found:", req.Username)
		h.recordLoginLog(r, 0, req.Username, 2, "用户不存在")
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logx.Error("Password mismatch for user:", req.Username)
		h.recordLoginLog(r, user.Id, req.Username, 2, "密码错误")
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		logx.Error("User is disabled:", req.Username)
		h.recordLoginLog(r, user.Id, req.Username, 2, "用户已禁用")
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户已禁用",
		})
		return
	}

	// 生成Token
	token, expiresAt, err := generateToken(user.Id, h.svcCtx.Config.JwtAuth)
	if err != nil {
		logx.Error("Failed to generate token:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 获取用户角色
	roles, err := h.svcCtx.RoleModel.FindByUserId(r.Context(), user.Id)
	if err != nil {
		logx.Error("Failed to get roles:", err)
	}

	roleIDs := make([]int64, 0)
	for _, role := range roles {
		roleIDs = append(roleIDs, role.Id)
	}

	// 记录登录成功日志
	h.recordLoginLog(r, user.Id, req.Username, 1, "登录成功")

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "登录成功",
		"data": map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
			"user_info": map[string]interface{}{
				"id":        user.Id,
				"username":  user.Username,
				"nickname":  user.Nickname.String,
				"avatar":    user.Avatar.String,
				"role_ids":  roleIDs,
			},
		},
	})
}

func generateToken(userId int64, jwtConfig config.JwtAuthConfig) (string, int64, error) {
	expiresAt := time.Now().Unix() + jwtConfig.AccessExpire
	claims := jwt.MapClaims{
		"uid":  userId,
		"exp":  expiresAt,
		"jti":  uuid.New().String(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.AccessSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

func (h *LoginHandler) recordLoginLog(r *http.Request, userId int64, username string, status int64, msg string) {
	log := &model.LoginLog{
		UserId:   sql.NullInt64{Int64: userId, Valid: userId > 0},
		Username: sql.NullString{String: username, Valid: username != ""},
		Ip:       sql.NullString{String: getClientIP(r), Valid: true},
		Status:   status,
		Msg:      sql.NullString{String: msg, Valid: true},
	}

	go func() {
		_, err := h.svcCtx.LoginLogModel.Insert(r.Context(), log)
		if err != nil {
			logx.Error("Failed to insert login log:", err)
		}
	}()
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
