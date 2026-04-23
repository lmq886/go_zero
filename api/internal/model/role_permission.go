package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionModel = (*defaultRolePermissionModel)(nil)

type (
	// RolePermission is an object of map RolePermission
	RolePermission struct {
		Id           int64     `db:"id"`
		RoleId       int64     `db:"role_id"`
		PermissionId int64     `db:"permission_id"`
		CreatedAt    time.Time `db:"created_at"`
	}

	// RolePermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionModel.
	RolePermissionModel interface {
		Insert(ctx context.Context, data *RolePermission) error
		DeleteByRoleId(ctx context.Context, roleId int64) error
		DeleteByPermissionId(ctx context.Context, permissionId int64) error
		FindPermissionIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
		BatchInsert(ctx context.Context, roleId int64, permissionIds []int64) error
	}

	defaultRolePermissionModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewRolePermissionModel returns a model for the database table.
func NewRolePermissionModel(conn sqlx.SqlConn) RolePermissionModel {
	return &defaultRolePermissionModel{
		conn:  conn,
		table: "role_permissions",
	}
}

func (m *defaultRolePermissionModel) Insert(ctx context.Context, data *RolePermission) error {
	query := fmt.Sprintf("insert into %s (role_id, permission_id, created_at) values ($1, $2, $3)", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.RoleId, data.PermissionId, time.Now())
	return err
}

func (m *defaultRolePermissionModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where role_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *defaultRolePermissionModel) DeleteByPermissionId(ctx context.Context, permissionId int64) error {
	query := fmt.Sprintf("delete from %s where permission_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, permissionId)
	return err
}

func (m *defaultRolePermissionModel) FindPermissionIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("select permission_id from %s where role_id = $1", m.table)
	var resp []int64
	err := m.conn.QueryRowsCtx(ctx, &resp, query, roleId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultRolePermissionModel) BatchInsert(ctx context.Context, roleId int64, permissionIds []int64) error {
	if len(permissionIds) == 0 {
		return nil
	}

	// 使用事务批量插入
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, permissionId := range permissionIds {
			query := fmt.Sprintf("insert into %s (role_id, permission_id, created_at) values ($1, $2, $3)", m.table)
			_, err := session.ExecCtx(ctx, query, roleId, permissionId, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
}
