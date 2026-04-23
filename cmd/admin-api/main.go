/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\cmd\admin-api\main.go
 * @Description: 管理后台API服务主入口
 */
package main

import (
	"flag"
	"fmt"

	"go_zero/api/internal/config"
	"go_zero/api/internal/handler"
	"go_zero/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// configFile 配置文件路径参数
// 默认值: etc/admin-api.yaml
var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

// main 程序主入口
func main() {
	// 解析命令行参数
	flag.Parse()

	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建HTTP服务器
	server := rest.MustNewServer(c.RestConf)
	// 确保服务器在退出时停止
	defer server.Stop()

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	// 注册所有API路由
	handler.RegisterHandlers(server, ctx)

	// 打印启动信息
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 启动服务器
	server.Start()
}
