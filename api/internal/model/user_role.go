package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRoleModel = (*defaultUserRoleModel)(nil)

type (
	// UserRole is an object of map UserRole
	UserRole struct {
		Id        int64     `db:"id"`
		UserId    int64     `db:"user_id"`
		RoleId    int64     `db:"role_id"`
		CreatedAt time.Time `db:"created_at"`
	}

	// UserRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRoleModel.
	UserRoleModel interface {
		Insert(ctx context.Context, data *UserRole) error
		DeleteByUserId(ctx context.Context, userId int64) error
		DeleteByRoleId(ctx context.Context, roleId int64) error
		FindRoleIdsByUserId(ctx context.Context, userId int64) ([]int64, error)
		BatchInsert(ctx context.Context, userId int64, roleIds []int64) error
	}

	defaultUserRoleModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewUserRoleModel returns a model for the database table.
func NewUserRoleModel(conn sqlx.SqlConn) UserRoleModel {
	return &defaultUserRoleModel{
		conn:  conn,
		table: "user_roles",
	}
}

func (m *defaultUserRoleModel) Insert(ctx context.Context, data *UserRole) error {
	query := fmt.Sprintf("insert into %s (user_id, role_id, created_at) values ($1, $2, $3)", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.RoleId, time.Now())
	return err
}

func (m *defaultUserRoleModel) DeleteByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("delete from %s where user_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *defaultUserRoleModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where role_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *defaultUserRoleModel) FindRoleIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select role_id from %s where user_id = $1", m.table)
	var resp []int64
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultUserRoleModel) BatchInsert(ctx context.Context, userId int64, roleIds []int64) error {
	if len(roleIds) == 0 {
		return nil
	}

	// 使用事务批量插入
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, roleId := range roleIds {
			query := fmt.Sprintf("insert into %s (user_id, role_id, created_at) values ($1, $2, $3)", m.table)
			_, err := session.ExecCtx(ctx, query, userId, roleId, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
}
