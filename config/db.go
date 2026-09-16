package config

import (
	"database/sql"
	"log"

	//_ "modernc.org/sqlite" // আন্ডারস্কোর দিয়ে শুধু ড্রাইভার রেজিস্টার করা
	_ "github.com/jackc/pgx/v5/stdlib" // pgx ড্রাইভার
)

func InitDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Connection error:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
}
