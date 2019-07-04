package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	// "github.com/gorilla/sessions"
)

var client *redis.Client
var store = sessions.NewCookieStore([]byte("victor-steven"))

var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //redis port
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/test", testGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// r.HandleFunc("/goodbye", goodbyeHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)

	//without mux
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":7000", nil)
}

func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	//get the first ten strings in redis from a string called "comments":
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	// fmt.Fprint(res, "Hello World!")
	templates.ExecuteTemplate(res, "index.html", comments)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
	//parse the form from the request body
	req.ParseForm()

	comment := req.PostForm.Get("comment")

	log.Println(comment)

	client.LPush("comments", comment)
	http.Redirect(res, req, "/", 302)

	// fmt.Fprint(res, "Hello World!")
	// templates.ExecuteTemplate(res, "index.html", comments)
}

func loginGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "login.html", nil)
}

func loginPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	session, _ := store.Get(req, "session")
	session.Values["username"] = username
	session.Save(req, res)
}

func testGetHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	untyped, ok := session.Values["username"]
	if !ok {
		return
	}
	username, ok := untyped.(string)
	if !ok {
		return
	}

	//write to the byte array of the username
	res.Write([]byte(username))
}

// func goodbyeHandler(res http.ResponseWriter, req *http.Request) {
// 	fmt.Fprint(res, "Goodbye World!")
// }
