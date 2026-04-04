package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// seedAdminAccount creates admin@deltaflare.com on first startup if it doesn't
// exist. Set ADMIN_PASSWORD to override the default password.
func seedAdminAccount() {
	const email = "admin@deltaflare.com"
	password := os.Getenv("ADMIN_PASSWORD")
	usingDefault := password == ""
	if usingDefault {
		password = "DFXadmin1!"
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("seed admin: %v", err)
		return
	}
	if _, err = db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, string(hash)); err != nil {
		log.Printf("seed admin: %v", err)
		return
	}
	if usingDefault {
		log.Printf("Seeded admin account: %s  password: %s  (set ADMIN_PASSWORD env var to override)", email, password)
	} else {
		log.Printf("Seeded admin account: %s", email)
	}
}

// fakeHash prevents user-enumeration via timing: we run bcrypt even when the
// email doesn't exist so the response time is indistinguishable from a real
// wrong-password attempt.
var fakeHash, _ = bcrypt.GenerateFromPassword([]byte("dfx_fake_timing_sentinel_x7k2"), 12)

// User is the authenticated user returned by session lookups.
type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if !strings.HasSuffix(email, "@deltaflare.com") {
		jsonError(w, "Only @deltaflare.com email addresses are allowed", 403)
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
	res, err := db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			jsonError(w, "An account with this email already exists", 409)
		} else {
			jsonError(w, "Server error", 500)
		}
		return
	}
	userID, _ := res.LastInsertId()
	token, err := createSession(userID)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	setSessionCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"user": map[string]interface{}{"id": userID, "email": email},
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))

	var userID int64
	var storedHash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).
		Scan(&userID, &storedHash)
	if err != nil {
		// Always run bcrypt to prevent timing-based user enumeration.
		bcrypt.CompareHashAndPassword(fakeHash, []byte(req.Password)) //nolint
		jsonError(w, "Invalid email or password", 401)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		jsonError(w, "Invalid email or password", 401)
		return
	}
	token, err := createSession(userID)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	setSessionCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"user": map[string]interface{}{"id": userID, "email": email},
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
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
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
