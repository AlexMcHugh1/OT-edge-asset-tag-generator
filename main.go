package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type PageData struct {
	UUID    string   `json:"uuid"`
	QRData  string   `json:"qr_data"`
	History []string `json:"history"`
}

var (
	history []string
	mu      sync.Mutex
)

func generateData() (string, string) {
	newUUID := "dfx-" + uuid.New().String()
	mu.Lock()
	history = append([]string{newUUID}, history...)
	if len(history) > 5 {
		history = history[:5]
	}
	mu.Unlock()

	png, _ := qrcode.Encode(newUUID, qrcode.Medium, 256)
	qrBase64 := base64.StdEncoding.EncodeToString(png)
	return newUUID, qrBase64
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	id, qr := generateData()
	mu.Lock()
	h := history
	mu.Unlock()

	json.NewEncoder(w).Encode(PageData{
		UUID:    id,
		QRData:  qr,
		History: h,
	})
}

func main() {
	generateData()

	http.HandleFunc("/api/generate", apiHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		h := history
		mu.Unlock()

		currentId := h[0]
		png, _ := qrcode.Encode(currentId, qrcode.Medium, 256)
		qr := base64.StdEncoding.EncodeToString(png)

		t := template.Must(template.New("phoenix").Parse(tmpl))
		t.Execute(w, PageData{UUID: currentId, QRData: qr, History: h})
	})

	fmt.Println("Phoenix DFX Tag Generator online at :9091")
	http.ListenAndServe(":9091", nil)
}

const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Phoenix | DFX Tag Generator</title>
    <style>
        :root { 
            --p-orange: #f26522; 
            --p-orange-hover: #ff7e42;
            --p-border: #e5e7eb; 
            --p-text: #374151; 
            --p-muted: #9ca3af; 
            --p-toast-bg: #e0f2fe; 
            --p-toast-text: #0369a1; 
        }
        body { font-family: "Inter", system-ui, sans-serif; background: #fff; color: var(--p-text); margin: 0; overflow-x: hidden; }
        header { height: 64px; border-bottom: 1px solid var(--p-border); display: flex; align-items: center; padding: 0 24px; }
        .brand { color: var(--p-orange); font-weight: 700; font-size: 20px; }
        .content { max-width: 600px; margin: 48px auto; padding: 0 24px; text-align: center; }
        .title { color: var(--p-orange); font-size: 12px; font-weight: 700; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 32px; }
        .card { border: 1px solid var(--p-border); border-radius: 8px; padding: 40px; margin-bottom: 40px; }
        .qr-frame { border: 1px solid var(--p-border); padding: 12px; border-radius: 4px; display: inline-block; margin-bottom: 24px; }
        .id-display { display: flex; align-items: center; background: #f9fafb; border: 1px solid var(--p-border); border-radius: 4px; padding: 10px 14px; margin-bottom: 24px; }
        code { flex-grow: 1; font-family: ui-monospace, monospace; font-size: 14px; text-align: left; }
        .btn-main { 
            background: var(--p-orange); 
            color: white; 
            border: none; 
            padding: 12px 40px; 
            border-radius: 4px; 
            font-weight: 600; 
            text-transform: uppercase; 
            cursor: pointer; 
            width: 100%; 
            font-size: 13px;
            transition: background 0.2s, transform 0.1s;
        }
        .btn-main:hover { background: var(--p-orange-hover); }
        .btn-main:active { transform: scale(0.98); }
        .icon-btn { background: none; border: none; cursor: pointer; padding: 6px; display: flex; align-items: center; border-radius: 4px; transition: color 0.2s; color: var(--p-text); }
        .icon-btn:hover { background: #f1f5f9; }
        .icon-btn svg { width: 18px; height: 18px; fill: none; stroke: currentColor; stroke-width: 2; }
        .icon-btn.copied { color: #10b981; }
        .history { border-top: 1px solid var(--p-border); padding-top: 32px; text-align: left; }
        .history-label { font-size: 11px; font-weight: 700; color: var(--p-muted); text-transform: uppercase; margin-bottom: 16px; display: block; text-align: center; }
        .history-row { display: flex; align-items: center; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #f3f4f6; }
        .h-id { font-family: ui-monospace, monospace; font-size: 13px; color: #4b5563; }
        #toast {
            position: fixed;
            bottom: 24px;
            left: 50%;
            transform: translateX(-50%) translateY(100px);
            background: var(--p-toast-bg);
            color: var(--p-toast-text);
            padding: 12px 24px;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 500;
            transition: transform 0.3s cubic-bezier(0.18, 0.89, 0.32, 1.28);
            box-shadow: 0 4px 12px rgba(0,0,0,0.05);
            z-index: 1000;
        }
        #toast.show { transform: translateX(-50%) translateY(0); }
    </style>
</head>
<body>
    <header><div class="brand">Phoenix</div></header>
    <div class="content">
        <div class="title">DFX Tag Generator</div>
        <div class="card">
            <div class="qr-frame"><img id="qrImg" src="data:image/png;base64,{{.QRData}}" width="180"></div>
            <div class="id-display">
                <code id="currentId">{{.UUID}}</code>
                <button class="icon-btn" onclick="handleCopy('currentId', this)">
                    <svg class="copy-icon" viewBox="0 0 24 24"><path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"></path><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"></path></svg>
                </button>
            </div>
            <button class="btn-main" onclick="generate()">Generate Tag</button>
        </div>
        <div class="history">
            <span class="history-label">Recent Assets</span>
            <div id="historyList">
                {{range .History}}
                <div class="history-row">
                    <span class="h-id">{{.}}</span>
                    <button class="icon-btn" onclick="handleCopyText('{{.}}', this)">
                        <svg class="copy-icon" viewBox="0 0 24 24"><path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"></path><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"></path></svg>
                    </button>
                </div>
                {{end}}
            </div>
        </div>
    </div>
    <div id="toast">✓ ID copied to clipboard</div>
    <script>
        const copySvg = '<path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"></path><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"></path>';
        const checkSvg = '<path d="M20 6L9 17L4 12" stroke-linecap="round" stroke-linejoin="round"/>';

        async function generate() {
            const res = await fetch('/api/generate');
            const data = await res.json();
            document.getElementById('qrImg').src = 'data:image/png;base64,' + data.qr_data;
            document.getElementById('currentId').innerText = data.uuid;
            const list = document.getElementById('historyList');
            let html = '';
            data.history.forEach(id => {
                html += '<div class="history-row"><span class="h-id">' + id + '</span>' +
                        '<button class="icon-btn" onclick="handleCopyText(\'' + id + '\', this)">' +
                        '<svg viewBox="0 0 24 24">' + copySvg + '</svg></button></div>';
            });
            list.innerHTML = html;
        }

        function handleCopy(id, btn) { handleCopyText(document.getElementById(id).innerText, btn); }
        function handleCopyText(text, btn) {
            navigator.clipboard.writeText(text);
            showToast();
            const svg = btn.querySelector('svg');
            const originalHtml = svg.innerHTML;
            btn.classList.add('copied');
            svg.innerHTML = checkSvg;
            setTimeout(() => {
                btn.classList.remove('copied');
                svg.innerHTML = originalHtml;
            }, 2000);
        }
        function showToast() {
            const t = document.getElementById('toast');
            t.classList.add('show');
            setTimeout(() => t.classList.remove('show'), 3000);
        }
    </script>
</body>
</html>`
