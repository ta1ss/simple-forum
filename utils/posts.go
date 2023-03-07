package utils

import (
	"fmt"
	database "forum/data"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func AddPost(title, body, fileName string, userID int, w http.ResponseWriter) (int, error) {
	flag := false
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	result, err := database.Exec("INSERT INTO posts (title, body, created, userID, media, flag) VALUES (?, ?, ?, ?, ?, ?)", title, body, currentTime, userID, fileName, flag)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(postID), nil
}
func AddCategories(postID int, categories []string, w http.ResponseWriter) {
	stmt, err := database.Prepare("INSERT INTO postcategories (postID, categoryID) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
	}
	for _, category := range categories {
		categoryID, err := strconv.Atoi(category)
		if err != nil {
			http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		}
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
	var createdby int
	err := database.QueryRow("SELECT userID FROM posts WHERE ID = ?", postID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}
	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		_, err := database.Exec("DELETE FROM comments WHERE postID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = database.Exec("DELETE FROM likes_dislikes WHERE postID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = database.Exec("DELETE FROM postcategories WHERE PostID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = database.Exec("DELETE FROM posts WHERE ID = ?", postID)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdatePost(postID, user int, newbody, newtitle string, w http.ResponseWriter) {
	var createdby int
	err := database.QueryRow("SELECT userID FROM posts WHERE ID = ?", postID).Scan(&createdby)
	if err != nil {
		fmt.Println(err)
	}
	userStatus := GetStatusByID(user)
	if createdby == user || userStatus == "mod" || userStatus == "admin" {
		if len(newbody) > 0 {
			_, err := database.Exec("UPDATE posts SET body = ? WHERE ID = ?", newbody, postID)
			if err != nil {
				fmt.Println(err)
			}
		}
		if len(newtitle) > 0 {
			_, err = database.Exec("UPDATE posts SET title = ? WHERE ID = ?", newtitle, postID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func AddPostFlag(id int) {
	_, err := database.Exec("UPDATE posts SET flag = ? WHERE ID = ?", true, id)
	if err != nil {
		fmt.Println(err)
	}
}

func RemovePostFlag(id int) {
	_, err := database.Exec("UPDATE posts SET flag = ? WHERE ID = ?", false, id)
	if err != nil {
		fmt.Println(err)
	}
}

func GetFlaggedPosts() []Posts {
	rows, err := database.Query("SELECT ID, title, body, created, userID, media, flag FROM posts WHERE flag = ?", true)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UserID, &post.Media, &post.Flag); err != nil {
			fmt.Println("failed to read POSTS")
		}
		post.Username = GetUsernameByID(post.UserID)
		post.Comments = getComments(post.ID)
		commentsinfo := "commentID"

		for i := range post.Comments {
			post.Comments[i].Likes_Dislikes = GetLikesDislikes(commentsinfo, post.Comments[i].ID)
			post.Comments[i].TotalLikes, post.Comments[i].TotalDislikes = CountCommentLikesDislikes(&post.Comments[i])
		}
		postsinfo := "postID"
		post.Likes_Dislikes = GetLikesDislikes(postsinfo, post.ID)
		post.TotalLikes, post.TotalDislikes = CountPostLikesDislikes(&post)
		posts = append(posts, post)
	}
	return posts
}
