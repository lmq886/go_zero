/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\types\types.go
 * @Description: API请求响应类型定义
 */
package types

// BaseRequest 基础请求结构体
type BaseRequest struct {
}

// BaseResponse 基础响应结构体
type BaseResponse struct {
	Code    int    `json:"code"`    // 响应码（0:成功，非0:失败）
	Message string `json:"message"` // 响应消息
}

// PageRequest 分页请求结构体
type PageRequest struct {
	Page     int `form:"page,default=1"`       // 页码（默认1）
	PageSize int `form:"page_size,default=10"` // 每页数量（默认10）
}

// PageResponse 分页响应结构体
type PageResponse struct {
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	List     interface{} `json:"list"`      // 数据列表
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username,optional"` // 用户名
	Password string `json:"password,optional"` // 密码
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Token     string `json:"token"`      // JWT Token
		ExpiresAt int64  `json:"expires_at"` // Token过期时间戳
		UserInfo  struct {
			ID       int64   `json:"id"`       // 用户ID
			Username string  `json:"username"` // 用户名
			Nickname string  `json:"nickname"` // 昵称
			Avatar   string  `json:"avatar"`   // 头像URL
			RoleIDs  []int64 `json:"role_ids"` // 角色ID列表
		} `json:"user_info"` // 用户信息
	} `json:"data"` // 响应数据
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username"`          // 用户名
	Password string `json:"password"`          // 密码
	Nickname string `json:"nickname,optional"` // 昵称
	Email    string `json:"email,optional"`    // 邮箱
	Phone    string `json:"phone,optional"`    // 手机号
}

// UserInfoResponse 用户信息响应结构体
type UserInfoResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		ID          int64            `json:"id"`          // 用户ID
		Username    string           `json:"username"`    // 用户名
		Nickname    string           `json:"nickname"`    // 昵称
		Avatar      string           `json:"avatar"`      // 头像URL
		Email       string           `json:"email"`       // 邮箱
		Phone       string           `json:"phone"`       // 手机号
		Status      int              `json:"status"`      // 状态
		CreatedAt   string           `json:"created_at"`  // 创建时间
		UpdatedAt   string           `json:"updated_at"`  // 更新时间
		Roles       []RoleInfo       `json:"roles"`       // 角色列表
		Permissions []PermissionInfo `json:"permissions"` // 权限列表
	} `json:"data"` // 响应数据
}

// CreateUserRequest 创建用户请求结构体
type CreateUserRequest struct {
	Username string  `json:"username"`                  // 用户名
	Password string  `json:"password"`                  // 密码
	Nickname string  `json:"nickname,optional"`         // 昵称
	Email    string  `json:"email,optional"`            // 邮箱
	Phone    string  `json:"phone,optional"`            // 手机号
	Avatar   string  `json:"avatar,optional"`           // 头像URL
	Status   int     `json:"status,optional,default=1"` // 状态（默认1:正常）
	RoleIDs  []int64 `json:"role_ids,optional"`         // 角色ID列表
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	ID       int64   `json:"id"`                // 用户ID
	Nickname string  `json:"nickname,optional"` // 昵称
	Email    string  `json:"email,optional"`    // 邮箱
	Phone    string  `json:"phone,optional"`    // 手机号
	Avatar   string  `json:"avatar,optional"`   // 头像URL
	Status   int     `json:"status,optional"`   // 状态
	RoleIDs  []int64 `json:"role_ids,optional"` // 角色ID列表
}

// UpdatePasswordRequest 更新密码请求结构体
type UpdatePasswordRequest struct {
	ID          int64  `json:"id"`                    // 用户ID
	OldPassword string `json:"old_password,optional"` // 旧密码
	NewPassword string `json:"new_password"`          // 新密码
}

// UserListRequest 用户列表请求结构体
type UserListRequest struct {
	PageRequest        // 分页参数
	Username    string `form:"username,optional"` // 用户名模糊查询
	Nickname    string `form:"nickname,optional"` // 昵称模糊查询
	Status      int    `form:"status,optional"`   // 状态查询
}

// UserListItem 用户列表项结构体
type UserListItem struct {
	ID        int64  `json:"id"`         // 用户ID
	Username  string `json:"username"`   // 用户名
	Nickname  string `json:"nickname"`   // 昵称
	Avatar    string `json:"avatar"`     // 头像URL
	Email     string `json:"email"`      // 邮箱
	Phone     string `json:"phone"`      // 手机号
	Status    int    `json:"status"`     // 状态
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// UserListResponse 用户列表响应结构体
type UserListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64          `json:"total"`     // 总记录数
		Page     int            `json:"page"`      // 当前页码
		PageSize int            `json:"page_size"` // 每页数量
		List     []UserListItem `json:"list"`      // 用户列表
	} `json:"data"` // 响应数据
}

// RoleInfo 角色信息结构体
type RoleInfo struct {
	ID          int64  `json:"id"`          // 角色ID
	Name        string `json:"name"`        // 角色名称
	Code        string `json:"code"`        // 角色编码
	Description string `json:"description"` // 角色描述
	Status      int    `json:"status"`      // 状态
	Sort        int    `json:"sort"`        // 排序
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
}

// CreateRoleRequest 创建角色请求结构体
type CreateRoleRequest struct {
	Name          string  `json:"name"`                      // 角色名称
	Code          string  `json:"code"`                      // 角色编码
	Description   string  `json:"description,optional"`      // 角色描述
	Status        int     `json:"status,optional,default=1"` // 状态（默认1:正常）
	Sort          int     `json:"sort,optional,default=0"`   // 排序（默认0）
	PermissionIDs []int64 `json:"permission_ids,optional"`   // 权限ID列表
}

// UpdateRoleRequest 更新角色请求结构体
type UpdateRoleRequest struct {
	ID            int64   `json:"id"`                      // 角色ID
	Name          string  `json:"name,optional"`           // 角色名称
	Description   string  `json:"description,optional"`    // 角色描述
	Status        int     `json:"status,optional"`         // 状态
	Sort          int     `json:"sort,optional"`           // 排序
	PermissionIDs []int64 `json:"permission_ids,optional"` // 权限ID列表
}

// RoleListRequest 角色列表请求结构体
type RoleListRequest struct {
	PageRequest        // 分页参数
	Name        string `form:"name,optional"`   // 角色名称模糊查询
	Status      int    `form:"status,optional"` // 状态查询
}

// RoleListResponse 角色列表响应结构体
type RoleListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64      `json:"total"`     // 总记录数
		Page     int        `json:"page"`      // 当前页码
		PageSize int        `json:"page_size"` // 每页数量
		List     []RoleInfo `json:"list"`      // 角色列表
	} `json:"data"` // 响应数据
}

// PermissionInfo 权限信息结构体
type PermissionInfo struct {
	ID        int64  `json:"id"`         // 权限ID
	Name      string `json:"name"`       // 权限名称
	Code      string `json:"code"`       // 权限编码
	Type      string `json:"type"`       // 权限类型
	ParentID  int64  `json:"parent_id"`  // 父权限ID
	Path      string `json:"path"`       // 路径
	Icon      string `json:"icon"`       // 图标
	Component string `json:"component"`  // 组件
	Status    int    `json:"status"`     // 状态
	Sort      int    `json:"sort"`       // 排序
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// CreatePermissionRequest 创建权限请求结构体
type CreatePermissionRequest struct {
	Name      string `json:"name"`                      // 权限名称
	Code      string `json:"code"`                      // 权限编码
	Type      string `json:"type"`                      // 权限类型
	ParentID  int64  `json:"parent_id,optional"`        // 父权限ID
	Path      string `json:"path,optional"`             // 路径
	Icon      string `json:"icon,optional"`             // 图标
	Component string `json:"component,optional"`        // 组件
	Status    int    `json:"status,optional,default=1"` // 状态（默认1:正常）
	Sort      int    `json:"sort,optional,default=0"`   // 排序（默认0）
}

// UpdatePermissionRequest 更新权限请求结构体
type UpdatePermissionRequest struct {
	ID        int64  `json:"id"`                 // 权限ID
	Name      string `json:"name,optional"`      // 权限名称
	Type      string `json:"type,optional"`      // 权限类型
	ParentID  int64  `json:"parent_id,optional"` // 父权限ID
	Path      string `json:"path,optional"`      // 路径
	Icon      string `json:"icon,optional"`      // 图标
	Component string `json:"component,optional"` // 组件
	Status    int    `json:"status,optional"`    // 状态
	Sort      int    `json:"sort,optional"`      // 排序
}

// PermissionListRequest 权限列表请求结构体
type PermissionListRequest struct {
	PageRequest        // 分页参数
	Name        string `form:"name,optional"`   // 权限名称模糊查询
	Type        string `form:"type,optional"`   // 权限类型查询
	Status      int    `form:"status,optional"` // 状态查询
}

// PermissionListResponse 权限列表响应结构体
type PermissionListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64            `json:"total"`     // 总记录数
		Page     int              `json:"page"`      // 当前页码
		PageSize int              `json:"page_size"` // 每页数量
		List     []PermissionInfo `json:"list"`      // 权限列表
	} `json:"data"` // 响应数据
}

// SystemConfigInfo 系统配置信息结构体
type SystemConfigInfo struct {
	ID        int64  `json:"id"`         // 配置ID
	Key       string `json:"key"`        // 配置键
	Value     string `json:"value"`      // 配置值
	Name      string `json:"name"`       // 配置名称
	Remark    string `json:"remark"`     // 配置备注
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// CreateSystemConfigRequest 创建系统配置请求结构体
type CreateSystemConfigRequest struct {
	Key    string `json:"key"`             // 配置键
	Value  string `json:"value"`           // 配置值
	Name   string `json:"name"`            // 配置名称
	Remark string `json:"remark,optional"` // 配置备注
// UpdateSystemConfigRequest 更新系统配置请求结构体
type UpdateSystemConfigRequest struct {
	ID     int64  `json:"id"`                   // 配置ID
	Key    string `json:"key,optional"`        // 配置键
	Value  string `json:"value,optional"`      // 配置值
	Name   string `json:"name,optional"`      // 配置名称
	Remark string `json:"remark,optional"`    // 配置备注
}
// SystemConfigListRequest 系统配置列表请求结构体
type SystemConfigListRequest struct {
	PageRequest                                  // 分页参数
	Key  string `form:"key,optional"`          // 配置键模糊查询
	Name string `form:"name,optional"`        // 配置名称模糊查询
}

// SystemConfigListResponse 系统配置列表响应结构体
type SystemConfigListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64                `json:"total"`    // 总记录数
		Page     int                  `json:"page"`     // 当前页码
		PageSize int                  `json:"page_size"`// 每页数量
		List     []SystemConfigInfo `json:"list"`     // 系统配置列表
	} `json:"data"` // 响应数据
}

// OperationLogInfo 操作日志信息结构体
type OperationLogInfo struct {
	ID            int64  `json:"id"`            // 日志ID
	UserID        int64  `json:"user_id"`      // 用户ID
	Username      string `json:"username"`    // 用户名
	Operation     string `json:"operation"`   // 操作描述
	Method        string `json:"method"`      // HTTP请求方法
	RequestURI    string `json:"request_uri"` // 请求URI
	RequestParams string `json:"request_params"` // 请求参数
	ResponseData  string `json:"response_data"` // 响应数据
	IP            string `json:"ip"`          // 客户端IP地址
	Location      string `json:"location"`    // IP所属地理位置
	Browser       string `json:"browser"`     // 客户端浏览器
	OS            string `json:"os"`          // 客户端操作系统
	Status        int    `json:"status"`      // 请求状态（0:失败，1:成功）
	ErrorMsg      string `json:"error_msg"`   // 错误信息
	Duration      int64  `json:"duration"`    // 请求处理时长（毫秒）
	CreatedAt     string `json:"created_at"`  // 日志创建时间
}

// OperationLogListRequest 操作日志列表请求结构体
type OperationLogListRequest struct {
	PageRequest                                  // 分页参数
	Username  string `form:"username,optional"`  // 用户名模糊查询
	Operation string `form:"operation,optional"` // 操作描述模糊查询
	Method    string `form:"method,optional"`    // 请求方法查询
	Status    int    `form:"status,optional"`    // 状态查询
	StartTime string `form:"start_time,optional"`// 开始时间
	EndTime   string `form:"end_time,optional"`  // 结束时间
}

// OperationLogListResponse 操作日志列表响应结构体
type OperationLogListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64               `json:"total"`    // 总记录数
		Page     int                 `json:"page"`     // 当前页码
		PageSize int                 `json:"page_size"`// 每页数量
		List     []OperationLogInfo `json:"list"`     // 操作日志列表
	} `json:"data"` // 响应数据
}

type LoginLogInfo struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Status    int    `json:"status"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"created_at"`
}

type LoginLogListRequest struct {
	PageRequest
	Username  string `form:"username,optional"`
	Status    int    `form:"status,optional"`
	StartTime string `form:"start_time,optional"`
	EndTime   string `form:"end_time,optional"`
}

type LoginLogListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total    int64          `json:"total"`
		Page     int            `json:"page"`
		PageSize int            `json:"page_size"`
		List     []LoginLogInfo `json:"list"`
	} `json:"data"`
}

// FileInfo 文件信息结构体
type FileInfo struct {
	ID           int64  `json:"id"`           // 文件ID
	Name         string `json:"name"`         // 存储文件名
	OriginalName string `json:"original_name"`// 原始文件名
	Path         string `json:"path"`         // 文件存储路径
	URL          string `json:"url"`          // 文件访问URL
	Size         int64  `json:"size"`         // 文件大小（字节）
	Type         string `json:"type"`         // 文件MIME类型
	Extension    string `json:"extension"`    // 文件扩展名
	MD5          string `json:"md5"`          // 文件MD5哈希值
	UserID       int64  `json:"user_id"`      // 上传用户ID
	CreatedAt    string `json:"created_at"`   // 上传时间
}

// UploadFileResponse 文件上传响应结构体
type UploadFileResponse struct {
	Code    int      `json:"code"`    // 响应码
	Message string   `json:"message"` // 响应消息
	Data    FileInfo `json:"data"`    // 文件信息
}

// FileListRequest 文件列表请求结构体
type FileListRequest struct {
	PageRequest                                  // 分页参数
	Name string `form:"name,optional"`        // 文件名模糊查询
	Type string `form:"type,optional"`        // 文件类型查询
}

// FileListResponse 文件列表响应结构体
type FileListResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    struct {
		Total    int64       `json:"total"`    // 总记录数
		Page     int         `json:"page"`     // 当前页码
		PageSize int         `json:"page_size"`// 每页数量
		List     []FileInfo `json:"list"`     // 文件列表
	} `json:"data"` // 响应数据
}
