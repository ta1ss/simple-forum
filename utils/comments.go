package utils

import (
	"fmt"
	database "forum/data"
	"net/http"
	"strings"
	"time"
)

type Comments struct {
	ID             int
	Body           string
	CreatedAt      string
	UserID         int
	PostID         int
	Username       string
	Likes_Dislikes []Likes_Dislikes
	TotalLikes     int
	TotalDislikes  int
	PostTitle      string
	PostBody       string
	PostMedia      string
	Flag           bool
}

func getComments(id int) []Comments {
	rows, err := database.Query("SELECT ID, body, creation, userID, postID, flag FROM comments WHERE postID = ?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.ID, &comment.Body, &comment.CreatedAt, &comment.UserID, &comment.PostID, &comment.Flag); err != nil {
			fmt.Println(err)
			fmt.Println("failed to read COMMENTS")
		}
		comment.Username = GetUsernameByID(comment.UserID)
		comments = append(comments, comment)
	}
	return comments
}

func AddComment(comment, postID string, user int, w http.ResponseWriter) {
	flag := false
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := database.Exec("INSERT INTO comments (body, creation, userID, postID, flag) VALUES (?, ?, ?, ?, ?)", comment, currentTime, user, postID, flag)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			_, err = database.Exec("CREATE TABLE comments (id INTEGER PRIMARY KEY, body, creation TIMESTAMP, userID INTEGER, postID INTEGER)")
			if err != nil {
				http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
			}
		}
	}
	var originalAuthor int
	var postTitle string
	commentAuhtor := GetUsernameByID(user)
	err = database.QueryRow("SELECT userID, title FROM posts WHERE ID =?", postID).Scan(&originalAuthor, &postTitle)
	if err != nil {
		fmt.Println(err)
	}
	text := "has commented on your post:" + postTitle
	action := commentAuhtor + " " + text
	if originalAuthor != user {
		_, err = database.Exec("INSERT INTO notifications (userID, action) VALUES (?, ?)", originalAuthor, action)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func DeleteComment(commentID, user int, w http.ResponseWriter) {
	var createdby int
	err := database.QueryRow("SELECT userID FROM comments WHERE ID = ?", commentID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}
	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		_, err := database.Exec("DELETE FROM likes_dislikes WHERE commentID = ?", commentID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = database.Exec("DELETE FROM comments WHERE ID = ?", commentID)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdateComment(commentID, user int, newcomment string, w http.ResponseWriter) {
	var createdby int
	err := database.QueryRow("SELECT userID FROM comments WHERE ID = ?", commentID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}

	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		if len(newcomment) > 0 {
			_, err := database.Exec("UPDATE comments SET body = ? WHERE ID = ?", newcomment, commentID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func AddCommentFlag(id int) {
	_, err := database.Exec("UPDATE comments SET flag = ? WHERE ID = ?", true, id)
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveCommentFlag(id int) {
	_, err := database.Exec("UPDATE comments SET flag = ? WHERE ID = ?", false, id)
	if err != nil {
		fmt.Println(err)
	}
}

func GetFlaggedComments() []Comments {
	rows, err := database.Query("SELECT ID, body, creation, userID, postID, flag FROM comments WHERE flag = ?", true)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.ID, &comment.Body, &comment.CreatedAt, &comment.UserID, &comment.PostID, &comment.Flag); err != nil {
			fmt.Println(err)
			fmt.Println("failed to read COMMENTS")
		}
		comment.Username = GetUsernameByID(comment.UserID)
		comments = append(comments, comment)
	}
	return comments
}
