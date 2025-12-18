package webserver

import (
	"database/sql"
	"embed"
	"github.com/TLop503/LogCrunch/structs"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

// embed html files in the binary for distribution.
//
//go:embed site/templates/*.html site/pages/*.html site/static/*
var templateFS embed.FS
var templates *template.Template

// StartRouter webserver server on specified addr (:8080, from main).
// Needs RO db passed in
func StartRouter(addr string, connList *structs.ConnectionList, db *sql.DB) {

	// register helper functions
	funcMap := template.FuncMap{
		"formatUnix": func(ts int64) string {
			loc, _ := time.LoadLocation("Local")
			return time.Unix(ts, 0).In(loc).Format("01-02 15:04:05")
		},
		"formatGoTime": func(t time.Time) string {
			loc, _ := time.LoadLocation("Local")
			return t.In(loc).Format("2006-01-02 15:04:05")
		},
	}

	// register templates
	var err error
	templates, err = template.New("").
		Funcs(funcMap).
		ParseFS(templateFS,
			"site/templates/*.html",
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
	mux.HandleFunc("/logs", serveLogPage(db))    // Needs to be RO
	mux.HandleFunc("/query", serveQueryPage(db)) // RO

	// Serve static files as subtree of fs
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
