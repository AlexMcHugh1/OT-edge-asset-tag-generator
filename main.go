package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type HistoryItem struct {
	UUID      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
}

type PageData struct {
	UUID    string        `json:"uuid"`
	QRData  string        `json:"qr_data"`
	History []HistoryItem `json:"history"`
}

var (
	history []HistoryItem
	mu      sync.Mutex
)

func generateData() (string, string) {
	newUUID := "dfx-" + uuid.New().String()
	ts := time.Now().Format("15:04")

	mu.Lock()
	history = append([]HistoryItem{{UUID: newUUID, Timestamp: ts}}, history...)
	if len(history) > 5 {
		history = history[:5]
	}
	mu.Unlock()

	png, _ := qrcode.Encode(newUUID, qrcode.Medium, 256)
	return newUUID, base64.StdEncoding.EncodeToString(png)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	id, qr := generateData()
	mu.Lock()
	h := history
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PageData{UUID: id, QRData: qr, History: h})
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "missing text param", 400)
		return
	}
	png, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "qr encode failed", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"qr_data": base64.StdEncoding.EncodeToString(png),
	})
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	generateData()
	backupDB() // take an immediate backup on every startup
	startBackupSchedule()
	startInactiveAccountCleanup()

	http.HandleFunc("GET /api/generate", apiHandler)
	http.HandleFunc("GET /api/qr", qrHandler)
	http.HandleFunc("GET /api/auth/me", meHandler)
	http.HandleFunc("POST /api/auth/register", registerHandler)
	http.HandleFunc("POST /api/auth/login", loginHandler)
	http.HandleFunc("POST /api/auth/logout", logoutHandler)
	http.HandleFunc("DELETE /api/auth/account", deleteAccountHandler)
	http.HandleFunc("GET /api/devices", listDevicesHandler)
	http.HandleFunc("POST /api/devices", createDeviceHandler)
	http.HandleFunc("PUT /api/devices/{id}", updateDeviceHandler)
	http.HandleFunc("DELETE /api/devices/{id}", deleteDeviceHandler)
	http.HandleFunc("POST /api/privacy/erasure-request", erasureRequestHandler)
	http.HandleFunc("GET /privacy", privacyHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id, qr := generateData()
		mu.Lock()
		h := history
		mu.Unlock()
		t := template.Must(template.New("dfx").Parse(tmpl))
		t.Execute(w, PageData{UUID: id, QRData: qr, History: h})
	})

	fmt.Println("DFX Tag Generator online at :9092")
	log.Fatal(http.ListenAndServe(":9092", nil))
}
