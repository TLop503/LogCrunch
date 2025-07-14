package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

// embed html files in the binary for distribution.
//
//go:embed site/templates/*.html
var templateFS embed.FS
var templates *template.Template

// Start web server on specified addr (:8080, from main).
func Start(addr string) {
	var err error
	templates, err = template.ParseFS(templateFS, "site/templates/*.html")
	if err != nil {
		log.Fatalf("error parsing embedded templates: %v", err)
	}

	// routing!
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveStatic("index"))

	go func() {
		log.Printf("Webserver running at %s\n", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
}

// serveStatic injects the given page (templateName) into layout.html using templates.
func serveStatic(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, templateName, nil)
		if err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
