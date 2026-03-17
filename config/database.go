package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB(cfg *DBConfig) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	fmt.Println("Database connected successfully")
	return conn
}
