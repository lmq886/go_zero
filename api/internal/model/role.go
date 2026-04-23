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

var roleFieldNames = builder.RawFieldNames(&Role{})
var roleRows = strings.Join(roleFieldNames, ",")

var _ RoleModel = (*customRoleModel)(nil)

type (
	// Role is an object of map Role
	Role struct {
		Id          int64          `db:"id"`
		Name        string         `db:"name"`
		Code        string         `db:"code"`
		Description sql.NullString `db:"description"`
		Status      int64          `db:"status"`
		Sort        int64          `db:"sort"`
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   time.Time      `db:"updated_at"`
		DeletedAt   sql.NullTime   `db:"deleted_at"`
	}

	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		FindByCodes(ctx context.Context, codes []string) ([]*Role, error)
		FindByUserId(ctx context.Context, userId int64) ([]*Role, error)
		FindPage(ctx context.Context, page, pageSize int64, name string, status int64) ([]*Role, int64, error)
	}

	customRoleModel struct {
		*defaultRoleModel
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn),
	}
}

func (m *customRoleModel) FindByCodes(ctx context.Context, codes []string) ([]*Role, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(codes))
	args := make([]interface{}, len(codes))
	for i, code := range codes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = code
	}

	query := fmt.Sprintf("select %s from %s where code in (%s) and deleted_at is null",
		roleRows, m.table, strings.Join(placeholders, ","))

	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customRoleModel) FindByUserId(ctx context.Context, userId int64) ([]*Role, error) {
	query := fmt.Sprintf(`select %s from %s r 
		inner join user_roles ur on r.id = ur.role_id 
		where ur.user_id = $1 and r.deleted_at is null`, roleRows, m.table)

	var resp []*Role
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

func (m *customRoleModel) FindPage(ctx context.Context, page, pageSize int64, name string, status int64) ([]*Role, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name like $%d", argIndex))
		args = append(args, "%"+name+"%")
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
		roleRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*Role
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

// roleModel is a default model implementation for table roles
type roleModel interface {
	Insert(ctx context.Context, data *Role) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*Role, error)
	Update(ctx context.Context, data *Role) error
	Delete(ctx context.Context, id int64) error
}

type defaultRoleModel struct {
	conn  sqlx.SqlConn
	table string
}

func newRoleModel(conn sqlx.SqlConn) *defaultRoleModel {
	return &defaultRoleModel{
		conn:  conn,
		table: "roles",
	}
}

func (m *defaultRoleModel) Insert(ctx context.Context, data *Role) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (name, code, description, status, sort, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Code, data.Description, data.Status, data.Sort, time.Now(), time.Now())
	return ret, err
}

func (m *defaultRoleModel) FindOne(ctx context.Context, id int64) (*Role, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", roleRows, m.table)
	var resp Role
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

func (m *defaultRoleModel) Update(ctx context.Context, data *Role) error {
	query := fmt.Sprintf("update %s set name = $1, code = $2, description = $3, status = $4, sort = $5, updated_at = $6 where id = $7", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Code, data.Description, data.Status, data.Sort, time.Now(), data.Id)
	return err
}

func (m *defaultRoleModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
