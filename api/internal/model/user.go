/*
 * @Author: 羡鱼
 * @Date: 2026-04-23 09:37:31
 * @FilePath: \go_zero\api\internal\model\user.go
 * @Description: 用户数据模型
 */
package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

// 用户表字段名称
var userFieldNames = builder.RawFieldNames(&User{})
var userRows = strings.Join(userFieldNames, ",")

var _ UserModel = (*customUserModel)(nil)

type (
	// User 用户结构体，对应数据库users表
	User struct {
		Id        int64          `db:"id"`         // 用户ID
		Username  string         `db:"username"`  // 用户名
		Password  string         `db:"password"`  // 密码（加密后）
		Nickname  sql.NullString `db:"nickname"`  // 昵称
		Avatar    sql.NullString `db:"avatar"`    // 头像URL
		Email     sql.NullString `db:"email"`     // 邮箱
		Phone     sql.NullString `db:"phone"`     // 手机号
		Status    int64          `db:"status"`    // 状态（1:正常，0:禁用）
		CreatedAt time.Time      `db:"created_at"`// 创建时间
		UpdatedAt time.Time      `db:"updated_at"`// 更新时间
		DeletedAt sql.NullTime   `db:"deleted_at"`// 删除时间（软删除）
	}

	// UserModel 用户数据模型接口
	UserModel interface {
		userModel
		FindOneByUsername(ctx context.Context, username string) (*User, error)
		FindPage(ctx context.Context, page, pageSize int64, username, nickname string, status int64) ([]*User, int64, error)
	}

	// customUserModel 自定义用户模型实现
	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel 创建用户数据模型实例
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

// FindOneByUsername 根据用户名查询用户
// 参数: ctx - 上下文
// 参数: username - 用户名
// 返回: *User - 用户信息
// 返回: error - 错误信息
func (m *customUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	query := fmt.Sprintf("select %s from %s where username = $1 and deleted_at is null limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindPage 分页查询用户列表
// 参数: ctx - 上下文
// 参数: page - 页码
// 参数: pageSize - 每页数量
// 参数: username - 用户名模糊查询
// 参数: nickname - 昵称模糊查询
// 参数: status - 状态查询
// 返回: []*User - 用户列表
// 返回: int64 - 总记录数
// 返回: error - 错误信息
func (m *customUserModel) FindPage(ctx context.Context, page, pageSize int64, username, nickname string, status int64) ([]*User, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// 构建查询条件
	if username != "" {
		conditions = append(conditions, fmt.Sprintf("username like $%d", argIndex))
		args = append(args, "%"+username+"%")
		argIndex++
	}

	if nickname != "" {
		conditions = append(conditions, fmt.Sprintf("nickname like $%d", argIndex))
		args = append(args, "%"+nickname+"%")
		argIndex++
	}

	if status > 0 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	// 软删除过滤
	conditions = append(conditions, "deleted_at is null")
	whereClause := strings.Join(conditions, " and ")

	// 1. 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 2. 查询数据
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where %s order by id desc limit $%d offset $%d",
		userRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*User
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, total, nil
	case sqlc.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

// userModel 默认用户模型接口
type userModel interface {
	Insert(ctx context.Context, data *User) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, id int64) error
}

// defaultUserModel 默认用户模型实现
type defaultUserModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string       // 表名
}

// newUserModel 创建默认用户模型实例
func newUserModel(conn sqlx.SqlConn) *defaultUserModel {
	return &defaultUserModel{
		conn:  conn,
		table: "users",
	}
}

// Insert 插入用户数据
// 参数: ctx - 上下文
// 参数: data - 用户数据
// 返回: sql.Result - 插入结果
// 返回: error - 错误信息
func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (username, password, nickname, avatar, email, phone, status, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Username, data.Password, data.Nickname, data.Avatar, data.Email, data.Phone, data.Status, time.Now(), time.Now())
	return ret, err
}

// FindOne 根据ID查询用户
// 参数: ctx - 上下文
// 参数: id - 用户ID
// 返回: *User - 用户信息
// 返回: error - 错误信息
func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新用户数据
// 参数: ctx - 上下文
// 参数: data - 用户数据
// 返回: error - 错误信息
func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := fmt.Sprintf("update %s set username = $1, password = $2, nickname = $3, avatar = $4, email = $5, phone = $6, status = $7, updated_at = $8 where id = $9", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Username, data.Password, data.Nickname, data.Avatar, data.Email, data.Phone, data.Status, time.Now(), data.Id)
	return err
}

// Delete 软删除用户
// 参数: ctx - 上下文
// 参数: id - 用户ID
// 返回: error - 错误信息
func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

// 缓存键前缀
var (
	cacheUserIdPrefix      = "cache:user:id:"
	cacheUserUsernamePrefix = "cache:user:username:"
)

// cacheUserIdKey 生成用户ID缓存键
func cacheUserIdKey(id int64) string {
	return fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
}

// cacheUserUsernameKey 生成用户名缓存键
func cacheUserUsernameKey(username string) string {
	return fmt.Sprintf("%s%s", cacheUserUsernamePrefix, username)
}

// ErrNotFound 数据未找到错误
var ErrNotFound = sqlc.ErrNotFound
