package utils

import (
	"fmt"
	database "forum/data"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func setSessionCookie(w http.ResponseWriter, username string) {
	sessionID, _ := uuid.NewV4()
	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID.String(),
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	_, err := database.Exec("DELETE FROM sessionIDs WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
	}
	_, err = database.Exec("INSERT INTO sessionIDs(sessionID, expires, username) VALUES (?,?,?)", sessionID.String(), expiration, username)
	if err != nil {
		fmt.Println(err)
	}
}

func checkSessionCookie(r *http.Request) bool {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return false
	}
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count == 1 {
		return true
	}
	return false
}

func userIDFromCookie(r *http.Request) int {
	cookie, _ := r.Cookie("sessionID")
	var user string
	err := database.QueryRow("SELECT username FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&user)
	if err != nil {
		fmt.Println(err)
	}
	var userID int
	err = database.QueryRow("SELECT ID FROM customer WHERE username = ?", user).Scan(&userID)
	if err != nil {
		fmt.Println(err)
	}
	return userID
}
