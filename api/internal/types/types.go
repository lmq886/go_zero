package types

// ==================== 通用类型定义 ====================

// PaginationReq 分页请求
type PaginationReq struct {
	Page     int `json:"page,optional" form:"page,optional"`         // 页码，默认 1
	PageSize int `json:"pageSize,optional" form:"pageSize,optional"` // 每页数量，默认 10
}

// PaginationResp 分页响应
type PaginationResp struct {
	Total    int64 `json:"total"`    // 总记录数
	Page     int   `json:"page"`     // 当前页码
	PageSize int   `json:"pageSize"` // 每页数量
}

// ==================== 认证模块类型定义 ====================

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginResp 登录响应
type LoginResp struct {
	Code         int      `json:"code"`         // 状态码
	Message      string   `json:"message"`      // 消息
	AccessToken  string   `json:"accessToken"`  // 访问令牌
	RefreshToken string   `json:"refreshToken"` // 刷新令牌
	ExpiresAt    int64    `json:"expiresAt"`    // 过期时间
	UserInfo     UserInfo `json:"userInfo"`     // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	Id          int64    `json:"id"`          // 用户ID
	Username    string   `json:"username"`    // 用户名
	Email       string   `json:"email"`       // 邮箱
	Phone       string   `json:"phone"`       // 手机号
	Nickname    string   `json:"nickname"`    // 昵称
	Avatar      string   `json:"avatar"`      // 头像
	Status      int      `json:"status"`      // 状态
	Roles       []string `json:"roles"`       // 角色列表
	Permissions []string `json:"permissions"` // 权限列表
	CreatedAt   int64    `json:"createdAt"`   // 创建时间
}

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username"`          // 用户名
	Password string `json:"password"`          // 密码
	Email    string `json:"email,optional"`    // 邮箱
	Phone    string `json:"phone,optional"`    // 手机号
	Nickname string `json:"nickname,optional"` // 昵称
}

// RegisterResp 注册响应
type RegisterResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
	UserId  int64  `json:"userId"`  // 用户ID
}

// LogoutReq 登出请求
type LogoutReq struct {
}

// LogoutResp 登出响应
type LogoutResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
}

// RefreshTokenReq 刷新令牌请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken"` // 刷新令牌
}

// RefreshTokenResp 刷新令牌响应
type RefreshTokenResp struct {
	Code        int    `json:"code"`        // 状态码
	Message     string `json:"message"`     // 消息
	AccessToken string `json:"accessToken"` // 访问令牌
	ExpiresAt   int64  `json:"expiresAt"`   // 过期时间
}

// ==================== 用户管理模块类型定义 ====================

// GetUserListReq 获取用户列表请求
type GetUserListReq struct {
	PaginationReq
	Username string `json:"username,optional" form:"username,optional"` // 用户名（模糊查询）
	Email    string `json:"email,optional" form:"email,optional"`       // 邮箱（模糊查询）
	Phone    string `json:"phone,optional" form:"phone,optional"`       // 手机号（模糊查询）
	Status   int    `json:"status,optional" form:"status,optional"`     // 状态
}

// GetUserListResp 获取用户列表响应
type GetUserListResp struct {
	Code       int            `json:"code"`       // 状态码
	Message    string         `json:"message"`    // 消息
	Pagination PaginationResp `json:"pagination"` // 分页信息
	Data       []UserListItem `json:"data"`       // 用户列表
}

// UserListItem 用户列表项
type UserListItem struct {
	Id        int64    `json:"id"`        // 用户ID
	Username  string   `json:"username"`  // 用户名
	Email     string   `json:"email"`     // 邮箱
	Phone     string   `json:"phone"`     // 手机号
	Nickname  string   `json:"nickname"`  // 昵称
	Avatar    string   `json:"avatar"`    // 头像
	Status    int      `json:"status"`    // 状态
	Roles     []string `json:"roles"`     // 角色列表
	CreatedAt int64    `json:"createdAt"` // 创建时间
	UpdatedAt int64    `json:"updatedAt"` // 更新时间
}

// GetUserByIdReq 根据ID获取用户请求
type GetUserByIdReq struct {
	Id int64 `json:"id,optional" path:"id"` // 用户ID
}

// GetUserByIdResp 根据ID获取用户响应
type GetUserByIdResp struct {
	Code    int        `json:"code"`    // 状态码
	Message string     `json:"message"` // 消息
	Data    UserDetail `json:"data"`    // 用户详情
}

// UserDetail 用户详情
type UserDetail struct {
	Id           int64      `json:"id"`           // 用户ID
	Username     string     `json:"username"`     // 用户名
	Email        string     `json:"email"`        // 邮箱
	Phone        string     `json:"phone"`        // 手机号
	Nickname     string     `json:"nickname"`     // 昵称
	Avatar       string     `json:"avatar"`       // 头像
	Gender       int        `json:"gender"`       // 性别
	Birthday     string     `json:"birthday"`     // 生日
	Address      string     `json:"address"`      // 地址
	Introduction string     `json:"introduction"` // 简介
	Status       int        `json:"status"`       // 状态
	Roles        []RoleInfo `json:"roles"`        // 角色列表
	CreatedAt    int64      `json:"createdAt"`    // 创建时间
	UpdatedAt    int64      `json:"updatedAt"`    // 更新时间
}

// RoleInfo 角色信息
type RoleInfo struct {
	Id          int64  `json:"id"`          // 角色ID
	Name        string `json:"name"`        // 角色名称
	Code        string `json:"code"`        // 角色编码
	Description string `json:"description"` // 角色描述
}

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username string  `json:"username"`          // 用户名
	Password string  `json:"password"`          // 密码
	Email    string  `json:"email,optional"`    // 邮箱
	Phone    string  `json:"phone,optional"`    // 手机号
	Nickname string  `json:"nickname,optional"` // 昵称
	Avatar   string  `json:"avatar,optional"`   // 头像
	Status   int     `json:"status,optional"`   // 状态
	RoleIds  []int64 `json:"roleIds,optional"`  // 角色ID列表
}

// CreateUserResp 创建用户响应
type CreateUserResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
	UserId  int64  `json:"userId"`  // 用户ID
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	Id           int64   `json:"id,optional" path:"id"` // 用户ID
	Email        string  `json:"email,optional"`        // 邮箱
	Phone        string  `json:"phone,optional"`        // 手机号
	Nickname     string  `json:"nickname,optional"`     // 昵称
	Avatar       string  `json:"avatar,optional"`       // 头像
	Gender       int     `json:"gender,optional"`       // 性别
	Birthday     string  `json:"birthday,optional"`     // 生日
	Address      string  `json:"address,optional"`      // 地址
	Introduction string  `json:"introduction,optional"` // 简介
	Status       int     `json:"status,optional"`       // 状态
	RoleIds      []int64 `json:"roleIds,optional"`      // 角色ID列表
}

// UpdateUserResp 更新用户响应
type UpdateUserResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
}

// DeleteUserReq 删除用户请求
type DeleteUserReq struct {
	Id int64 `json:"id,optional" path:"id"` // 用户ID
}

// DeleteUserResp 删除用户响应
type DeleteUserResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	Id          int64  `json:"id,optional" path:"id"` // 用户ID
	NewPassword string `json:"newPassword"`           // 新密码
}

// ResetPasswordResp 重置密码响应
type ResetPasswordResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
}

// UpdateProfileReq 更新个人资料请求
type UpdateProfileReq struct {
	Nickname     string `json:"nickname,optional"`     // 昵称
	Avatar       string `json:"avatar,optional"`       // 头像
	Gender       int    `json:"gender,optional"`       // 性别
	Birthday     string `json:"birthday,optional"`     // 生日
	Address      string `json:"address,optional"`      // 地址
	Introduction string `json:"introduction,optional"` // 简介
}

// UpdateProfileResp 更新个人资料响应
type UpdateProfileResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
}

// GetProfileReq 获取个人资料请求
type GetProfileReq struct {
}

// GetProfileResp 获取个人资料响应
type GetProfileResp struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    UserProfile `json:"data"`    // 个人资料
}

// UserProfile 个人资料
type UserProfile struct {
	Id           int64    `json:"id"`           // 用户ID
	Username     string   `json:"username"`     // 用户名
	Email        string   `json:"email"`        // 邮箱
	Phone        string   `json:"phone"`        // 手机号
	Nickname     string   `json:"nickname"`     // 昵称
	Avatar       string   `json:"avatar"`       // 头像
	Gender       int      `json:"gender"`       // 性别
	Birthday     string   `json:"birthday"`     // 生日
	Address      string   `json:"address"`      // 地址
	Introduction string   `json:"introduction"` // 简介
	Roles        []string `json:"roles"`        // 角色列表
	Permissions  []string `json:"permissions"`  // 权限列表
	CreatedAt    int64    `json:"createdAt"`    // 创建时间
	UpdatedAt    int64    `json:"updatedAt"`    // 更新时间
}
