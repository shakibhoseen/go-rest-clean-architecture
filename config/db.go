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

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	return db
}
