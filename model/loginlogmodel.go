package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LoginLogModel 登录日志模型接口
// 定义登录日志相关的数据库操作方法
type LoginLogModel interface {
	// 插入新登录日志
	Insert(ctx context.Context, data *LoginLog) error
	// 根据ID查找登录日志
	FindOne(ctx context.Context, id int64) (*LoginLog, error)
	// 删除登录日志
	Delete(ctx context.Context, id int64) error
	// 清空登录日志
	Clear(ctx context.Context) error
	// 获取登录日志列表
	FindList(ctx context.Context, page, pageSize int, userId int64, username string, status int, startTime, endTime int64) ([]*LoginLog, int64, error)
}

// LoginLog 登录日志结构体
// 对应数据库中的 login_logs 表
type LoginLog struct {
	Id        int64  `db:"id"`         // 日志ID
	UserId    int64  `db:"user_id"`    // 用户ID
	Username  string `db:"username"`   // 用户名
	Ip        string `db:"ip"`         // IP地址
	Location  string `db:"location"`   // 地理位置
	Browser   string `db:"browser"`    // 浏览器
	Os        string `db:"os"`         // 操作系统
	UserAgent string `db:"user_agent"` // 用户代理
	Status    int64  `db:"status"`     // 状态（1-成功，0-失败）
	Message   string `db:"message"`    // 消息
	CreatedAt int64  `db:"created_at"` // 创建时间
}

// defaultLoginLogModel 默认登录日志模型实现
type defaultLoginLogModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewLoginLogModel 创建登录日志模型实例
// 参数 conn: 数据库连接
// 返回值: 登录日志模型接口
func NewLoginLogModel(conn sqlx.SqlConn) LoginLogModel {
	return &defaultLoginLogModel{
		conn:  conn,
		table: "login_logs",
	}
}

// Insert 插入新登录日志
// 参数 ctx: 上下文
// 参数 data: 登录日志数据
// 返回值: 错误信息
func (m *defaultLoginLogModel) Insert(ctx context.Context, data *LoginLog) error {
	// 设置创建时间
	now := time.Now().Unix()
	data.CreatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		user_id, username, ip, location, browser, os, user_agent, status, message, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, m.table)

	// 执行插入操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.UserId, data.Username, data.Ip, data.Location, data.Browser,
		data.Os, data.UserAgent, data.Status, data.Message, data.CreatedAt,
	)
	return err
}

// FindOne 根据ID查找登录日志
// 参数 ctx: 上下文
// 参数 id: 日志ID
// 返回值: 登录日志信息和错误信息
func (m *defaultLoginLogModel) FindOne(ctx context.Context, id int64) (*LoginLog, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, user_id, username, ip, location, browser, os, user_agent, status, message, created_at
	FROM %s WHERE id = $1`, m.table)

	// 执行查询
	var resp LoginLog
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

// Delete 删除登录日志
// 参数 ctx: 上下文
// 参数 id: 日志ID
// 返回值: 错误信息
func (m *defaultLoginLogModel) Delete(ctx context.Context, id int64) error {
	// 构建删除SQL
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, m.table)

	// 执行删除操作
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// Clear 清空登录日志
// 参数 ctx: 上下文
// 返回值: 错误信息
func (m *defaultLoginLogModel) Clear(ctx context.Context) error {
	// 构建清空SQL
	query := fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY`, m.table)

	// 执行清空操作
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}

// FindList 获取登录日志列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 userId: 用户ID
// 参数 username: 用户名（模糊查询）
// 参数 status: 状态
// 参数 startTime: 开始时间
// 参数 endTime: 结束时间
// 返回值: 登录日志列表、总记录数和错误信息
func (m *defaultLoginLogModel) FindList(ctx context.Context, page, pageSize int, userId int64, username string, status int, startTime, endTime int64) ([]*LoginLog, int64, error) {
	// 构建基础查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// 添加用户ID查询条件
	if userId > 0 {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userId)
		argIndex++
	}

	// 添加用户名查询条件
	if username != "" {
		whereClause += fmt.Sprintf(" AND username LIKE $%d", argIndex)
		args = append(args, "%"+username+"%")
		argIndex++
	}

	// 添加状态查询条件
	if status != 0 {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// 添加开始时间查询条件
	if startTime > 0 {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, startTime)
		argIndex++
	}

	// 添加结束时间查询条件
	if endTime > 0 {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, endTime)
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
		id, user_id, username, ip, location, browser, os, user_agent, status, message, created_at
	FROM %s %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*LoginLog
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
