package main

import (
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		username       TEXT    NOT NULL UNIQUE,
		password_hash  TEXT    NOT NULL,
		created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}
	// Idempotent migrations for existing deployments.
	db.Exec(`ALTER TABLE users RENAME COLUMN email TO username`)
	db.Exec(`ALTER TABLE users ADD COLUMN last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP`)
	// Strip @domain suffix from any usernames that are still stored as emails.
	db.Exec(`UPDATE users SET username = substr(username, 1, instr(username, '@') - 1) WHERE username LIKE '%@%'`)
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		token      TEXT     PRIMARY KEY,
		user_id    INTEGER  NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS erasure_requests (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		username   TEXT    NOT NULL,
		reason     TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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
		SELECT u.id, u.username
		FROM users u
		JOIN sessions s ON s.user_id = u.id
		WHERE s.token = ? AND s.expires_at > datetime('now')
	`, token).Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func deleteSession(token string) {
	db.Exec("DELETE FROM sessions WHERE token = ?", token)
}

// updateLastActive records the current time as the user's last activity.
// Called on login and registration so the 12-month retention clock resets.
func updateLastActive(userID int64) {
	db.Exec("UPDATE users SET last_active_at = datetime('now') WHERE id = ?", userID)
}

// deleteInactiveUsers permanently removes accounts that have had no activity
// for more than 12 months, fulfilling the UK GDPR retention policy.
func deleteInactiveUsers() {
	db.Exec(`DELETE FROM users
		WHERE COALESCE(last_active_at, created_at) < datetime('now', '-12 months')`)
}

// backupDB creates a clean snapshot via VACUUM INTO, compresses it with gzip,
// and skips saving if the database content hasn't changed since the last
// backup. Backups are stored in /backups/ next to the main database file.
// Only the 7 most recent are kept.
func backupDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dfx.db"
	}
	backupDir := filepath.Join(filepath.Dir(dbPath), "backups")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		log.Printf("backup: could not create backup directory: %v", err)
		return
	}

	// VACUUM INTO a temp file so we get a consistent snapshot.
	tmp := filepath.Join(backupDir, ".tmp.db")
	os.Remove(tmp)
	if _, err := db.Exec("VACUUM INTO ?", tmp); err != nil {
		log.Printf("backup: VACUUM INTO failed: %v", err)
		return
	}
	defer os.Remove(tmp)

	// Hash the snapshot and compare to the most recent backup.
	newHash, err := sha256File(tmp)
	if err != nil {
		log.Printf("backup: hash failed: %v", err)
		return
	}
	if newHash == latestBackupHash(backupDir) {
		log.Printf("backup: no changes, skipping")
		return
	}

	// Compress and write the new backup.
	dest := filepath.Join(backupDir, fmt.Sprintf("dfx-%s.db.gz", time.Now().Format("2006-01-02T15-04-05")))
	if err := gzipFile(tmp, dest); err != nil {
		log.Printf("backup: compress failed: %v", err)
		return
	}
	log.Printf("backup: wrote %s", dest)
	pruneBackups(backupDir, 7)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// latestBackupHash returns the SHA256 of the most recent backup by
// decompressing it, so we can compare it to the current snapshot.
func latestBackupHash(dir string) string {
	entries, _ := os.ReadDir(dir)
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".db.gz") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	if len(files) == 0 {
		return ""
	}
	sort.Strings(files)
	latest := files[len(files)-1]

	f, err := os.Open(latest)
	if err != nil {
		return ""
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return ""
	}
	defer gr.Close()
	h := sha256.New()
	io.Copy(h, gr) //nolint
	return hex.EncodeToString(h.Sum(nil))
}

func gzipFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	gz := gzip.NewWriter(out)
	if _, err := io.Copy(gz, in); err != nil {
		return err
	}
	return gz.Close()
}

// pruneBackups removes old backups keeping only the most recent n files.
func pruneBackups(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".db.gz") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)
	for len(files) > keep {
		os.Remove(files[0])
		files = files[1:]
	}
}

// startBackupSchedule runs a daily backup at 02:00, one hour before the
// inactive-account cleanup, so there is always a recent backup before any
// deletions occur.
func startBackupSchedule() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
			time.Sleep(time.Until(next))
			backupDB()
		}
	}()
}
