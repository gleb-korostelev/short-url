package impl

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func InitDB() (db.DatabaseI, error) {
	сonnection, err := pgx.Connect(context.Background(), config.DBDSN)
	if err != nil {
		logger.Infof("Unable to connect to database: %v\n", err)
		return nil, err
	}
	logger.Infof("Connected to database.")
	data := &Database{Conn: сonnection}
	return data, nil
}

func (db *Database) GetConn(ctx context.Context) *pgx.Conn {
	return db.Conn
}

func (db *Database) Close() error {
	err := db.Conn.Close(context.Background())
	if err != nil {
		logger.Fatalf("internal error %v", err)
		return err
	}
	return nil
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.Conn.Ping(ctx)
	if err != nil {
		logger.Fatalf("internal error %v", err)
		return err
	}
	return nil
}

func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := db.Conn.Exec(ctx, query, args...)
	if err != nil {
		logger.Fatalf("internal error %v", err)
		return err
	}
	return nil
}

func (db *Database) Query(ctx context.Context, query string, args ...interface{}) pgx.Rows {
	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		logger.Fatalf("internal error %v", err)
		return nil
	}
	return rows
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Conn.QueryRow(ctx, query, args...)
}
