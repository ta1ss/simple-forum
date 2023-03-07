package database

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func ConnectDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	fmt.Println("Database connection established")
}

func CloseDB() {
	db.Close()
	fmt.Println("Database connection closed")
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Prepare(query string) (*sql.Stmt, error) {
	return db.Prepare(query)
}
