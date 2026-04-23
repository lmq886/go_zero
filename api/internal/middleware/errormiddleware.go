package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ErrorMiddleware 错误处理中间件
// 用于统一处理请求过程中的错误，包括 panic 恢复
// 遵循企业级系统错误处理最佳实践
type ErrorMiddleware struct {
}

// NewErrorMiddleware 创建错误处理中间件实例
// 返回值: 错误处理中间件实例
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *ErrorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 捕获 panic
		defer func() {
			if err := recover(); err != nil {
				// 记录错误和堆栈信息
				logx.Errorf("Panic recovered: %v\nStack: %s", err, debug.Stack())

				// 返回统一的错误响应
				m.writeErrorResponse(w, http.StatusInternalServerError, "服务器内部错误，请稍后重试")
			}
		}()

		// 创建响应捕获器
		responseRecorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// 继续处理请求
		next(responseRecorder, r)

		// 检查响应状态码，对错误状态码进行统一处理
		if responseRecorder.status >= 400 {
			logx.Errorf("Request failed: method=%s, path=%s, status=%d",
				r.Method, r.URL.Path, responseRecorder.status)
		}
	}
}

// writeErrorResponse 写入错误响应
// 参数 w: HTTP 响应写入器
// 参数 code: HTTP 状态码
// 参数 message: 错误消息
func (m *ErrorMiddleware) writeErrorResponse(w http.ResponseWriter, code int, message string) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	// 构建错误响应
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	// 序列化并写入响应
	json.NewEncoder(w).Encode(response)
}

// ErrorResponse 错误响应结构体
// 定义统一的错误响应格式
type ErrorResponse struct {
	// 状态码
	Code int `json:"code"`
	// 错误消息
	Message string `json:"message"`
	// 详细错误信息（可选）
	Details string `json:"details,omitempty"`
}

// ResponseMiddleware 响应包装中间件
// 用于统一包装响应格式，将所有成功响应包装为统一格式
// 遵循企业级系统响应格式最佳实践
type ResponseMiddleware struct {
}

// NewResponseMiddleware 创建响应包装中间件实例
// 返回值: 响应包装中间件实例
func NewResponseMiddleware() *ResponseMiddleware {
	return &ResponseMiddleware{}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *ResponseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 创建响应捕获器
		responseRecorder := &responseBodyCapture{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           make([]byte, 0),
		}

		// 继续处理请求
		next(responseRecorder, r)

		// 如果响应状态码不是 200，不进行包装
		if responseRecorder.status != http.StatusOK {
			return
		}

		// 如果响应内容类型不是 JSON，不进行包装
		contentType := responseRecorder.Header().Get("Content-Type")
		if !isJSONContentType(contentType) {
			return
		}

		// 包装响应
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(responseRecorder.status)

		// 尝试解析原始响应
		var rawResponse interface{}
		if len(responseRecorder.body) > 0 {
			if err := json.Unmarshal(responseRecorder.body, &rawResponse); err != nil {
				// 如果解析失败，直接返回原始响应
				w.Write(responseRecorder.body)
				return
			}
		}

		// 检查是否已经是统一格式的响应
		if isAlreadyFormatted(rawResponse) {
			w.Write(responseRecorder.body)
			return
		}

		// 包装为统一格式
		formattedResponse := SuccessResponse{
			Code:    http.StatusOK,
			Message: "操作成功",
			Data:    rawResponse,
		}

		// 序列化并写入响应
		json.NewEncoder(w).Encode(formattedResponse)
	}
}

// isJSONContentType 检查内容类型是否为 JSON
// 参数 contentType: 内容类型
// 返回值: 是否为 JSON 内容类型
func isJSONContentType(contentType string) bool {
	return contentType == "application/json" ||
		contentType == "application/json; charset=utf-8"
}

// isAlreadyFormatted 检查响应是否已经是统一格式
// 参数 response: 响应数据
// 返回值: 是否已经是统一格式
func isAlreadyFormatted(response interface{}) bool {
	if response == nil {
		return false
	}

	// 检查是否是 map 类型
	respMap, ok := response.(map[string]interface{})
	if !ok {
		return false
	}

	// 检查是否包含 code 字段
	_, hasCode := respMap["code"]
	_, hasMessage := respMap["message"]

	return hasCode && hasMessage
}

// SuccessResponse 成功响应结构体
// 定义统一的成功响应格式
type SuccessResponse struct {
	// 状态码
	Code int `json:"code"`
	// 消息
	Message string `json:"message"`
	// 数据
	Data interface{} `json:"data,omitempty"`
}

// PagedResponse 分页响应结构体
// 定义统一的分页响应格式
type PagedResponse struct {
	// 状态码
	Code int `json:"code"`
	// 消息
	Message string `json:"message"`
	// 数据
	Data interface{} `json:"data"`
	// 分页信息
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo 分页信息结构体
// 定义分页信息的格式
type PaginationInfo struct {
	// 当前页码
	Page int `json:"page"`
	// 每页数量
	PageSize int `json:"pageSize"`
	// 总记录数
	Total int64 `json:"total"`
	// 总页数
	TotalPages int `json:"totalPages"`
}

// responseBodyCapture 响应体捕获器
// 用于包装 HTTP 响应写入器，捕获响应状态码和响应体
type responseBodyCapture struct {
	http.ResponseWriter
	status int
	body   []byte
}

// Write 重写 Write 方法
// 参数 b: 字节数组
// 返回值: 写入的字节数和错误
func (r *responseBodyCapture) Write(b []byte) (int, error) {
	// 保存到缓冲区
	r.body = append(r.body, b...)
	// 返回原始写入的字节数
	return len(b), nil
}

// WriteHeader 重写 WriteHeader 方法
// 参数 code: HTTP 状态码
func (r *responseBodyCapture) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// NewSuccessResponse 创建成功响应
// 参数 data: 响应数据
// 返回值: 成功响应
func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
	}
}

// NewPagedResponse 创建分页响应
// 参数 data: 响应数据
// 参数 page: 当前页码
// 参数 pageSize: 每页数量
// 参数 total: 总记录数
// 返回值: 分页响应
func NewPagedResponse(data interface{}, page, pageSize int, total int64) *PagedResponse {
	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResponse{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
		Pagination: PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// NewErrorResponse 创建错误响应
// 参数 code: 状态码
// 参数 message: 错误消息
// 返回值: 错误响应
func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// WriteJSONResponse 写入 JSON 响应
// 参数 w: HTTP 响应写入器
// 参数 response: 响应数据
func WriteJSONResponse(w http.ResponseWriter, response interface{}) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// 序列化并写入响应
	json.NewEncoder(w).Encode(response)
}

// WriteError 写入错误响应
// 参数 w: HTTP 响应写入器
// 参数 code: HTTP 状态码
// 参数 message: 错误消息
func WriteError(w http.ResponseWriter, code int, message string) {
	// 使用 httpx 写入错误
	httpx.Error(w, code, message)
}
