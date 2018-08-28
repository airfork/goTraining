package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

var db *sql.DB
var err error

// todo struct contains the todo text and the id of the todo
type todo struct {
	Text string `json:"todo"`
	ID   string `json:"id"`
	Done bool   `json:"completed"`
}

func main() {
	// Replace password and database with your password for the mysql
	// and the database respectively
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/database")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	http.Handle("/", http.FileServer(http.Dir("views")))
	http.HandleFunc("/api/todos", index)
	http.HandleFunc("/api/todos/", api)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

// Code for post and get at /todos/api
func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Creates unique id for the todo
		u, err := uuid.NewV4()
		check(err)
		// Gets the text of the todo, and inserts it, the id,
		// and completed = false into DB
		t := r.FormValue("name")
		stmt, err := db.Prepare(`INSERT INTO todos VALUES ("` + t + `", "` + u.String() + `", "0");`)
		check(err)
		defer stmt.Close()
		// Execute prepare statement
		_, err = stmt.Exec()
		check(err)
		todo := &todo{t, u.String(), false}
		j, err := json.Marshal(todo)
		check(err)
		w.Write(j)
		// http.Redirect(w, r, "/api/todos/"+u.String(), http.StatusSeeOther)

		// Any other type of request acts as a get request
	} else {
		// Get rows that match Query, everything
		rows, err := db.Query(`SELECT * FROM todos;`)
		check(err)
		defer rows.Close()

		// Variable declaration
		var text string
		var id string
		var complete bool
		// Slice of *todo, used to print final JSON
		bs := make([]*todo, 0)

		// Iterate through found rows
		for rows.Next() {
			// Store values into variables
			err = rows.Scan(&text, &id, &complete)
			check(err)
			// Create a todo struct with the elements from db
			todo := &todo{text, id, complete}
			// add elements to slice
			bs = append(bs, todo)
		}
		// Turn slice into JSON and then write it
		// to the response bosy
		t, err := json.Marshal(bs)
		check(err)
		w.Write(t)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	// check to see if valid ID length, if not, redirect
	if len(getID(r.URL.String())) != 36 {
		http.Redirect(w, r, "/api/todos", http.StatusTemporaryRedirect)
	}

	// Displays JSON for that specific todo
	if r.Method == "GET" {
		w.Write(findByID(r))
	}

	// How to update todos
	if r.Method == "PUT" {
		data := r.FormValue("completed")
		text := r.FormValue("name")
		u := getID(r.URL.String())

		if text == "" {
			stmt, err := db.Prepare(`UPDATE todos SET completed="` + data + `" WHERE id="` + u + `";`)
			check(err)
			defer stmt.Close()

			_, err = stmt.Exec()
			check(err)
			w.Write(findByID(r))
			// This is actually unecessary
			// From the actual page, the user can only change
			// completed from false to true and vice-versa
		} else {
			stmt, err := db.Prepare(`UPDATE todos SET todo="` + text + `" WHERE id="` + u + `";`)
			check(err)
			defer stmt.Close()

			_, err = stmt.Exec()
			check(err)
			w.Write(findByID(r))
		}
	}

	// Deletes todo with that specific ID
	if r.Method == "DELETE" {
		// Get id from request URL
		u := getID(r.URL.String())
		// Delete from table
		stmt, err := db.Prepare(`DELETE FROM todos WHERE id="` + u + `";`)
		check(err)
		defer stmt.Close()

		_, err = stmt.Exec()
		check(err)
		http.Redirect(w, r, "/api/todos", http.StatusTemporaryRedirect)
	}
}

// Some basic error checking
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Takes URL and extracts the id
func getID(u string) string {
	i := strings.Split(u, "")
	i = i[11:]
	return strings.Join(i, "")
}

// Takes the id from the request and finds that element
// Returns the element as a byte slice
func findByID(r *http.Request) []byte {
	// Declarations
	var text string
	var id string
	var complete bool
	var s []byte

	// Get id from url
	u := getID(r.URL.String())
	// Find element with that id
	rows, err := db.Query(`SELECT * FROM todos WHERE id="` + u + `";`)
	check(err)
	defer rows.Close()

	// Iterate through found elements(only one)
	// Save to struct and Marshal that struct into JSON
	for rows.Next() {
		err = rows.Scan(&text, &id, &complete)
		check(err)
		todo := &todo{text, id, complete}
		s, err = json.Marshal(todo)
		check(err)
	}
	return s
}

// Run first time to create table with todos and IDS
// Call it inside func main with http.HandleFunc("/route", create)

// func create(w http.ResponseWriter, r *http.Request) {
//
// 	stmt, err := db.Prepare(`CREATE TABLE todos (todo VARCHAR(160), id VARCHAR(45), completed TINYINT);`)
// 	check(err)
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec()
// 	check(err)
// }
