package repository

import (
	"database/sql"
	"log"
	"time"
)

var DB *sql.DB

func InitDB(dsn string) error {
	var err error
	for i := 0; i < 30; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err == nil {
			err = DB.Ping()
		}
		if err == nil {
			DB.SetMaxOpenConns(20)
			DB.SetMaxIdleConns(5)
			DB.SetConnMaxLifetime(5 * time.Minute)
			log.Println("[DB] Connected")
			return nil
		}
		log.Printf("[DB] Waiting... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return err
}