# GoZero 企业级后台管理系统

## 项目概述

本项目是一个基于 GoZero 框架构建的企业级后台管理系统 API 服务，采用 RESTful 架构设计，遵循 GoZero 最佳实践。

### 技术栈

- **Go 版本**: 1.25.0
- **框架**: GoZero v1.10.1
- **数据库**: PostgreSQL
- **认证**: JWT (JSON Web Token)
- **ORM**: GoZero SQLX
- **日志**: GoZero Logx

### 功能模块

| 模块 | 功能描述 | 状态 |
|------|----------|------|
| 认证模块 | 登录、注册、登出、刷新令牌 | ✅ 已实现 |
| 健康检查 | 服务健康状态检查 | ✅ 已实现 |
| 用户管理 | 用户 CRUD、重置密码、个人资料 | 📋 待实现 |
| 角色管理 | 角色 CRUD、权限分配 | 📋 待实现 |
| 权限管理 | 权限 CRUD | 📋 待实现 |
| 菜单管理 | 菜单 CRUD、用户菜单 | 📋 待实现 |
| 日志管理 | 操作日志、登录日志 | 📋 待实现 |
| 系统配置 | 配置项 CRUD | 📋 待实现 |

### 项目结构

```
go_zero/
├── api/                          # API 层
│   ├── api.go                    # 服务入口文件（main 函数）
│   ├── admin.api                 # API 定义文件（goctl 模板）
│   ├── etc/
│   │   └── api.yaml              # 服务配置文件
│   ├── internal/
│   │   ├── config/               # 配置定义
│   │   │   └── config.go
│   │   ├── handler/              # HTTP 处理器
│   │   │   ├── healthhandler.go
│   │   │   └── loginhandler.go
│   │   ├── logic/                # 业务逻辑
│   │   │   └── loginlogic.go
│   │   ├── middleware/           # 中间件
│   │   │   ├── corsmiddleware.go      # CORS 跨域
│   │   │   ├── errormiddleware.go     # 错误处理
│   │   │   ├── jwtauthmiddleware.go   # JWT 认证
│   │   │   ├── operationlogmiddleware.go # 操作日志
│   │   │   ├── permissionmiddleware.go  # 权限验证
│   │   │   └── ratelimitmiddleware.go   # 限流
│   │   ├── svc/                  # 服务上下文（依赖注入）
│   │   │   └── servicecontext.go
│   │   └── types/                # 类型定义（请求/响应结构）
│   │       └── types.go
│   └── routes/                   # 路由注册
│       └── routes.go
├── model/                        # 数据模型层
│   ├── usermodel.go              # 用户模型
│   ├── rolemodel.go              # 角色模型
│   ├── permissionmodel.go        # 权限模型
│   ├── menumodel.go              # 菜单模型
│   ├── operationlogmodel.go      # 操作日志模型
│   ├── loginlogmodel.go          # 登录日志模型
│   ├── configmodel.go            # 系统配置模型
│   ├── errors.go                 # 错误定义
│   └── schema.sql                # 数据库 schema
├── go.mod
├── go.sum
├── README.md
└── rename.md                     # 本文档
```

---

## 运行方式

### 前置条件

1. **Go 环境**: Go 1.25.0 或更高版本
2. **PostgreSQL 数据库**: 版本 12+（可选，用于完整功能）

### 方式一：命令行运行

#### 1. 进入项目目录

```bash
cd d:\code\work\go_zero
```

#### 2. 下载依赖

```bash
go mod tidy
```

#### 3. 编译项目

```bash
go build -o admin-api.exe ./api
```

#### 4. 运行服务

**注意：必须在 `api` 目录下运行，或指定配置文件路径**

```bash
# 方式 A：进入 api 目录运行（推荐）
cd api
..\admin-api.exe

# 方式 B：使用 -f 参数指定配置文件路径
.\admin-api.exe -f api/etc/api.yaml
```

### 方式二：GoLand 运行配置

#### 正确的配置步骤：

1. 打开 GoLand，打开项目 `d:\code\work\go_zero`

2. 点击右上角 **Edit Configurations...**

3. 点击 **+** 号，选择 **Go Build**

4. 填写配置：
   - **Name**: `admin-api`
   - **Run kind**: `File`
   - **File path**: `d:\code\work\go_zero\api\api.go`
   - **Output directory**: （可选）
   - **Working directory**: `d:\code\work\go_zero\api`  ⚠️ **关键配置**
   - **Environment**: （可选，如数据库环境变量）
   - **Go tool arguments**: （可选）
   - **Program arguments**: （可选，如 `-f etc/api.yaml`，默认就是这个）

5. 点击 **Apply** → **OK**

6. 点击运行按钮 ▶️

### 方式三：使用 go run 直接运行

```bash
# 必须在 api 目录下运行
cd d:\code\work\go_zero\api
go run api.go
```

---

## 常见问题

### Q1: 在 GoLand 中运行没有反应/卡住

**可能原因 1：工作目录配置错误**

配置文件 `etc/api.yaml` 是相对于运行时工作目录的。如果工作目录不是 `api` 目录，程序会找不到配置文件。

**解决方法**：
- 确保 **Working directory** 设置为 `d:\code\work\go_zero\api`
- 或者在 **Program arguments** 中添加 `-f api/etc/api.yaml`

**可能原因 2：数据库连接超时**

配置文件默认连接 PostgreSQL 数据库（`localhost:5432`）。如果没有启动数据库，服务会在连接数据库时卡住。

**解决方法**：
- 方案 A：启动 PostgreSQL 数据库，并创建 `admin_system` 数据库
- 方案 B：暂时注释掉数据库相关代码（仅用于测试 API 框架）

**检查是否有数据库**：
```bash
# 检查 PostgreSQL 服务状态
# Windows: 服务中查看 postgresql 服务
```

### Q2: 编译错误：undefined: xxx

**原因**：代码中引用了未实现的函数或类型。

**当前状态**：项目已修复所有编译错误，可以正常编译。

### Q3: 运行时错误：cannot connect to database

**原因**：无法连接到 PostgreSQL 数据库。

**解决方法**：
1. 确保 PostgreSQL 服务已启动
2. 检查 `api/etc/api.yaml` 中的数据库配置是否正确
3. 创建数据库 `admin_system`
4. 执行 `model/schema.sql` 初始化表结构

### Q4: 如何测试 API

服务启动成功后，可以使用以下方式测试：

**健康检查**：
```bash
curl http://localhost:8080/health
```

**预期响应**：
```json
{
  "status": "ok",
  "version": "v1.0.0",
  "uptime": "service is running"
}
```

**登录接口**：
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

---

## Swagger 文档使用说明

项目已集成 Swagger API 文档，可以通过 Web 界面可视化查看和测试所有 API 接口。

### 访问 Swagger UI

服务启动后，访问以下地址：

| 地址 | 描述 |
|------|------|
| http://localhost:8080/swagger/ | Swagger UI 界面 |
| http://localhost:8080/swagger.json | OpenAPI 规范文档（JSON 格式） |
| http://localhost:8080/swagger.yaml | OpenAPI 规范文档（YAML 格式） |

### 使用 Swagger UI 测试接口

1. **打开 Swagger UI**：访问 http://localhost:8080/swagger/

2. **认证接口测试**：
   - 找到 `/api/v1/auth/login` 接口
   - 点击 `Try it out`
   - 输入请求体：
     ```json
     {
       "username": "admin",
       "password": "admin123"
     }
     ```
   - 点击 `Execute`，获取 `access_token`

3. **设置认证 Token**：
   - 点击页面右上角的 `Authorize` 按钮
   - 在 `Value` 字段输入：`Bearer <你的access_token>`
   - 点击 `Authorize`，然后点击 `Close`

4. **测试需要认证的接口**：
   - 找到需要认证的接口（如 `/api/v1/auth/logout`）
   - 点击 `Try it out`
   - 点击 `Execute`，查看响应

### Swagger UI 功能特性

- **接口可视化**：所有接口按模块分组展示
- **参数说明**：每个接口的请求参数、响应结构都有详细说明
- **在线测试**：直接在浏览器中调用接口，查看请求和响应
- **认证支持**：支持 Bearer Token 认证，可测试需要登录的接口
- **Schema 查看**：查看所有请求/响应的数据结构定义

### 文档文件位置

```
go_zero/
├── api/
│   └── internal/
│       └── handler/
│           └── swaggerhandler.go    # Swagger UI 处理器
├── docs/
│   └── swagger.yaml                  # OpenAPI 规范文档
```

---

## API 接口列表

### 公开接口（无需认证）

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/refresh` | 刷新令牌 |
| GET | `/health` | 健康检查 |

### 需认证接口（需 Bearer Token）

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/logout` | 用户登出 |

---

## 配置说明

### 配置文件位置

`api/etc/api.yaml`

### 主要配置项

```yaml
# 服务配置
Name: admin-api          # 服务名称
Host: 0.0.0.0            # 监听地址
Port: 8080               # 监听端口
Mode: debug              # 运行模式 (debug/test/release)

# 日志配置
Log:
  Level: info            # 日志级别
  Mode: console          # 输出模式
  Encoding: plain        # 编码格式

# 数据库配置
DataSource:
  Type: postgres         # 数据库类型
  Host: localhost        # 主机
  Port: 5432             # 端口
  Database: admin_system # 数据库名
  Username: postgres     # 用户名
  Password: password     # 密码

# JWT 配置
Auth:
  AccessSecret: gozero-admin-api-access-secret-key-2024
  AccessExpire: 7200     # 2小时
  RefreshSecret: gozero-admin-api-refresh-secret-key-2024
  RefreshExpire: 604800  # 7天

# 限流配置
RateLimit:
  Enabled: true
  RequestsPerSecond: 100
  Burst: 200

# CORS 配置
CORS:
  Enabled: true
  AllowOrigins:
    - "*"
```

---

## 下一步计划

1. **实现用户管理模块**：用户 CRUD、重置密码、个人资料
2. **实现角色管理模块**：角色 CRUD、权限分配
3. **实现权限管理模块**：权限 CRUD
4. **实现菜单管理模块**：菜单 CRUD、用户菜单
5. **实现日志管理模块**：操作日志、登录日志
6. **实现系统配置模块**：配置项 CRUD
7. **添加单元测试**
8. **添加 Docker 部署支持**

---

## 联系方式

如有问题，请查看代码注释或参考 GoZero 官方文档。

- GoZero 官方文档: https://go-zero.dev/
- GoZero GitHub: https://github.com/zeromicro/go-zero
