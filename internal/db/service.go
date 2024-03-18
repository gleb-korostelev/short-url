package db

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func InitDB() *Database {
	сonnection, err := pgx.Connect(context.Background(), config.DBDSN)
	if err != nil {
		logger.Infof("Unable to connect to database: %v\n", err)
		return nil
	}
	logger.Infof("Connected to database.")
	return &Database{Conn: сonnection}
}

func (db *Database) Close() {
	db.Conn.Close(context.Background())
}
