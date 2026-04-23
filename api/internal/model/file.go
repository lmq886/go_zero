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

var fileFieldNames = builder.RawFieldNames(&File{})
var fileRows = strings.Join(fileFieldNames, ",")

var _ FileModel = (*customFileModel)(nil)

type (
	// File is an object of map File
	File struct {
		Id           int64          `db:"id"`
		Name         string         `db:"name"`
		OriginalName string        `db:"original_name"`
		Path         string         `db:"path"`
		Url          sql.NullString `db:"url"`
		Size         int64          `db:"size"`
		Type         string         `db:"type"`
		Extension    sql.NullString `db:"extension"`
		Md5          sql.NullString `db:"md5"`
		UserId       sql.NullInt64  `db:"user_id"`
		CreatedAt    time.Time      `db:"created_at"`
		DeletedAt    sql.NullTime   `db:"deleted_at"`
	}

	// FileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileModel.
	FileModel interface {
		fileModel
		FindByMd5(ctx context.Context, md5 string) (*File, error)
		FindPage(ctx context.Context, page, pageSize int64, name, typ string) ([]*File, int64, error)
	}

	customFileModel struct {
		*defaultFileModel
	}
)

// NewFileModel returns a model for the database table.
func NewFileModel(conn sqlx.SqlConn) FileModel {
	return &customFileModel{
		defaultFileModel: newFileModel(conn),
	}
}

func (m *customFileModel) FindByMd5(ctx context.Context, md5 string) (*File, error) {
	query := fmt.Sprintf("select %s from %s where md5 = $1 and deleted_at is null limit 1", fileRows, m.table)
	var resp File
	err := m.conn.QueryRowCtx(ctx, &resp, query, md5)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFileModel) FindPage(ctx context.Context, page, pageSize int64, name, typ string) ([]*File, int64, error) {
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
	query := fmt.Sprintf("select %s from %s where %s order by created_at desc limit $%d offset $%d",
		fileRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*File
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

// fileModel is a default model implementation for table files
type fileModel interface {
	Insert(ctx context.Context, data *File) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*File, error)
	Update(ctx context.Context, data *File) error
	Delete(ctx context.Context, id int64) error
}

type defaultFileModel struct {
	conn  sqlx.SqlConn
	table string
}

func newFileModel(conn sqlx.SqlConn) *defaultFileModel {
	return &defaultFileModel{
		conn:  conn,
		table: "files",
	}
}

func (m *defaultFileModel) Insert(ctx context.Context, data *File) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (name, original_name, path, url, size, type, extension, md5, user_id, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.OriginalName, data.Path, data.Url, data.Size, data.Type, data.Extension, data.Md5, data.UserId, time.Now())
	return ret, err
}

func (m *defaultFileModel) FindOne(ctx context.Context, id int64) (*File, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and deleted_at is null limit 1", fileRows, m.table)
	var resp File
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

func (m *defaultFileModel) Update(ctx context.Context, data *File) error {
	query := fmt.Sprintf("update %s set name = $1, original_name = $2, path = $3, url = $4, size = $5, type = $6, extension = $7, md5 = $8, user_id = $9, updated_at = $10 where id = $11", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.OriginalName, data.Path, data.Url, data.Size, data.Type, data.Extension, data.Md5, data.UserId, time.Now(), data.Id)
	return err
}

func (m *defaultFileModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
