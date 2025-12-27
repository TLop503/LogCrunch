package webserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/TLop503/LogCrunch/server/db/users"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "logcrunch_session"
	sessionDuration   = 2 * time.Hour
)

// LoginRequest represents the JSON body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PasswordUpdateRequest represents the JSON body for password update
type PasswordUpdateRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// APIResponse is a generic JSON response
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the list
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colonIdx := strings.LastIndex(ip, ":"); colonIdx != -1 {
		ip = ip[:colonIdx]
	}
	return ip
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// serveLoginPage serves the login page
func serveLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "login", nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// servePasswordChangePage serves the password change page
func servePasswordChangePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "password-change", nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// authMiddleware checks for valid session and redirects to login if not authenticated
func authMiddleware(userDb *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			clientIP := getClientIP(r)
			session, err := users.ValidateSession(userDb, cookie.Value, clientIP)
			if err != nil || session == nil {
				// Clear invalid cookie
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   -1,
				})
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Check if user requires password change
			user, err := users.GetUserByID(userDb, session.UserID)
			if err != nil || user == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if user.RequiresPasswordChange {
				// Redirect to password change page unless already there
				http.Redirect(w, r, "/password-change", http.StatusFound)
				return
			}

			// Session is valid and no password change required, proceed
			next.ServeHTTP(w, r)
		})
	}
}

// passwordChangeAuthMiddleware allows access only if logged in (for password change page)
func passwordChangeAuthMiddleware(userDb *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			clientIP := getClientIP(r)
			session, err := users.ValidateSession(userDb, cookie.Value, clientIP)
			if err != nil || session == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Session is valid, proceed (don't check password change requirement)
			next.ServeHTTP(w, r)
		})
	}
}

// handleLogin handles POST /api/auth/login
func handleLogin(userDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Username and password are required",
			})
			return
		}

		// Look up user
		user, err := users.GetUserByUsername(userDB, req.Username)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Internal server error",
			})
			return
		}

		if user == nil {
			// User not found - use same error as wrong password to prevent enumeration
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid username or password",
			})
			return
		}

		// Check if user is active
		if !user.IsActive {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Account is disabled",
			})
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid username or password",
			})
			return
		}

		// Get client IP for session
		clientIP := getClientIP(r)

		// Create session
		session, err := users.CreateSession(userDB, user.ID, clientIP, sessionDuration)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Failed to create session",
			})
			return
		}

		// Update last login
		users.UpdateLastLogin(userDB, user.ID, clientIP)

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    session.ID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteStrictMode,
			Expires:  session.ExpiresAt,
		})

		// Check if password change is required
		if user.RequiresPasswordChange {
			writeJSON(w, http.StatusOK, APIResponse{
				Success: true,
				Message: "Login successful. Password change required.",
			})
			return
		}

		writeJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Login successful",
		})
	}
}

// handleLogout handles POST /api/auth/logout
func handleLogout(userDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err == nil && cookie.Value != "" {
			// Delete session from database
			users.DeleteSession(userDB, cookie.Value)
		}

		// Clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		writeJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Logged out successfully",
		})
	}
}

// handlePasswordUpdate handles POST /api/auth/password
func handlePasswordUpdate(userDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session from cookie
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Not authenticated",
			})
			return
		}

		// Validate session
		clientIP := getClientIP(r)
		session, err := users.ValidateSession(userDB, cookie.Value, clientIP)
		if err != nil || session == nil {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid or expired session",
			})
			return
		}

		// Parse request body
		var req PasswordUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		// Validate input
		if req.CurrentPassword == "" || req.NewPassword == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Current password and new password are required",
			})
			return
		}

		// Validate new password strength (minimum 8 characters)
		if len(req.NewPassword) < 8 {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "New password must be at least 8 characters",
			})
			return
		}

		// Get user
		user, err := users.GetUserByID(userDB, session.UserID)
		if err != nil || user == nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Failed to retrieve user",
			})
			return
		}

		// Verify current password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Current password is incorrect",
			})
			return
		}

		// Hash new password
		newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Failed to hash password",
			})
			return
		}

		// Update password
		if err := users.UpdatePassword(userDB, user.ID, string(newHash)); err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Failed to update password",
			})
			return
		}

		// Invalidate all other sessions for security
		users.DeleteAllUserSessions(userDB, user.ID)

		// Create a new session for the current user
		newSession, err := users.CreateSession(userDB, user.ID, clientIP, sessionDuration)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Password updated but failed to create new session",
			})
			return
		}

		// Set new session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    newSession.ID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  newSession.ExpiresAt,
		})

		writeJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Password updated successfully",
		})
	}
}

// handleSessionCheck handles GET /api/auth/check - returns current auth status
func handleSessionCheck(userDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"authenticated": false,
			})
			return
		}

		clientIP := getClientIP(r)
		session, err := users.ValidateSession(userDB, cookie.Value, clientIP)
		if err != nil || session == nil {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"authenticated": false,
			})
			return
		}

		// Get user info
		user, err := users.GetUserByID(userDB, session.UserID)
		if err != nil || user == nil {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"authenticated": false,
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated":            true,
			"username":                 user.Username,
			"can_create_users":         user.CanCreateUsers,
			"requires_password_change": user.RequiresPasswordChange,
		})
	}
}
