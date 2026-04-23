package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OperationLogModel 操作日志模型接口
// 定义操作日志相关的数据库操作方法
type OperationLogModel interface {
	// 插入新操作日志
	Insert(ctx context.Context, data *OperationLog) (sql.Result, error)
	// 根据ID查找操作日志
	FindOne(ctx context.Context, id int64) (*OperationLog, error)
	// 删除操作日志
	Delete(ctx context.Context, id int64) error
	// 清空操作日志
	Clear(ctx context.Context) error
	// 获取操作日志列表
	FindList(ctx context.Context, page, pageSize int, userId int64, username, module, operation, method string, status int, startTime, endTime int64) ([]*OperationLog, int64, error)
}

// OperationLog 操作日志结构体
// 对应数据库中的 operation_logs 表
type OperationLog struct {
	Id           int64          `db:"id"`            // 日志ID
	UserId       int64          `db:"user_id"`       // 用户ID
	Username     string         `db:"username"`      // 用户名
	Module       string         `db:"module"`        // 模块
	Operation    string         `db:"operation"`     // 操作
	Method       string         `db:"method"`        // 请求方法
	Path         string         `db:"path"`          // 请求路径
	Status       int64          `db:"status"`        // 状态（1-成功，0-失败）
	Ip           string         `db:"ip"`            // IP地址
	UserAgent    string         `db:"user_agent"`    // 用户代理
	RequestData  sql.NullString `db:"request_data"`  // 请求数据
	ResponseData sql.NullString `db:"response_data"` // 响应数据
	ErrorMsg     sql.NullString `db:"error_msg"`     // 错误信息
	Duration     int64          `db:"duration"`      // 执行时长（毫秒）
	CreatedAt    int64          `db:"created_at"`    // 创建时间
}

// defaultOperationLogModel 默认操作日志模型实现
type defaultOperationLogModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewOperationLogModel 创建操作日志模型实例
// 参数 conn: 数据库连接
// 返回值: 操作日志模型接口
func NewOperationLogModel(conn sqlx.SqlConn) OperationLogModel {
	return &defaultOperationLogModel{
		conn:  conn,
		table: "operation_logs",
	}
}

// Insert 插入新操作日志
// 参数 ctx: 上下文
// 参数 data: 操作日志数据
// 返回值: 执行结果和错误信息
func (m *defaultOperationLogModel) Insert(ctx context.Context, data *OperationLog) (sql.Result, error) {
	// 设置创建时间
	now := time.Now().Unix()
	data.CreatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		user_id, username, module, operation, method, path, status,
		ip, user_agent, request_data, response_data, error_msg, duration, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.UserId, data.Username, data.Module, data.Operation, data.Method, data.Path, data.Status,
		data.Ip, data.UserAgent, data.RequestData, data.ResponseData, data.ErrorMsg, data.Duration, data.CreatedAt,
	)
}

// FindOne 根据ID查找操作日志
// 参数 ctx: 上下文
// 参数 id: 日志ID
// 返回值: 操作日志信息和错误信息
func (m *defaultOperationLogModel) FindOne(ctx context.Context, id int64) (*OperationLog, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, user_id, username, module, operation, method, path, status,
		ip, user_agent, request_data, response_data, error_msg, duration, created_at
	FROM %s WHERE id = $1`, m.table)

	// 执行查询
	var resp OperationLog
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

// Delete 删除操作日志
// 参数 ctx: 上下文
// 参数 id: 日志ID
// 返回值: 错误信息
func (m *defaultOperationLogModel) Delete(ctx context.Context, id int64) error {
	// 构建删除SQL
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, m.table)

	// 执行删除操作
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// Clear 清空操作日志
// 参数 ctx: 上下文
// 返回值: 错误信息
func (m *defaultOperationLogModel) Clear(ctx context.Context) error {
	// 构建清空SQL
	query := fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY`, m.table)

	// 执行清空操作
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}

// FindList 获取操作日志列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 userId: 用户ID
// 参数 username: 用户名（模糊查询）
// 参数 module: 模块
// 参数 operation: 操作
// 参数 method: 请求方法
// 参数 status: 状态
// 参数 startTime: 开始时间
// 参数 endTime: 结束时间
// 返回值: 操作日志列表、总记录数和错误信息
func (m *defaultOperationLogModel) FindList(ctx context.Context, page, pageSize int, userId int64, username, module, operation, method string, status int, startTime, endTime int64) ([]*OperationLog, int64, error) {
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

	// 添加模块查询条件
	if module != "" {
		whereClause += fmt.Sprintf(" AND module = $%d", argIndex)
		args = append(args, module)
		argIndex++
	}

	// 添加操作查询条件
	if operation != "" {
		whereClause += fmt.Sprintf(" AND operation LIKE $%d", argIndex)
		args = append(args, "%"+operation+"%")
		argIndex++
	}

	// 添加请求方法查询条件
	if method != "" {
		whereClause += fmt.Sprintf(" AND method = $%d", argIndex)
		args = append(args, method)
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
		id, user_id, username, module, operation, method, path, status,
		ip, user_agent, request_data, response_data, error_msg, duration, created_at
	FROM %s %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*OperationLog
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
