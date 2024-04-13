package main

import (
	"fmt"
	"log"
	"net"
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

	port := 8080
	serverSocket := fmt.Sprintf("%s:%d", hostIP(), port)
	log.Printf("Server started at http://%s/share", serverSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func hostIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	hostLocalIP := conn.LocalAddr().(*net.UDPAddr)

	return hostLocalIP.IP.String()
}
