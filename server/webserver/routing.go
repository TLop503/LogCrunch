package webserver

import (
	"database/sql"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/TLop503/LogCrunch/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// embed html files in the binary for distribution.
//
//go:embed website_content/templates/*.html website_content/pages/*.html website_content/static/*
var templateFS embed.FS

// templates holds all parsed templates with helper functions
var templates *template.Template

// initTemplates parses and registers all templates with helper functions
func initTemplates() error {
	var err error
	templates, err = template.New("").
		Funcs(helperFuncMap()).
		ParseFS(templateFS,
			"website_content/templates/*.html",
			"website_content/pages/*.html",
		)
	return err
}

// helperFuncMap returns a mapping of helper functions for templates
func helperFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatUnix": func(ts int64) string {
			loc, _ := time.LoadLocation("Local")
			return time.Unix(ts, 0).In(loc).Format("01-02 15:04:05")
		},
		"formatGoTime": func(t time.Time) string {
			loc, _ := time.LoadLocation("Local")
			return t.In(loc).Format("2006-01-02 15:04:05")
		},
	}
}

// setupRoutes configures all application routes
func setupRoutes(r *chi.Mux, connList *structs.ConnectionList, logDb *sql.DB, userDb *sql.DB) {
	// Middleware
	r.Use(middleware.Logger)

	// Static content - serve from embedded filesystem
	staticFS, err := fs.Sub(templateFS, "website_content/static")
	if err != nil {
		log.Fatalf("error creating static file system: %v", err)
	}
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Public routes (no auth required)
	r.Get("/login", serveLoginPage())
	r.Post("/api/auth/login", handleLogin(userDb))

	// Protected routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(userDb))

		// Pages
		r.Get("/", servePage("index", nil))
		r.Get("/connections", serveConnectionsPage(connList))
		r.Get("/logs", serveLogPage(logDb))
		r.Get("/query", serveQueryPage(logDb))

		// API endpoints
		r.Post("/alias", handleAliasSet(connList))
		r.Get("/alias/edit", handleAliasEditForm(connList, templates))

		// Auth API endpoints (require existing session)
		r.Post("/api/auth/logout", handleLogout(userDb))
		r.Post("/api/auth/password", handlePasswordUpdate(userDb))
		r.Get("/api/auth/check", handleSessionCheck(userDb))
	})
}

// StartRouter starts the webserver on the specified address
func StartRouter(addr string, connList *structs.ConnectionList, logDb *sql.DB, userDb *sql.DB) {
	// Initialize templates
	if err := initTemplates(); err != nil {
		log.Fatalf("error parsing embedded templates: %v", err)
	}

	// Setup router
	r := chi.NewRouter()
	setupRoutes(r, connList, logDb, userDb)

	// Start server
	log.Printf("Starting webserver at %s\n", addr)
	go func() {
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
}
