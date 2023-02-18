package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Categories":       CategoriesHandler(),
		"Session":          SessionHandler(r),
		"hasNotifications": false,
		"currentUserID":    0,
		"userType":         "user",
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/" {
		tmpl := template.Must(template.ParseFiles("templates/home.html"))
		tmpl.Execute(w, data)
	} else {
		notFoundHandler(w, r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	if isLoggedIn {
		userIDint := userIDFromCookie(r)
		userID := strconv.Itoa(userIDint)

		data := map[string]interface{}{
			"IsLoggedIn":       isLoggedIn,
			"User":             GetProfile(userID, userIDint),
			"hasNotifications": false,
			"currentUserID":    userIDFromCookie(r),
			"refreshPageURL":   r.URL.Path,
		}
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		if r.URL.Path == "/profile" {
			tmpl := template.Must(template.ParseFiles("templates/profile.html"))
			tmpl.Execute(w, data)
		}
	} else {
		data := map[string]interface{}{
			"IsLoggedIn": isLoggedIn,
		}
		if r.URL.Path == "/profile" {
			tmpl := template.Must(template.ParseFiles("templates/profile.html"))
			tmpl.Execute(w, data)
		}
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Categories":       CategoriesHandler(),
		"Session":          SessionHandler(r),
		"hasNotifications": false,
		"currentUserID":    0,
		"Users":            getUsers(),
		"userType":         "user",
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/users" {
		tmpl := template.Must(template.ParseFiles("templates/users.html"))
		tmpl.Execute(w, data)
	} else {
		notFoundHandler(w, r)
	}
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" {
		data := false
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, data)
	}
}

func RegisterDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		email := r.PostFormValue("email")
		username := r.PostFormValue("username")
		hashedpassword, _ := HashPassword(r.PostFormValue("password"))

		db, err := sql.Open("sqlite3", "./data/database.db")
		if err != nil {
			http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// send the users info to db
		status := "user"
		_, err = db.Exec("INSERT INTO customer (email, username, password, status) VALUES (?, ?, ?, ?)", email, username, hashedpassword, status)

		// database gives error, return
		if err != nil {
			fmt.Println("User Failed To Register")
			fmt.Println(err)
			fmt.Println("---------")
			if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
				data := true
				tmpl := template.Must(template.ParseFiles("templates/register.html"))
				tmpl.Execute(w, data)
			} else {
				http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			// else store user and redirect to login page
			//Terminal Output
			fmt.Println("---------------------------------")
			fmt.Println("User Registration Successful")
			fmt.Println("Email:", email)
			fmt.Println("Usernme:", username)
			fmt.Println("Hashed Password:", hashedpassword)
			fmt.Println("---------------------------------")
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Session":          SessionHandler(r),
		"hasNotifications": false,
		"currentUserID":    0,
		"userType":         "user",
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/create-post" {
		tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
		tmpl.Execute(w, data)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		data := false
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, data)
	}
}

func LoginDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		db, err := sql.Open("sqlite3", "file:data/database.db")
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
		if Login(db, username, password) {
			setSessionCookie(w, username)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			data := true
			tmpl := template.Must(template.ParseFiles("templates/login.html"))
			tmpl.Execute(w, data)
		}
		//Terminal Output
		fmt.Printf("User %v - Login Successful\n", username)
		fmt.Println("---------------------------------")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session ID/username from the request's cookies
	sessionID, _ := r.Cookie("sessionID")
	// Delete the user's session from the sessions table in the database
	db, err := sql.Open("sqlite3", "file:data/database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessionIDs WHERE sessionID = ?", sessionID.Value)
	if err != nil {
		fmt.Println(err)
	}

	// Redirect the user to the home page or a login page
	http.Redirect(w, r, "/", http.StatusFound)
	fmt.Println("User Logged Out")
}

func SentCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}
		comment := r.PostFormValue("comment")
		user := userIDFromCookie(r)
		postID := r.PostFormValue("postID")
		//takes the page address from html in order to return the user back to where they posted comment
		refreshpage := r.PostFormValue("url")
		AddComment(comment, postID, user, w)
		http.Redirect(w, r, refreshpage, http.StatusSeeOther)
		//Terminal Output
		fmt.Println("Comment:", comment)
		fmt.Println("---------------------------------")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ReceiveLikesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		itemID := r.FormValue("item_id")
		userIDint := userIDFromCookie(r)
		userID := strconv.Itoa(userIDint)
		action := r.FormValue("action")
		item_type := r.FormValue("item_type")

		db, err := sql.Open("sqlite3", "file:data/database.db")
		if err != nil {
			http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		//check if like/dislike went on post or comment
		var col string
		if item_type == "post" {
			col = "posts"
			where := "postID"
			AddLikesDislikes(itemID, userID, col, action, item_type, where)
		} else {
			col = "comments"
			where := "commentID"
			AddLikesDislikes(itemID, userID, col, action, item_type, where)
		}
		refreshpage := r.PostFormValue("url")
		http.Redirect(w, r, refreshpage, http.StatusSeeOther)
		//Terminal Output
		fmt.Printf("%v %vd \n", strings.Title(item_type), strings.Title(action))
		fmt.Println("---------------------------------")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
		}
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		userID := userIDFromCookie(r)
		categories := r.Form["topic[]"]
		file, header, err := r.FormFile("fileInput")
		fileName := "false"
		if err == nil {
			fileName = header.Filename
			AddMedia(fileName, file)
		}
		postID, err := AddPost(title, body, fileName, userID, w)
		if err != nil {
			fmt.Println(err)
		}
		AddCategories(postID, categories, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		//Terminal Output
		fmt.Println("Post Creation Successful")
		fmt.Println("Title:", title)
		fmt.Println("Body:", body)
		fmt.Println("---------------------------------")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		url := r.PostFormValue("url")
		commentID, _ := strconv.Atoi(r.PostFormValue("item_id"))
		user := userIDFromCookie(r)
		DeleteComment(commentID, user, w)
		http.Redirect(w, r, url, http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		url := r.PostFormValue("url")
		postID, _ := strconv.Atoi(r.PostFormValue("item_id"))
		user := userIDFromCookie(r)
		DeletePost(postID, user, w)
		http.Redirect(w, r, url, http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ModifyCommentHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	commentID, _ := strconv.Atoi(r.PostFormValue("id"))
	user := userIDFromCookie(r)
	newcomment := r.PostFormValue("textarea")
	url := r.PostFormValue("url")
	UpdateComment(commentID, user, newcomment, w)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func ModifyPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	postID, _ := strconv.Atoi(r.PostFormValue("id"))
	user := userIDFromCookie(r)
	url := r.PostFormValue("url")
	newbody := r.PostFormValue("postarea")
	newtitle := r.PostFormValue("titlearea")
	UpdatePost(postID, user, newbody, newtitle, w)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCookie(r)
	notifications := GetNotifications(userID)
	json.NewEncoder(w).Encode(notifications)
}

func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCookie(r)
	DeleteNotifications(userID)
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	url := r.PostFormValue("url")
	http.Redirect(w, r, url, http.StatusSeeOther)

}

func modPageHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Posts":            GetContent(5),
		"Session":          SessionHandler(r),
		"currentUserID":    0,
		"hasNotifications": false,
		"userType":         "user",
		"FlaggedPosts":     GetFlaggedPosts(),
		"FlaggedComments":  GetFlaggedComments(),
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/mod" {
		tmpl := template.Must(template.ParseFiles("templates/mod.html"))
		tmpl.Execute(w, data)
	}
}

func userStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(r.PostFormValue("id"))
		status := r.PostFormValue("status")
		db, err := sql.Open("sqlite3", "file:data/database.db")
		if err != nil {
			http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()
		fmt.Println(id, status)
		_, err = db.Exec("UPDATE customer SET status = ? WHERE ID = ?", status, id)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/content/"):]
	id_int, _ := strconv.Atoi(id)
	isLoggedIn := checkSessionCookie(r)
	title := GetCategoryName(id_int)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Posts":            GetContent(id_int),
		"Session":          SessionHandler(r),
		"currentUserID":    0,
		"hasNotifications": false,
		"userType":         "user",
		"pageTitle":        title,
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	tmpl := template.Must(template.ParseFiles("templates/content.html"))
	tmpl.Execute(w, data)
}

func CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
		}
		title := "Ticket"
		body := r.PostFormValue("body")
		userID := userIDFromCookie(r)
		var categories []string
		categories = append(categories, "5")
		file, header, err := r.FormFile("fileInput")
		fileName := "false"
		if err == nil {
			fileName = header.Filename
			AddMedia(fileName, file)
		}
		postID, err := AddPost(title, body, fileName, userID, w)
		if err != nil {
			fmt.Println(err)
		}
		AddCategories(postID, categories, w)
		http.Redirect(w, r, "/mod", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func NewTicketPageHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Session":          SessionHandler(r),
		"hasNotifications": false,
		"currentUserID":    0,
		"userType":         "user",
		"refreshPageURL":   r.URL.Path,
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/create-ticket" {
		tmpl := template.Must(template.ParseFiles("templates/create-ticket.html"))
		tmpl.Execute(w, data)
	}
}

func AddFlagHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(requestData.Id)
	if requestData.Type == "post" {
		AddPostFlag(id)
	}
	if requestData.Type == "comment" {
		AddCommentFlag(id)
	}
}

func RemoveFlagHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(requestData.Id)
	if requestData.Type == "post" {
		RemovePostFlag(id)
	}
	if requestData.Type == "comment" {
		RemoveCommentFlag(id)
	}
}

func manageCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := checkSessionCookie(r)
	data := map[string]interface{}{
		"IsLoggedIn":       isLoggedIn,
		"Categories":       CategoriesHandler(),
		"Session":          SessionHandler(r),
		"hasNotifications": false,
		"currentUserID":    0,
		"userType":         "user",
	}
	if isLoggedIn {
		if len(GetNotifications(userIDFromCookie(r))) > 0 {
			data["hasNotifications"] = true
		}
		data["currentUserID"] = userIDFromCookie(r)
		data["userType"] = GetStatusByID(userIDFromCookie(r))
	}
	if r.URL.Path == "/manage-cats" {
		tmpl := template.Must(template.ParseFiles("templates/manage-cats.html"))
		tmpl.Execute(w, data)
	} else {
		notFoundHandler(w, r)
	}
}

func deleteCategorieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		catID, _ := strconv.Atoi(r.PostFormValue("item_id"))
		deleteCategorie(catID)
		http.Redirect(w, r, "https://localhost/manage-cats", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func createCategorieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}

		title := r.FormValue("title")
		body := r.FormValue("body")

		createCategorie(title, body)
		fmt.Println("Categorie Creation Successful")
		fmt.Println("Title:", title)
		fmt.Println("Body:", body)
		fmt.Println("---------------------------------")
		http.Redirect(w, r, "https://localhost/manage-cats", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
