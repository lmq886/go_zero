package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"go_zero/api/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

// RateLimitMiddleware 限流中间件
// 用于限制请求速率，防止系统被过多请求压垮
// 采用令牌桶算法实现限流
type RateLimitMiddleware struct {
	config         config.RateLimitConfig
	mu             sync.Mutex
	tokens         int
	lastRefillTime time.Time
}

// NewRateLimitMiddleware 创建限流中间件实例
// 参数 config: 限流配置
// 返回值: 限流中间件实例
func NewRateLimitMiddleware(config config.RateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		config:         config,
		tokens:         config.Burst,
		lastRefillTime: time.Now(),
	}
}

// refillTokens 补充令牌
func (m *RateLimitMiddleware) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(m.lastRefillTime)
	if elapsed > 0 {
		newTokens := int(elapsed.Seconds() * float64(m.config.RequestsPerSecond))
		if newTokens > 0 {
			m.tokens = min(m.config.Burst, m.tokens+newTokens)
			m.lastRefillTime = now
		}
	}
}

// min 返回两个数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// tryAcquire 尝试获取令牌
// 返回值: 是否成功获取，需要等待的时间
func (m *RateLimitMiddleware) tryAcquire() (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.refillTokens()

	if m.tokens > 0 {
		m.tokens--
		return true, 0
	}

	if m.config.RequestsPerSecond <= 0 {
		return false, 0
	}

	waitTime := time.Second / time.Duration(m.config.RequestsPerSecond)
	return false, waitTime
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.config.Enabled {
			next(w, r)
			return
		}

		ok, delay := m.tryAcquire()
		if !ok {
			logx.Errorf("Rate limit exceeded: too many requests")
			m.writeRateLimitError(w, delay)
			return
		}

		if delay > 0 {
			if delay > time.Second {
				logx.Errorf("Rate limit exceeded: delay too long (%v)", delay)
				m.writeRateLimitError(w, delay)
				return
			}
			time.Sleep(delay)
		}

		startTime := time.Now()
		wrappedWriter := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next(wrappedWriter, r)

		duration := time.Since(startTime)
		logx.Infof("Request processed: method=%s, path=%s, status=%d, duration=%v",
			r.Method, r.URL.Path, wrappedWriter.status, duration)
	}
}

// writeRateLimitError 写入限流错误响应
// 参数 w: HTTP 响应写入器
// 参数 delay: 建议的等待时间
func (m *RateLimitMiddleware) writeRateLimitError(w http.ResponseWriter, delay time.Duration) {
	if delay > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(delay.Seconds())+1))
	}

	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(m.config.Burst))
	w.Header().Set("X-RateLimit-Remaining", "0")

	WriteError(w, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
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
type IPBasedRateLimitMiddleware struct {
	config   config.RateLimitConfig
	limiters map[string]*RateLimitMiddleware
	mu       sync.RWMutex
}

// NewIPBasedRateLimitMiddleware 创建基于 IP 的限流中间件实例
// 参数 config: 限流配置
// 返回值: 基于 IP 的限流中间件实例
func NewIPBasedRateLimitMiddleware(config config.RateLimitConfig) *IPBasedRateLimitMiddleware {
	return &IPBasedRateLimitMiddleware{
		config:   config,
		limiters: make(map[string]*RateLimitMiddleware),
	}
}

// getOrCreateLimiter 获取或创建限流器
// 参数 clientIP: 客户端 IP 地址
// 返回值: 限流器实例
func (m *IPBasedRateLimitMiddleware) getOrCreateLimiter(clientIP string) *RateLimitMiddleware {
	m.mu.RLock()
	limiter, ok := m.limiters[clientIP]
	m.mu.RUnlock()

	if ok {
		return limiter
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	limiter, ok = m.limiters[clientIP]
	if ok {
		return limiter
	}

	limiter = NewRateLimitMiddleware(m.config)
	m.limiters[clientIP] = limiter
	return limiter
}

// Handle 中间件处理函数
// 参数 next: 下一个处理函数
// 返回值: 处理函数
func (m *IPBasedRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.config.Enabled {
			next(w, r)
			return
		}

		clientIP := m.getClientIP(r)
		limiter := m.getOrCreateLimiter(clientIP)

		ok, delay := limiter.tryAcquire()
		if !ok {
			logx.Errorf("IP rate limit exceeded for %s: too many requests", clientIP)
			m.writeRateLimitError(w, clientIP, delay)
			return
		}

		if delay > 0 {
			if delay > time.Second {
				logx.Errorf("IP rate limit exceeded for %s: delay too long (%v)", clientIP, delay)
				m.writeRateLimitError(w, clientIP, delay)
				return
			}
			time.Sleep(delay)
		}

		next(w, r)
	}
}

// getClientIP 获取客户端 IP 地址
// 参数 r: HTTP 请求
// 返回值: 客户端 IP 地址
func (m *IPBasedRateLimitMiddleware) getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	return r.RemoteAddr
}

// writeRateLimitError 写入限流错误响应
// 参数 w: HTTP 响应写入器
// 参数 clientIP: 客户端 IP 地址
// 参数 delay: 建议的等待时间
func (m *IPBasedRateLimitMiddleware) writeRateLimitError(w http.ResponseWriter, clientIP string, delay time.Duration) {
	if delay > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(delay.Seconds())+1))
	}

	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(m.config.Burst))
	w.Header().Set("X-RateLimit-Remaining", "0")

	WriteError(w, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
}
