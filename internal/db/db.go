package db

import (
	"database/sql"
	"log"

	"github.com/Siddheshk02/coupon-system/internal/config"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Init(cfg config.Config) {
	var err error
	Conn, err = sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	if err := Conn.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	log.Println("Connected to database")
}
