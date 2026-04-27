package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserModel 用户模型接口
// 定义用户相关的数据库操作方法
type UserModel interface {
	// 插入新用户
	Insert(ctx context.Context, data *User) (sql.Result, error)
	// 根据ID查找用户
	FindOne(ctx context.Context, id int64) (*User, error)
	// 根据用户名查找用户
	FindOneByUsername(ctx context.Context, username string) (*User, error)
	// 根据邮箱查找用户
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	// 更新用户信息
	Update(ctx context.Context, data *User) error
	// 删除用户（软删除）
	Delete(ctx context.Context, id int64) error
	// 获取用户列表
	FindList(ctx context.Context, page, pageSize int, username, email, phone string, status int) ([]*User, int64, error)
	// 获取用户的角色列表
	FindUserRoles(ctx context.Context, userId int64) ([]*Role, error)
	// 获取用户的权限列表
	FindUserPermissions(ctx context.Context, userId int64) ([]*Permission, error)
	// 为用户分配角色
	AssignRoles(ctx context.Context, userId int64, roleIds []int64) error
}

// User 用户结构体
// 对应数据库中的 users 表
type User struct {
	Id           int64          `db:"id"`           // 用户ID
	Username     string         `db:"username"`     // 用户名
	Password     string         `db:"password"`     // 密码（加密存储）
	Email        sql.NullString `db:"email"`        // 邮箱
	Phone        sql.NullString `db:"phone"`        // 手机号
	Nickname     sql.NullString `db:"nickname"`     // 昵称
	Avatar       sql.NullString `db:"avatar"`       // 头像URL
	Gender       sql.NullInt64  `db:"gender"`       // 性别（0-未知，1-男，2-女）
	Birthday     sql.NullString `db:"birthday"`     // 生日
	Address      sql.NullString `db:"address"`      // 地址
	Introduction sql.NullString `db:"introduction"` // 简介
	Status       int64          `db:"status"`       // 状态（0-禁用，1-启用）
	LastLoginAt  sql.NullInt64  `db:"last_login_at"` // 最后登录时间
	LastLoginIp  sql.NullString `db:"last_login_ip"` // 最后登录IP
	CreatedAt    int64          `db:"created_at"`   // 创建时间
	UpdatedAt    int64          `db:"updated_at"`   // 更新时间
	DeletedAt    sql.NullInt64  `db:"deleted_at"`   // 删除时间（软删除）
}

// defaultUserModel 默认用户模型实现
type defaultUserModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewUserModel 创建用户模型实例
// 参数 conn: 数据库连接
// 返回值: 用户模型接口
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &defaultUserModel{
		conn:  conn,
		table: "users",
	}
}

// Insert 插入新用户
// 参数 ctx: 上下文
// 参数 data: 用户数据
// 返回值: 执行结果和错误信息
func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	// 设置创建时间和更新时间
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		username, password, email, phone, nickname, avatar, 
		gender, birthday, address, introduction, status, 
		last_login_at, last_login_ip, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.Username, data.Password, data.Email, data.Phone, data.Nickname, data.Avatar,
		data.Gender, data.Birthday, data.Address, data.Introduction, data.Status,
		data.LastLoginAt, data.LastLoginIp, data.CreatedAt, data.UpdatedAt,
	)
}

// FindOne 根据ID查找用户
// 参数 ctx: 上下文
// 参数 id: 用户ID
// 返回值: 用户信息和错误信息
func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, username, password, email, phone, nickname, avatar,
		gender, birthday, address, introduction, status,
		last_login_at, last_login_ip, created_at, updated_at, deleted_at
	FROM %s WHERE id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneByUsername 根据用户名查找用户
// 参数 ctx: 上下文
// 参数 username: 用户名
// 返回值: 用户信息和错误信息
func (m *defaultUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, username, password, email, phone, nickname, avatar,
		gender, birthday, address, introduction, status,
		last_login_at, last_login_ip, created_at, updated_at, deleted_at
	FROM %s WHERE username = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneByEmail 根据邮箱查找用户
// 参数 ctx: 上下文
// 参数 email: 邮箱
// 返回值: 用户信息和错误信息
func (m *defaultUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, username, password, email, phone, nickname, avatar,
		gender, birthday, address, introduction, status,
		last_login_at, last_login_ip, created_at, updated_at, deleted_at
	FROM %s WHERE email = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新用户信息
// 参数 ctx: 上下文
// 参数 data: 用户数据
// 返回值: 错误信息
func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	// 设置更新时间
	data.UpdatedAt = time.Now().Unix()

	// 构建更新SQL
	query := fmt.Sprintf(`UPDATE %s SET
		email = $1, phone = $2, nickname = $3, avatar = $4,
		gender = $5, birthday = $6, address = $7, introduction = $8,
		status = $9, last_login_at = $10, last_login_ip = $11, updated_at = $12
	WHERE id = $13 AND deleted_at IS NULL`, m.table)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.Email, data.Phone, data.Nickname, data.Avatar,
		data.Gender, data.Birthday, data.Address, data.Introduction,
		data.Status, data.LastLoginAt, data.LastLoginIp, data.UpdatedAt,
		data.Id,
	)
	return err
}

// Delete 删除用户（软删除）
// 参数 ctx: 上下文
// 参数 id: 用户ID
// 返回值: 错误信息
func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	// 构建软删除SQL
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`, m.table)

	// 执行软删除操作
	_, err := m.conn.ExecCtx(ctx, query, time.Now().Unix(), id)
	return err
}

// FindList 获取用户列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 username: 用户名（模糊查询）
// 参数 email: 邮箱（模糊查询）
// 参数 phone: 手机号（模糊查询）
// 参数 status: 状态
// 返回值: 用户列表、总记录数和错误信息
func (m *defaultUserModel) FindList(ctx context.Context, page, pageSize int, username, email, phone string, status int) ([]*User, int64, error) {
	// 构建基础查询条件
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	// 添加用户名查询条件
	if username != "" {
		whereClause += fmt.Sprintf(" AND username LIKE $%d", argIndex)
		args = append(args, "%"+username+"%")
		argIndex++
	}

	// 添加邮箱查询条件
	if email != "" {
		whereClause += fmt.Sprintf(" AND email LIKE $%d", argIndex)
		args = append(args, "%"+email+"%")
		argIndex++
	}

	// 添加手机号查询条件
	if phone != "" {
		whereClause += fmt.Sprintf(" AND phone LIKE $%d", argIndex)
		args = append(args, "%"+phone+"%")
		argIndex++
	}

	// 添加状态查询条件
	if status != 0 {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// 查询总记录数
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建列表查询SQL
	listQuery := fmt.Sprintf(`SELECT 
		id, username, password, email, phone, nickname, avatar,
		gender, birthday, address, introduction, status,
		last_login_at, last_login_ip, created_at, updated_at, deleted_at
	FROM %s %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*User
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// FindUserRoles 获取用户的角色列表
// 参数 ctx: 上下文
// 参数 userId: 用户ID
// 返回值: 角色列表和错误信息
func (m *defaultUserModel) FindUserRoles(ctx context.Context, userId int64) ([]*Role, error) {
	// 构建查询SQL
	query := `SELECT 
		r.id, r.name, r.code, r.description, r.status, r.sort, r.created_at, r.updated_at, r.deleted_at
	FROM roles r
	INNER JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1 AND r.deleted_at IS NULL
	ORDER BY r.sort ASC`

	// 执行查询
	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FindUserPermissions 获取用户的权限列表
// 参数 ctx: 上下文
// 参数 userId: 用户ID
// 返回值: 权限列表和错误信息
func (m *defaultUserModel) FindUserPermissions(ctx context.Context, userId int64) ([]*Permission, error) {
	// 构建查询SQL（通过用户->角色->权限获取权限列表）
	query := `SELECT DISTINCT
		p.id, p.name, p.code, p.type, p.parent_id, p.description, 
		p.path, p.method, p.status, p.sort, p.created_at, p.updated_at, p.deleted_at
	FROM permissions p
	INNER JOIN role_permissions rp ON p.id = rp.permission_id
	INNER JOIN user_roles ur ON rp.role_id = ur.role_id
	WHERE ur.user_id = $1 AND p.deleted_at IS NULL
	ORDER BY p.sort ASC`

	// 执行查询
	var resp []*Permission
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AssignRoles 为用户分配角色
// 参数 ctx: 上下文
// 参数 userId: 用户ID
// 参数 roleIds: 角色ID列表
// 返回值: 错误信息
func (m *defaultUserModel) AssignRoles(ctx context.Context, userId int64, roleIds []int64) error {
	// 使用事务处理
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先删除用户已有的角色关联
		deleteQuery := `DELETE FROM user_roles WHERE user_id = $1`
		_, err := session.ExecCtx(ctx, deleteQuery, userId)
		if err != nil {
			return err
		}

		// 如果没有要分配的角色，直接返回
		if len(roleIds) == 0 {
			return nil
		}

		// 插入新的角色关联
		insertQuery := `INSERT INTO user_roles (user_id, role_id, created_at) VALUES ($1, $2, $3)`
		now := time.Now().Unix()

		for _, roleId := range roleIds {
			_, err = session.ExecCtx(ctx, insertQuery, userId, roleId, now)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
