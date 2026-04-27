package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// PermissionModel 权限模型接口
// 定义权限相关的数据库操作方法
type PermissionModel interface {
	// 插入新权限
	Insert(ctx context.Context, data *Permission) (sql.Result, error)
	// 根据ID查找权限
	FindOne(ctx context.Context, id int64) (*Permission, error)
	// 根据编码查找权限
	FindOneByCode(ctx context.Context, code string) (*Permission, error)
	// 更新权限信息
	Update(ctx context.Context, data *Permission) error
	// 删除权限（软删除）
	Delete(ctx context.Context, id int64) error
	// 获取权限列表
	FindList(ctx context.Context, page, pageSize int, name, code, permType string) ([]*Permission, int64, error)
	// 获取所有权限（用于树形结构）
	FindAll(ctx context.Context) ([]*Permission, error)
	// 检查权限是否被角色使用
	CheckPermissionUsed(ctx context.Context, permissionId int64) (bool, error)
}

// Permission 权限结构体
// 对应数据库中的 permissions 表
type Permission struct {
	Id          int64          `db:"id"`          // 权限ID
	Name        string         `db:"name"`        // 权限名称
	Code        string         `db:"code"`        // 权限编码（唯一标识）
	Type        string         `db:"type"`        // 权限类型（directory-目录，menu-菜单，api-接口）
	ParentId    int64          `db:"parent_id"`   // 父权限ID
	Description sql.NullString `db:"description"` // 权限描述
	Path        sql.NullString `db:"path"`        // 接口路径
	Method      sql.NullString `db:"method"`      // 请求方法
	Status      int64          `db:"status"`      // 状态（0-禁用，1-启用）
	Sort        int64          `db:"sort"`        // 排序
	CreatedAt   int64          `db:"created_at"`  // 创建时间
	UpdatedAt   int64          `db:"updated_at"`  // 更新时间
	DeletedAt   sql.NullInt64  `db:"deleted_at"`  // 删除时间（软删除）
}

// defaultPermissionModel 默认权限模型实现
type defaultPermissionModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewPermissionModel 创建权限模型实例
// 参数 conn: 数据库连接
// 返回值: 权限模型接口
func NewPermissionModel(conn sqlx.SqlConn) PermissionModel {
	return &defaultPermissionModel{
		conn:  conn,
		table: "permissions",
	}
}

// Insert 插入新权限
// 参数 ctx: 上下文
// 参数 data: 权限数据
// 返回值: 执行结果和错误信息
func (m *defaultPermissionModel) Insert(ctx context.Context, data *Permission) (sql.Result, error) {
	// 设置创建时间和更新时间
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.Name, data.Code, data.Type, data.ParentId, data.Description,
		data.Path, data.Method, data.Status, data.Sort, data.CreatedAt, data.UpdatedAt,
	)
}

// FindOne 根据ID查找权限
// 参数 ctx: 上下文
// 参数 id: 权限ID
// 返回值: 权限信息和错误信息
func (m *defaultPermissionModel) FindOne(ctx context.Context, id int64) (*Permission, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Permission
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

// FindOneByCode 根据编码查找权限
// 参数 ctx: 上下文
// 参数 code: 权限编码
// 返回值: 权限信息和错误信息
func (m *defaultPermissionModel) FindOneByCode(ctx context.Context, code string) (*Permission, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE code = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Permission
	err := m.conn.QueryRowCtx(ctx, &resp, query, code)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新权限信息
// 参数 ctx: 上下文
// 参数 data: 权限数据
// 返回值: 错误信息
func (m *defaultPermissionModel) Update(ctx context.Context, data *Permission) error {
	// 设置更新时间
	data.UpdatedAt = time.Now().Unix()

	// 构建更新SQL
	query := fmt.Sprintf(`UPDATE %s SET
		name = $1, description = $2, path = $3, method = $4, status = $5, sort = $6, updated_at = $7
	WHERE id = $8 AND deleted_at IS NULL`, m.table)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.Name, data.Description, data.Path, data.Method, data.Status, data.Sort, data.UpdatedAt,
		data.Id,
	)
	return err
}

// Delete 删除权限（软删除）
// 参数 ctx: 上下文
// 参数 id: 权限ID
// 返回值: 错误信息
func (m *defaultPermissionModel) Delete(ctx context.Context, id int64) error {
	// 使用事务处理
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先删除角色权限关联
		deleteQuery := `DELETE FROM role_permissions WHERE permission_id = $1`
		_, err := session.ExecCtx(ctx, deleteQuery, id)
		if err != nil {
			return err
		}

		// 软删除权限
		deletePermissionQuery := fmt.Sprintf(`UPDATE %s SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`, m.table)
		_, err = session.ExecCtx(ctx, deletePermissionQuery, time.Now().Unix(), id)
		return err
	})
}

// FindList 获取权限列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 name: 权限名称（模糊查询）
// 参数 code: 权限编码（模糊查询）
// 参数 permType: 权限类型
// 返回值: 权限列表、总记录数和错误信息
func (m *defaultPermissionModel) FindList(ctx context.Context, page, pageSize int, name, code, permType string) ([]*Permission, int64, error) {
	// 构建基础查询条件
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	// 添加名称查询条件
	if name != "" {
		whereClause += fmt.Sprintf(" AND name LIKE $%d", argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	// 添加编码查询条件
	if code != "" {
		whereClause += fmt.Sprintf(" AND code LIKE $%d", argIndex)
		args = append(args, "%"+code+"%")
		argIndex++
	}

	// 添加类型查询条件
	if permType != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, permType)
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
		id, name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at, deleted_at
	FROM %s %s ORDER BY sort ASC, created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*Permission
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// FindAll 获取所有权限（用于树形结构）
// 参数 ctx: 上下文
// 返回值: 权限列表和错误信息
func (m *defaultPermissionModel) FindAll(ctx context.Context) ([]*Permission, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, code, type, parent_id, description, path, method, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE deleted_at IS NULL ORDER BY sort ASC, created_at DESC`, m.table)

	// 执行查询
	var resp []*Permission
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckPermissionUsed 检查权限是否被角色使用
// 参数 ctx: 上下文
// 参数 permissionId: 权限ID
// 返回值: 是否被使用和错误信息
func (m *defaultPermissionModel) CheckPermissionUsed(ctx context.Context, permissionId int64) (bool, error) {
	// 构建查询SQL
	query := `SELECT COUNT(*) FROM role_permissions WHERE permission_id = $1`

	// 执行查询
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, permissionId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
