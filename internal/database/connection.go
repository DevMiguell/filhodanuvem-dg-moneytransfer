package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateConnection() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL is not set, cannot connect to db")
	}
	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, err
}
