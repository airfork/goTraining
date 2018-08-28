package main

import (
	"html/template"
	"net/http"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName string
	// You should never do this, as you should not store password on server
	// Password string

	// This allows password to be encrypted, much more secure
	Password []byte
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/info", info)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/info", http.StatusSeeOther)
		return
	}
	login(w, r)
}

func info(w http.ResponseWriter, r *http.Request) {
	u, err := getUser(w, r)
	if err != nil {
		http.Error(w, "Please sign in again", http.StatusForbidden)
		return
	}
	tpl.ExecuteTemplate(w, "userinfo.gohtml", u)

}

func getUser(w http.ResponseWriter, r *http.Request) (user, error) {
	var u user
	c, err := r.Cookie("snickerdoodle")
	if err != nil {
		return u, err
	}

	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}
	return u, nil
}

func isLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("snickerdoodle")
	if err != nil {
		return false
	}

	u := dbSessions[c.Value]
	_, ok := dbUsers[u]
	return ok
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// get form values
		un := r.FormValue("username")
		p := r.FormValue("password")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")

		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "snickerdoodle",
			Value: sID.String(),
		}
		http.SetCookie(w, c)

		dbSessions[c.Value] = un

		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// store user in dbUsers
		u := user{un, bs, f, l}
		dbUsers[un] = u

		// redirect
		http.Redirect(w, r, "/info", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}
