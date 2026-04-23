package main

import (
	"flag"
	"fmt"

	"d:\code\work\go_zero\api\internal\config"
	"d:\code\work\go_zero\api\internal\handler"
	"d:\code\work\go_zero\api\internal\svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// main 服务入口函数
// 遵循 GoZero 最佳实践，使用 flag 解析命令行参数，加载配置，初始化服务
func main() {
	// 定义命令行参数
	// -f: 配置文件路径，默认为 etc/api.yaml
	var configFile = flag.String("f", "etc/api.yaml", "the config file")
	flag.Parse()

	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 配置日志
	// 使用配置文件中的日志配置
	logx.MustSetup(c.Log)

	// 创建服务上下文
	// 初始化数据库连接、模型实例等依赖
	serverCtx := svc.NewServiceContext(c)

	// 创建 REST 服务器
	// 使用配置文件中的服务配置
	server := rest.MustNewServer(c.RestConf)
	// 确保在服务退出时关闭服务器
	defer server.Stop()

	// 注册路由
	// 将所有 HTTP 处理器注册到服务器
	handler.RegisterHandlers(server, serverCtx)

	// 打印服务启动信息
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("Service mode: %s\n", c.Mode)
	fmt.Printf("Log level: %s\n", c.Log.Level)

	// 启动服务器
	// 这是一个阻塞调用，直到服务停止
	server.Start()
}
