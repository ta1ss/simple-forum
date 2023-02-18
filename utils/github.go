package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// <!-- GITHUB CLIEND ID -->
// <!-- de61cc110026ec3b4a3a -->
// <!-- CLIENT SECRET -->
// <!-- 92fd55bf30bbc5570ea775e67247a0c74c631244 -->

var (
	//oAuth2.0 Configurations
	githubOauthConfig = &oauth2.Config{
		ClientID:     "de61cc110026ec3b4a3a",
		ClientSecret: "92fd55bf30bbc5570ea775e67247a0c74c631244",
		RedirectURL:  "https://localhost/ghcallback",
		Scopes:       []string{"user:email", "user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	githuboauthStateString = "random"
)

func ghLogIn(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(githuboauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func ghCallBack(w http.ResponseWriter, r *http.Request) {

	// do the state and code security checks with Google
	state := r.FormValue("state")
	if state != githuboauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", githuboauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Printf("Failed to get user info with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	// Unmarshal the response body into a struct
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"login"`
	}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		log.Fatal(err)
	}

	// make the necessary changes in db if needed
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	password, _ := HashPassword("earthisflat")
	username := userInfo.Name
	email := userInfo.Email
	status := "user"

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM customer WHERE email = ?", email).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO customer (email, username, password, status) VALUES (?, ?, ?, ?) ", email, username, password, status)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = db.Exec("UPDATE customer SET username = ? WHERE email = ?", username, email)
		if err != nil {
			fmt.Println(err)
		}
	}
	setSessionCookie(w, username)
	http.Redirect(w, r, "/", http.StatusFound)

	//Terminal Output
	fmt.Printf("User %v - Github-Login Successful\n", username)
	fmt.Println("---------------------------------")
}
