/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\middleware\operation_log.go
 * @Description: 操作日志中间件，用于记录用户操作行为
 */
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

// OperationLogMiddleware 操作日志中间件结构体
type OperationLogMiddleware struct {
	SvcCtx *svc.ServiceContext // 服务上下文，用于数据库操作
}

// NewOperationLogMiddleware 创建操作日志中间件实例
func NewOperationLogMiddleware(svcCtx *svc.ServiceContext) *OperationLogMiddleware {
	return &OperationLogMiddleware{SvcCtx: svcCtx}
}

// responseRecorder 自定义响应记录器，用于捕获响应内容
type responseRecorder struct {
	http.ResponseWriter // 原始响应写入器
	body *bytes.Buffer   // 用于存储响应内容
}

// Write 重写Write方法，同时将响应内容写入缓冲区
// 参数: b - 响应字节数据
// 返回: int - 写入的字节数
// 返回: error - 错误信息
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Handle 中间件处理函数
// 参数: next - 下一个处理函数
// 返回: http.HandlerFunc - 包装后的处理函数
func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 记录请求开始时间
		startTime := time.Now()

		// 2. 读取并记录请求体参数
		var requestParams string
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				requestParams = string(body)
				// 重置请求体，以便后续处理
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}

		// 3. 创建响应记录器，用于捕获响应内容
		recorder := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		// 4. 执行下一个中间件或处理函数
		next(recorder, r)

		// 5. 计算请求处理时长（毫秒）
		duration := time.Since(startTime).Milliseconds()

		// 6. 从Context中获取用户ID
		userId := GetUserId(r.Context())

		// 7. 获取用户名
		var username sql.NullString
		if userId > 0 {
			user, err := m.SvcCtx.UserModel.FindOne(r.Context(), userId)
			if err == nil {
				username = sql.NullString{String: user.Username, Valid: true}
			}
		}

		// 8. 构建操作日志记录
		log := &model.OperationLog{
			UserId:        sql.NullInt64{Int64: userId, Valid: userId > 0},
			Username:      username,
			Operation:     getOperationDesc(r.Method, r.URL.Path), // 操作描述
			Method:        r.Method,        // HTTP方法
			RequestUri:    r.URL.Path,      // 请求路径
			RequestParams: sql.NullString{String: requestParams, Valid: requestParams != ""}, // 请求参数
			ResponseData:  sql.NullString{String: recorder.body.String(), Valid: recorder.body.Len() > 0}, // 响应数据
			Ip:            sql.NullString{String: GetClientIP(r), Valid: true}, // 客户端IP
			Status:        1,              // 操作状态（1:成功）
			Duration:      duration,       // 处理时长
		}

		// 9. 异步记录操作日志到数据库
		go func() {
			_, err := m.SvcCtx.OperationLogModel.Insert(r.Context(), log)
			if err != nil {
				logx.Error("Failed to insert operation log:", err)
			}
		}()
	}
}

// GetClientIP 获取客户端真实IP地址
// 参数: r - HTTP请求对象
// 返回: string - 客户端IP地址
func GetClientIP(r *http.Request) string {
	// 优先从X-Forwarded-For获取（代理场景）
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		// 从X-Real-IP获取
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		// 从RemoteAddr获取
		ip = r.RemoteAddr
	}
	// 如果是多个代理，取第一个IP
	if idx := bytes.IndexByte([]byte(ip), ','); idx != -1 {
		ip = string([]byte(ip)[:idx])
	}
	return ip
}

// getOperationDesc 根据请求方法和路径生成操作描述
// 参数: method - HTTP方法
// 参数: path - 请求路径
// 返回: string - 操作描述
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
		// 默认返回方法和路径
		return method + " " + path
	}
}
