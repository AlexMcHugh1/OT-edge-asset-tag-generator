package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type PageData struct {
	UUID    string
	QRData  string
	History []string
}

var (
	history []string
	mu      sync.Mutex
)

func handler(w http.ResponseWriter, r *http.Request) {
	newUUID := "dfx-" + uuid.New().String()

	// Update History (Last 5)
	mu.Lock()
	history = append([]string{newUUID}, history...)
	if len(history) > 5 {
		history = history[:5]
	}
	currentHistory := make([]string, len(history))
	copy(currentHistory, history)
	mu.Unlock()

	png, _ := qrcode.Encode(newUUID, qrcode.Medium, 256)
	qrBase64 := base64.StdEncoding.EncodeToString(png)

	data := PageData{
		UUID:    newUUID,
		QRData:  qrBase64,
		History: currentHistory,
	}

	tmpl := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Phoenix | ID Generator</title>
        <style>
            body { font-family: -apple-system, sans-serif; background-color: #ffffff; color: #333; margin: 0; display: flex; flex-direction: column; align-items: center; min-height: 100vh; }
            header { width: 100%; padding: 20px; display: flex; align-items: center; border-bottom: 1px solid #eee; margin-bottom: 40px; }
            .logo { color: #e67e22; font-weight: bold; font-size: 24px; text-decoration: none; }
            .container { text-align: center; max-width: 600px; width: 90%; }
            h2 { color: #e67e22; text-transform: uppercase; font-size: 14px; letter-spacing: 1px; margin-bottom: 30px; }
            .main-card { background: white; border: 1px solid #e0e0e0; border-radius: 12px; padding: 40px; margin-bottom: 30px; transition: box-shadow 0.3s; }
            .main-card:hover { box-shadow: 0 4px 20px rgba(0,0,0,0.05); }
            .qr-box { border: 1px solid #eee; padding: 15px; border-radius: 8px; display: inline-block; margin-bottom: 20px; }
            .uuid-text { font-family: monospace; background: #f9f9f9; padding: 10px 20px; border-radius: 5px; border: 1px solid #eee; font-size: 16px; margin-bottom: 20px; display: block; }
            .btn { background: #e67e22; color: white; border: none; padding: 12px 30px; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 600; text-transform: uppercase; }
            
            .history-section { width: 100%; margin-top: 40px; border-top: 1px solid #eee; padding-top: 20px; }
            .history-title { font-size: 12px; color: #999; text-transform: uppercase; margin-bottom: 15px; }
            .history-item { font-family: monospace; font-size: 13px; color: #666; padding: 8px 0; border-bottom: 1px dotted #eee; text-align: left; }
        </style>
    </head>
    <body>
        <header><div class="logo">Phoenix</div></header>
        <div class="container">
            <h2>ID Generator</h2>
            <div class="main-card">
                <div class="qr-box">
                    <img src="data:image/png;base64,{{.QRData}}" width="200">
                </div>
                <span class="uuid-text">{{.UUID}}</span>
                <button class="btn" onclick="window.location.reload()">Generate</button>
            </div>

            <div class="history-section">
                <div class="history-title">Recent Assets</div>
                {{range .History}}
                <div class="history-item">{{.}}</div>
                {{end}}
            </div>
        </div>
    </body>
    </html>`

	t := template.Must(template.New("phoenix").Parse(tmpl))
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Phoenix ID Gen starting on :9091...")
	http.ListenAndServe(":9091", nil)
}
