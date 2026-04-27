package model

import "errors"

// 定义模型层常用的错误
var (
	// ErrNotFound 记录未找到错误
	// 当查询数据库时，如果没有找到对应的记录，返回此错误
	ErrNotFound = errors.New("record not found")

	// ErrDuplicateKey 重复键错误
	// 当插入数据时，如果违反唯一约束，返回此错误
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrInvalidParam 无效参数错误
	// 当传入的参数无效时，返回此错误
	ErrInvalidParam = errors.New("invalid parameter")

	// ErrNoRowsAffected 没有行受影响错误
	// 当执行更新或删除操作时，如果没有任何行受到影响，返回此错误
	ErrNoRowsAffected = errors.New("no rows affected")

	// ErrConnectionFailed 数据库连接失败错误
	// 当无法连接到数据库时，返回此错误
	ErrConnectionFailed = errors.New("database connection failed")

	// ErrTransactionFailed 事务失败错误
	// 当事务执行失败时，返回此错误
	ErrTransactionFailed = errors.New("transaction failed")
)
