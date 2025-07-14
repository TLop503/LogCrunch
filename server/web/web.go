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
func Start(addr string, connList *structs.ConnectionList) {
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
	mux.HandleFunc("/connections", serveConnectionsPage(connList))
	mux.HandleFunc("/alias", handleAliasSet(connList))
	mux.HandleFunc("/alias/edit", handleAliasEditForm(connList, templates))

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
func serveConnectionsPage(connList *structs.ConnectionList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Safely read from the connections list
		connList.RLock()
		connections := make([]*structs.Connection, 0, len(connList.Connections))
		for _, conn := range connList.Connections {
			connections = append(connections, conn)
		}
		connList.RUnlock()

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

func handleAliasEditForm(connList *structs.ConnectionList, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "Missing IP parameter", http.StatusBadRequest)
			return
		}

		connList.RLock()
		conn, ok := connList.Connections[ip]
		connList.RUnlock()

		if !ok {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		}

		conn.Lock()
		data := struct {
			RemoteAddr string
			Alias      string
		}{
			RemoteAddr: conn.RemoteAddr,
			Alias:      conn.Alias,
		}
		conn.Unlock()

		err := tmpl.ExecuteTemplate(w, "alias-edit", data)
		if err != nil {
			http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleAliasSet(connList *structs.ConnectionList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ip := r.FormValue("ip")
		alias := r.FormValue("alias")

		if len(alias) > 32 {
			http.Error(w, "Alias too long (max 32 chars)", http.StatusBadRequest)
			return
		}

		connList.RLock()
		conn, ok := connList.Connections[ip]
		connList.RUnlock()

		if !ok {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		}

		conn.Lock()
		conn.Alias = alias
		conn.Unlock()

		http.Redirect(w, r, "/connections", http.StatusSeeOther)
	}
}
