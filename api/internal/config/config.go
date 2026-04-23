/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\config\config.go
 * @Description: 系统配置定义结构体
 */
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

// Config 系统配置结构体
// 包含了所有模块的配置信息
type Config struct {
	rest.RestConf // 基础HTTP服务配置（端口、地址、日志等）
	
	// 数据库配置
	DB struct {
		DataSource   string // 数据库连接字符串
		MaxOpenConns int    // 最大打开连接数
		MaxIdleConns int    // 最大空闲连接数
		LogLevel     string // 日志级别
	}
	
	// 缓存配置
	Cache cache.CacheConf
	
	// JWT认证配置
	JwtAuth struct {
		AccessSecret string // JWT签名密钥
		AccessExpire int64  // Token过期时间（秒）
	}
	
	// 文件上传配置
	Upload struct {
		MaxSize     int64  // 最大文件大小（字节）
		AllowTypes  string // 允许上传的文件类型（逗号分隔）
		SavePath    string // 文件保存路径
	}
}
