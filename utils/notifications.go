package utils

import (
	"database/sql"
	"fmt"
)

// gets all the notifications for a certain userID
func GetNotifications(id int) []string {
	var text []string
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT action FROM notifications where userID =?", id)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			fmt.Println(err)
		}
		text = append(text, t)
	}
	return text
}

func DeleteNotifications(id int) {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("DELETE FROM notifications WHERE userID = ?", id)
	if err != nil {
		fmt.Println(err)
	}
}
