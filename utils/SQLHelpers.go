package utils

// ! DB TODO
import (
	"database/sql"
	"fmt"
)

func GetUsernameByID(id int) string {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}

	var username string
	err = db.QueryRow("SELECT username FROM customer WHERE ID = ?", id).Scan(&username)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return username
}

func GetStatusByID(id int) string {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}

	var status string
	err = db.QueryRow("SELECT status FROM customer WHERE ID = ?", id).Scan(&status)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return status
}

func GetCategoryName(id int) string {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	var title string
	err = db.QueryRow("SELECT title FROM categories WHERE ID = ?", id).Scan(&title)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return title
}
