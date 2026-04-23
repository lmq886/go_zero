package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"d:\code\work\go_zero\api\internal\config"

	"github.com/zeromicro/go-zero/core/logx"
)

// CORSMiddleware CORS 跨域中间件
// 用于处理跨域资源共享请求
// 遵循 W3C CORS 规范
type CORSMiddleware struct {
	// 配置信息
	config config.CORSConfig
}

// NewCORSMiddleware 创建 CORS 中间件实例
// 参数 config: CORS 配置
// 返回值: CORS 中间件实例
func NewCORSMiddleware(config config.CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *CORSMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 如果 CORS 未启用，直接处理请求
		if !m.config.Enabled {
			next(w, r)
			return
		}

		// 获取请求来源
		origin := r.Header.Get("Origin")

		// 处理预检请求（OPTIONS 方法）
		if r.Method == http.MethodOptions {
			m.handlePreflight(w, r, origin)
			return
		}

		// 处理实际请求
		m.handleActualRequest(w, r, origin, next)
	}
}

// handlePreflight 处理预检请求
// 参数 w: HTTP 响应写入器
// 参数 r: HTTP 请求
// 参数 origin: 请求来源
func (m *CORSMiddleware) handlePreflight(w http.ResponseWriter, r *http.Request, origin string) {
	// 检查来源是否允许
	if !m.isOriginAllowed(origin) {
		logx.Errorf("Preflight request from origin %s is not allowed", origin)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 设置 CORS 响应头
	m.setCORSHeaders(w, origin)

	// 检查请求方法是否允许
	requestMethod := r.Header.Get("Access-Control-Request-Method")
	if !m.isMethodAllowed(requestMethod) {
		logx.Errorf("Preflight request method %s is not allowed", requestMethod)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 检查请求头是否允许
	requestHeaders := r.Header.Get("Access-Control-Request-Headers")
	if !m.areHeadersAllowed(requestHeaders) {
		logx.Errorf("Preflight request headers %s are not allowed", requestHeaders)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 设置预检请求特定的响应头
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.AllowMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", m.normalizeHeaders(requestHeaders))
	w.Header().Set("Access-Control-Max-Age", strconv.Itoa(m.config.MaxAge))

	// 返回 200 OK 或 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// handleActualRequest 处理实际请求
// 参数 w: HTTP 响应写入器
// 参数 r: HTTP 请求
// 参数 origin: 请求来源
// 参数 next: 下一个处理函数
func (m *CORSMiddleware) handleActualRequest(w http.ResponseWriter, r *http.Request, origin string, next http.HandlerFunc) {
	// 检查来源是否允许
	if !m.isOriginAllowed(origin) {
		logx.Errorf("Request from origin %s is not allowed", origin)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 设置 CORS 响应头
	m.setCORSHeaders(w, origin)

	// 继续处理请求
	next(w, r)
}

// isOriginAllowed 检查来源是否允许
// 参数 origin: 请求来源
// 返回值: 是否允许
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	// 如果没有来源，允许（非浏览器请求）
	if origin == "" {
		return true
	}

	// 检查是否允许所有来源
	for _, allowedOrigin := range m.config.AllowOrigins {
		if allowedOrigin == "*" {
			return true
		}
		// 支持子域名匹配，例如 *.example.com 匹配 api.example.com
		if m.matchOrigin(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

// matchOrigin 匹配来源
// 支持通配符匹配
// 参数 allowedOrigin: 允许的来源
// 参数 origin: 请求来源
// 返回值: 是否匹配
func (m *CORSMiddleware) matchOrigin(allowedOrigin, origin string) bool {
	// 完全匹配
	if allowedOrigin == origin {
		return true
	}

	// 通配符匹配
	// 例如：*.example.com 匹配 api.example.com
	if strings.HasPrefix(allowedOrigin, "*.") {
		suffix := strings.TrimPrefix(allowedOrigin, "*")
		if strings.HasSuffix(origin, suffix) {
			return true
		}
	}

	return false
}

// isMethodAllowed 检查请求方法是否允许
// 参数 method: 请求方法
// 返回值: 是否允许
func (m *CORSMiddleware) isMethodAllowed(method string) bool {
	for _, allowedMethod := range m.config.AllowMethods {
		if strings.EqualFold(allowedMethod, method) {
			return true
		}
	}
	return false
}

// areHeadersAllowed 检查请求头是否允许
// 参数 headers: 请求头字符串
// 返回值: 是否允许
func (m *CORSMiddleware) areHeadersAllowed(headers string) bool {
	// 如果没有请求头，允许
	if headers == "" {
		return true
	}

	// 分割请求头
	requestHeaders := strings.Split(headers, ",")

	// 检查每个请求头
	for _, header := range requestHeaders {
		header = strings.TrimSpace(header)
		header = strings.ToLower(header)

		// 检查是否在允许列表中
		found := false
		for _, allowedHeader := range m.config.AllowHeaders {
			if strings.ToLower(allowedHeader) == header {
				found = true
				break
			}
			// 支持通配符
			if allowedHeader == "*" {
				found = true
				break
			}
		}

		if !found {
			logx.Errorf("Header %s is not allowed", header)
			return false
		}
	}

	return true
}

// setCORSHeaders 设置 CORS 响应头
// 参数 w: HTTP 响应写入器
// 参数 origin: 请求来源
func (m *CORSMiddleware) setCORSHeaders(w http.ResponseWriter, origin string) {
	// 设置允许的来源
	// 如果配置了 *，则返回 *，否则返回请求的来源
	if m.hasWildcardOrigin() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Add("Vary", "Origin")
	}

	// 设置是否允许凭证
	if m.config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// 设置暴露的响应头
	if len(m.config.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.config.ExposeHeaders, ", "))
	}
}

// hasWildcardOrigin 检查是否允许所有来源
// 返回值: 是否允许所有来源
func (m *CORSMiddleware) hasWildcardOrigin() bool {
	for _, origin := range m.config.AllowOrigins {
		if origin == "*" {
			return true
		}
	}
	return false
}

// normalizeHeaders 规范化请求头
// 参数 headers: 请求头字符串
// 返回值: 规范化的请求头字符串
func (m *CORSMiddleware) normalizeHeaders(headers string) string {
	// 如果没有请求头，返回空字符串
	if headers == "" {
		return ""
	}

	// 分割请求头
	headerList := strings.Split(headers, ",")

	// 清理每个请求头
	var normalizedHeaders []string
	for _, header := range headerList {
		header = strings.TrimSpace(header)
		if header != "" {
			normalizedHeaders = append(normalizedHeaders, header)
		}
	}

	return strings.Join(normalizedHeaders, ", ")
}
