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

	mu.Lock()
	// Insert at the top (index 0)
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
        <title>Phoenix | DFX Tag Generator</title>
        <style>
            :root { --phoenix-orange: #f26522; --border-color: #e5e7eb; --text-main: #374151; --text-muted: #9ca3af; }
            body { font-family: "Inter", -apple-system, sans-serif; background-color: #ffffff; color: var(--text-main); margin: 0; padding: 0; }
            
            /* Header */
            header { width: 100%; height: 64px; display: flex; align-items: center; padding: 0 24px; border-bottom: 1px solid var(--border-color); box-sizing: border-box; }
            .brand { color: var(--phoenix-orange); font-weight: 700; font-size: 20px; letter-spacing: -0.5px; }

            .content { max-width: 800px; margin: 48px auto; padding: 0 24px; text-align: center; }
            .section-title { color: var(--phoenix-orange); font-size: 12px; font-weight: 700; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 32px; }

            /* Generator Card */
            .gen-card { border: 1px solid var(--border-color); border-radius: 8px; padding: 40px; margin-bottom: 48px; display: flex; flex-direction: column; align-items: center; }
            .qr-wrapper { border: 1px solid var(--border-color); padding: 16px; border-radius: 4px; margin-bottom: 24px; }
            .display-box { display: flex; align-items: center; background: #f9fafb; border: 1px solid var(--border-color); border-radius: 4px; padding: 8px 12px; width: 100%; max-width: 440px; margin-bottom: 24px; }
            .display-box code { flex-grow: 1; font-family: monospace; font-size: 14px; text-align: left; overflow: hidden; text-overflow: ellipsis; }
            
            /* Buttons */
            .copy-btn { background: none; border: none; cursor: pointer; color: var(--text-muted); font-size: 16px; padding: 4px; display: flex; align-items: center; transition: color 0.2s; }
            .copy-btn:hover { color: var(--phoenix-orange); }
            .btn-primary { background: var(--phoenix-orange); color: white; border: none; padding: 10px 32px; border-radius: 4px; font-weight: 600; font-size: 13px; text-transform: uppercase; cursor: pointer; transition: opacity 0.2s; }
            .btn-primary:hover { opacity: 0.9; }

            /* History Table */
            .history-table { width: 100%; border-top: 1px solid var(--border-color); margin-top: 48px; }
            .history-header { font-size: 11px; color: var(--text-muted); text-transform: uppercase; padding: 16px 0; text-align: center; font-weight: 600; }
            .history-row { display: flex; align-items: center; justify-content: space-between; padding: 12px 16px; border-bottom: 1px solid #f3f4f6; transition: background 0.2s; }
            .history-row:hover { background: #fafafa; }
            .history-id { font-family: monospace; font-size: 13px; color: #4b5563; }
        </style>
    </head>
    <body>
        <header><div class="brand">Phoenix</div></header>
        
        <div class="content">
            <div class="section-title">Phoenix DFX Tag Generator</div>
            
            <div class="gen-card">
                <div class="qr-wrapper">
                    <img src="data:image/png;base64,{{.QRData}}" width="160" height="160">
                </div>
                
                <div class="display-box">
                    <code id="mainId">{{.UUID}}</code>
                    <button class="copy-btn" onclick="copyId('mainId')" title="Copy">📋</button>
                </div>

                <form method="GET">
                    <button type="submit" class="btn-primary">Generate</button>
                </form>
            </div>

            <div class="history-table">
                <div class="history-header">Recent Assets</div>
                {{range .History}}
                <div class="history-row">
                    <span class="history-id">{{.}}</span>
                    <button class="copy-btn" onclick="copyText('{{.}}')" title="Copy">📋</button>
                </div>
                {{end}}
            </div>
        </div>

        <script>
            function copyId(elementId) {
                const text = document.getElementById(elementId).innerText;
                copyText(text);
            }

            function copyText(text) {
                navigator.clipboard.writeText(text).then(() => {
                    // Optional: Add Phoenix-style toast here
                });
            }
        </script>
    </body>
    </html>`

	t := template.Must(template.New("phoenix").Parse(tmpl))
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Phoenix DFX Tag Generator starting on :9091...")
	http.ListenAndServe(":9091", nil)
}
