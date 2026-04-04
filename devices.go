package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Device represents a saved tag-to-device mapping.
type Device struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	UserEmail    string `json:"user_email"`
	Tag          string `json:"tag"`
	DeviceName   string `json:"device_name"`
	SerialNumber string `json:"serial_number"`
	Environment  string `json:"environment"`
	Location     string `json:"location"`
	IsGlobal     bool   `json:"is_global"`
	CreatedAt    string `json:"created_at"`
}

func listDevicesHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Device{})
		return
	}
	rows, err := db.Query(`
		SELECT d.id, d.user_id, u.email, d.tag, d.device_name, d.serial_number,
		       d.environment, d.location, d.is_global, d.created_at
		FROM devices d
		JOIN users u ON u.id = d.user_id
		WHERE d.is_global = 1 OR d.user_id = ?
		ORDER BY d.created_at DESC
	`, user.ID)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	defer rows.Close()
	devices := []Device{}
	for rows.Next() {
		var d Device
		var isGlobal int
		if err := rows.Scan(&d.ID, &d.UserID, &d.UserEmail, &d.Tag,
			&d.DeviceName, &d.SerialNumber, &d.Environment, &d.Location, &isGlobal, &d.CreatedAt); err != nil {
			continue
		}
		d.IsGlobal = isGlobal == 1
		devices = append(devices, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	if user == nil {
		jsonError(w, "Authentication required", 401)
		return
	}
	var req struct {
		Tag          string `json:"tag"`
		DeviceName   string `json:"device_name"`
		SerialNumber string `json:"serial_number"`
		Environment  string `json:"environment"`
		Location     string `json:"location"`
		IsGlobal     bool   `json:"is_global"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	req.Tag = strings.TrimSpace(req.Tag)
	req.DeviceName = strings.TrimSpace(req.DeviceName)
	req.SerialNumber = strings.TrimSpace(req.SerialNumber)
	req.Environment = strings.TrimSpace(req.Environment)
	req.Location = strings.TrimSpace(req.Location)

	if req.Tag == "" || req.DeviceName == "" || req.Environment == "" || req.Location == "" {
		jsonError(w, "All fields except serial number are required", 400)
		return
	}
	if len(req.DeviceName) > 100 || len(req.Environment) > 100 || len(req.Location) > 100 || len(req.Tag) > 200 {
		jsonError(w, "Field value too long", 400)
		return
	}
	isGlobal := 0
	if req.IsGlobal {
		isGlobal = 1
	}
	res, err := db.Exec(`
		INSERT INTO devices (user_id, tag, device_name, serial_number, environment, location, is_global)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, user.ID, req.Tag, req.DeviceName, req.SerialNumber, req.Environment, req.Location, isGlobal)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	id, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id})
}

func updateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	if user == nil {
		jsonError(w, "Authentication required", 401)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		jsonError(w, "Invalid device ID", 400)
		return
	}
	var ownerID int64
	var isGlobal int
	if err := db.QueryRow("SELECT user_id, is_global FROM devices WHERE id = ?", id).
		Scan(&ownerID, &isGlobal); err != nil {
		jsonError(w, "Device not found", 404)
		return
	}
	if ownerID != user.ID && isGlobal == 0 {
		jsonError(w, "Forbidden", 403)
		return
	}
	var req struct {
		DeviceName   string `json:"device_name"`
		SerialNumber string `json:"serial_number"`
		Environment  string `json:"environment"`
		Location     string `json:"location"`
		IsGlobal     bool   `json:"is_global"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	req.DeviceName = strings.TrimSpace(req.DeviceName)
	req.SerialNumber = strings.TrimSpace(req.SerialNumber)
	req.Environment = strings.TrimSpace(req.Environment)
	req.Location = strings.TrimSpace(req.Location)

	if req.DeviceName == "" || req.Environment == "" || req.Location == "" {
		jsonError(w, "All fields except serial number are required", 400)
		return
	}
	if len(req.DeviceName) > 100 || len(req.Environment) > 100 || len(req.Location) > 100 {
		jsonError(w, "Field value too long", 400)
		return
	}
	newIsGlobal := 0
	if req.IsGlobal {
		newIsGlobal = 1
	}
	_, err = db.Exec(`
		UPDATE devices SET device_name = ?, serial_number = ?, environment = ?, location = ?, is_global = ?
		WHERE id = ?
	`, req.DeviceName, req.SerialNumber, req.Environment, req.Location, newIsGlobal, id)
	if err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	if user == nil {
		jsonError(w, "Authentication required", 401)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		jsonError(w, "Invalid device ID", 400)
		return
	}
	var ownerID int64
	var isGlobal int
	if err := db.QueryRow("SELECT user_id, is_global FROM devices WHERE id = ?", id).
		Scan(&ownerID, &isGlobal); err != nil {
		jsonError(w, "Device not found", 404)
		return
	}
	if ownerID != user.ID && isGlobal == 0 {
		jsonError(w, "Forbidden", 403)
		return
	}
	db.Exec("DELETE FROM devices WHERE id = ?", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
