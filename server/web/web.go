package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/TLop503/LogCrunch/structs"
)

// embed html files in the binary for distribution.
//
//go:embed site/templates/*.html site/pages/*.html site/static/*
var templateFS embed.FS
var templates *template.Template

// Start web server on specified addr (:8080, from main).
func Start(addr string, conns *structs.ConnectionList) {
	var err error
	templates, err = template.ParseFS(templateFS,
		"site/templates/base.html",
		"site/templates/navbar.html",
		"site/pages/*.html",
	)
	if err != nil {
		log.Fatalf("error parsing embedded templates: %v", err)
	}

	// routing!
	mux := http.NewServeMux()
	mux.HandleFunc("/", servePage("index", nil))
	mux.HandleFunc("/connections", serveConnectionsPage(conns))

	// Serve static files
	staticFS, err := fs.Sub(templateFS, "site/static")
	if err != nil {
		log.Fatalf("error creating static filesystem: %v", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	go func() {
		log.Printf("Webserver running at %s\n", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
}

// serveConnectionsPage safely reads the connections list and serves the connections page
func serveConnectionsPage(conns *structs.ConnectionList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Safely read from the connections list
		conns.RLock()
		connections := make([]*structs.Connection, 0, len(conns.Connections))
		for _, conn := range conns.Connections {
			connections = append(connections, conn)
		}
		conns.RUnlock()

		err := templates.ExecuteTemplate(w, "connections", connections)
		if err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// servePage renders the given template with the provided data.
func servePage(templateName string, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, templateName, data)
		if err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
