package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"d:\code\work\go_zero\api\internal\config"
	"d:\code\work\go_zero\api\internal\svc"
	"d:\code\work\go_zero\model"

	"github.com/zeromicro/go-zero/core/logx"
)

// OperationLogMiddleware 操作日志中间件
// 用于记录用户的操作日志，包括请求信息、响应信息等
// 遵循企业级系统审计要求
type OperationLogMiddleware struct {
	// 服务上下文
	svcCtx *svc.ServiceContext
	// 配置信息
	config config.OperationLogConfig
}

// NewOperationLogMiddleware 创建操作日志中间件实例
// 参数 svcCtx: 服务上下文
// 返回值: 操作日志中间件实例
func NewOperationLogMiddleware(svcCtx *svc.ServiceContext) *OperationLogMiddleware {
	return &OperationLogMiddleware{
		svcCtx: svcCtx,
		config: svcCtx.Config.OperationLog,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 如果操作日志未启用，直接处理请求
		if !m.config.Enabled {
			next(w, r)
			return
		}

		// 检查是否在排除路径中
		if m.isExcludedPath(r.URL.Path) {
			next(w, r)
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体（如果配置了记录请求数据）
		var requestBody []byte
		if m.config.RecordRequest {
			requestBody = m.readRequestBody(r)
		}

		// 创建响应捕获器
		responseRecorder := &responseCapture{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           bytes.NewBufferString(""),
		}

		// 继续处理请求
		next(responseRecorder, r)

		// 计算处理时间
		duration := time.Since(startTime).Milliseconds()

		// 异步记录操作日志
		go m.logOperation(r, requestBody, responseRecorder, startTime, duration)
	}
}

// isExcludedPath 检查路径是否在排除列表中
// 参数 path: 请求路径
// 返回值: 是否在排除列表中
func (m *OperationLogMiddleware) isExcludedPath(path string) bool {
	for _, excludedPath := range m.config.ExcludePaths {
		// 精确匹配
		if path == excludedPath {
			return true
		}
		// 前缀匹配（如果排除路径以 / 结尾）
		if strings.HasSuffix(excludedPath, "/") && strings.HasPrefix(path, excludedPath) {
			return true
		}
	}
	return false
}

// readRequestBody 读取请求体
// 参数 r: HTTP 请求
// 返回值: 请求体内容
func (m *OperationLogMiddleware) readRequestBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}

	// 读取请求体
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logx.Errorf("Failed to read request body: %v", err)
		return nil
	}

	// 重新设置请求体，以便后续处理
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

// logOperation 记录操作日志
// 参数 r: HTTP 请求
// 参数 requestBody: 请求体内容
// 参数 responseRecorder: 响应捕获器
// 参数 startTime: 开始时间
// 参数 duration: 处理时间（毫秒）
func (m *OperationLogMiddleware) logOperation(
	r *http.Request,
	requestBody []byte,
	responseRecorder *responseCapture,
	startTime time.Time,
	duration int64,
) {
	// 从上下文中获取用户信息
	ctx := r.Context()
	userId, _ := GetUserIdFromContext(ctx)
	username, _ := GetUsernameFromContext(ctx)

	// 如果没有用户信息，可能是未登录的请求，跳过日志记录
	if userId == 0 {
		return
	}

	// 解析模块和操作
	module, operation := m.parseModuleAndOperation(r.URL.Path, r.Method)

	// 构建操作日志
	log := &model.OperationLog{
		UserId:    userId,
		Username:  username,
		Module:    module,
		Operation: operation,
		Method:    r.Method,
		Path:      r.URL.Path,
		Status:    int64(responseRecorder.status),
		Ip:        m.getClientIP(r),
		UserAgent: r.UserAgent(),
		Duration:  duration,
	}

	// 记录请求数据（如果配置了）
	if m.config.RecordRequest && len(requestBody) > 0 {
		// 限制请求数据长度，避免存储过大
		maxLength := 4000
		if len(requestBody) > maxLength {
			requestBody = requestBody[:maxLength]
		}
		log.RequestData = stringToNullString(string(requestBody))
	}

	// 记录响应数据（如果配置了）
	if m.config.RecordResponse && responseRecorder.body.Len() > 0 {
		// 限制响应数据长度，避免存储过大
		maxLength := 4000
		responseBody := responseRecorder.body.Bytes()
		if len(responseBody) > maxLength {
			responseBody = responseBody[:maxLength]
		}
		log.ResponseData = stringToNullString(string(responseBody))
	}

	// 如果响应状态码 >= 400，记录错误信息
	if responseRecorder.status >= 400 {
		log.ErrorMsg = stringToNullString("请求失败，状态码: " + string(rune(responseRecorder.status)))
	}

	// 插入数据库
	_, err := m.svcCtx.OperationLogModel.Insert(context.Background(), log)
	if err != nil {
		logx.Errorf("Failed to insert operation log: %v", err)
	}
}

// parseModuleAndOperation 解析模块和操作
// 参数 path: 请求路径
// 参数 method: 请求方法
// 返回值: 模块和操作
func (m *OperationLogMiddleware) parseModuleAndOperation(path, method string) (string, string) {
	// 移除前缀 /api/v1/
	path = strings.TrimPrefix(path, "/api/v1/")

	// 分割路径
	parts := strings.Split(path, "/")

	// 默认模块和操作
	module := "unknown"
	operation := "unknown"

	// 解析模块
	if len(parts) > 0 {
		module = parts[0]
	}

	// 解析操作
	switch method {
	case http.MethodGet:
		operation = "查询"
		// 检查是否是列表查询
		if len(parts) == 1 || (len(parts) == 2 && !isNumeric(parts[1])) {
			operation = "列表查询"
		}
	case http.MethodPost:
		operation = "创建"
		// 检查是否是登录
		if module == "auth" && strings.Contains(path, "login") {
			operation = "登录"
		}
		if module == "auth" && strings.Contains(path, "register") {
			operation = "注册"
		}
		if module == "auth" && strings.Contains(path, "logout") {
			operation = "登出"
		}
		if module == "auth" && strings.Contains(path, "refresh") {
			operation = "刷新令牌"
		}
	case http.MethodPut:
		operation = "更新"
	case http.MethodDelete:
		operation = "删除"
	case http.MethodPatch:
		operation = "部分更新"
	}

	return module, operation
}

// isNumeric 检查字符串是否是数字
// 参数 s: 字符串
// 返回值: 是否是数字
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// getClientIP 获取客户端 IP 地址
// 参数 r: HTTP 请求
// 返回值: 客户端 IP 地址
func (m *OperationLogMiddleware) getClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 头获取
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从 X-Real-IP 头获取
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// 从 RemoteAddr 获取
	return strings.TrimSpace(r.RemoteAddr)
}

// stringToNullString 将字符串转换为 sql.NullString
// 参数 s: 字符串
// 返回值: sql.NullString
func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// responseCapture 响应捕获器
// 用于包装 HTTP 响应写入器，捕获响应状态码和响应体
type responseCapture struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

// Write 重写 Write 方法
// 参数 b: 字节数组
// 返回值: 写入的字节数和错误
func (r *responseCapture) Write(b []byte) (int, error) {
	// 写入到缓冲区
	r.body.Write(b)
	// 写入到原始响应
	return r.ResponseWriter.Write(b)
}

// WriteHeader 重写 WriteHeader 方法
// 参数 code: HTTP 状态码
func (r *responseCapture) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
