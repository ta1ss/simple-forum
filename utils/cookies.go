package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func setSessionCookie(w http.ResponseWriter, username string) {
	//uses UUID to generate a unique sessionID for the user
	sessionID, _ := uuid.NewV4()
	expiration := time.Now().Add(24 * time.Hour)
	//Creates the cookie with HTTP only, meaning our program will only
	//send the cookie to http requests, so for example JS cant see it.
	//Secure means the cookie will be sent only over secure connections
	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID.String(),
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	//adds the cookie to our database for later validation
	db, _ := sql.Open("sqlite3", "file:data/database.db")
	//first remove any existing sessions for the user
	_, err := db.Exec("DELETE FROM sessionIDs WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("INSERT INTO sessionIDs(sessionID, expires, username) VALUES (?,?,?)", sessionID.String(), expiration, username)
	if err != nil {
		fmt.Println(err)
	}
}

func checkSessionCookie(r *http.Request) bool {
	//if there is no such cookie, return false
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return false
	}
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	//SQL will the return the number of matches it got from cookie table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count == 1 {
		return true
	}
	return false
}

// function to get the userID from cookie value
func userIDFromCookie(r *http.Request) int {
	cookie, _ := r.Cookie("sessionID")
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	var user string
	err = db.QueryRow("SELECT username FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&user)
	if err != nil {
		fmt.Println(err)
	}
	var userID int
	err = db.QueryRow("SELECT ID FROM customer WHERE username = ?", user).Scan(&userID)

	if err != nil {
		fmt.Println(err)
	}
	return userID
}
