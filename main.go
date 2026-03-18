package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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

		current := h[0]
		png, _ := qrcode.Encode(current.UUID, qrcode.Medium, 256)
		qr := base64.StdEncoding.EncodeToString(png)

		t := template.Must(template.New("dfx").Parse(tmpl))
		t.Execute(w, PageData{UUID: current.UUID, QRData: qr, History: h})
	})

	fmt.Println("DFX Tag Generator online at :9092")
	http.ListenAndServe(":9092", nil)
}

const tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>DFX Tag Generator</title>
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
        body { font-family: "Inter", system-ui, sans-serif; background: #fff; color: var(--p-text); margin: 0; -webkit-tap-highlight-color: transparent; }
        
        /* Adjusted content margin now that header is gone */
        .content { max-width: 600px; margin: 48px auto; padding: 0 16px; text-align: center; }
        .title { color: var(--p-orange); font-size: 14px; font-weight: 800; text-transform: uppercase; letter-spacing: 1.5px; margin-bottom: 28px; }
        
        .card { border: 1px solid var(--p-border); border-radius: 8px; padding: 24px; margin-bottom: 32px; position: relative; }
        .qr-frame { border: 1px solid var(--p-border); padding: 20px; border-radius: 4px; display: inline-block; margin-bottom: 20px; background: #fff; position: relative; cursor: pointer; }
        .qr-frame img { width: 180px; height: 180px; display: block; }
        
        .fs-btn { position: absolute; top: 4px; right: 4px; background: rgba(255,255,255,0.7); backdrop-filter: blur(4px); border: 1px solid var(--p-border); border-radius: 4px; padding: 4px; cursor: pointer; }
        .fs-btn svg { width: 14px; height: 14px; stroke: var(--p-muted); stroke-width: 2.5; fill: none; }

        .id-display { display: flex; align-items: center; background: #f9fafb; border: 1px solid var(--p-border); border-radius: 4px; padding: 10px 14px; margin-bottom: 20px; cursor: pointer; transition: background 0.3s, border-color 0.3s; }
        .id-display:hover { background: #f1f5f9; }
        .copy-flash { background: var(--p-toast-bg) !important; border-color: #7dd3fc !important; }
        code { flex-grow: 1; font-family: ui-monospace, monospace; font-size: 13px; text-align: left; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; pointer-events: none; }
        
        .btn-main { background: var(--p-orange); color: white; border: none; padding: 14px; border-radius: 4px; font-weight: 600; text-transform: uppercase; cursor: pointer; width: 100%; font-size: 13px; transition: background 0.2s, transform 0.1s; }
        .btn-main:hover { background: var(--p-orange-hover); }
        .btn-main:active { transform: scale(0.98); }
        
        .history { border-top: 1px solid var(--p-border); padding-top: 24px; text-align: left; }
        .history-label { font-size: 10px; font-weight: 700; color: var(--p-muted); text-transform: uppercase; margin-bottom: 12px; display: block; text-align: center; }
        .history-row { display: flex; align-items: center; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #f3f4f6; gap: 8px; cursor: pointer; }
        .history-row:hover { background: #fafafa; }
        .h-info { display: flex; flex-direction: column; pointer-events: none; }
        .h-id { font-family: ui-monospace, monospace; font-size: 12px; color: #4b5563; }
        .h-ts { font-size: 10px; color: var(--p-muted); font-weight: 500; }
        
        .icon-btn { background: none; border: none; cursor: pointer; padding: 8px; color: var(--p-text); }
        .icon-btn svg { width: 18px; height: 18px; fill: none; stroke: currentColor; stroke-width: 2; }
        
        #modal { display: none; position: fixed; inset: 0; background: #fff; z-index: 2000; flex-direction: column; align-items: center; justify-content: center; padding: 20px; }
        #modal.show { display: flex; }
        #modal img { width: 90vw; height: 90vw; max-width: 500px; }
        .close-fs { position: absolute; top: 20px; right: 20px; font-size: 32px; cursor: pointer; }
        
        #toast { position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%) translateY(100px); background: var(--p-toast-bg); color: var(--p-toast-text); padding: 12px 24px; border-radius: 8px; font-size: 14px; font-weight: 500; transition: transform 0.3s; box-shadow: 0 4px 12px rgba(0,0,0,0.05); z-index: 3000; }
        #toast.show { transform: translateX(-50%) translateY(0); }
    </style>
</head>
<body>
    <div class="content">
        <div class="title">DFX Tag Generator</div>
        <div class="card">
            <div class="qr-frame" onclick="openFullscreen()">
                <img id="qrImg" src="data:image/png;base64,{{.QRData}}">
                <div class="fs-btn">
                    <svg viewBox="0 0 24 24"><path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/></svg>
                </div>
            </div>
            <div class="id-display" id="mainIdBox" onclick="handleCopy('currentId', this)">
                <code id="currentId">{{.UUID}}</code>
                <div class="icon-btn">
                    <svg viewBox="0 0 24 24"><path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"></path><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"></path></svg>
                </div>
            </div>
            <button class="btn-main" onclick="generate()">Generate Tag</button>
        </div>
        <div class="history">
            <span class="history-label">Recent History</span>
            <div id="historyList"></div>
        </div>
    </div>
    <div id="modal" onclick="closeFullscreen()">
        <div class="close-fs">&times;</div>
        <img id="fsImg">
    </div>
    <div id="toast">✓ ID copied to clipboard</div>
    <script>
        const copySvg = '<path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"></path><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"></path>';
        const checkSvg = '<path d="M20 6L9 17L4 12" stroke-linecap="round" stroke-linejoin="round"/>';

        function saveHistory(data) { localStorage.setItem('dfx_history_v2', JSON.stringify(data)); renderHistory(data); }
        function loadHistory() { 
            const saved = localStorage.getItem('dfx_history_v2');
            if (saved) renderHistory(JSON.parse(saved));
            else renderHistory([]);
        }

        function renderHistory(items) {
            const list = document.getElementById('historyList');
            let h = "";
            items.forEach(item => {
                h += '<div class="history-row" onclick="handleCopyText(\'' + item.uuid + '\', this)">' +
                     '<div class="h-info"><span class="h-ts">' + item.timestamp + '</span>' +
                     '<span class="h-id">' + item.uuid + '</span></div>' +
                     '<div class="icon-btn"><svg viewBox="0 0 24 24">' + copySvg + '</svg></div></div>';
            });
            list.innerHTML = h;
        }

        async function generate() {
            const res = await fetch('/api/generate');
            const data = await res.json();
            document.getElementById('qrImg').src = 'data:image/png;base64,' + data.qr_data;
            document.getElementById('currentId').innerText = data.uuid;
            saveHistory(data.history);
        }

        function handleCopy(id, el) { handleCopyText(document.getElementById(id).innerText, el); }
        function handleCopyText(text, el) {
            navigator.clipboard.writeText(text);
            showToast();
            el.classList.add('copy-flash');
            const svg = el.querySelector('svg');
            const original = svg.innerHTML;
            svg.innerHTML = checkSvg;
            setTimeout(() => { 
                el.classList.remove('copy-flash'); 
                svg.innerHTML = original;
            }, 1000);
        }

        function openFullscreen() {
            document.getElementById('fsImg').src = document.getElementById('qrImg').src;
            document.getElementById('modal').classList.add('show');
        }
        function closeFullscreen() { document.getElementById('modal').classList.remove('show'); }
        function showToast() {
            const t = document.getElementById('toast');
            t.classList.add('show');
            setTimeout(() => t.classList.remove('show'), 3000);
        }
        window.onload = loadHistory;
    </script>
</body>
</html>`
