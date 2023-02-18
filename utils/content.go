package utils

import (
	"database/sql"
	"fmt"
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

func SessionHandler(r *http.Request) []Sessions {
	var sessions []Sessions
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return sessions
	}
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT username FROM sessionIDs WHERE sessionID =?", cookie.Value)
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
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT ID, title, description FROM categories")
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
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO categories (title, description) VALUES (?, ?)", title, body)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteCategorie(catID int) {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM categories WHERE ID = ?", catID)
	if err != nil {
		fmt.Println(err)
	}
}

func GetContent(categoryID int) []Posts {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query(`SELECT posts.ID, posts.title, posts.body, posts.created, posts.userID, posts.media, posts.flag
                         FROM postcategories
                         JOIN posts ON postcategories.PostID = posts.ID
                         WHERE postcategories.CategoryID = ?`, categoryID)
	if err != nil {
		fmt.Print(err)
	}

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

func getUsers() []Users {
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT ID, username, email, status FROM customer")
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
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}

	var users []Users
	var user Users

	// get the userdata
	err = db.QueryRow("SELECT ID, username, email FROM customer where ID = ?", usrID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// get posts
	rows, err := db.Query("SELECT ID, title, body, created, media FROM posts WHERE userID = ?", userIDint)
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

	// get the user's comments
	commentsRows, err := db.Query("SELECT body, creation, postID FROM comments WHERE userID = ?", usrID)
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
		err = db.QueryRow("SELECT title, body, media FROM posts WHERE ID = ?", comment.PostID).Scan(&postTitle, &postBody, &postMedia)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		comment.PostTitle = postTitle
		comment.PostBody = postBody
		comment.PostMedia = postMedia
		user.Comments = append(user.Comments, comment)
	}

	// get user's liked posts
	likedPostRows, err := db.Query("SELECT postID FROM likes_dislikes WHERE userID = ? and type = 'like' and item_type = 'posts'", usrID)
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
		// get liked post title
		var title string
		err := db.QueryRow("SELECT title FROM posts WHERE ID = ?", likedPost.PostID).Scan(&title)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		likedPost.Title = title
		user.Likes = append(user.Likes, likedPost)
	}

	// get user's Disliked posts
	DislikedPostRows, err := db.Query("SELECT postID FROM likes_dislikes WHERE userID = ? and type = 'dislike' and item_type = 'posts'", usrID)
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
		// get liked post title
		var title string
		err := db.QueryRow("SELECT title FROM posts WHERE ID = ?", likedPost.PostID).Scan(&title)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		likedPost.Title = title
		user.DisLikes = append(user.DisLikes, likedPost)
	}

	// if len(user.Likes) > 0 {
	// 	fmt.Println("liked posts found")
	// } else {
	// 	fmt.Println("No liked found")
	// }

	users = append(users, user)
	return users
}
