package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ConfigModel 系统配置模型接口
// 定义系统配置相关的数据库操作方法
type ConfigModel interface {
	// 插入新配置
	Insert(ctx context.Context, data *Config) (sql.Result, error)
	// 根据ID查找配置
	FindOne(ctx context.Context, id int64) (*Config, error)
	// 根据键查找配置
	FindOneByKey(ctx context.Context, key string) (*Config, error)
	// 更新配置信息
	Update(ctx context.Context, data *Config) error
	// 删除配置（软删除）
	Delete(ctx context.Context, key string) error
	// 获取配置列表
	FindList(ctx context.Context, page, pageSize int, key, name, group string) ([]*Config, int64, error)
	// 根据分组获取配置
	FindByGroup(ctx context.Context, group string) ([]*Config, error)
	// 批量获取配置值
	FindValuesByKeys(ctx context.Context, keys []string) (map[string]string, error)
}

// Config 系统配置结构体
// 对应数据库中的 configs 表
type Config struct {
	Id          int64          `db:"id"`           // 配置ID
	Key         string         `db:"key"`          // 配置键（唯一标识）
	Name        string         `db:"name"`         // 配置名称
	Value       string         `db:"value"`        // 配置值
	Type        string         `db:"type"`         // 配置类型（string, number, boolean, json）
	GroupName   string         `db:"group_name"`   // 配置分组
	Description sql.NullString `db:"description"`  // 配置描述
	Status      int64          `db:"status"`       // 状态（0-禁用，1-启用）
	CreatedAt   int64          `db:"created_at"`   // 创建时间
	UpdatedAt   int64          `db:"updated_at"`   // 更新时间
	DeletedAt   sql.NullInt64  `db:"deleted_at"`   // 删除时间（软删除）
}

// defaultConfigModel 默认系统配置模型实现
type defaultConfigModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewConfigModel 创建系统配置模型实例
// 参数 conn: 数据库连接
// 返回值: 系统配置模型接口
func NewConfigModel(conn sqlx.SqlConn) ConfigModel {
	return &defaultConfigModel{
		conn:  conn,
		table: "configs",
	}
}

// Insert 插入新配置
// 参数 ctx: 上下文
// 参数 data: 配置数据
// 返回值: 执行结果和错误信息
func (m *defaultConfigModel) Insert(ctx context.Context, data *Config) (sql.Result, error) {
	// 设置创建时间和更新时间
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		key, name, value, type, group_name, description, status, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.Key, data.Name, data.Value, data.Type, data.GroupName,
		data.Description, data.Status, data.CreatedAt, data.UpdatedAt,
	)
}

// FindOne 根据ID查找配置
// 参数 ctx: 上下文
// 参数 id: 配置ID
// 返回值: 配置信息和错误信息
func (m *defaultConfigModel) FindOne(ctx context.Context, id int64) (*Config, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, key, name, value, type, group_name, description, status, created_at, updated_at, deleted_at
	FROM %s WHERE id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Config
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

// FindOneByKey 根据键查找配置
// 参数 ctx: 上下文
// 参数 key: 配置键
// 返回值: 配置信息和错误信息
func (m *defaultConfigModel) FindOneByKey(ctx context.Context, key string) (*Config, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, key, name, value, type, group_name, description, status, created_at, updated_at, deleted_at
	FROM %s WHERE key = $1 AND deleted_at IS NULL AND status = 1`, m.table)

	// 执行查询
	var resp Config
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新配置信息
// 参数 ctx: 上下文
// 参数 data: 配置数据
// 返回值: 错误信息
func (m *defaultConfigModel) Update(ctx context.Context, data *Config) error {
	// 设置更新时间
	data.UpdatedAt = time.Now().Unix()

	// 构建更新SQL
	query := fmt.Sprintf(`UPDATE %s SET
		name = $1, value = $2, description = $3, status = $4, updated_at = $5
	WHERE key = $6 AND deleted_at IS NULL`, m.table)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.Name, data.Value, data.Description, data.Status, data.UpdatedAt,
		data.Key,
	)
	return err
}

// Delete 删除配置（软删除）
// 参数 ctx: 上下文
// 参数 key: 配置键
// 返回值: 错误信息
func (m *defaultConfigModel) Delete(ctx context.Context, key string) error {
	// 构建软删除SQL
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = $1 WHERE key = $2 AND deleted_at IS NULL`, m.table)

	// 执行软删除操作
	_, err := m.conn.ExecCtx(ctx, query, time.Now().Unix(), key)
	return err
}

// FindList 获取配置列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 key: 配置键（模糊查询）
// 参数 name: 配置名称（模糊查询）
// 参数 group: 配置分组
// 返回值: 配置列表、总记录数和错误信息
func (m *defaultConfigModel) FindList(ctx context.Context, page, pageSize int, key, name, group string) ([]*Config, int64, error) {
	// 构建基础查询条件
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	// 添加键查询条件
	if key != "" {
		whereClause += fmt.Sprintf(" AND key LIKE $%d", argIndex)
		args = append(args, "%"+key+"%")
		argIndex++
	}

	// 添加名称查询条件
	if name != "" {
		whereClause += fmt.Sprintf(" AND name LIKE $%d", argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	// 添加分组查询条件
	if group != "" {
		whereClause += fmt.Sprintf(" AND group_name = $%d", argIndex)
		args = append(args, group)
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
		id, key, name, value, type, group_name, description, status, created_at, updated_at, deleted_at
	FROM %s %s ORDER BY group_name ASC, created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*Config
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// FindByGroup 根据分组获取配置
// 参数 ctx: 上下文
// 参数 group: 配置分组
// 返回值: 配置列表和错误信息
func (m *defaultConfigModel) FindByGroup(ctx context.Context, group string) ([]*Config, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, key, name, value, type, group_name, description, status, created_at, updated_at, deleted_at
	FROM %s WHERE group_name = $1 AND deleted_at IS NULL AND status = 1 ORDER BY created_at ASC`, m.table)

	// 执行查询
	var resp []*Config
	err := m.conn.QueryRowsCtx(ctx, &resp, query, group)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FindValuesByKeys 批量获取配置值
// 参数 ctx: 上下文
// 参数 keys: 配置键列表
// 返回值: 配置键值对和错误信息
func (m *defaultConfigModel) FindValuesByKeys(ctx context.Context, keys []string) (map[string]string, error) {
	// 如果没有键，直接返回空映射
	if len(keys) == 0 {
		return map[string]string{}, nil
	}

	// 构建查询SQL（使用 IN 子句）
	// 注意：这里为了简单，我们使用字符串拼接，实际项目中应该使用更安全的方式
	// 或者查询所有配置后在内存中过滤
	query := fmt.Sprintf(`SELECT key, value FROM %s WHERE deleted_at IS NULL AND status = 1`, m.table)

	// 执行查询
	type configKeyValue struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}

	var resp []*configKeyValue
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	// 构建结果映射
	result := make(map[string]string)
	// 先创建一个键的集合用于快速查找
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for _, item := range resp {
		if keySet[item.Key] {
			result[item.Key] = item.Value
		}
	}

	return result, nil
}
