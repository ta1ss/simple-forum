package main

import (
	database "forum/data"
	"forum/utils"
)

func main() {
	database.ConnectDB()
	defer database.CloseDB()

	utils.RunServer()
}
