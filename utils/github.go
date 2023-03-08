package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "forum/data"

	"golang.org/x/oauth2"
)

var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     "de61cc110026ec3b4a3a",
		ClientSecret: "92fd55bf30bbc5570ea775e67247a0c74c631244",
		RedirectURL:  "https://176.112.158.14:8443/ghcallback",
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

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"login"`
	}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
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
