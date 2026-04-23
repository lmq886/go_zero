package middleware

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
)

type OperationLogMiddleware struct {
	SvcCtx *svc.ServiceContext
}

func NewOperationLogMiddleware(svcCtx *svc.ServiceContext) *OperationLogMiddleware {
	return &OperationLogMiddleware{SvcCtx: svcCtx}
}

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 读取请求体
		var requestParams string
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				requestParams = string(body)
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}

		// 记录响应
		recorder := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		// 执行下一个中间件
		next(recorder, r)

		duration := time.Since(startTime).Milliseconds()

		// 获取用户ID
		userId := GetUserId(r.Context())

		// 获取用户名
		var username sql.NullString
		if userId > 0 {
			user, err := m.SvcCtx.UserModel.FindOne(r.Context(), userId)
			if err == nil {
				username = sql.NullString{String: user.Username, Valid: true}
			}
		}

		// 记录操作日志
		log := &model.OperationLog{
			UserId:        sql.NullInt64{Int64: userId, Valid: userId > 0},
			Username:      username,
			Operation:     getOperationDesc(r.Method, r.URL.Path),
			Method:        r.Method,
			RequestUri:    r.URL.Path,
			RequestParams: sql.NullString{String: requestParams, Valid: requestParams != ""},
			ResponseData:  sql.NullString{String: recorder.body.String(), Valid: recorder.body.Len() > 0},
			Ip:            sql.NullString{String: GetClientIP(r), Valid: true},
			Status:        1,
			Duration:      duration,
		}

		go func() {
			_, err := m.SvcCtx.OperationLogModel.Insert(r.Context(), log)
			if err != nil {
				logx.Error("Failed to insert operation log:", err)
			}
		}()
	}
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	// 如果是多个代理，取第一个IP
	if idx := bytes.IndexByte([]byte(ip), ','); idx != -1 {
		ip = string([]byte(ip)[:idx])
	}
	return ip
}

func getOperationDesc(method, path string) string {
	switch {
	case method == "POST" && path == "/api/v1/auth/login":
		return "用户登录"
	case method == "POST" && path == "/api/v1/auth/register":
		return "用户注册"
	case method == "POST" && path == "/api/v1/users":
		return "新增用户"
	case method == "PUT" && path == "/api/v1/users":
		return "更新用户"
	case method == "DELETE" && path == "/api/v1/users":
		return "删除用户"
	case method == "POST" && path == "/api/v1/roles":
		return "新增角色"
	case method == "PUT" && path == "/api/v1/roles":
		return "更新角色"
	case method == "DELETE" && path == "/api/v1/roles":
		return "删除角色"
	case method == "POST" && path == "/api/v1/permissions":
		return "新增权限"
	case method == "PUT" && path == "/api/v1/permissions":
		return "更新权限"
	case method == "DELETE" && path == "/api/v1/permissions":
		return "删除权限"
	case method == "POST" && path == "/api/v1/configs":
		return "新增配置"
	case method == "PUT" && path == "/api/v1/configs":
		return "更新配置"
	case method == "DELETE" && path == "/api/v1/configs":
		return "删除配置"
	case method == "POST" && path == "/api/v1/files/upload":
		return "上传文件"
	case method == "DELETE" && path == "/api/v1/files":
		return "删除文件"
	default:
		return method + " " + path
	}
}
