package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DB struct {
		DataSource string
		MaxOpenConns int
		MaxIdleConns int
		LogLevel string
	}
	Cache cache.CacheConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Upload struct {
		MaxSize   int64
		AllowTypes string
		SavePath  string
	}
}
