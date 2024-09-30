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
		fmt.Printf("error connecting to db: %v\n", err)
		return
	}
	defer db.Close()

	sqlFilePath := "migrations/1_init.sql"
	schemaSQL, err := os.ReadFile(sqlFilePath)
	if err != nil {
		fmt.Printf("error reading sql file: %v\n", err)
		return
	}

	_, err = db.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		fmt.Printf("failed to execute schema SQL: %v\n", err)
		return
	}

	fmt.Println("Schema created successfully!")
}
