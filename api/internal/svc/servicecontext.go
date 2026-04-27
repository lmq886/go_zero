package svc

import (
	"fmt"

	"go_zero/api/internal/config"
	"go_zero/model"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 服务上下文
// 用于在处理请求时传递依赖项，如数据库连接、模型实例等
// 遵循 GoZero 最佳实践，通过依赖注入方式管理依赖
type ServiceContext struct {
	// 配置信息
	Config config.Config

	// 数据库连接
	DB sqlx.SqlConn

	// 模型实例
	// 用户模型
	UserModel model.UserModel
	// 角色模型
	RoleModel model.RoleModel
	// 权限模型
	PermissionModel model.PermissionModel
	// 菜单模型
	MenuModel model.MenuModel
	// 操作日志模型
	OperationLogModel model.OperationLogModel
	// 登录日志模型
	LoginLogModel model.LoginLogModel
	// 系统配置模型
	ConfigModel model.ConfigModel
}

// NewServiceContext 创建服务上下文实例
// 参数 c: 配置信息
// 返回值: 服务上下文实例
// 遵循 GoZero 最佳实践，在服务启动时初始化所有依赖
func NewServiceContext(c config.Config) *ServiceContext {
	// 构建数据库连接字符串
	dsn := buildDSN(c.DataSource)

	// 创建数据库连接
	// 使用 postgres 驱动，sqlx 提供了更方便的数据库操作
	conn := postgres.New(dsn)

	return &ServiceContext{
		Config:            c,
		DB:                conn,
		UserModel:         model.NewUserModel(conn),
		RoleModel:         model.NewRoleModel(conn),
		PermissionModel:   model.NewPermissionModel(conn),
		MenuModel:         model.NewMenuModel(conn),
		OperationLogModel: model.NewOperationLogModel(conn),
		LoginLogModel:     model.NewLoginLogModel(conn),
		ConfigModel:       model.NewConfigModel(conn),
	}
}

// buildDSN 构建数据库连接字符串
// 参数 config: 数据库配置
// 返回值: PostgreSQL DSN 连接字符串
// 格式: host=localhost port=5432 user=postgres password=password dbname=admin_system sslmode=disable TimeZone=Asia/Shanghai
func buildDSN(config config.DataSourceConfig) string {
	// 构建 PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
		config.TimeZone,
	)

	return dsn
}
