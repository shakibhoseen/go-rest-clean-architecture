package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // আন্ডারস্কোর দিয়ে শুধু ড্রাইভার রেজিস্টার করা
)

func InitDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
}
