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

var operationLogFieldNames = builder.RawFieldNames(&OperationLog{})
var operationLogRows = strings.Join(operationLogFieldNames, ",")

var _ OperationLogModel = (*defaultOperationLogModel)(nil)

type (
	// OperationLog is an object of map OperationLog
	OperationLog struct {
		Id            int64          `db:"id"`
		UserId        sql.NullInt64  `db:"user_id"`
		Username      sql.NullString `db:"username"`
		Operation     string         `db:"operation"`
		Method        string         `db:"method"`
		RequestUri    string         `db:"request_uri"`
		RequestParams sql.NullString `db:"request_params"`
		ResponseData  sql.NullString `db:"response_data"`
		Ip            sql.NullString `db:"ip"`
		Location      sql.NullString `db:"location"`
		Browser       sql.NullString `db:"browser"`
		Os            sql.NullString `db:"os"`
		Status        int64          `db:"status"`
		ErrorMsg      sql.NullString `db:"error_msg"`
		Duration      int64          `db:"duration"`
		CreatedAt     time.Time      `db:"created_at"`
	}

	// OperationLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOperationLogModel.
	OperationLogModel interface {
		Insert(ctx context.Context, data *OperationLog) (sql.Result, error)
		Delete(ctx context.Context, id int64) error
		FindPage(ctx context.Context, page, pageSize int64, username, operation, method string, status int64, startTime, endTime string) ([]*OperationLog, int64, error)
	}

	defaultOperationLogModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewOperationLogModel returns a model for the database table.
func NewOperationLogModel(conn sqlx.SqlConn) OperationLogModel {
	return &defaultOperationLogModel{
		conn:  conn,
		table: "operation_logs",
	}
}

func (m *defaultOperationLogModel) Insert(ctx context.Context, data *OperationLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (user_id, username, operation, method, request_uri, request_params, response_data, ip, location, browser, os, status, error_msg, duration, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Username, data.Operation, data.Method, data.RequestUri, data.RequestParams, data.ResponseData, data.Ip, data.Location, data.Browser, data.Os, data.Status, data.ErrorMsg, data.Duration, time.Now())
	return ret, err
}

func (m *defaultOperationLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultOperationLogModel) FindPage(ctx context.Context, page, pageSize int64, username, operation, method string, status int64, startTime, endTime string) ([]*OperationLog, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if username != "" {
		conditions = append(conditions, fmt.Sprintf("username like $%d", argIndex))
		args = append(args, "%"+username+"%")
		argIndex++
	}

	if operation != "" {
		conditions = append(conditions, fmt.Sprintf("operation like $%d", argIndex))
		args = append(args, "%"+operation+"%")
		argIndex++
	}

	if method != "" {
		conditions = append(conditions, fmt.Sprintf("method = $%d", argIndex))
		args = append(args, method)
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
		operationLogRows, m.table, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var resp []*OperationLog
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
