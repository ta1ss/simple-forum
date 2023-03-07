package utils

import (
	"database/sql"
	"fmt"
	database "forum/data"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func login(username, password string) bool {
	var dbUsername, dbPassword string
	err := database.QueryRow("SELECT username, password FROM customer WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
	}
	if CheckPasswordHash(password, dbPassword) {
		return true
	}
	return false
}
