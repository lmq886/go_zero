/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\handler\login_handler.go
 * @Description: 登录接口控制器
 */
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

// LoginHandler 登录接口控制器结构体
type LoginHandler struct {
	svcCtx *svc.ServiceContext // 服务上下文
}

// NewLoginHandler 创建登录接口控制器实例
func NewLoginHandler(svcCtx *svc.ServiceContext) *LoginHandler {
	return &LoginHandler{svcCtx: svcCtx}
}

// ServeHTTP 处理登录请求
// 参数: w - HTTP响应写入器
// 参数: r - HTTP请求对象
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.LoginRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	// 2. 根据用户名查询用户
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

	// 3. 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logx.Error("Password mismatch for user:", req.Username)
		h.recordLoginLog(r, user.Id, req.Username, 2, "密码错误")
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户名或密码错误",
		})
		return
	}

	// 4. 检查用户状态是否正常
	if user.Status != 1 {
		logx.Error("User is disabled:", req.Username)
		h.recordLoginLog(r, user.Id, req.Username, 2, "用户已禁用")
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "用户已禁用",
		})
		return
	}

	// 5. 生成JWT Token
	token, expiresAt, err := generateToken(user.Id, h.svcCtx.Config.JwtAuth)
	if err != nil {
		logx.Error("Failed to generate token:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 6. 获取用户角色信息
	roles, err := h.svcCtx.RoleModel.FindByUserId(r.Context(), user.Id)
	if err != nil {
		logx.Error("Failed to get roles:", err)
	}

	roleIDs := make([]int64, 0)
	for _, role := range roles {
		roleIDs = append(roleIDs, role.Id)
	}

	// 7. 记录登录成功日志
	h.recordLoginLog(r, user.Id, req.Username, 1, "登录成功")

	// 8. 返回登录成功响应
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

// generateToken 生成JWT Token
// 参数: userId - 用户ID
// 参数: jwtConfig - JWT配置
// 返回: string - Token字符串
// 返回: int64 - 过期时间戳
// 返回: error - 错误信息
func generateToken(userId int64, jwtConfig config.JwtAuthConfig) (string, int64, error) {
	expiresAt := time.Now().Unix() + jwtConfig.AccessExpire
	claims := jwt.MapClaims{
		"uid":  userId,       // 用户ID
		"exp":  expiresAt,    // 过期时间
		"jti":  uuid.New().String(), // JWT ID
		"iat":  time.Now().Unix(),   // 签发时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.AccessSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// recordLoginLog 记录登录日志
// 参数: r - HTTP请求对象
// 参数: userId - 用户ID
// 参数: username - 用户名
// 参数: status - 登录状态（1:成功，2:失败）
// 参数: msg - 登录信息
func (h *LoginHandler) recordLoginLog(r *http.Request, userId int64, username string, status int64, msg string) {
	log := &model.LoginLog{
		UserId:   sql.NullInt64{Int64: userId, Valid: userId > 0},
		Username: sql.NullString{String: username, Valid: username != ""},
		Ip:       sql.NullString{String: getClientIP(r), Valid: true},
		Status:   status,
		Msg:      sql.NullString{String: msg, Valid: true},
	}

	// 异步记录日志
	go func() {
		_, err := h.svcCtx.LoginLogModel.Insert(r.Context(), log)
		if err != nil {
			logx.Error("Failed to insert login log:", err)
		}
	}()
}

// getClientIP 获取客户端IP地址
// 参数: r - HTTP请求对象
// 返回: string - 客户端IP地址
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
