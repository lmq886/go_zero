-- 企业级后台管理系统数据库模式
-- 数据库：PostgreSQL
-- 作者：GoZero Enterprise Team
-- 版本：v1.0.0

-- ==================== 系统配置 ====================
-- 创建扩展（如果需要）
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==================== 核心表 ====================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(20),
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    gender INTEGER DEFAULT 0,
    birthday VARCHAR(20),
    address VARCHAR(255),
    introduction TEXT,
    status INTEGER NOT NULL DEFAULT 1,
    last_login_at BIGINT,
    last_login_ip VARCHAR(50),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    status INTEGER NOT NULL DEFAULT 1,
    sort INTEGER NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    parent_id BIGINT DEFAULT 0,
    description VARCHAR(255),
    path VARCHAR(255),
    method VARCHAR(20),
    status INTEGER NOT NULL DEFAULT 1,
    sort INTEGER NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT
);

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    parent_id BIGINT DEFAULT 0,
    path VARCHAR(255),
    component VARCHAR(255),
    icon VARCHAR(100),
    title VARCHAR(100) NOT NULL,
    redirect VARCHAR(255),
    hidden BOOLEAN NOT NULL DEFAULT FALSE,
    always_show BOOLEAN NOT NULL DEFAULT FALSE,
    permission VARCHAR(100),
    status INTEGER NOT NULL DEFAULT 1,
    sort INTEGER NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    UNIQUE(user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    UNIQUE(role_id, permission_id)
);

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS role_menus (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    menu_id BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    UNIQUE(role_id, menu_id)
);

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(50) NOT NULL,
    module VARCHAR(50) NOT NULL,
    operation VARCHAR(100) NOT NULL,
    method VARCHAR(20) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    ip VARCHAR(50) NOT NULL,
    user_agent VARCHAR(500) NOT NULL,
    request_data TEXT,
    response_data TEXT,
    error_msg TEXT,
    duration BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL
);

-- 登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(50) NOT NULL,
    ip VARCHAR(50) NOT NULL,
    location VARCHAR(100),
    browser VARCHAR(100),
    os VARCHAR(100),
    user_agent VARCHAR(500) NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    message VARCHAR(255),
    created_at BIGINT NOT NULL
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS configs (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    type VARCHAR(20) NOT NULL,
    group_name VARCHAR(50) NOT NULL DEFAULT 'default',
    description VARCHAR(255),
    status INTEGER NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT
);

-- ==================== 索引创建 ====================

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 角色表索引
CREATE INDEX IF NOT EXISTS idx_roles_code ON roles(code);
CREATE INDEX IF NOT EXISTS idx_roles_status ON roles(status);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);

-- 权限表索引
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code);
CREATE INDEX IF NOT EXISTS idx_permissions_parent_id ON permissions(parent_id);
CREATE INDEX IF NOT EXISTS idx_permissions_type ON permissions(type);
CREATE INDEX IF NOT EXISTS idx_permissions_status ON permissions(status);
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions(deleted_at);

-- 菜单表索引
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_type ON menus(type);
CREATE INDEX IF NOT EXISTS idx_menus_status ON menus(status);
CREATE INDEX IF NOT EXISTS idx_menus_deleted_at ON menus(deleted_at);

-- 用户角色关联表索引
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- 角色权限关联表索引
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

-- 角色菜单关联表索引
CREATE INDEX IF NOT EXISTS idx_role_menus_role_id ON role_menus(role_id);
CREATE INDEX IF NOT EXISTS idx_role_menus_menu_id ON role_menus(menu_id);

-- 操作日志表索引
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_username ON operation_logs(username);
CREATE INDEX IF NOT EXISTS idx_operation_logs_module ON operation_logs(module);
CREATE INDEX IF NOT EXISTS idx_operation_logs_status ON operation_logs(status);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);

-- 登录日志表索引
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_username ON login_logs(username);
CREATE INDEX IF NOT EXISTS idx_login_logs_status ON login_logs(status);
CREATE INDEX IF NOT EXISTS idx_login_logs_created_at ON login_logs(created_at);

-- 系统配置表索引
CREATE INDEX IF NOT EXISTS idx_configs_key ON configs(key);
CREATE INDEX IF NOT EXISTS idx_configs_group_name ON configs(group_name);
CREATE INDEX IF NOT EXISTS idx_configs_status ON configs(status);
CREATE INDEX IF NOT EXISTS idx_configs_deleted_at ON configs(deleted_at);

-- ==================== 初始化数据 ====================

-- 插入默认管理员用户（密码：admin123，使用 bcrypt 加密）
INSERT INTO users (username, password, email, phone, nickname, status, created_at, updated_at)
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5E', 'admin@example.com', '13800138000', '系统管理员', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT)
ON CONFLICT (username) DO NOTHING;

-- 插入默认角色
INSERT INTO roles (name, code, description, status, sort, created_at, updated_at)
VALUES 
('超级管理员', 'super_admin', '拥有系统所有权限', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('管理员', 'admin', '拥有大部分管理权限', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('普通用户', 'user', '普通用户权限', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT)
ON CONFLICT (code) DO NOTHING;

-- 为管理员用户分配超级管理员角色
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE NOT EXISTS (SELECT 1 FROM user_roles WHERE user_id = 1 AND role_id = 1);

-- 插入默认菜单
INSERT INTO menus (name, type, parent_id, path, component, icon, title, redirect, hidden, always_show, permission, status, sort, created_at, updated_at)
VALUES 
-- 一级菜单
('系统管理', 'directory', 0, '/system', '', 'Setting', '系统管理', '', false, true, '', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('用户管理', 'menu', 1, '/system/user', 'system/user/index', 'User', '用户管理', '', false, false, 'system:user:list', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('角色管理', 'menu', 1, '/system/role', 'system/role/index', 'UserFilled', '角色管理', '', false, false, 'system:role:list', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('菜单管理', 'menu', 1, '/system/menu', 'system/menu/index', 'Menu', '菜单管理', '', false, false, 'system:menu:list', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('权限管理', 'menu', 1, '/system/permission', 'system/permission/index', 'Key', '权限管理', '', false, false, 'system:permission:list', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('日志管理', 'directory', 0, '/log', '', 'Document', '日志管理', '', false, true, '', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('操作日志', 'menu', 6, '/log/operation', 'log/operation/index', 'List', '操作日志', '', false, false, 'log:operation:list', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('登录日志', 'menu', 6, '/log/login', 'log/login/index', 'User', '登录日志', '', false, false, 'log:login:list', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('系统配置', 'directory', 0, '/config', '', 'Tools', '系统配置', '', false, true, '', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('配置管理', 'menu', 9, '/config/list', 'config/list/index', 'Setting', '配置管理', '', false, false, 'config:list', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT);

-- 插入默认权限
INSERT INTO permissions (name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at)
VALUES 
-- 用户管理权限
('用户列表', 'system:user:list', 'api', 0, '查看用户列表', '/api/v1/users', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('用户详情', 'system:user:detail', 'api', 0, '查看用户详情', '/api/v1/users/:id', 'GET', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('创建用户', 'system:user:create', 'api', 0, '创建用户', '/api/v1/users', 'POST', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('更新用户', 'system:user:update', 'api', 0, '更新用户', '/api/v1/users/:id', 'PUT', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除用户', 'system:user:delete', 'api', 0, '删除用户', '/api/v1/users/:id', 'DELETE', 1, 5, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('重置密码', 'system:user:reset-password', 'api', 0, '重置用户密码', '/api/v1/users/:id/reset-password', 'POST', 1, 6, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
-- 角色管理权限
('角色列表', 'system:role:list', 'api', 0, '查看角色列表', '/api/v1/roles', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('角色详情', 'system:role:detail', 'api', 0, '查看角色详情', '/api/v1/roles/:id', 'GET', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('创建角色', 'system:role:create', 'api', 0, '创建角色', '/api/v1/roles', 'POST', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('更新角色', 'system:role:update', 'api', 0, '更新角色', '/api/v1/roles/:id', 'PUT', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除角色', 'system:role:delete', 'api', 0, '删除角色', '/api/v1/roles/:id', 'DELETE', 1, 5, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('分配权限', 'system:role:assign-permissions', 'api', 0, '为角色分配权限', '/api/v1/roles/:id/permissions', 'POST', 1, 6, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
-- 菜单管理权限
('菜单列表', 'system:menu:list', 'api', 0, '查看菜单列表', '/api/v1/menus', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('菜单详情', 'system:menu:detail', 'api', 0, '查看菜单详情', '/api/v1/menus/:id', 'GET', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('创建菜单', 'system:menu:create', 'api', 0, '创建菜单', '/api/v1/menus', 'POST', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('更新菜单', 'system:menu:update', 'api', 0, '更新菜单', '/api/v1/menus/:id', 'PUT', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除菜单', 'system:menu:delete', 'api', 0, '删除菜单', '/api/v1/menus/:id', 'DELETE', 1, 5, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
-- 权限管理权限
('权限列表', 'system:permission:list', 'api', 0, '查看权限列表', '/api/v1/permissions', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('权限详情', 'system:permission:detail', 'api', 0, '查看权限详情', '/api/v1/permissions/:id', 'GET', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('创建权限', 'system:permission:create', 'api', 0, '创建权限', '/api/v1/permissions', 'POST', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('更新权限', 'system:permission:update', 'api', 0, '更新权限', '/api/v1/permissions/:id', 'PUT', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除权限', 'system:permission:delete', 'api', 0, '删除权限', '/api/v1/permissions/:id', 'DELETE', 1, 5, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
-- 日志管理权限
('操作日志列表', 'log:operation:list', 'api', 0, '查看操作日志列表', '/api/v1/logs/operation', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除操作日志', 'log:operation:delete', 'api', 0, '删除操作日志', '/api/v1/logs/operation/:id', 'DELETE', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('清空操作日志', 'log:operation:clear', 'api', 0, '清空操作日志', '/api/v1/logs/operation', 'DELETE', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('登录日志列表', 'log:login:list', 'api', 0, '查看登录日志列表', '/api/v1/logs/login', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除登录日志', 'log:login:delete', 'api', 0, '删除登录日志', '/api/v1/logs/login/:id', 'DELETE', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('清空登录日志', 'log:login:clear', 'api', 0, '清空登录日志', '/api/v1/logs/login', 'DELETE', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
-- 系统配置权限
('配置列表', 'config:list', 'api', 0, '查看配置列表', '/api/v1/configs', 'GET', 1, 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('配置详情', 'config:detail', 'api', 0, '查看配置详情', '/api/v1/configs/:key', 'GET', 1, 2, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('创建配置', 'config:create', 'api', 0, '创建配置', '/api/v1/configs', 'POST', 1, 3, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('更新配置', 'config:update', 'api', 0, '更新配置', '/api/v1/configs/:key', 'PUT', 1, 4, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('删除配置', 'config:delete', 'api', 0, '删除配置', '/api/v1/configs/:key', 'DELETE', 1, 5, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT);

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 1, id, EXTRACT(EPOCH FROM NOW())::BIGINT
FROM permissions
WHERE NOT EXISTS (
    SELECT 1 FROM role_permissions rp 
    WHERE rp.role_id = 1 AND rp.permission_id = permissions.id
);

-- 为超级管理员角色分配所有菜单
INSERT INTO role_menus (role_id, menu_id, created_at)
SELECT 1, id, EXTRACT(EPOCH FROM NOW())::BIGINT
FROM menus
WHERE NOT EXISTS (
    SELECT 1 FROM role_menus rm 
    WHERE rm.role_id = 1 AND rm.menu_id = menus.id
);

-- 插入默认系统配置
INSERT INTO configs (key, name, value, type, group_name, description, status, created_at, updated_at)
VALUES 
('site.name', '站点名称', '企业级后台管理系统', 'string', 'site', '网站名称', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('site.logo', '站点Logo', '', 'string', 'site', '网站Logo地址', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('site.description', '站点描述', '基于GoZero框架开发的企业级后台管理系统', 'string', 'site', '网站描述', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('system.register.enabled', '是否允许注册', 'true', 'boolean', 'system', '是否允许用户注册', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('system.login.maxRetries', '登录最大重试次数', '5', 'number', 'system', '登录失败最大重试次数', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('system.login.lockDuration', '登录锁定时长(分钟)', '30', 'number', 'system', '登录失败锁定时长', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('jwt.accessToken.expire', '访问令牌有效期(秒)', '7200', 'number', 'jwt', 'JWT访问令牌有效期', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('jwt.refreshToken.expire', '刷新令牌有效期(秒)', '604800', 'number', 'jwt', 'JWT刷新令牌有效期', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('log.operation.enabled', '操作日志开关', 'true', 'boolean', 'log', '是否开启操作日志记录', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('log.login.enabled', '登录日志开关', 'true', 'boolean', 'log', '是否开启登录日志记录', 1, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT)
ON CONFLICT (key) DO NOTHING;
