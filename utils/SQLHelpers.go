package utils

import (
	"fmt"
	database "forum/data"
)

func GetUsernameByID(id int) string {
	var username string
	err := database.QueryRow("SELECT username FROM customer WHERE ID = ?", id).Scan(&username)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return username
}

func GetStatusByID(id int) string {
	var status string
	err := database.QueryRow("SELECT status FROM customer WHERE ID = ?", id).Scan(&status)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return status
}

func GetCategoryName(id int) string {
	var title string
	err := database.QueryRow("SELECT title FROM categories WHERE ID = ?", id).Scan(&title)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return title
}
