package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var loginLogFieldNames = builder.RawFieldNames(&LoginLog{})
var loginLogRows = strings.Join(loginLogFieldNames, ",")

var _ LoginLogModel = (*defaultLoginLogModel)(nil)

type (
	// LoginLog is an object of map LoginLog
	LoginLog struct {
		Id        int64          `db:"id"`
		UserId    sql.NullInt64  `db:"user_id"`
		Username  sql.NullString `db:"username"`
		Ip        sql.NullString `db:"ip"`
		Location  sql.NullString `db:"location"`
		Browser   sql.NullString `db:"browser"`
		Os        sql.NullString `db:"os"`
		Status    int64          `db:"status"`
		Msg       sql.NullString `db:"msg"`
		CreatedAt time.Time      `db:"created_at"`
	}

	// LoginLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginLogModel.
	LoginLogModel interface {
		Insert(ctx context.Context, data *LoginLog) (sql.Result, error)
		Delete(ctx context.Context, id int64) error
		FindPage(ctx context.Context, page, pageSize int64, username string, status int64, startTime, endTime string) ([]*LoginLog, int64, error)
	}

	defaultLoginLogModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewLoginLogModel returns a model for the database table.
func NewLoginLogModel(conn sqlx.SqlConn) LoginLogModel {
	return &defaultLoginLogModel{
		conn:  conn,
		table: "login_logs",
	}
}

func (m *defaultLoginLogModel) Insert(ctx context.Context, data *LoginLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (user_id, username, ip, location, browser, os, status, msg, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Username, data.Ip, data.Location, data.Browser, data.Os, data.Status, data.Msg, time.Now())
	return ret, err
}

func (m *defaultLoginLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLoginLogModel) FindPage(ctx context.Context, page, pageSize int64, username string, status int64, startTime, endTime string) ([]*LoginLog, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if username != "" {
		conditions = append(conditions, fmt.Sprintf("username like $%d", argIndex))
		args = append(args, "%"+username+"%")
		argIndex++
	}

	if status > 0 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if startTime != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, startTime)
		argIndex++
	}

	if endTime != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, endTime)
		argIndex++
	}

	var whereClause string
	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " and ")
	} else {
		whereClause = "1=1"
	}

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
		loginLogRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*LoginLog
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
