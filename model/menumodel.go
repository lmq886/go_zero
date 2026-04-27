package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MenuModel 菜单模型接口
// 定义菜单相关的数据库操作方法
type MenuModel interface {
	// 插入新菜单
	Insert(ctx context.Context, data *Menu) (sql.Result, error)
	// 根据ID查找菜单
	FindOne(ctx context.Context, id int64) (*Menu, error)
	// 更新菜单信息
	Update(ctx context.Context, data *Menu) error
	// 删除菜单（软删除）
	Delete(ctx context.Context, id int64) error
	// 获取菜单列表
	FindList(ctx context.Context, page, pageSize int, name, menuType string) ([]*Menu, int64, error)
	// 获取所有菜单（用于树形结构）
	FindAll(ctx context.Context) ([]*Menu, error)
	// 获取用户的菜单列表
	FindUserMenus(ctx context.Context, userId int64) ([]*Menu, error)
	// 检查菜单是否有子菜单
	CheckHasChildren(ctx context.Context, parentId int64) (bool, error)
	// 检查菜单是否被角色使用
	CheckMenuUsed(ctx context.Context, menuId int64) (bool, error)
}

// Menu 菜单结构体
// 对应数据库中的 menus 表
type Menu struct {
	Id          int64          `db:"id"`           // 菜单ID
	Name        string         `db:"name"`         // 菜单名称
	Type        string         `db:"type"`         // 菜单类型（directory-目录，menu-菜单，button-按钮）
	ParentId    int64          `db:"parent_id"`    // 父菜单ID
	Path        sql.NullString `db:"path"`         // 路由路径
	Component   sql.NullString `db:"component"`    // 组件路径
	Icon        sql.NullString `db:"icon"`         // 图标
	Title       string         `db:"title"`        // 标题（用于国际化）
	Redirect    sql.NullString `db:"redirect"`     // 重定向路径
	Hidden      bool           `db:"hidden"`       // 是否隐藏
	AlwaysShow  bool           `db:"always_show"`  // 是否总是显示
	Permission  sql.NullString `db:"permission"`   // 权限编码
	Status      int64          `db:"status"`       // 状态（0-禁用，1-启用）
	Sort        int64          `db:"sort"`         // 排序
	CreatedAt   int64          `db:"created_at"`   // 创建时间
	UpdatedAt   int64          `db:"updated_at"`   // 更新时间
	DeletedAt   sql.NullInt64  `db:"deleted_at"`   // 删除时间（软删除）
}

// defaultMenuModel 默认菜单模型实现
type defaultMenuModel struct {
	conn  sqlx.SqlConn // 数据库连接
	table string        // 表名
}

// NewMenuModel 创建菜单模型实例
// 参数 conn: 数据库连接
// 返回值: 菜单模型接口
func NewMenuModel(conn sqlx.SqlConn) MenuModel {
	return &defaultMenuModel{
		conn:  conn,
		table: "menus",
	}
}

// Insert 插入新菜单
// 参数 ctx: 上下文
// 参数 data: 菜单数据
// 返回值: 执行结果和错误信息
func (m *defaultMenuModel) Insert(ctx context.Context, data *Menu) (sql.Result, error) {
	// 设置创建时间和更新时间
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	// 构建插入SQL
	query := fmt.Sprintf(`INSERT INTO %s (
		name, type, parent_id, path, component, icon, title, redirect,
		hidden, always_show, permission, status, sort, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, m.table)

	// 执行插入操作
	return m.conn.ExecCtx(ctx, query,
		data.Name, data.Type, data.ParentId, data.Path, data.Component,
		data.Icon, data.Title, data.Redirect, data.Hidden, data.AlwaysShow,
		data.Permission, data.Status, data.Sort, data.CreatedAt, data.UpdatedAt,
	)
}

// FindOne 根据ID查找菜单
// 参数 ctx: 上下文
// 参数 id: 菜单ID
// 返回值: 菜单信息和错误信息
func (m *defaultMenuModel) FindOne(ctx context.Context, id int64) (*Menu, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, type, parent_id, path, component, icon, title, redirect,
		hidden, always_show, permission, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var resp Menu
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

// Update 更新菜单信息
// 参数 ctx: 上下文
// 参数 data: 菜单数据
// 返回值: 错误信息
func (m *defaultMenuModel) Update(ctx context.Context, data *Menu) error {
	// 设置更新时间
	data.UpdatedAt = time.Now().Unix()

	// 构建更新SQL
	query := fmt.Sprintf(`UPDATE %s SET
		name = $1, parent_id = $2, path = $3, component = $4, icon = $5,
		title = $6, redirect = $7, hidden = $8, always_show = $9, permission = $10,
		status = $11, sort = $12, updated_at = $13
	WHERE id = $14 AND deleted_at IS NULL`, m.table)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		data.Name, data.ParentId, data.Path, data.Component, data.Icon,
		data.Title, data.Redirect, data.Hidden, data.AlwaysShow, data.Permission,
		data.Status, data.Sort, data.UpdatedAt, data.Id,
	)
	return err
}

// Delete 删除菜单（软删除）
// 参数 ctx: 上下文
// 参数 id: 菜单ID
// 返回值: 错误信息
func (m *defaultMenuModel) Delete(ctx context.Context, id int64) error {
	// 使用事务处理
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先删除角色菜单关联
		deleteQuery := `DELETE FROM role_menus WHERE menu_id = $1`
		_, err := session.ExecCtx(ctx, deleteQuery, id)
		if err != nil {
			return err
		}

		// 软删除菜单
		deleteMenuQuery := fmt.Sprintf(`UPDATE %s SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`, m.table)
		_, err = session.ExecCtx(ctx, deleteMenuQuery, time.Now().Unix(), id)
		return err
	})
}

// FindList 获取菜单列表
// 参数 ctx: 上下文
// 参数 page: 页码
// 参数 pageSize: 每页数量
// 参数 name: 菜单名称（模糊查询）
// 参数 menuType: 菜单类型
// 返回值: 菜单列表、总记录数和错误信息
func (m *defaultMenuModel) FindList(ctx context.Context, page, pageSize int, name, menuType string) ([]*Menu, int64, error) {
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

	// 添加类型查询条件
	if menuType != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, menuType)
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
		id, name, type, parent_id, path, component, icon, title, redirect,
		hidden, always_show, permission, status, sort, created_at, updated_at, deleted_at
	FROM %s %s ORDER BY sort ASC, created_at DESC LIMIT $%d OFFSET $%d`, 
		m.table, whereClause, argIndex, argIndex+1)

	// 添加分页参数
	args = append(args, pageSize, offset)

	// 执行查询
	var resp []*Menu
	err = m.conn.QueryRowsCtx(ctx, &resp, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// FindAll 获取所有菜单（用于树形结构）
// 参数 ctx: 上下文
// 返回值: 菜单列表和错误信息
func (m *defaultMenuModel) FindAll(ctx context.Context) ([]*Menu, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT 
		id, name, type, parent_id, path, component, icon, title, redirect,
		hidden, always_show, permission, status, sort, created_at, updated_at, deleted_at
	FROM %s WHERE deleted_at IS NULL AND status = 1 ORDER BY sort ASC, created_at DESC`, m.table)

	// 执行查询
	var resp []*Menu
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FindUserMenus 获取用户的菜单列表
// 参数 ctx: 上下文
// 参数 userId: 用户ID
// 返回值: 菜单列表和错误信息
func (m *defaultMenuModel) FindUserMenus(ctx context.Context, userId int64) ([]*Menu, error) {
	// 构建查询SQL（通过用户->角色->菜单获取菜单列表）
	query := `SELECT DISTINCT
		m.id, m.name, m.type, m.parent_id, m.path, m.component, m.icon, m.title, m.redirect,
		m.hidden, m.always_show, m.permission, m.status, m.sort, m.created_at, m.updated_at, m.deleted_at
	FROM menus m
	INNER JOIN role_menus rm ON m.id = rm.menu_id
	INNER JOIN user_roles ur ON rm.role_id = ur.role_id
	WHERE ur.user_id = $1 AND m.deleted_at IS NULL AND m.status = 1
	ORDER BY m.sort ASC, m.created_at DESC`

	// 执行查询
	var resp []*Menu
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckHasChildren 检查菜单是否有子菜单
// 参数 ctx: 上下文
// 参数 parentId: 父菜单ID
// 返回值: 是否有子菜单和错误信息
func (m *defaultMenuModel) CheckHasChildren(ctx context.Context, parentId int64) (bool, error) {
	// 构建查询SQL
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE parent_id = $1 AND deleted_at IS NULL`, m.table)

	// 执行查询
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, parentId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckMenuUsed 检查菜单是否被角色使用
// 参数 ctx: 上下文
// 参数 menuId: 菜单ID
// 返回值: 是否被使用和错误信息
func (m *defaultMenuModel) CheckMenuUsed(ctx context.Context, menuId int64) (bool, error) {
	// 构建查询SQL
	query := `SELECT COUNT(*) FROM role_menus WHERE menu_id = $1`

	// 执行查询
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, menuId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
