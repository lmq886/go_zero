package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

// Config 服务配置结构体
// 对应配置文件 api.yaml 中的配置项
// 遵循 GoZero 最佳实践，使用结构体映射配置
type Config struct {
	// 基础服务配置
	rest.RestConf

	// 数据库配置
	DataSource DataSourceConfig

	// JWT 认证配置
	Auth AuthConfig

	// 缓存配置（可选）
	CacheRedis cache.CacheConf `json:",optional"`

	// 限流配置
	RateLimit RateLimitConfig

	// 熔断配置
	CircuitBreaker CircuitBreakerConfig

	// CORS 跨域配置
	CORS CORSConfig

	// 操作日志配置
	OperationLog OperationLogConfig

	// 登录日志配置
	LoginLog LoginLogConfig

	// 系统配置
	System SystemConfig
}

// DataSourceConfig 数据库配置结构体
// 定义数据库连接相关的配置项
type DataSourceConfig struct {
	// 数据库类型（mysql, postgres, sqlite）
	Type string `json:",default=postgres"`
	// 数据库主机地址
	Host string `json:",default=localhost"`
	// 数据库端口
	Port int `json:",default=5432"`
	// 数据库名称
	Database string
	// 数据库用户名
	Username string
	// 数据库密码
	Password string
	// 字符集
	Charset string `json:",default=utf8mb4"`
	// 最大空闲连接数
	MaxIdleConns int `json:",default=10"`
	// 最大打开连接数
	MaxOpenConns int `json:",default=100"`
	// 连接最大生命周期（秒）
	MaxLifetime int `json:",default=3600"`
	// 是否启用 SSL（PostgreSQL 专用）
	SSLMode string `json:",default=disable"`
	// 时区
	TimeZone string `json:",default=Asia/Shanghai"`
}

// AuthConfig JWT 认证配置结构体
// 定义 JWT 令牌相关的配置项
type AuthConfig struct {
	// 访问令牌密钥
	AccessSecret string
	// 访问令牌有效期（秒）
	AccessExpire int64 `json:",default=7200"`
	// 刷新令牌密钥
	RefreshSecret string
	// 刷新令牌有效期（秒）
	RefreshExpire int64 `json:",default=604800"`
	// 发行者
	Issuer string `json:",default=admin-api"`
	// 受众
	Audience string `json:",default=admin-web"`
}

// RateLimitConfig 限流配置结构体
// 定义请求限流相关的配置项
type RateLimitConfig struct {
	// 是否启用限流
	Enabled bool `json:",default=true"`
	// 每秒允许的请求数
	RequestsPerSecond int `json:",default=100"`
	// 突发请求数
	Burst int `json:",default=200"`
}

// CircuitBreakerConfig 熔断配置结构体
// 定义服务熔断相关的配置项
type CircuitBreakerConfig struct {
	// 是否启用熔断
	Enabled bool `json:",default=true"`
	// 熔断窗口大小（秒）
	Window int `json:",default=10"`
	// 熔断阈值（错误率百分比）
	Threshold int `json:",default=50"`
	// 最小请求数
	MinRequests int `json:",default=10"`
	// 半开状态允许的请求数
	HalfOpenRequests int `json:",default=5"`
	// 冷却时间（秒）
	CoolDown int `json:",default=5"`
}

// CORSConfig 跨域配置结构体
// 定义 CORS 相关的配置项
type CORSConfig struct {
	// 是否启用 CORS
	Enabled bool `json:",default=true"`
	// 允许的来源
	AllowOrigins []string `json:",default=[\"*\"]"`
	// 允许的请求方法
	AllowMethods []string `json:",default=[\"GET\",\"POST\",\"PUT\",\"DELETE\",\"OPTIONS\",\"PATCH\"]"`
	// 允许的请求头
	AllowHeaders []string `json:",default=[\"Content-Type\",\"Authorization\",\"X-Requested-With\",\"Accept\",\"Origin\"]"`
	// 暴露的响应头
	ExposeHeaders []string `json:",optional"`
	// 是否允许凭证
	AllowCredentials bool `json:",default=true"`
	// 预检请求的有效期（秒）
	MaxAge int `json:",default=86400"`
}

// OperationLogConfig 操作日志配置结构体
// 定义操作日志记录相关的配置项
type OperationLogConfig struct {
	// 是否启用操作日志
	Enabled bool `json:",default=true"`
	// 是否记录请求数据
	RecordRequest bool `json:",default=true"`
	// 是否记录响应数据
	RecordResponse bool `json:",default=true"`
	// 排除的路径（不记录日志）
	ExcludePaths []string `json:",optional"`
}

// LoginLogConfig 登录日志配置结构体
// 定义登录日志记录相关的配置项
type LoginLogConfig struct {
	// 是否启用登录日志
	Enabled bool `json:",default=true"`
}

// SystemConfig 系统配置结构体
// 定义系统级别的配置项
type SystemConfig struct {
	// 超级管理员角色编码
	SuperAdminRoleCode string `json:",default=super_admin"`
	// 默认密码（用于重置密码）
	DefaultPassword string `json:",default=admin123"`
	// 密码最小长度
	PasswordMinLength int `json:",default=6"`
	// 密码最大长度
	PasswordMaxLength int `json:",default=20"`
	// 登录失败最大重试次数
	LoginMaxRetries int `json:",default=5"`
	// 登录失败锁定时长（分钟）
	LoginLockDuration int `json:",default=30"`
	// 是否允许用户注册
	AllowRegister bool `json:",default=true"`
}
