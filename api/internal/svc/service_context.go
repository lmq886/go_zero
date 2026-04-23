/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\svc\service_context.go
 * @Description: 服务上下文，用于管理数据库连接和数据模型实例
 */
package svc

import (
	"go_zero/api/internal/config"
	"go_zero/api/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 服务上下文结构体
// 用于集中管理所有数据模型实例和配置
// 避免在每个handler中重复创建数据库连接和模型
type ServiceContext struct {
	Config               config.Config         // 系统配置
	UserModel            model.UserModel       // 用户数据模型
	RoleModel            model.RoleModel       // 角色数据模型
	PermissionModel      model.PermissionModel // 权限数据模型
	UserRoleModel        model.UserRoleModel   // 用户角色关联模型
	RolePermissionModel  model.RolePermissionModel // 角色权限关联模型
	SystemConfigModel    model.SystemConfigModel   // 系统配置数据模型
	OperationLogModel    model.OperationLogModel   // 操作日志数据模型
	LoginLogModel        model.LoginLogModel       // 登录日志数据模型
	FileModel            model.FileModel           // 文件数据模型
}

// NewServiceContext 创建服务上下文实例
// 参数: c - 系统配置
// 返回: *ServiceContext - 服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	// 创建PostgreSQL数据库连接
	conn := sqlx.NewSqlConn("postgres", c.DB.DataSource)

	// 初始化所有数据模型实例
	return &ServiceContext{
		Config:              c,
		UserModel:           model.NewUserModel(conn),
		RoleModel:           model.NewRoleModel(conn),
		PermissionModel:     model.NewPermissionModel(conn),
		UserRoleModel:       model.NewUserRoleModel(conn),
		RolePermissionModel: model.NewRolePermissionModel(conn),
		SystemConfigModel:   model.NewSystemConfigModel(conn),
		OperationLogModel:   model.NewOperationLogModel(conn),
		LoginLogModel:       model.NewLoginLogModel(conn),
		FileModel:           model.NewFileModel(conn),
	}
}
