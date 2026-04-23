package handler

import (
	"github.com/zeromicro/go-zero/rest"

	"go_zero/api/internal/middleware"
	"go_zero/api/internal/svc"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext) {
	// 认证中间件
	jwtAuth := middleware.NewJwtAuthMiddleware(svcCtx.Config)
	permissionAuth := middleware.NewPermissionAuthMiddleware(svcCtx)

	// Auth
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/auth/login",
		Handler: NewLoginHandler(svcCtx).ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/auth/register",
		Handler: NewRegisterHandler(svcCtx).ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/auth/logout",
		Handler: NewLogoutHandler(svcCtx).ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/auth/userinfo",
		Handler: NewGetUserInfoHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
		},
	})

	// User
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/users",
		Handler: NewCreateUserHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpPut,
		Path:    "/api/v1/users/:id",
		Handler: NewUpdateUserHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/users/:id",
		Handler: NewDeleteUserHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/users/:id",
		Handler: NewGetUserHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/users",
		Handler: NewListUsersHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpPut,
		Path:    "/api/v1/users/:id/password",
		Handler: NewUpdatePasswordHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})

	// Role
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/roles",
		Handler: NewCreateRoleHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpPut,
		Path:    "/api/v1/roles/:id",
		Handler: NewUpdateRoleHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/roles/:id",
		Handler: NewDeleteRoleHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/roles/:id",
		Handler: NewGetRoleHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/roles",
		Handler: NewListRolesHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})

	// Permission
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/permissions",
		Handler: NewCreatePermissionHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpPut,
		Path:    "/api/v1/permissions/:id",
		Handler: NewUpdatePermissionHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/permissions/:id",
		Handler: NewDeletePermissionHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/permissions/:id",
		Handler: NewGetPermissionHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/permissions",
		Handler: NewListPermissionsHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})

	// System Config
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/configs",
		Handler: NewCreateSystemConfigHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpPut,
		Path:    "/api/v1/configs/:id",
		Handler: NewUpdateSystemConfigHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/configs/:id",
		Handler: NewDeleteSystemConfigHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/configs/:id",
		Handler: NewGetSystemConfigHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/configs",
		Handler: NewListSystemConfigsHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})

	// Log
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/operation-logs",
		Handler: NewListOperationLogsHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/operation-logs/:id",
		Handler: NewDeleteOperationLogHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/login-logs",
		Handler: NewListLoginLogsHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/login-logs/:id",
		Handler: NewDeleteLoginLogHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})

	// File
	server.AddRoute(rest.Route{
		Method:  httpPost,
		Path:    "/api/v1/files/upload",
		Handler: NewUploadFileHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/files/:id",
		Handler: NewGetFileHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpGet,
		Path:    "/api/v1/files",
		Handler: NewListFilesHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
	server.AddRoute(rest.Route{
		Method:  httpDelete,
		Path:    "/api/v1/files/:id",
		Handler: NewDeleteFileHandler(svcCtx).ServeHTTP,
		Middleware: []rest.Middleware{
			jwtAuth.Handle,
			permissionAuth.Handle,
		},
	})
}

const (
	httpGet    = "GET"
	httpPost   = "POST"
	httpPut    = "PUT"
	httpDelete = "DELETE"
)
