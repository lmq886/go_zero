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
)

var permissionFieldNames = builder.RawFieldNames(&Permission{})
var permissionRows = strings.Join(permissionFieldNames, ",")

var _ PermissionModel = (*customPermissionModel)(nil)

type (
	// Permission is an object of map Permission
	Permission struct {
		Id        int64          `db:"id"`
		Name      string         `db:"name"`
		Code      string         `db:"code"`
		Type      string         `db:"type"`
		ParentId  int64          `db:"parent_id"`
		Path      sql.NullString `db:"path"`
		Icon      sql.NullString `db:"icon"`
		Component sql.NullString `db:"component"`
		Status    int64          `db:"status"`
		Sort      int64          `db:"sort"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
		DeletedAt sql.NullTime   `db:"deleted_at"`
	}

	// PermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionModel.
	PermissionModel interface {
		permissionModel
		FindByUserId(ctx context.Context, userId int64) ([]*Permission, error)
		FindByRoleId(ctx context.Context, roleId int64) ([]*Permission, error)
		FindPage(ctx context.Context, page, pageSize int64, name, typ string, status int64) ([]*Permission, int64, error)
	}

	customPermissionModel struct {
		*defaultPermissionModel
	}
)

// NewPermissionModel returns a model for the database table.
func NewPermissionModel(conn sqlx.SqlConn) PermissionModel {
	return &customPermissionModel{
		defaultPermissionModel: newPermissionModel(conn),
	}
}

func (m *customPermissionModel) FindByUserId(ctx context.Context, userId int64) ([]*Permission, error) {
	query := fmt.Sprintf(`select distinct p.%s from %s p 
		inner join role_permissions rp on p.id = rp.permission_id 
		inner join user_roles ur on rp.role_id = ur.role_id 
		where ur.user_id = $1 and p.deleted_at is null`, permissionRows, m.table)

	var resp []*Permission
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customPermissionModel) FindByRoleId(ctx context.Context, roleId int64) ([]*Permission, error) {
	query := fmt.Sprintf(`select p.%s from %s p 
		inner join role_permissions rp on p.id = rp.permission_id 
		where rp.role_id = $1 and p.deleted_at is null`, permissionRows, m.table)

	var resp []*Permission
	err := m.conn.QueryRowsCtx(ctx, &resp, query, roleId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customPermissionModel) FindPage(ctx context.Context, page, pageSize int64, name, typ string, status int64) ([]*Permission, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name like $%d", argIndex))
		args = append(args, "%"+name+"%")
		argIndex++
	}

	if typ != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, typ)
		argIndex++
	}

	if status > 0 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	conditions = append(conditions, "deleted_at is null")
	whereClause := strings.Join(conditions, " and ")

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where %s order by sort asc, id desc limit $%d offset $%d",
		permissionRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*Permission
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

// permissionModel is a default model implementation for table permissions
type permissionModel interface {
	Insert(ctx context.Context, data *Permission) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*Permission, error)
	Update(ctx context.Context, data *Permission) error
	Delete(ctx context.Context, id int64) error
}

type defaultPermissionModel struct {
	conn  sqlx.SqlConn
	table string
}

func newPermissionModel(conn sqlx.SqlConn) *defaultPermissionModel {
	return &defaultPermissionModel{
		conn:  conn,
		table: "permissions",
	}
}

func (m *defaultPermissionModel) Insert(ctx context.Context, data *Permission) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (name, code, type, parent_id, path, icon, component, status, sort, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Code, data.Type, data.ParentId, data.Path, data.Icon, data.Component, data.Status, data.Sort, time.Now(), time.Now())
	return ret, err
}

func (m *defaultPermissionModel) FindOne(ctx context.Context, id int64) (*Permission, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", permissionRows, m.table)
	var resp Permission
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

func (m *defaultPermissionModel) Update(ctx context.Context, data *Permission) error {
	query := fmt.Sprintf("update %s set name = $1, code = $2, type = $3, parent_id = $4, path = $5, icon = $6, component = $7, status = $8, sort = $9, updated_at = $10 where id = $11", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Code, data.Type, data.ParentId, data.Path, data.Icon, data.Component, data.Status, data.Sort, time.Now(), data.Id)
	return err
}

func (m *defaultPermissionModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
