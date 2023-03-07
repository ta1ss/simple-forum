package utils

import (
	"fmt"
	database "forum/data"
)

func GetNotifications(id int) []string {
	var text []string
	rows, err := database.Query("SELECT action FROM notifications where userID =?", id)
	if err != nil {
		fmt.Println(err)
	}

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
	_, err := database.Exec("DELETE FROM notifications WHERE userID = ?", id)
	if err != nil {
		fmt.Println(err)
	}
}
