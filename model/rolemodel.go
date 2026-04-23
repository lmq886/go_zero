package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// RoleModel 角色模型接口
// 定义角色相关的数据库操作方法
type RoleModel interface {
	// 插入新角色
	Insert(ctx context.Context, data *Role) (sql.Result, error)
	// 根据ID查找角色
	FindOne(ctx context.Context, id int64) (*Role, error)
	// 根据编码查找角色
	FindOneByCode(ctx context.Context, code string) (*Role, error)
	// 更新角色信息
	Update(ctx context.Context, data *Role) error
	// 删除角色（软删除）
	Delete(ctx context.Context, id int64) error
	// 获取角色列表
	FindList(ctx context.Context, page, pageSize int, name, code string) ([]*Role, int64, error)
	// 获取角色的权限列表
	FindRolePermissions(ctx context.Context, roleId int64) ([]*Permission, error)
	// 为角色分配权限
	AssignPermissions(ctx context.Context, roleId int64, permissionIds []int64) error
	// 检查角色是否被用户使用
	CheckRoleUsed(ctx context.Context, roleId int64) (bool, error)
}

// Role 角色结构体
// 对应数据库中的 roles 表
type Role struct {
	Id          int64          `db:"id"`          // 角色ID
	Name        string         `db:"name"`        // 角色名称
	Code        string         `db:"code"`        // 角色编码（唯一标识）
	Description sql.NullString `db:"description"` // 角色描述
	Status      int64          `db:"status"`      // 状态（0-禁用，1-启用）
	Sort        int64          `db:"sort"`        // 排序
	CreatedAt   int64          `db:"created_at"`  // 创建时间
	UpdatedAt   int64          `db:"updated_at"`  // 更新时间
	DeletedAt   sql.NullInt64  `db:"deleted_at"`  // 删除时间（软删除）
}

// defaultRoleModel 默认角色模型实现
type defaultRoleModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewRoleModel 创建角色模型实例
// 参数 conn: 数据库连接
// 返回值: 角色模型接口
func NewRoleModel(conn sqlx.SqlConn) RoleModel {
	return &defaultRoleModel{
		conn:  conn,
		table: "roles",
	}
}

// Insert 插入新角色
// 参数 ctx: 上下文
// 参数 data: 角色数据
// 返回值: 执行结果和错误信息
func (m *defaultRoleModel) Insert(ctx context.Context, data *Role) (sql.Result, error) {
	// 设置创建时间和更新时间
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		name, code, description, status, sort, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.Name, data.Code, data.Description, data.Status, data.Sort,
		data.CreatedAt, data.UpdatedAt,
	)
}

// FindOne 根据ID查找角色
// 参数 ctx: 上下文
// 参数 id: 角色ID
// 返回值: 角色信息和错误信息
func (m *defaultRoleModel) FindOne(ctx context.Context, id int64) (*Role, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, code, description, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Role
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

// FindOneByCode 根据编码查找角色
// 参数 ctx: 上下文
// 参数 code: 角色编码
// 返回值: 角色信息和错误信息
func (m *defaultRoleModel) FindOneByCode(ctx context.Context, code string) (*Role, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, code, description, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE code = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Role
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

// Update 更新角色信息
// 参数 ctx: 上下文
// 参数 data: 角色数据
// 返回值: 错误信息
func (m *defaultRoleModel) Update(ctx context.Context, data *Role) error {
	// 设置更新时间
	data.UpdatedAt = time.Now().Unix()

	// 构建更新SQL
	query := fmt.Sprintf(`UPDATE %s SET
		name = $1, description = $2, status = $3, sort = $4, updated_at = $5
	WHERE id = $6 AND deleted_at IS NULL`, m.table)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.Name, data.Description, data.Status, data.Sort, data.UpdatedAt,
		data.Id,
	)
	return err
}

// Delete 删除角色（软删除）
// 参数 ctx: 上下文
// 参数 id: 角色ID
// 返回值: 错误信息
func (m *defaultRoleModel) Delete(ctx context.Context, id int64) error {
	// 使用事务处理
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先删除角色权限关联
		deleteRolePermissionsQuery := `DELETE FROM role_permissions WHERE role_id = $1`
		_, err := session.ExecCtx(ctx, deleteRolePermissionsQuery, id)
		if err != nil {
			return err
		}

		// 删除角色菜单关联
		deleteRoleMenusQuery := `DELETE FROM role_menus WHERE role_id = $1`
		_, err = session.ExecCtx(ctx, deleteRoleMenusQuery, id)
		if err != nil {
			return err
		}

		// 删除用户角色关联
		deleteUserRolesQuery := `DELETE FROM user_roles WHERE role_id = $1`
		_, err = session.ExecCtx(ctx, deleteUserRolesQuery, id)
		if err != nil {
			return err
		}

		// 软删除角色
		deleteRoleQuery := fmt.Sprintf(`UPDATE %s SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`, m.table)
		_, err = session.ExecCtx(ctx, deleteRoleQuery, time.Now().Unix(), id)
		return err
	})
}

// FindList 获取角色列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 name: 角色名称（模糊查询）
// 参数 code: 角色编码（模糊查询）
// 返回值: 角色列表、总记录数和错误信息
func (m *defaultRoleModel) FindList(ctx context.Context, page, pageSize int, name, code string) ([]*Role, int64, error) {
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
		id, name, code, description, status, sort, created_at, updated_at, deleted_at
	FROM %s %s ORDER BY sort ASC, created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*Role
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// FindRolePermissions 获取角色的权限列表
// 参数 ctx: 上下文
// 参数 roleId: 角色ID
// 返回值: 权限列表和错误信息
func (m *defaultRoleModel) FindRolePermissions(ctx context.Context, roleId int64) ([]*Permission, error) {
	// 构建查询SQL
	query := `SELECT 
		p.id, p.name, p.code, p.type, p.parent_id, p.description, 
		p.path, p.method, p.status, p.sort, p.created_at, p.updated_at, p.deleted_at
	FROM permissions p
	INNER JOIN role_permissions rp ON p.id = rp.permission_id
	WHERE rp.role_id = $1 AND p.deleted_at IS NULL
	ORDER BY p.sort ASC`

	// 执行查询
	var resp []*Permission
	err := m.conn.QueryRowsCtx(ctx, &resp, query, roleId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AssignPermissions 为角色分配权限
// 参数 ctx: 上下文
// 参数 roleId: 角色ID
// 参数 permissionIds: 权限ID列表
// 返回值: 错误信息
func (m *defaultRoleModel) AssignPermissions(ctx context.Context, roleId int64, permissionIds []int64) error {
	// 使用事务处理
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先删除角色已有的权限关联
		deleteQuery := `DELETE FROM role_permissions WHERE role_id = $1`
		_, err := session.ExecCtx(ctx, deleteQuery, roleId)
		if err != nil {
			return err
		}

		// 如果没有要分配的权限，直接返回
		if len(permissionIds) == 0 {
			return nil
		}

		// 插入新的权限关联
		insertQuery := `INSERT INTO role_permissions (role_id, permission_id, created_at) VALUES ($1, $2, $3)`
		now := time.Now().Unix()

		for _, permissionId := range permissionIds {
			_, err = session.ExecCtx(ctx, insertQuery, roleId, permissionId, now)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// CheckRoleUsed 检查角色是否被用户使用
// 参数 ctx: 上下文
// 参数 roleId: 角色ID
// 返回值: 是否被使用和错误信息
func (m *defaultRoleModel) CheckRoleUsed(ctx context.Context, roleId int64) (bool, error) {
	// 构建查询SQL
	query := `SELECT COUNT(*) FROM user_roles WHERE role_id = $1`

	// 执行查询
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, roleId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
