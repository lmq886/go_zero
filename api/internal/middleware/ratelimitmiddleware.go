package middleware

import (
	"net/http"
	"strconv"
	"time"

	"d:\code\work\go_zero\api\internal\config"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RateLimitMiddleware 限流中间件
// 用于限制请求速率，防止系统被过多请求压垮
// 采用令牌桶算法实现限流
type RateLimitMiddleware struct {
	// 配置信息
	config config.RateLimitConfig
	// 限流器
	limiter *limit.TokenLimiter
}

// NewRateLimitMiddleware 创建限流中间件实例
// 参数 config: 限流配置
// 返回值: 限流中间件实例
func NewRateLimitMiddleware(config config.RateLimitConfig) *RateLimitMiddleware {
	// 创建令牌桶限流器
	// rate: 每秒生成的令牌数
	// capacity: 令牌桶容量（突发请求数）
	limiter := limit.NewTokenLimiter(
		config.RequestsPerSecond,
		config.Burst,
	)

	return &RateLimitMiddleware{
		config:  config,
		limiter: limiter,
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 如果限流未启用，直接处理请求
		if !m.config.Enabled {
			next(w, r)
			return
		}

		// 尝试获取令牌
		// reserve: 预约令牌
		// ok: 是否成功获取
		reserve, ok := m.limiter.Reserve()
		if !ok {
			// 令牌桶已满，无法预约
			logx.Errorf("Rate limit exceeded: too many requests")
			m.writeRateLimitError(w, 0)
			return
		}

		// 检查是否需要等待
		delay := reserve.Delay()
		if delay > 0 {
			// 需要等待，检查等待时间是否可接受
			// 如果等待时间超过 1 秒，直接返回限流错误
			if delay > time.Second {
				logx.Errorf("Rate limit exceeded: delay too long (%v)", delay)
				reserve.Cancel()
				m.writeRateLimitError(w, delay)
				return
			}

			// 等待指定时间
			time.Sleep(delay)
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 包装响应写入器以捕获状态码
		wrappedWriter := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// 继续处理请求
		next(wrappedWriter, r)

		// 记录请求处理时间
		duration := time.Since(startTime)
		logx.Infof("Request processed: method=%s, path=%s, status=%d, duration=%v",
			r.Method, r.URL.Path, wrappedWriter.status, duration)
	}
}

// writeRateLimitError 写入限流错误响应
// 参数 w: HTTP 响应写入器
// 参数 delay: 建议的等待时间
func (m *RateLimitMiddleware) writeRateLimitError(w http.ResponseWriter, delay time.Duration) {
	// 设置 Retry-After 响应头
	if delay > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(delay.Seconds())+1))
	}

	// 设置 X-RateLimit-Limit 响应头
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(m.config.Burst))

	// 设置 X-RateLimit-Remaining 响应头
	w.Header().Set("X-RateLimit-Remaining", "0")

	// 返回 429 Too Many Requests 错误
	httpx.Error(w, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
}

// statusRecorder 状态记录器
// 用于包装 HTTP 响应写入器，捕获响应状态码
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader 重写 WriteHeader 方法
// 参数 code: HTTP 状态码
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// IPBasedRateLimitMiddleware 基于 IP 的限流中间件
// 用于对每个 IP 地址进行独立的限流
// 适用于防止单个 IP 发送过多请求
type IPBasedRateLimitMiddleware struct {
	// 配置信息
	config config.RateLimitConfig
	// IP 限流器映射
	limiters map[string]*limit.TokenLimiter
}

// NewIPBasedRateLimitMiddleware 创建基于 IP 的限流中间件实例
// 参数 config: 限流配置
// 返回值: 基于 IP 的限流中间件实例
func NewIPBasedRateLimitMiddleware(config config.RateLimitConfig) *IPBasedRateLimitMiddleware {
	return &IPBasedRateLimitMiddleware{
		config:   config,
		limiters: make(map[string]*limit.TokenLimiter),
	}
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *IPBasedRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 如果限流未启用，直接处理请求
		if !m.config.Enabled {
			next(w, r)
			return
		}

		// 获取客户端 IP 地址
		clientIP := m.getClientIP(r)

		// 获取或创建该 IP 的限流器
		limiter := m.getOrCreateLimiter(clientIP)

		// 尝试获取令牌
		reserve, ok := limiter.Reserve()
		if !ok {
			// 令牌桶已满，无法预约
			logx.Errorf("IP rate limit exceeded for %s: too many requests", clientIP)
			m.writeRateLimitError(w, clientIP, 0)
			return
		}

		// 检查是否需要等待
		delay := reserve.Delay()
		if delay > 0 {
			// 如果等待时间超过 1 秒，直接返回限流错误
			if delay > time.Second {
				logx.Errorf("IP rate limit exceeded for %s: delay too long (%v)", clientIP, delay)
				reserve.Cancel()
				m.writeRateLimitError(w, clientIP, delay)
				return
			}

			// 等待指定时间
			time.Sleep(delay)
		}

		// 继续处理请求
		next(w, r)
	}
}

// getClientIP 获取客户端 IP 地址
// 参数 r: HTTP 请求
// 返回值: 客户端 IP 地址
func (m *IPBasedRateLimitMiddleware) getClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 头获取
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		ips := splitAndTrim(xff, ",")
		if len(ips) > 0 {
			return ips[0]
		}
	}

	// 尝试从 X-Real-IP 头获取
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// 从 RemoteAddr 获取
	return r.RemoteAddr
}

// splitAndTrim 分割字符串并去除空格
// 参数 s: 字符串
// 参数 sep: 分隔符
// 返回值: 分割后的字符串切片
func splitAndTrim(s, sep string) []string {
	var result []string
	parts := split(s, sep)
	for _, part := range parts {
		trimmed := trim(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// split 分割字符串
// 参数 s: 字符串
// 参数 sep: 分隔符
// 返回值: 分割后的字符串切片
func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// trim 去除字符串两端的空格
// 参数 s: 字符串
// 返回值: 去除空格后的字符串
func trim(s string) string {
	start := 0
	end := len(s) - 1

	// 去除开头空格
	for start <= end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// 去除结尾空格
	for end >= start && (s[end] == ' ' || s[end] == '\t' || s[end] == '\n' || s[end] == '\r') {
		end--
	}

	if start > end {
		return ""
	}

	return s[start : end+1]
}

// getOrCreateLimiter 获取或创建限流器
// 参数 clientIP: 客户端 IP 地址
// 返回值: 限流器实例
func (m *IPBasedRateLimitMiddleware) getOrCreateLimiter(clientIP string) *limit.TokenLimiter {
	// 检查是否已存在
	if limiter, ok := m.limiters[clientIP]; ok {
		return limiter
	}

	// 创建新的限流器
	limiter := limit.NewTokenLimiter(
		m.config.RequestsPerSecond,
		m.config.Burst,
	)

	// 存储到映射中
	m.limiters[clientIP] = limiter

	return limiter
}

// writeRateLimitError 写入限流错误响应
// 参数 w: HTTP 响应写入器
// 参数 clientIP: 客户端 IP 地址
// 参数 delay: 建议的等待时间
func (m *IPBasedRateLimitMiddleware) writeRateLimitError(w http.ResponseWriter, clientIP string, delay time.Duration) {
	// 设置 Retry-After 响应头
	if delay > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(delay.Seconds())+1))
	}

	// 设置 X-RateLimit-Limit 响应头
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(m.config.Burst))

	// 设置 X-RateLimit-Remaining 响应头
	w.Header().Set("X-RateLimit-Remaining", "0")

	// 返回 429 Too Many Requests 错误
	httpx.Error(w, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
}
