package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
)

type Content struct {
	mu      sync.RWMutex
	content string
}

func newContent(text string) *Content {
	return &Content{content: text}
}

func (c *Content) get() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.content
}

func (c *Content) set(text string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.content = text
}

func main() {
	// setup template
	tmpl := template.Must(template.ParseFiles("static/index.html"))

	// serve static files from the /static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// create instance of content
	content := newContent("Paste text here...")

	http.HandleFunc("/share", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			serveWebPage(tmpl, w, content)
		case "POST":
			saveContent(w, r, content)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	port := 8080 // TODO: make a CLI arg
	log.Printf("Server started and listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func serveWebPage(tmpl *template.Template, w http.ResponseWriter,
	content *Content) {
	// get the current text
	currentText := content.get()
	data := struct{ Content string }{Content: currentText}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveContent(w http.ResponseWriter, r *http.Request, content *Content) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	content.set(r.FormValue("content"))

	// redirect to a new URL to prevent form resubmission issues
	http.Redirect(w, r, "/share", http.StatusSeeOther)
}
