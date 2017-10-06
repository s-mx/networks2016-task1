package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"html/template"

	"github.com/gorilla/mux"

	"networks2016-task1/session"
	auth "networks2016-task1/users"
	"networks2016-task1/api"
	"unicode"
)

var users	*auth.Users

var templates	map[string]*template.Template
var port *int
var postgresConfigPath *string

func initTemplates() {
	templates = make(map[string]*template.Template)
	temp := template.Must(template.ParseFiles("templates/base.html", "templates/default.html"))
	templates["default"] = temp
	temp = template.Must(template.ParseFiles("templates/base.html", "templates/logged.html"))
	templates["logged"] = temp
	/*temp = template.Must(template.ParseFiles("templates/base.html", "templates/fact.html"))
	templates["fact"] = temp*/
}

func init() {
	port = flag.Int("port", 12000, "server port")
	postgresConfigPath = flag.String("postgres-config",
									 "configs/postgres_test.json",
		                             "path to postgres config")
	flag.Parse()

	var err error
	session.Manager, err = session.NewSessionManager("app1")
	if err != nil {
		log.Fatal(err)
	}

	users, err = auth.NewUsers()
	if err != nil {
		log.Fatal(err)
	}

	initTemplates()
	auth.InitAuth(*postgresConfigPath)
}

func indexPage(responseWriter http.ResponseWriter, request *http.Request) {
	if session.Manager.IsLogged(responseWriter, request) {
		curSession, _ := session.Manager.GetSession(request)
		templates["logged"].ExecuteTemplate(responseWriter, "base", curSession.User)
	} else {
		templates["default"].ExecuteTemplate(responseWriter, "base", nil)
	}
}

func ValidateString(str string) bool {
	if len(str) == 0 {
		return false
	}

	for _, elem := range str {
		if !unicode.IsDigit(elem) && !unicode.IsLetter(elem) {
			return false
		}
	}

	return true
}

func ValidateData(login, password string) (bool) {
	return ValidateString(login) && ValidateString(password)
}

func loginPage(responseWriter http.ResponseWriter, request *http.Request) {
	if session.Manager.IsLogged(responseWriter, request) {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	if request.Method == "GET" {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	request.ParseForm()
	form := request.Form
	if len(form["login"]) == 0 || len(form["password"]) == 0 {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	login := form["login"][0]
	password := form["password"][0]

	if ValidateData(login, password) == false {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	existance, err := users.Db.CheckExistance(login, password)
	if err != nil {
		log.Fatal(err)
	}

	if existance == false {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	user, _ := users.Db.GetUser(login, password)
	session.Manager.StartSession(responseWriter, request, user)
}

func logout(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "Get" {
		http.Redirect(responseWriter, request, "/", 302)
		return
	} else {
		if session.Manager.IsLogged(responseWriter, request) {
			session.Manager.Logout(responseWriter, request)
		}

		http.Redirect(responseWriter, request, "/", 302)
		return
	}
}

func registerAction(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		http.Redirect(responseWriter, request, "/", 302)
		return
	} else {
		request.ParseForm()

		login := request.Form["login"][0]
		password := request.Form["password"][0]
		ValidateData(login, password)

		exist, _ := users.Db.CheckExistance(login, password)
		if exist {
			http.Redirect(responseWriter, request, "/", 302)
			return
		}

		user, _ := users.Db.AddUser(login, password)
		session.Manager.StartSession(responseWriter, request, user)
	}
}

func factPage(responseWriter http.ResponseWriter, request *http.Request) {
	if session.Manager.IsLogged(responseWriter, request) == false {
		http.Redirect(responseWriter, request, "/", 302)
		return
	}

	curSession, _ := session.Manager.GetSession(request)
	user := curSession.User
	temp := template.Must(template.ParseFiles("templates/base.html", "templates/fact.html"))
	temp.ExecuteTemplate(responseWriter, "base", user)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", loginPage)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/register", registerAction)
	router.HandleFunc("/fact", factPage)
	router.HandleFunc("/api/fact/{number}", api.Factorize)

	log.Println("Listening...")
	http.ListenAndServe(":" + strconv.Itoa(*port), router)
}
