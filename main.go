package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type PageData struct {
	UUID      string
	QRData    string
	IsChecked bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	useSuffix := r.URL.Query().Get("monsters") == "true"

	// Generate Base UUID
	rawUUID := uuid.New().String()
	newID := "dfx-" + rawUUID

	// Apply "2319" logic: Replace the last 4 characters if checked
	if useSuffix {
		newID = newID[:len(newID)-4] + "2319"
	}

	// Generate QR Code
	png, _ := qrcode.Encode(newID, qrcode.Medium, 256)
	qrBase64 := base64.StdEncoding.EncodeToString(png)

	data := PageData{
		UUID:      newID,
		QRData:    qrBase64,
		IsChecked: useSuffix,
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>DFX Generator</title>
		<style>
			body { font-family: -apple-system, sans-serif; background-color: #f4f7f6; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
			.card { background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; max-width: 450px; width: 90%; }
			code { background: #eee; padding: 10px; border-radius: 4px; display: block; margin: 1.5rem 0; word-break: break-all; font-size: 0.95rem; border: 1px solid #ccc; }
			.options { margin-bottom: 1.5rem; font-size: 0.9rem; color: #555; }
			button { background: #3498db; color: white; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-size: 1rem; width: 100%; }
			button:hover { background: #2980b9; }
			input[type="checkbox"] { transform: scale(1.2); margin-right: 8px; vertical-align: middle; }
		</style>
	</head>
	<body>
		<div class="card">
			<h1>DFX ID Generator</h1>
			<img src="data:image/png;base64,{{.QRData}}" alt="QR Code">
			<code>{{.UUID}}</code>
			
			<div class="options">
				<label>
					<input type="checkbox" id="suffixCheck" {{if .IsChecked}}checked{{end}} onchange="toggleSuffix()">
					White Sock Incident (End in 2319)
				</label>
			</div>

			<button onclick="window.location.reload()">Generate New ID</button>
		</div>

		<script>
			function toggleSuffix() {
				const isChecked = document.getElementById('suffixCheck').checked;
				// Reload the page with the parameter to let the backend handle the ID generation
				window.location.href = isChecked ? "/?monsters=true" : "/";
			}
		</script>
	</body>
	</html>`

	t := template.Must(template.New("webpage").Parse(tmpl))
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server starting on :9091...")
	http.ListenAndServe(":9091", nil)
}