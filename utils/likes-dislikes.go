package utils

import (
	"database/sql"
	"fmt"
	database "forum/data"
	"strconv"
)

type Likes_Dislikes struct {
	ID        int
	PostID    sql.NullInt64
	CommentID sql.NullInt64
	UserID    int
	Type      string
	Title     string
	Likes     string
	DisLikes  string
}

func AddLikesDislikes(itemID, userID, col, action, item_type, where string) {
	voted, err := HasUserVotedPost(where, itemID, userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !voted {
		_, err = database.Exec("INSERT INTO likes_dislikes ("+where+", userID, type, item_type) VALUES (?, ?, ?, ?) ", itemID, userID, action, col)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// update the current count
		_, err = database.Exec("UPDATE likes_dislikes SET type = ? WHERE "+where+" = ? and userID = ? and item_type = ? ", action, itemID, userID, col)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//adds a notification line to the database
	var originalAuthor int
	user, _ := strconv.Atoi(userID)
	commentAuhtor := GetUsernameByID(user)
	var title string
	if item_type == "post" {
		err = database.QueryRow("SELECT userID, title FROM posts WHERE ID =?", itemID).Scan(&originalAuthor, &title)
		if err != nil {
			fmt.Println(err)
		}
	}
	if item_type == "comment" {
		err = database.QueryRow("SELECT userID, body FROM comments WHERE ID =?", itemID).Scan(&originalAuthor, &title)
		if err != nil {
			fmt.Println(err)
		}
	}
	text := commentAuhtor + " has " + action + "d your " + item_type + ": " + title
	if originalAuthor != user {
		_, err = database.Exec("INSERT INTO notifications (userID, action) VALUES (?, ?)", originalAuthor, text)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func HasUserVotedPost(where, postID, userID string) (bool, error) {
	var count int
	err := database.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE "+where+" = ? and userID = ? ", postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetLikesDislikes(where string, postID int) []Likes_Dislikes {
	rows, err := database.Query("SELECT ID, postID, commentID, userID, type FROM likes_dislikes WHERE "+where+" = ? ", postID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var likes_dislikes []Likes_Dislikes
	for rows.Next() {
		var like_dislike Likes_Dislikes
		if err := rows.Scan(&like_dislike.ID, &like_dislike.PostID, &like_dislike.CommentID, &like_dislike.UserID, &like_dislike.Type); err != nil {
			fmt.Println(err)
			fmt.Println("failed to read likes_dislikes")
		}
		likes_dislikes = append(likes_dislikes, like_dislike)
	}
	return likes_dislikes
}

func CountPostLikesDislikes(post *Posts) (int, int) {
	var likes, dislikes int

	for _, v := range post.Likes_Dislikes {
		if v.Type == "like" {
			likes++
		} else if v.Type == "dislike" {
			dislikes++
		}
	}
	return likes, dislikes
}

func CountCommentLikesDislikes(comment *Comments) (int, int) {
	var likes, dislikes int

	for _, v := range comment.Likes_Dislikes {
		if v.Type == "like" {
			likes++
		} else if v.Type == "dislike" {
			dislikes++
		}
	}
	return likes, dislikes
}
