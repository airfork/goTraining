package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// Page struct defines a page as consiting of a Title and a Body
type Page struct {
	Title string
	// Needs to be a byte slice for methods used to work
	Body []byte
}

// Parses templates into a single *Template
// Can then use ExecuteTemplate to render a specific template
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// Returns a regexp
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// Renders a specific template from the template variable
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		// Send an error to the ResponseWriter, the code of which is
		// http.StatusInternalServerError, aka code 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// Gets body data from the textarea tag with name of body
	// Returns this data as a string
	body := r.FormValue("body")
	// Creates Page struct with data
	// []byte(body) converts the string body data into a byte slice
	p := &Page{Title: title, Body: []byte(body)}
	// Calls save function that all Page types have
	err := p.save()
	if err != nil {
		// Send an error to the ResponseWriter, the code of which is
		// http.StatusInternalServerError, aka code 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirects to the view page for that title
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the page title from the request
		// Call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
