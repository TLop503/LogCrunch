package webserver

import (
	"github.com/TLop503/LogCrunch/structs"
	"html/template"
	"net/http"
)

// handleAliasEditForm renders a form to edit the alias for a given connection.
// expects a GET request with the `ip` parameter in the query string.
// form is rendered using the "alias-edit" template.
func handleAliasEditForm(connList *structs.ConnectionList, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get ip addr from query (request)
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "Missing IP parameter", http.StatusBadRequest)
			return
		}

		// lookup the connection by IP
		connList.RLock()
		conn, ok := connList.Connections[ip]
		connList.RUnlock()

		if !ok {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		}

		// safely extract the alias and remote address under lock
		conn.Lock()
		data := struct {
			RemoteAddr string
			Alias      string
		}{
			RemoteAddr: conn.RemoteAddr,
			Alias:      conn.Alias,
		}
		conn.Unlock()

		// render the alias-edit template with the extracted connection info
		err := tmpl.ExecuteTemplate(w, "alias-edit", data)
		if err != nil {
			http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// handleAliasSet processes the submitted alias form and updates the in-memory connection alias.
// expects a POST request with `ip` and `alias` fields.
func handleAliasSet(connList *structs.ConnectionList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// extract form values
		ip := r.FormValue("ip")
		alias := r.FormValue("alias")

		// validate alias length (max 32 characters)
		if len(alias) > 32 {
			http.Error(w, "Alias too long (max 32 chars)", http.StatusBadRequest)
			return
		}

		// lookup the connection by IP
		connList.RLock()
		conn, ok := connList.Connections[ip]
		connList.RUnlock()

		if !ok {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		}

		// safely update the alias
		conn.Lock()
		conn.Alias = alias
		conn.Unlock()

		// re-dir back to the connection list after saving
		http.Redirect(w, r, "/connections", http.StatusSeeOther)
	}
}
