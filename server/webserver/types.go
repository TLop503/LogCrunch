package webserver

import "time"

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
