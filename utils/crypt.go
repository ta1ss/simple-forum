package utils

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	// built in function, that takes the regular password and compares it against stored has
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(db *sql.DB, username, password string) bool {
	var dbUsername, dbPassword string
	err := db.QueryRow("SELECT username, password FROM customer WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
	}
	// compare password to hashed passwhord in db
	if CheckPasswordHash(password, dbPassword) {
		return true
	}
	return false
}
