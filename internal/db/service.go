// Package db provides the database interface and utilities for handling
// connections and operations with PostgreSQL using the pgx library.
package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB defines the interface for database operations that any DB implementation must satisfy.
// It provides methods for executing queries and commands, and for managing the database connection.
type DB interface {
	// Close terminates the database connection.
	Close() error

	// Ping tests the connectivity with the database.
	// It returns an error if the database is not reachable.
	Ping(ctx context.Context) error

	// Exec executes a SQL query without returning any rows.
	// The args are for any placeholder parameters in the query.
	// It returns a CommandTag, which reports the outcome of the query (like rows affected).
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)

	// Query executes a SQL query that returns rows, typically SELECT.
	// The args are for any placeholder parameters in the query.
	// It returns Rows, which is an iterator for reading fetched data.
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)

	// QueryRow executes a SQL query that is expected to return at most one row.
	// The args are for any placeholder parameters in the query.
	// QueryRow returns a Row, which is a lazy loader for fetching the single row.
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row

	// GetConn returns a connection pool object which can be used to execute multiple queries
	// with a single database connection. It is useful for transactions.
	GetConn(ctx context.Context) *pgxpool.Pool
}
