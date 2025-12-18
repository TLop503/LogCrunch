package webserver

import (
	"database/sql"
	"github.com/TLop503/LogCrunch/server/db"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"net/http"
)

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

func serveQueryPage(dbase *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data []structs.Log
			err  error
		)

		switch r.Method {
		case http.MethodGet:
			// last 50 by default!
			data, err = db.MostRecent50(dbase)
			if err != nil {
				http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
			// run the user-provided query
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Bad form data", http.StatusBadRequest)
				return
			}
			query := r.FormValue("query")
			data, err = db.RunQuery(dbase, query)
			if err != nil {
				http.Error(w, "Query failed", http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// render template with whatever data we got
		err = templates.ExecuteTemplate(w, "query", data)
		if err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// serveLogPage renders the contents of the intake log file
func serveLogPage(dbase *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.MostRecent50(dbase)
		if err != nil {
			http.Error(w, "Failed to parse log intake file", http.StatusInternalServerError)
			return
		}

		// Try passing a pointer to the struct
		err = templates.ExecuteTemplate(w, "logs", data)
		if err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
