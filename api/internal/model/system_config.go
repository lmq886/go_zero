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

var systemConfigFieldNames = builder.RawFieldNames(&SystemConfig{})
var systemConfigRows = strings.Join(systemConfigFieldNames, ",")

var _ SystemConfigModel = (*customSystemConfigModel)(nil)

type (
	// SystemConfig is an object of map SystemConfig
	SystemConfig struct {
		Id        int64          `db:"id"`
		Key       string         `db:"key"`
		Value     string         `db:"value"`
		Name      string         `db:"name"`
		Remark    sql.NullString `db:"remark"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
		DeletedAt sql.NullTime   `db:"deleted_at"`
	}

	// SystemConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSystemConfigModel.
	SystemConfigModel interface {
		systemConfigModel
		FindByKey(ctx context.Context, key string) (*SystemConfig, error)
		FindPage(ctx context.Context, page, pageSize int64, key, name string) ([]*SystemConfig, int64, error)
	}

	customSystemConfigModel struct {
		*defaultSystemConfigModel
	}
)

// NewSystemConfigModel returns a model for the database table.
func NewSystemConfigModel(conn sqlx.SqlConn) SystemConfigModel {
	return &customSystemConfigModel{
		defaultSystemConfigModel: newSystemConfigModel(conn),
	}
}

func (m *customSystemConfigModel) FindByKey(ctx context.Context, key string) (*SystemConfig, error) {
	query := fmt.Sprintf("select %s from %s where key = $1 and deleted_at is null limit 1", systemConfigRows, m.table)
	var resp SystemConfig
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customSystemConfigModel) FindPage(ctx context.Context, page, pageSize int64, key, name string) ([]*SystemConfig, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if key != "" {
		conditions = append(conditions, fmt.Sprintf("key like $%d", argIndex))
		args = append(args, "%"+key+"%")
		argIndex++
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name like $%d", argIndex))
		args = append(args, "%"+name+"%")
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
		systemConfigRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*SystemConfig
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

// systemConfigModel is a default model implementation for table system_configs
type systemConfigModel interface {
	Insert(ctx context.Context, data *SystemConfig) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*SystemConfig, error)
	Update(ctx context.Context, data *SystemConfig) error
	Delete(ctx context.Context, id int64) error
}

type defaultSystemConfigModel struct {
	conn  sqlx.SqlConn
	table string
}

func newSystemConfigModel(conn sqlx.SqlConn) *defaultSystemConfigModel {
	return &defaultSystemConfigModel{
		conn:  conn,
		table: "system_configs",
	}
}

func (m *defaultSystemConfigModel) Insert(ctx context.Context, data *SystemConfig) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (key, value, name, remark, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Key, data.Value, data.Name, data.Remark, time.Now(), time.Now())
	return ret, err
}

func (m *defaultSystemConfigModel) FindOne(ctx context.Context, id int64) (*SystemConfig, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", systemConfigRows, m.table)
	var resp SystemConfig
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

func (m *defaultSystemConfigModel) Update(ctx context.Context, data *SystemConfig) error {
	query := fmt.Sprintf("update %s set key = $1, value = $2, name = $3, remark = $4, updated_at = $5 where id = $6", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Key, data.Value, data.Name, data.Remark, time.Now(), data.Id)
	return err
}

func (m *defaultSystemConfigModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
