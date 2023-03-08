package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ipCounter = make(map[string]int)

func RunServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", rateLimiter(HomeHandler))

	//login related
	mux.HandleFunc("/login", rateLimiter(LoginPageHandler))
	//Google Log-in
	mux.HandleFunc("/glogin", rateLimiter(gLogIn))
	mux.HandleFunc("/callback", rateLimiter(gCallBack))
	//Github Log-in
	mux.HandleFunc("/ghlogin", rateLimiter(ghLogIn))
	mux.HandleFunc("/ghcallback", rateLimiter(ghCallBack))

	//moderation related
	mux.HandleFunc("/users", rateLimiter(usersHandler))
	mux.HandleFunc("/mod", rateLimiter(modPageHandler))
	mux.HandleFunc("/update-user-status", rateLimiter(userStatusHandler))
	mux.HandleFunc("/add-flag", rateLimiter(AddFlagHandler))
	mux.HandleFunc("/remove-flag", rateLimiter(RemoveFlagHandler))

	mux.HandleFunc("/logout", rateLimiter(LogOutHandler))
	mux.HandleFunc("/register", rateLimiter(RegisterPageHandler))
	mux.HandleFunc("/registerdata", rateLimiter(RegisterDataHandler))
	mux.HandleFunc("/log-in-data", rateLimiter(LoginDataHandler))
	mux.HandleFunc("/profile", rateLimiter(ProfileHandler))
	mux.HandleFunc("/manage-cats", rateLimiter(manageCategoriesHandler))
	mux.HandleFunc("/delete-categorie", rateLimiter(deleteCategorieHandler))
	mux.HandleFunc("/create-categorie", rateLimiter(createCategorieHandler))

	//interaction related
	mux.HandleFunc("/like", rateLimiter(ReceiveLikesHandler))
	mux.HandleFunc("/commentfield", rateLimiter(SentCommentHandler))
	mux.HandleFunc("/create-data", rateLimiter(NewPostHandler))
	mux.HandleFunc("/create-post", rateLimiter(CreatePostPageHandler))
	mux.HandleFunc("/create-ticket", rateLimiter(NewTicketPageHandler))
	mux.HandleFunc("/ticket-data", rateLimiter(CreateTicketHandler))
	mux.HandleFunc("/delete-comment", rateLimiter(DeleteCommentHandler))
	mux.HandleFunc("/modify-comment", rateLimiter(ModifyCommentHandler))
	mux.HandleFunc("/delete-post", rateLimiter(DeletePostHandler))
	mux.HandleFunc("/modify-post", rateLimiter(ModifyPostHandler))

	mux.HandleFunc("/api/notifications", rateLimiter(NotificationHandler))
	mux.HandleFunc("/delete-notifications", rateLimiter(DeleteNotificationHandler))

	fileServer := http.FileServer(http.Dir("templates/UI/"))
	mux.Handle("/UI/style.css", http.StripPrefix("/UI/", fileServer))
	//serving media
	fs := http.FileServer(http.Dir("templates/media/"))
	mux.Handle("/media/", http.StripPrefix("/media/", fs))
	fserv := http.FileServer(http.Dir("templates/media/"))
	mux.Handle("/content/media/", http.StripPrefix("/content/media/", fserv))
	//serving js script
	jsserv := http.FileServer(http.Dir("./templates"))
	mux.Handle("/content/js.js", http.StripPrefix("/content/", jsserv))
	mux.Handle("/js.js", http.FileServer(http.Dir("./templates")))

	mux.HandleFunc("/content/", rateLimiter(contentHandler))

	go func() {
		if err := http.ListenAndServe(":8081", http.HandlerFunc(redirect)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Set timeouts
	server := &http.Server{
		Addr:         ":8443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	cert := "utils/cert/server.crt"
	key := "utils/cert/server.key"

	fmt.Println("Server started on port 8443")
	log.Fatal(server.ListenAndServeTLS(cert, key))
}

// Redirect all http -> https
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

// middleware checks how many times the user with this IP has made requests
func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if count, ok := ipCounter[ip]; ok {
			if count >= 100 {
				http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
				return
			}
			ipCounter[ip] = count + 1
		} else {
			ipCounter[ip] = 1
		}

		// start a goroutine to reset the count for the IP address after 1 minute
		go func() {
			time.Sleep(time.Minute)
			delete(ipCounter, ip)
		}()
		next.ServeHTTP(w, r)
	})
}
