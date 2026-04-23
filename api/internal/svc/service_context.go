package svc

import (
	"go_zero/api/internal/config"
	"go_zero/api/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	UserModel            model.UserModel
	RoleModel            model.RoleModel
	PermissionModel      model.PermissionModel
	UserRoleModel        model.UserRoleModel
	RolePermissionModel  model.RolePermissionModel
	SystemConfigModel    model.SystemConfigModel
	OperationLogModel    model.OperationLogModel
	LoginLogModel        model.LoginLogModel
	FileModel            model.FileModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DB.DataSource)

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
