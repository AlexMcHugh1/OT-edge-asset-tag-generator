package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// fakeHash prevents user-enumeration via timing: we run bcrypt even when the
// username doesn't exist so the response time is indistinguishable from a real
// wrong-password attempt.
var fakeHash, _ = bcrypt.GenerateFromPassword([]byte("dfx_fake_timing_sentinel_x7k2"), 12)

// User is the authenticated user returned by session lookups.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// registrationKey reads the key from /data/registration.key (same volume as
// the database). Returns "" if the file doesn't exist — meaning registration
// is closed.
func registrationKey() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dfx.db"
	}
	keyFile := filepath.Join(filepath.Dir(dbPath), "registration.key")
	b, err := os.ReadFile(keyFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	regKey := registrationKey()
	if regKey == "" {
		jsonError(w, "Registration is closed", 403)
		return
	}
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		InviteCode  string `json:"invite_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	if req.InviteCode != regKey {
		jsonError(w, "Invalid invite code", 403)
		return
	}
	username := strings.TrimSpace(strings.ToLower(req.Username))
	// Strip domain suffix if someone submits a full email address.
	if idx := strings.Index(username, "@"); idx != -1 {
		username = username[:idx]
	}
	if username == "" {
		jsonError(w, "Username is required", 400)
		return
	}
	if len(username) > 64 {
		jsonError(w, "Username must be under 64 characters", 400)
		return
	}
	if len(req.Password) < 8 {
		jsonError(w, "Password must be at least 8 characters", 400)
		return
	}
	if len(req.Password) > 72 {
		jsonError(w, "Password must be under 72 characters", 400)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	res, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			jsonError(w, "That username is already taken", 409)
		} else {
			jsonError(w, "Server error", 500)
		}
		return
	}
	userID, _ := res.LastInsertId()
	updateLastActive(userID)
	token, err := createSession(userID)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	setSessionCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":   true,
		"user": map[string]any{"id": userID, "username": username},
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	username := strings.TrimSpace(strings.ToLower(req.Username))
	if idx := strings.Index(username, "@"); idx != -1 {
		username = username[:idx]
	}

	var userID int64
	var storedHash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).
		Scan(&userID, &storedHash)
	if err != nil {
		// Always run bcrypt to prevent timing-based user enumeration.
		bcrypt.CompareHashAndPassword(fakeHash, []byte(req.Password)) //nolint
		jsonError(w, "Invalid username or password", 401)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		jsonError(w, "Invalid username or password", 401)
		return
	}
	updateLastActive(userID)
	token, err := createSession(userID)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	setSessionCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":   true,
		"user": map[string]any{"id": userID, "username": username},
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("dfx_session"); err == nil {
		deleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "dfx_session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	w.Header().Set("Content-Type", "application/json")
	if user == nil {
		json.NewEncoder(w).Encode(map[string]any{"user": nil})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"user": user})
}

// setSessionCookie writes a long-lived, HttpOnly session cookie.
func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "dfx_session",
		Value:    token,
		MaxAge:   30 * 24 * 3600, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

// userFromRequest extracts the authenticated user from the session cookie.
// Returns nil if there is no valid session.
func userFromRequest(r *http.Request) *User {
	cookie, err := r.Cookie("dfx_session")
	if err != nil {
		return nil
	}
	user, err := getUserFromSession(cookie.Value)
	if err != nil {
		return nil
	}
	return user
}

// jsonError writes a JSON error response.
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
