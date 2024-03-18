package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DatabaseI interface {
	Close() error
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...interface{}) pgx.Rows
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetConn(ctx context.Context) *pgx.Conn
}
