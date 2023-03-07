package utils

import (
	"context"
	"encoding/json"
	"fmt"
	database "forum/data"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "789924993075-pmhol8b8bi3ft24ln190643pjbkcqbeo.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "https://localhost/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random"
)

func gLogIn(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func gCallBack(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("Failed to get user info with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Fatal(err)
	}

	password, _ := HashPassword("earthisflat")
	username := userInfo.Name
	email := userInfo.Email
	status := "user"

	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM customer WHERE email = ?", email).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count == 0 {
		_, err = database.Exec("INSERT INTO customer (email, username, password, status) VALUES (?, ?, ?, ?) ", email, username, password, status)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = database.Exec("UPDATE customer SET username = ? WHERE email = ?", username, email)
		if err != nil {
			fmt.Println(err)
		}
	}
	setSessionCookie(w, username)
	http.Redirect(w, r, "/", http.StatusFound)
}
