-- 管理后台数据库初始化脚本
-- 数据库: PostgreSQL
-- 版本: 1.0.0

-- 创建枚举类型
CREATE TYPE permission_type AS ENUM ('menu', 'button', 'api');
CREATE TYPE log_status AS ENUM ('success', 'fail');

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    email VARCHAR(100),
    phone VARCHAR(20),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    status SMALLINT DEFAULT 1,
    sort INT DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    type permission_type DEFAULT 'menu',
    parent_id BIGINT DEFAULT 0,
    path VARCHAR(255),
    icon VARCHAR(100),
    component VARCHAR(255),
    status SMALLINT DEFAULT 1,
    sort INT DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    remark VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    username VARCHAR(50),
    operation VARCHAR(100) NOT NULL,
    method VARCHAR(10) NOT NULL,
    request_uri VARCHAR(255) NOT NULL,
    request_params TEXT,
    response_data TEXT,
    ip VARCHAR(50),
    location VARCHAR(100),
    browser VARCHAR(100),
    os VARCHAR(50),
    status SMALLINT DEFAULT 1,
    error_msg TEXT,
    duration BIGINT DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    username VARCHAR(50),
    ip VARCHAR(50),
    location VARCHAR(100),
    browser VARCHAR(100),
    os VARCHAR(50),
    status SMALLINT DEFAULT 1,
    msg VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 文件表
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    url VARCHAR(500),
    size BIGINT NOT NULL,
    type VARCHAR(100) NOT NULL,
    extension VARCHAR(20),
    md5 VARCHAR(32),
    user_id BIGINT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_roles_code ON roles(code);
CREATE INDEX IF NOT EXISTS idx_permissions_parent_id ON permissions(parent_id);
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_system_configs_key ON system_configs(key);
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_created_at ON login_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_files_md5 ON files(md5);
CREATE INDEX IF NOT EXISTS idx_files_user_id ON files(user_id);

-- 初始化数据
-- 插入超级管理员角色
INSERT INTO roles (name, code, description, status, sort) VALUES
('超级管理员', 'super_admin', '拥有系统所有权限', 1, 1)
ON CONFLICT (code) DO NOTHING;

-- 插入默认权限
INSERT INTO permissions (name, code, type, parent_id, path, icon, component, status, sort) VALUES
('系统管理', 'system', 'menu', 0, '/system', 'Setting', 'Layout', 1, 1),
('用户管理', 'system:user', 'menu', 1, '/system/user', 'User', 'system/user/index', 1, 1),
('用户列表', 'system:user:list', 'api', 2, '', '', '', 1, 1),
('用户新增', 'system:user:add', 'api', 2, '', '', '', 1, 2),
('用户编辑', 'system:user:edit', 'api', 2, '', '', '', 1, 3),
('用户删除', 'system:user:delete', 'api', 2, '', '', '', 1, 4),
('角色管理', 'system:role', 'menu', 1, '/system/role', 'Shield', 'system/role/index', 1, 2),
('角色列表', 'system:role:list', 'api', 7, '', '', '', 1, 1),
('角色新增', 'system:role:add', 'api', 7, '', '', '', 1, 2),
('角色编辑', 'system:role:edit', 'api', 7, '', '', '', 1, 3),
('角色删除', 'system:role:delete', 'api', 7, '', '', '', 1, 4),
('权限管理', 'system:permission', 'menu', 1, '/system/permission', 'Key', 'system/permission/index', 1, 3),
('权限列表', 'system:permission:list', 'api', 12, '', '', '', 1, 1),
('权限新增', 'system:permission:add', 'api', 12, '', '', '', 1, 2),
('权限编辑', 'system:permission:edit', 'api', 12, '', '', '', 1, 3),
('权限删除', 'system:permission:delete', 'api', 12, '', '', '', 1, 4),
('系统配置', 'system:config', 'menu', 1, '/system/config', 'Tool', 'system/config/index', 1, 4),
('配置列表', 'system:config:list', 'api', 17, '', '', '', 1, 1),
('配置新增', 'system:config:add', 'api', 17, '', '', '', 1, 2),
('配置编辑', 'system:config:edit', 'api', 17, '', '', '', 1, 3),
('配置删除', 'system:config:delete', 'api', 17, '', '', '', 1, 4),
('日志管理', 'log', 'menu', 0, '/log', 'Document', 'Layout', 1, 2),
('操作日志', 'log:operation', 'menu', 22, '/log/operation', 'List', 'log/operation/index', 1, 1),
('操作日志列表', 'log:operation:list', 'api', 23, '', '', '', 1, 1),
('操作日志删除', 'log:operation:delete', 'api', 23, '', '', '', 1, 2),
('登录日志', 'log:login', 'menu', 22, '/log/login', 'Login', 'log/login/index', 1, 2),
('登录日志列表', 'log:login:list', 'api', 26, '', '', '', 1, 1),
('登录日志删除', 'log:login:delete', 'api', 26, '', '', '', 1, 2),
('文件管理', 'file', 'menu', 0, '/file', 'Folder', 'Layout', 1, 3),
('文件列表', 'file:list', 'api', 29, '', '', '', 1, 1),
('文件上传', 'file:upload', 'api', 29, '', '', '', 1, 2),
('文件删除', 'file:delete', 'api', 29, '', '', '', 1, 3)
ON CONFLICT (code) DO NOTHING;

-- 给超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 插入默认系统配置
INSERT INTO system_configs (key, value, name, remark) VALUES
('site.name', '管理后台系统', '站点名称', '系统站点名称'),
('site.logo', '', '站点Logo', '系统站点Logo'),
('site.description', '基于gozero的管理后台系统', '站点描述', '系统站点描述'),
('upload.max_size', '10485760', '上传文件最大大小', '单位：字节，默认10MB'),
('upload.allow_types', 'jpg,jpeg,png,gif,doc,docx,xls,xlsx,pdf,txt', '允许上传的文件类型', '文件扩展名，多个用逗号分隔'),
('jwt.secret', 'gozero-admin-secret-key-2024', 'JWT密钥', 'JWT签名密钥'),
('jwt.expire', '86400', 'JWT过期时间', '单位：秒，默认24小时')
ON CONFLICT (key) DO NOTHING;
