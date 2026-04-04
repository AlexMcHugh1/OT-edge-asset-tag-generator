package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() error {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "./dfx.db"
	}
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		email         TEXT    NOT NULL UNIQUE,
		password_hash TEXT    NOT NULL,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		token      TEXT     PRIMARY KEY,
		user_id    INTEGER  NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS devices (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id       INTEGER NOT NULL,
		tag           TEXT    NOT NULL,
		device_name   TEXT    NOT NULL,
		serial_number TEXT, -- Added this line
		environment   TEXT    NOT NULL,
		location      TEXT    NOT NULL,
		is_global     INTEGER NOT NULL DEFAULT 0,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	return nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func createSession(userID int64) (string, error) {
	token, err := newSessionToken()
	if err != nil {
		return "", err
	}
	expires := time.Now().Add(30 * 24 * time.Hour)
	_, err = db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expires,
	)
	return token, err
}

func getUserFromSession(token string) (*User, error) {
	var u User
	err := db.QueryRow(`
		SELECT u.id, u.email
		FROM users u
		JOIN sessions s ON s.user_id = u.id
		WHERE s.token = ? AND s.expires_at > datetime('now')
	`, token).Scan(&u.ID, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func deleteSession(token string) {
	db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
