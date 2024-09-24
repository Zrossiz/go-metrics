package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	DBDSN := os.Getenv("DATABASE_DSN")

	db, err := pgxpool.Connect(context.Background(), DBDSN)
	if err != nil {
		fmt.Errorf("error connect to db: %v", err)
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		metric_type TEXT NOT NULL,
		name TEXT NOT NULL,
		value DOUBLE PRECISION,
		delta BIGINT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics (name);`)

	if err != nil {
		fmt.Errorf("failed to create table: %w", err)
	}
}
