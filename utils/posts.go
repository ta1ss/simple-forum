package utils

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func AddPost(title, body, fileName string, userID int, w http.ResponseWriter) (int, error) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
		return 0, err
	}
	defer db.Close()
	flag := false
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	result, err := db.Exec("INSERT INTO posts (title, body, created, userID, media, flag) VALUES (?, ?, ?, ?, ?, ?)", title, body, currentTime, userID, fileName, flag)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return 0, err
	}
	// Get the ID of the last inserted row
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(postID), nil
}
func AddCategories(postID int, categories []string, w http.ResponseWriter) {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	// Prepare the INSERT statement
	stmt, err := db.Prepare("INSERT INTO postcategories (postID, categoryID) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
	}
	// Execute the INSERT statement for each category
	for _, category := range categories {
		// Convert the category string to an int
		categoryID, err := strconv.Atoi(category)
		if err != nil {
			http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		}
		// Insert the postID and categoryID into the postcategories table
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func AddMedia(filename string, file multipart.File) {
	filePath := "templates/media/" + filename
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Println(err)
	}
}

func DeletePost(postID, user int, w http.ResponseWriter) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
	var createdby int
	err = db.QueryRow("SELECT userID FROM posts WHERE ID = ?", postID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}
	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		_, err := db.Exec("DELETE FROM comments WHERE postID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("DELETE FROM likes_dislikes WHERE postID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("DELETE FROM postcategories WHERE PostID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("DELETE FROM posts WHERE ID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdatePost(postID, user int, newbody, newtitle string, w http.ResponseWriter) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
	var createdby int
	err = db.QueryRow("SELECT userID FROM posts WHERE ID = ?", postID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}
	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		if len(newbody) > 0 {
			_, err := db.Exec("UPDATE posts SET body = ? WHERE ID = ?", newbody, postID)
			if err != nil {
				fmt.Println(err)
			}
		}
		if len(newtitle) > 0 {
			_, err = db.Exec("UPDATE posts SET title = ? WHERE ID = ?", newtitle, postID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// function to add a moderation flag for posts
func AddPostFlag(id int) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("UPDATE posts SET flag = ? WHERE ID = ?", true, id)
	if err != nil {
		fmt.Println(err)
	}
}

// function to remove a moderation flag for posts
func RemovePostFlag(id int) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("UPDATE posts SET flag = ? WHERE ID = ?", false, id)
	if err != nil {
		fmt.Println(err)
	}
}

// function to get all flagged posts
func GetFlaggedPosts() []Posts {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT ID, title, body, created, userID, media, flag FROM posts WHERE flag = ?", true)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	// get the all POSTS
	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UserID, &post.Media, &post.Flag); err != nil {
			fmt.Println("failed to read POSTS")
		}

		// Get the username of the author of the post
		post.Username = GetUsernameByID(post.UserID)

		// get the COMMENTS(with likes) of the POST
		post.Comments = GetComments(post.ID)
		commentsinfo := "commentID"

		for i := range post.Comments {
			post.Comments[i].Likes_Dislikes = GetLikesDislikes(commentsinfo, post.Comments[i].ID)
			post.Comments[i].TotalLikes, post.Comments[i].TotalDislikes = CountCommentLikesDislikes(&post.Comments[i])
		}

		// get likes of the POST
		postsinfo := "postID"
		post.Likes_Dislikes = GetLikesDislikes(postsinfo, post.ID)
		// Count the likes/dislikes
		post.TotalLikes, post.TotalDislikes = CountPostLikesDislikes(&post)

		posts = append(posts, post)
	}
	return posts
}
