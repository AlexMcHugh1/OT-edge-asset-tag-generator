package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type PageData struct {
	UUID   string
	QRData string
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. Generate standard dfx- UUID
	newUUID := "dfx-" + uuid.New().String()

	// 2. Generate QR Code
	png, _ := qrcode.Encode(newUUID, qrcode.Medium, 256)
	qrBase64 := base64.StdEncoding.EncodeToString(png)

	data := PageData{UUID: newUUID, QRData: qrBase64}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>DFX ID Generator</title>
		<style>
			body { font-family: -apple-system, sans-serif; background-color: #f4f7f6; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
			.card { background: white; padding: 2.5rem; border-radius: 12px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); text-align: center; max-width: 480px; width: 95%; }
			h1 { color: #2c3e50; font-size: 1.6rem; margin-top: 0; margin-bottom: 2rem; font-weight: 600; }
			.qr-container { display: flex; justify-content: center; margin-bottom: 2rem; border: 1px solid #ddd; padding: 12px; border-radius: 8px; max-width: 256px; margin-left: auto; margin-right: auto; }
			.uuid-display { display: flex; align-items: center; justify-content: center; gap: 8px; background: #eee; padding: 12px; border-radius: 6px; border: 1px solid #ccc; margin-bottom: 2rem; }
			code { word-break: break-all; font-size: 0.95rem; color: #333; flex-grow: 1; text-align: left; }
			.copy-btn { background: #e0e0e0; border: none; padding: 6px; border-radius: 4px; cursor: pointer; transition: background 0.2s; color: #666; font-size: 1.1rem; display: flex; align-items: center; justify-content: center; }
			.copy-btn:hover { background: #d0d0d0; }
			.primary-btn { background: #3498db; color: white; border: none; padding: 14px; border-radius: 8px; cursor: pointer; font-size: 1.05rem; width: 100%; transition: background 0.2s; }
			.primary-btn:hover { background: #2980b9; }
			.copy-feedback { display: none; margin-top: -1.5rem; margin-bottom: 1.5rem; color: #27ae60; font-size: 0.85rem; font-weight: 500; }
		</style>
	</head>
	<body>
		<div class="card">
			<h1>DFX ID Generator</h1>
			
			<div class="qr-container">
				<img src="data:image/png;base64,{{.QRData}}" alt="QR Code" width="256" height="256">
			</div>

			<div class="uuid-display">
				<code id="uuidValue">{{.UUID}}</code>
				<button class="copy-btn" onclick="copyToClipboard()" title="Copy ID">📋</button>
			</div>

			<div class="copy-feedback" id="copyFeedback">✓ Copied to clipboard!</div>

			<button class="primary-btn" onclick="window.location.reload()">Generate New ID</button>
		</div>

		<script>
			function copyToClipboard() {
				const uuidText = document.getElementById('uuidValue').innerText;
				const feedback = document.getElementById('copyFeedback');

				navigator.clipboard.writeText(uuidText).then(() => {
					// Show the "Copied" message
					feedback.style.display = 'block';

					// Hide it after 2 seconds
					setTimeout(() => {
						feedback.style.display = 'none';
					}, 2000);
				}).catch(err => {
					console.error('Failed to copy: ', err);
					feedback.innerText = '❌ Copy failed';
					feedback.style.display = 'block';
					setTimeout(() => {
						feedback.style.display = 'none';
						feedback.innerText = '✓ Copied to clipboard!'; // Reset text
					}, 2000);
				});
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