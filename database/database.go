package database

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() (*PostQreSQLCon, error) {

	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres sslmode=require password=aymaan132 host=database.c9oigwacc6k7.ap-south-1.rds.amazonaws.com")
	if err != nil {
		slog.Error("Failed to open database connection", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Successfully connected to the database")

	return &PostQreSQLCon{
		dbCon: db,
	}, nil
}
