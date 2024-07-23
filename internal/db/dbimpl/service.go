// Package dbimpl implements the database interface defined in the db package for PostgreSQL.
// It utilizes the pgx library to establish and manage the database connection pool.
package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/db"
	"github.com/gleb-korostelev/short-url/tools/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database wraps a pgxpool.Pool to manage and interact with the database.
type Database struct {
	Conn *pgxpool.Pool // Conn is a connection pool for PostgreSQL.
}

// InitDB initializes and returns a new instance of Database.
// It establishes a connection pool using the DSN provided in the configuration,
// logs the connection status, and initializes necessary database tables.
func InitDB() (db.DB, error) {
	connection, err := pgxpool.New(context.Background(), config.DBDSN)
	if err != nil {
		logger.Infof("Unable to connect to database: %v", err)
		return nil, err
	}
	logger.Infof("Connected to database.")
	data := &Database{Conn: connection}

	err = InitializeTables(data)
	if err != nil {
		logger.Infof("Failed to initialize tables: %v", err)
		return nil, err
	}
	return data, nil
}

// GetConn retrieves the database connection pool.
func (db *Database) GetConn(ctx context.Context) *pgxpool.Pool {
	return db.Conn
}

// Close terminates the database connection pool.
func (db *Database) Close() error {
	db.Conn.Close()
	return nil
}

// Ping verifies the connection to the database by attempting to contact the server.
func (db *Database) Ping(ctx context.Context) error {
	err := db.Conn.Ping(ctx)
	if err != nil {
		logger.Errorf("Internal error: %v", err)
		return err
	}
	return nil
}

// Exec executes a SQL command (insert, update, delete, etc.) and returns the result tag.
func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tg, err := db.Conn.Exec(ctx, query, args...)
	if err != nil {
		logger.Errorf("Internal error: %v", err)
		return pgconn.CommandTag{}, err
	}
	return tg, nil
}

// Query executes a SQL query and returns the rows as a pgx.Rows object.
func (db *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		logger.Errorf("Internal error: %v", err)
		return nil, err
	}
	return rows, nil
}

// QueryRow executes a SQL query that is expected to return at most one row.
func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Conn.QueryRow(ctx, query, args...)
}
