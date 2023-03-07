package utils

import (
	"fmt"
	database "forum/data"
	"net/http"
)

type Sessions struct {
	Username string
}

type Categories struct {
	ID          int
	Title       string
	Description string
}

type Posts struct {
	ID             int
	Title          string
	Body           string
	CreatedAt      string
	Comments       []Comments
	UserID         int
	Media          string
	Username       string
	Likes_Dislikes []Likes_Dislikes
	TotalLikes     int
	TotalDislikes  int
	Flag           bool
}

type Users struct {
	ID       int
	Username string
	Email    string
	Posts    []Posts
	Comments []Comments
	Likes    []Likes_Dislikes
	DisLikes []Likes_Dislikes
	Status   string
}

func SessionHandler(r *http.Request) []Sessions {
	var sessions []Sessions
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return sessions
	}
	rows, err := database.Query("SELECT username FROM sessionIDs WHERE sessionID =?", cookie.Value)
	if err != nil {
		fmt.Print(err)
	}
	for rows.Next() {
		var session Sessions
		if err := rows.Scan(&session.Username); err != nil {
			fmt.Println("failed to read")
		}
		sessions = append(sessions, session)
	}
	return sessions
}

func CategoriesHandler() []Categories {
	rows, err := database.Query("SELECT ID, title, description FROM categories")
	if err != nil {
		fmt.Print(err)
	}
	var categories []Categories

	for rows.Next() {
		var category Categories
		if err := rows.Scan(&category.ID, &category.Title, &category.Description); err != nil {
			fmt.Println("failed to read")
		}
		categories = append(categories, category)
	}
	return categories
}

func createCategorie(title, body string) {
	_, err := database.Exec("INSERT INTO categories (title, description) VALUES (?, ?)", title, body)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteCategorie(catID int) {
	_, err := database.Exec("DELETE FROM categories WHERE ID = ?", catID)
	if err != nil {
		fmt.Println(err)
	}
}

func GetContent(categoryID int) []Posts {
	rows, err := database.Query(`SELECT posts.ID, posts.title, posts.body, posts.created, posts.userID, posts.media, posts.flag
                         FROM postcategories
                         JOIN posts ON postcategories.PostID = posts.ID
                         WHERE postcategories.CategoryID = ?`, categoryID)
	if err != nil {
		fmt.Print(err)
	}

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

func getUsers() []Users {
	rows, err := database.Query("SELECT ID, username, email, status FROM customer")
	if err != nil {
		fmt.Print(err)
	}
	var users []Users

	for rows.Next() {
		var user Users
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Status); err != nil {
			fmt.Println("failed to read")
			fmt.Println(err)
		}
		users = append(users, user)
	}
	return users
}

func GetProfile(usrID string, userIDint int) []Users {
	var users []Users
	var user Users

	err := database.QueryRow("SELECT ID, username, email FROM customer where ID = ?", usrID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rows, err := database.Query("SELECT ID, title, body, created, media FROM posts WHERE userID = ?", userIDint)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.Media); err != nil {
			fmt.Println(err)
			return nil
		}
		user.Posts = append(user.Posts, post)
	}

	commentsRows, err := database.Query("SELECT body, creation, postID FROM comments WHERE userID = ?", usrID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer commentsRows.Close()

	for commentsRows.Next() {
		var comment Comments
		if err := commentsRows.Scan(&comment.Body, &comment.CreatedAt, &comment.PostID); err != nil {
			fmt.Println(err)
			return nil
		}
		var postTitle string
		var postBody string
		var postMedia string
		err = database.QueryRow("SELECT title, body, media FROM posts WHERE ID = ?", comment.PostID).Scan(&postTitle, &postBody, &postMedia)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		comment.PostTitle = postTitle
		comment.PostBody = postBody
		comment.PostMedia = postMedia
		user.Comments = append(user.Comments, comment)
	}

	likedPostRows, err := database.Query("SELECT postID FROM likes_dislikes WHERE userID = ? and type = 'like' and item_type = 'posts'", usrID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer likedPostRows.Close()
	for likedPostRows.Next() {
		var likedPost Likes_Dislikes
		if err := likedPostRows.Scan(&likedPost.PostID); err != nil {
			fmt.Println(err)
			return nil
		}

		var title string
		err := database.QueryRow("SELECT title FROM posts WHERE ID = ?", likedPost.PostID).Scan(&title)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		likedPost.Title = title
		user.Likes = append(user.Likes, likedPost)
	}

	DislikedPostRows, err := database.Query("SELECT postID FROM likes_dislikes WHERE userID = ? and type = 'dislike' and item_type = 'posts'", usrID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer DislikedPostRows.Close()
	for DislikedPostRows.Next() {
		var likedPost Likes_Dislikes
		if err := DislikedPostRows.Scan(&likedPost.PostID); err != nil {
			fmt.Println(err)
			return nil
		}
		var title string
		err := database.QueryRow("SELECT title FROM posts WHERE ID = ?", likedPost.PostID).Scan(&title)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		likedPost.Title = title
		user.DisLikes = append(user.DisLikes, likedPost)
	}
	users = append(users, user)
	return users
}
