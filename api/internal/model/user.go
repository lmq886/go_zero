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
	"github.com/zeromicro/go-zero/core/stringx"
)

var userFieldNames = builder.RawFieldNames(&User{})
var userRows = strings.Join(userFieldNames, ",")

var _ UserModel = (*customUserModel)(nil)

type (
	// User is an object of map User
	User struct {
		Id        int64          `db:"id"`
		Username  string         `db:"username"`
		Password  string         `db:"password"`
		Nickname  sql.NullString `db:"nickname"`
		Avatar    sql.NullString `db:"avatar"`
		Email     sql.NullString `db:"email"`
		Phone     sql.NullString `db:"phone"`
		Status    int64          `db:"status"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
		DeletedAt sql.NullTime   `db:"deleted_at"`
	}

	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindOneByUsername(ctx context.Context, username string) (*User, error)
		FindPage(ctx context.Context, page, pageSize int64, username, nickname string, status int64) ([]*User, int64, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	query := fmt.Sprintf("select %s from %s where username = $1 and deleted_at is null limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) FindPage(ctx context.Context, page, pageSize int64, username, nickname string, status int64) ([]*User, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if username != "" {
		conditions = append(conditions, fmt.Sprintf("username like $%d", argIndex))
		args = append(args, "%"+username+"%")
		argIndex++
	}

	if nickname != "" {
		conditions = append(conditions, fmt.Sprintf("nickname like $%d", argIndex))
		args = append(args, "%"+nickname+"%")
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
	query := fmt.Sprintf("select %s from %s where %s order by id desc limit $%d offset $%d",
		userRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*User
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

// userModel is a default model implementation for table users
type userModel interface {
	Insert(ctx context.Context, data *User) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, id int64) error
}

type defaultUserModel struct {
	conn  sqlx.SqlConn
	table string
}

func newUserModel(conn sqlx.SqlConn) *defaultUserModel {
	return &defaultUserModel{
		conn:  conn,
		table: "users",
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (username, password, nickname, avatar, email, phone, status, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Username, data.Password, data.Nickname, data.Avatar, data.Email, data.Phone, data.Status, time.Now(), time.Now())
	return ret, err
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", userRows, m.table)
	var resp User
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

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := fmt.Sprintf("update %s set username = $1, password = $2, nickname = $3, avatar = $4, email = $5, phone = $6, status = $7, updated_at = $8 where id = $9", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Username, data.Password, data.Nickname, data.Avatar, data.Email, data.Phone, data.Status, time.Now(), data.Id)
	return err
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

var (
	cacheUserIdPrefix      = "cache:user:id:"
	cacheUserUsernamePrefix = "cache:user:username:"
)

func cacheUserIdKey(id int64) string {
	return fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
}

func cacheUserUsernameKey(username string) string {
	return fmt.Sprintf("%s%s", cacheUserUsernamePrefix, username)
}

// ErrNotFound is an alias of sqlc.ErrNotFound
var ErrNotFound = sqlc.ErrNotFound
