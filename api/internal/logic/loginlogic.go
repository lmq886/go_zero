package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go_zero/api/internal/middleware"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
	"go_zero/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

// LoginLogic 登录业务逻辑
// 处理用户登录的核心业务逻辑
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic 创建登录逻辑实例
// 参数 ctx: 上下文
// 参数 svcCtx: 服务上下文
// 返回值: 登录逻辑实例
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 执行登录逻辑
// 参数 req: 登录请求
// 返回值: 登录响应和错误信息
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	// 1. 验证参数
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	// 2. 根据用户名查找用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		l.Errorf("Failed to find user by username: %v", err)
		return nil, errors.New("登录失败，请稍后重试")
	}

	// 3. 验证用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用，请联系管理员")
	}

	// 4. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 记录登录失败日志
		l.recordLoginLog(user.Id, user.Username, 0, "密码错误")
		return nil, errors.New("用户名或密码错误")
	}

	// 5. 获取用户角色
	roles, err := l.svcCtx.UserModel.FindUserRoles(l.ctx, user.Id)
	if err != nil {
		l.Errorf("Failed to find user roles: %v", err)
		return nil, errors.New("登录失败，请稍后重试")
	}

	// 6. 获取用户权限
	permissions, err := l.svcCtx.UserModel.FindUserPermissions(l.ctx, user.Id)
	if err != nil {
		l.Errorf("Failed to find user permissions: %v", err)
		return nil, errors.New("登录失败，请稍后重试")
	}

	// 7. 构建角色和权限列表
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	// 8. 生成 JWT 令牌
	accessToken, refreshToken, expiresAt, err := middleware.GenerateToken(
		user.Id,
		user.Username,
		roleCodes,
		permissionCodes,
		l.svcCtx.Config.Auth,
	)
	if err != nil {
		l.Errorf("Failed to generate token: %v", err)
		return nil, errors.New("登录失败，请稍后重试")
	}

	// 9. 更新用户最后登录信息
	now := time.Now().Unix()
	user.LastLoginAt = sql.NullInt64{Int64: now, Valid: true}
	user.LastLoginIp = sql.NullString{String: l.getClientIP(), Valid: true}
	if err := l.svcCtx.UserModel.Update(l.ctx, user); err != nil {
		l.Errorf("Failed to update user last login info: %v", err)
		// 不影响登录流程，只记录错误
	}

	// 10. 记录登录成功日志
	l.recordLoginLog(user.Id, user.Username, 1, "登录成功")

	// 11. 构建响应
	return &types.LoginResp{
		Code:         200,
		Message:      "登录成功",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UserInfo: types.UserInfo{
			Id:          user.Id,
			Username:    user.Username,
			Email:       nullStringToString(user.Email),
			Phone:       nullStringToString(user.Phone),
			Nickname:    nullStringToString(user.Nickname),
			Avatar:      nullStringToString(user.Avatar),
			Status:      int(user.Status),
			Roles:       roleCodes,
			Permissions: permissionCodes,
			CreatedAt:   user.CreatedAt,
		},
	}, nil
}

// recordLoginLog 记录登录日志
// 参数 userId: 用户ID
// 参数 username: 用户名
// 参数 status: 状态（1-成功，0-失败）
// 参数 message: 消息
func (l *LoginLogic) recordLoginLog(userId int64, username string, status int, message string) {
	// 如果登录日志未启用，直接返回
	if !l.svcCtx.Config.LoginLog.Enabled {
		return
	}

	// 构建登录日志
	log := &model.LoginLog{
		UserId:    userId,
		Username:  username,
		Ip:        l.getClientIP(),
		UserAgent: l.getUserAgent(),
		Status:    int64(status),
		Message:   message,
	}

	// 插入数据库
	if err := l.svcCtx.LoginLogModel.Insert(l.ctx, log); err != nil {
		l.Errorf("Failed to insert login log: %v", err)
	}
}

// getClientIP 获取客户端 IP 地址
// 返回值: 客户端 IP 地址
func (l *LoginLogic) getClientIP() string {
	// 从上下文中获取 IP（实际项目中应该从请求中获取）
	return "127.0.0.1"
}

// getUserAgent 获取用户代理
// 返回值: 用户代理字符串
func (l *LoginLogic) getUserAgent() string {
	// 从上下文中获取用户代理（实际项目中应该从请求中获取）
	return "Unknown"
}

// nullStringToString 将 sql.NullString 转换为 string
// 参数 ns: sql.NullString
// 返回值: 字符串值
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// RegisterLogic 注册业务逻辑
// 处理用户注册的核心业务逻辑
type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 创建注册逻辑实例
// 参数 ctx: 上下文
// 参数 svcCtx: 服务上下文
// 返回值: 注册逻辑实例
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 执行注册逻辑
// 参数 req: 注册请求
// 返回值: 注册响应和错误信息
func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	// 1. 检查是否允许注册
	if !l.svcCtx.Config.System.AllowRegister {
		return nil, errors.New("系统暂不开放注册功能")
	}

	// 2. 验证参数
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	// 3. 验证密码长度
	if len(req.Password) < l.svcCtx.Config.System.PasswordMinLength ||
		len(req.Password) > l.svcCtx.Config.System.PasswordMaxLength {
		return nil, errors.New("密码长度必须在 6-20 个字符之间")
	}

	// 4. 检查用户名是否已存在
	existingUser, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("Failed to check username existence: %v", err)
		return nil, errors.New("注册失败，请稍后重试")
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 5. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("Failed to hash password: %v", err)
		return nil, errors.New("注册失败，请稍后重试")
	}

	// 6. 构建用户数据
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Status:   1,
	}

	// 设置可选字段
	if req.Email != "" {
		user.Email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.Phone != "" {
		user.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	if req.Nickname != "" {
		user.Nickname = sql.NullString{String: req.Nickname, Valid: true}
	}

	// 7. 插入用户
	result, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		l.Errorf("Failed to insert user: %v", err)
		return nil, errors.New("注册失败，请稍后重试")
	}

	// 8. 获取用户ID
	userId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Failed to get last insert id: %v", err)
		return nil, errors.New("注册失败，请稍后重试")
	}

	// 9. 为新用户分配默认角色（普通用户）
	defaultRole, err := l.svcCtx.RoleModel.FindOneByCode(l.ctx, "user")
	if err != nil && err != model.ErrNotFound {
		l.Errorf("Failed to find default role: %v", err)
		// 不影响注册流程，只记录错误
	}
	if defaultRole != nil {
		if err := l.svcCtx.UserModel.AssignRoles(l.ctx, userId, []int64{defaultRole.Id}); err != nil {
			l.Errorf("Failed to assign default role: %v", err)
			// 不影响注册流程，只记录错误
		}
	}

	// 10. 返回成功响应
	return &types.RegisterResp{
		Code:    200,
		Message: "注册成功",
		UserId:  userId,
	}, nil
}

// LogoutLogic 登出业务逻辑
// 处理用户登出的核心业务逻辑
type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLogoutLogic 创建登出逻辑实例
// 参数 ctx: 上下文
// 参数 svcCtx: 服务上下文
// 返回值: 登出逻辑实例
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Logout 执行登出逻辑
// 参数 req: 登出请求
// 返回值: 登出响应和错误信息
func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.LogoutResp, error) {
	// 从上下文中获取用户ID
	userId, ok := middleware.GetUserIdFromContext(l.ctx)
	if !ok {
		return nil, errors.New("未授权")
	}

	// 记录登出日志
	l.Infof("User %d logged out", userId)

	// 返回成功响应
	return &types.LogoutResp{
		Code:    200,
		Message: "登出成功",
	}, nil
}

// RefreshTokenLogic 刷新令牌业务逻辑
// 处理令牌刷新的核心业务逻辑
type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRefreshTokenLogic 创建刷新令牌逻辑实例
// 参数 ctx: 上下文
// 参数 svcCtx: 服务上下文
// 返回值: 刷新令牌逻辑实例
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RefreshToken 执行令牌刷新逻辑
// 参数 req: 刷新令牌请求
// 返回值: 刷新令牌响应和错误信息
func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (*types.RefreshTokenResp, error) {
	// 1. 验证参数
	if req.RefreshToken == "" {
		return nil, errors.New("刷新令牌不能为空")
	}

	// 2. 解析刷新令牌
	userId, _, err := middleware.ParseRefreshToken(req.RefreshToken, l.svcCtx.Config.Auth)
	if err != nil {
		l.Errorf("Failed to parse refresh token: %v", err)
		return nil, errors.New("刷新令牌无效或已过期")
	}

	// 3. 验证用户是否存在
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("Failed to find user: %v", err)
		return nil, errors.New("刷新令牌失败，请稍后重试")
	}

	// 4. 验证用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用，请联系管理员")
	}

	// 5. 获取用户角色
	roles, err := l.svcCtx.UserModel.FindUserRoles(l.ctx, user.Id)
	if err != nil {
		l.Errorf("Failed to find user roles: %v", err)
		return nil, errors.New("刷新令牌失败，请稍后重试")
	}

	// 6. 获取用户权限
	permissions, err := l.svcCtx.UserModel.FindUserPermissions(l.ctx, user.Id)
	if err != nil {
		l.Errorf("Failed to find user permissions: %v", err)
		return nil, errors.New("刷新令牌失败，请稍后重试")
	}

	// 7. 构建角色和权限列表
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	// 8. 生成新的访问令牌
	accessToken, _, expiresAt, err := middleware.GenerateToken(
		user.Id,
		user.Username,
		roleCodes,
		permissionCodes,
		l.svcCtx.Config.Auth,
	)
	if err != nil {
		l.Errorf("Failed to generate new access token: %v", err)
		return nil, errors.New("刷新令牌失败，请稍后重试")
	}

	// 9. 返回成功响应
	return &types.RefreshTokenResp{
		Code:        200,
		Message:     "刷新令牌成功",
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}
