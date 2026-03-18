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

const tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>DFX Tag Generator</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,300;9..40,400;9..40,500;9..40,600;9..40,700;9..40,800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        *,*::before,*::after{box-sizing:border-box;margin:0;padding:0}

        :root {
            --bg: #fafafa;
            --bg-card: #ffffff;
            --bg-header: rgba(255,255,255,0.88);
            --bg-input: #f5f5f5;
            --text: #1a1a1a;
            --text-sec: #555;
            --text-muted: #999;
            --border: #e5e7eb;
            --border-hover: #d0d0d0;
            --accent: #e85d2a;
            --accent-hover: #ff7636;
            --accent-light: rgba(232,93,42,0.08);
            --shadow-card: 0 1px 4px rgba(0,0,0,0.04), 0 0 0 1px rgba(0,0,0,0.04);
            --toast-bg: #e0f2fe;
            --toast-text: #0369a1;
            --modal-bg: rgba(255,255,255,0.96);
            --qr-border: #e5e7eb;
            --qr-bg: #ffffff;
            --hist-hover: #f7f7f8;
            --hist-border: #f3f4f6;
            --copy-flash-bg: #e0f2fe;
            --copy-flash-border: #7dd3fc;
            --icon-stroke: #444;
            --tr: 0.35s;
            color-scheme: light;
        }

        [data-theme="dark"] {
            --bg: #1a1a1e;
            --bg-card: #242428;
            --bg-header: rgba(26,26,30,0.9);
            --bg-input: #2a2a2f;
            --text: #e8e8ea;
            --text-sec: #a0a0a5;
            --text-muted: #666669;
            --border: #323236;
            --border-hover: #434348;
            --accent: #f06b35;
            --accent-hover: #ff8a50;
            --accent-light: rgba(240,107,53,0.1);
            --shadow-card: 0 1px 4px rgba(0,0,0,0.2), 0 0 0 1px rgba(255,255,255,0.04);
            --toast-bg: #1e3a5f;
            --toast-text: #7dd3fc;
            --modal-bg: rgba(26,26,30,0.96);
            --qr-border: #323236;
            --qr-bg: #ffffff;
            --hist-hover: #2a2a2f;
            --hist-border: #323236;
            --copy-flash-bg: rgba(30,58,95,0.4);
            --copy-flash-border: #2563eb;
            --icon-stroke: #b0b0b4;
            color-scheme: dark;
        }

        html { font-size: 16px; }

        body {
            font-family: 'DM Sans', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg);
            color: var(--text);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            -webkit-font-smoothing: antialiased;
            transition: background var(--tr) ease, color var(--tr) ease;
            -webkit-tap-highlight-color: transparent;
        }

        /* ── Header ── */
        .header {
            position: sticky;
            top: 0;
            z-index: 100;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 1.25rem;
            height: 56px;
            background: var(--bg-header);
            backdrop-filter: blur(20px) saturate(1.4);
            -webkit-backdrop-filter: blur(20px) saturate(1.4);
            border-bottom: 1px solid var(--border);
            transition: background var(--tr) ease, border-color var(--tr) ease;
        }

        .header-left {
            display: flex;
            align-items: center;
            gap: 0.6rem;
        }

        .dfx-logo { width: 34px; height: 34px; flex-shrink: 0; }
        .dfx-logo svg { width: 100%; height: 100%; display: block; }

        .app-title {
            font-weight: 700;
            font-size: 0.88rem;
            color: var(--accent);
            letter-spacing: -0.01em;
            transition: color var(--tr) ease;
        }

        /* ── Theme Toggle: minimalist sun/moon ── */
        .theme-btn {
            width: 36px;
            height: 36px;
            border: none;
            background: transparent;
            cursor: pointer;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 8px;
            color: var(--text-muted);
            transition: color 0.25s ease, background 0.25s ease;
            position: relative;
        }

        .theme-btn:hover { color: var(--accent); background: var(--accent-light); }

        .theme-btn svg {
            width: 20px;
            height: 20px;
            position: absolute;
            transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease;
        }

        .theme-btn .icon-sun { opacity: 1; transform: rotate(0deg) scale(1); }
        .theme-btn .icon-moon { opacity: 0; transform: rotate(-60deg) scale(0.5); }

        [data-theme="dark"] .theme-btn .icon-sun { opacity: 0; transform: rotate(90deg) scale(0.5); }
        [data-theme="dark"] .theme-btn .icon-moon { opacity: 1; transform: rotate(0deg) scale(1); }

        /* ── Content ── */
        .content {
            flex: 1;
            max-width: 560px;
            width: 100%;
            margin: 0 auto;
            padding: 2.5rem 1rem 3rem;
            text-align: center;
        }

        .page-title {
            color: var(--accent);
            font-size: 0.72rem;
            font-weight: 800;
            text-transform: uppercase;
            letter-spacing: 0.15em;
            margin-bottom: 1.75rem;
            transition: color var(--tr) ease;
        }

        .card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 14px;
            padding: 1.75rem 1.5rem;
            margin-bottom: 2rem;
            box-shadow: var(--shadow-card);
            transition: background var(--tr) ease, border-color var(--tr) ease, box-shadow var(--tr) ease;
        }

        .qr-frame {
            border: 1.5px solid var(--qr-border);
            padding: 1rem;
            border-radius: 10px;
            display: inline-block;
            margin-bottom: 1.25rem;
            background: var(--qr-bg);
            position: relative;
            cursor: pointer;
            transition: border-color 0.25s ease, box-shadow 0.25s ease;
        }

        .qr-frame:hover { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-light); }
        .qr-frame img { width: 180px; height: 180px; display: block; border-radius: 4px; }

        .fs-btn {
            position: absolute;
            top: 8px;
            right: 8px;
            background: rgba(255,255,255,0.85);
            backdrop-filter: blur(6px);
            border: 1px solid rgba(0,0,0,0.06);
            border-radius: 6px;
            padding: 5px;
            cursor: pointer;
            transition: background 0.2s ease;
        }

        .fs-btn:hover { background: rgba(255,255,255,0.95); }
        .fs-btn svg { width: 14px; height: 14px; stroke: #666; stroke-width: 2.5; fill: none; display: block; }

        .id-display {
            display: flex;
            align-items: center;
            background: var(--bg-input);
            border: 1.5px solid var(--border);
            border-radius: 10px;
            padding: 0.7rem 0.9rem;
            margin-bottom: 1.25rem;
            cursor: pointer;
            transition: background 0.25s ease, border-color 0.25s ease;
            gap: 0.5rem;
        }

        .id-display:hover { border-color: var(--border-hover); }
        .id-display.copy-flash { background: var(--copy-flash-bg) !important; border-color: var(--copy-flash-border) !important; }

        .id-display code {
            flex: 1;
            font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.8rem;
            text-align: left;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            pointer-events: none;
            color: var(--text-sec);
            transition: color var(--tr) ease;
        }

        .btn-main {
            background: var(--accent);
            color: #fff;
            border: none;
            padding: 0.85rem;
            border-radius: 10px;
            font-family: inherit;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            cursor: pointer;
            width: 100%;
            font-size: 0.8rem;
            transition: background 0.2s ease, transform 0.1s ease, box-shadow 0.2s ease;
        }

        .btn-main:hover { background: var(--accent-hover); box-shadow: 0 4px 16px rgba(232,93,42,0.25); }
        .btn-main:active { transform: scale(0.98); }

        /* ── History ── */
        .history {
            border-top: 1px solid var(--border);
            padding-top: 1.5rem;
            text-align: left;
            transition: border-color var(--tr) ease;
        }

        .history-label {
            font-size: 0.62rem;
            font-weight: 700;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.12em;
            margin-bottom: 0.75rem;
            display: block;
            text-align: center;
            transition: color var(--tr) ease;
        }

        .history-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0.75rem 0.5rem;
            border-bottom: 1px solid var(--hist-border);
            gap: 0.5rem;
            cursor: pointer;
            border-radius: 8px;
            margin: 0 -0.5rem;
            transition: background 0.2s ease, border-color var(--tr) ease;
        }

        .history-row:hover { background: var(--hist-hover); }
        .history-row:last-child { border-bottom: none; }
        .history-row.copy-flash { background: var(--copy-flash-bg); }

        .h-info { display: flex; flex-direction: column; pointer-events: none; gap: 2px; min-width: 0; flex: 1; }

        .h-id {
            font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.75rem;
            color: var(--text-sec);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            transition: color var(--tr) ease;
        }

        .h-ts {
            font-size: 0.62rem;
            color: var(--text-muted);
            font-weight: 600;
            letter-spacing: 0.02em;
            transition: color var(--tr) ease;
        }

        .icon-btn {
            background: none;
            border: none;
            cursor: pointer;
            padding: 6px;
            color: var(--text-muted);
            border-radius: 6px;
            transition: color 0.2s ease, background 0.2s ease;
            flex-shrink: 0;
        }

        .icon-btn:hover { color: var(--accent); background: var(--accent-light); }
        .icon-btn svg { width: 16px; height: 16px; fill: none; stroke: currentColor; stroke-width: 2; display: block; }

        /* ── Modal ── */
        #modal {
            display: none;
            position: fixed;
            inset: 0;
            background: var(--modal-bg);
            backdrop-filter: blur(24px);
            -webkit-backdrop-filter: blur(24px);
            z-index: 2000;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 20px;
            transition: background var(--tr) ease;
        }

        #modal.show { display: flex; }

        #modal img {
            width: 80vw;
            max-width: 420px;
            max-height: 420px;
            object-fit: contain;
            aspect-ratio: 1/1;
            image-rendering: pixelated;
            border-radius: 12px;
            box-shadow: 0 8px 40px rgba(0,0,0,0.1);
        }

        .close-fs {
            position: absolute;
            top: 1rem;
            right: 1.25rem;
            width: 40px;
            height: 40px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            background: var(--bg-card);
            border: 1px solid var(--border);
            cursor: pointer;
            font-size: 1.25rem;
            color: var(--text-muted);
            transition: all 0.2s ease;
        }

        .close-fs:hover { color: var(--text); border-color: var(--border-hover); }

        /* ── Toast ── */
        #toast {
            position: fixed;
            bottom: 1.5rem;
            left: 50%;
            transform: translateX(-50%) translateY(80px);
            background: var(--toast-bg);
            color: var(--toast-text);
            padding: 0.65rem 1.25rem;
            border-radius: 10px;
            font-size: 0.8rem;
            font-weight: 600;
            opacity: 0;
            transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1),
                        opacity 0.35s ease,
                        background var(--tr) ease,
                        color var(--tr) ease;
            box-shadow: 0 4px 20px rgba(0,0,0,0.12);
            z-index: 3000;
            white-space: nowrap;
            pointer-events: none;
        }

        #toast.show { transform: translateX(-50%) translateY(0); opacity: 1; }

        /* ── Animations ── */
        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(14px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .page-title { animation: fadeInUp 0.4s ease both; }
        .card { animation: fadeInUp 0.45s ease 0.05s both; }
        .history { animation: fadeInUp 0.45s ease 0.1s both; }

        @media (max-width: 768px) {
            .content { padding: 2rem 1rem 2.5rem; }
        }

        @media (max-width: 480px) {
            .header { padding: 0 0.75rem; height: 50px; }
            .dfx-logo { width: 28px; height: 28px; }
            .app-title { font-size: 0.82rem; }
            .theme-btn { width: 32px; height: 32px; }
            .theme-btn svg { width: 18px; height: 18px; }
            .content { padding: 1.5rem 0.85rem 2rem; }
            .card { padding: 1.25rem 1rem; border-radius: 12px; }
            .qr-frame { padding: 0.75rem; }
            .qr-frame img { width: 150px; height: 150px; }
            .id-display { padding: 0.6rem 0.75rem; }
            .id-display code { font-size: 0.7rem; }
            .btn-main { padding: 0.75rem; font-size: 0.75rem; border-radius: 8px; }
            .h-id { font-size: 0.68rem; }
            .h-ts { font-size: 0.58rem; }
            .page-title { font-size: 0.65rem; margin-bottom: 1.25rem; }
        }

        @media (max-width: 340px) {
            .qr-frame img { width: 130px; height: 130px; }
            .id-display code { font-size: 0.62rem; }
        }
    </style>
</head>
<body>

    <header class="header">
        <div class="header-left">
            <div class="dfx-logo">
                <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="6" y="6" width="36" height="36" rx="4" stroke="var(--icon-stroke)" stroke-width="1.8"/>
                    <rect x="10" y="10" width="10" height="10" rx="1.5" stroke="var(--accent)" stroke-width="1.8"/>
                    <rect x="13" y="13" width="4" height="4" rx="0.5" fill="var(--accent)"/>
                    <rect x="28" y="10" width="10" height="10" rx="1.5" stroke="var(--accent)" stroke-width="1.8"/>
                    <rect x="31" y="13" width="4" height="4" rx="0.5" fill="var(--accent)"/>
                    <rect x="10" y="28" width="10" height="10" rx="1.5" stroke="var(--accent)" stroke-width="1.8"/>
                    <rect x="13" y="31" width="4" height="4" rx="0.5" fill="var(--accent)"/>
                    <rect x="24" y="24" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.5"/>
                    <rect x="29" y="24" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.35"/>
                    <rect x="34" y="29" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.5"/>
                    <rect x="29" y="34" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.35"/>
                    <rect x="24" y="29" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.25"/>
                    <rect x="34" y="34" width="3" height="3" rx="0.5" fill="var(--icon-stroke)" opacity="0.4"/>
                </svg>
            </div>
            <span class="app-title">DFX Tag Generator</span>
        </div>

        <button class="theme-btn" id="themeToggle" aria-label="Toggle dark mode">
            <svg class="icon-sun" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="5"/>
                <line x1="12" y1="1" x2="12" y2="3"/>
                <line x1="12" y1="21" x2="12" y2="23"/>
                <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
                <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
                <line x1="1" y1="12" x2="3" y2="12"/>
                <line x1="21" y1="12" x2="23" y2="12"/>
                <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
                <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <svg class="icon-moon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
        </button>
    </header>

    <div class="content">
        <div class="page-title">DFX Tag Generator</div>
        <div class="card">
            <div class="qr-frame" onclick="openFullscreen()">
                <img id="qrImg" src="data:image/png;base64,{{.QRData}}" alt="QR Code">
                <div class="fs-btn">
                    <svg viewBox="0 0 24 24"><path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/></svg>
                </div>
            </div>
            <div class="id-display" id="mainIdBox" onclick="handleCopy('currentId', this)">
                <code id="currentId">{{.UUID}}</code>
                <div class="icon-btn">
                    <svg viewBox="0 0 24 24"><path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"/><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"/></svg>
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
        <img id="fsImg" alt="QR Code Fullscreen">
    </div>

    <div id="toast">Copied to clipboard</div>

    <script>
        /* ── Theme ── */
        var html = document.documentElement;
        var toggleBtn = document.getElementById('themeToggle');

        function setTheme(t) {
            html.setAttribute('data-theme', t);
            localStorage.setItem('dfx-theme', t);
        }

        (function() {
            var manual = localStorage.getItem('dfx-theme-manual');
            var saved = localStorage.getItem('dfx-theme');
            if (manual && saved) {
                setTheme(saved);
            } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                setTheme('dark');
            } else {
                setTheme('light');
            }
        })();

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
            if (!localStorage.getItem('dfx-theme-manual')) {
                setTheme(e.matches ? 'dark' : 'light');
            }
        });

        toggleBtn.addEventListener('click', function() {
            var next = html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
            setTheme(next);
            localStorage.setItem('dfx-theme-manual', 'true');
        });

        /* ── SVG icons for copy state ── */
        var checkSvg = '<path d="M20 6L9 17L4 12" stroke-linecap="round" stroke-linejoin="round"/>';
        var copySvg = '<path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"/><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"/>';

        /* ─────────────────────────────────────────────
           Toast: single global timer.
           Every call resets the countdown to a full
           2.5 s so rapid copies never cut it short.
           ───────────────────────────────────────────── */
        var toastTimer = null;
        var toastEl = document.getElementById('toast');

        function showToast() {
            if (toastTimer) clearTimeout(toastTimer);

            // If already visible just reset the timer — no flicker
            if (!toastEl.classList.contains('show')) {
                toastEl.classList.add('show');
            }

            toastTimer = setTimeout(function() {
                toastEl.classList.remove('show');
                toastTimer = null;
            }, 2500);
        }

        /* ─────────────────────────────────────────────
           Flash: EXCLUSIVE — only one element can be
           flashed at a time. A new copy cancels the
           previous flash immediately.
           ───────────────────────────────────────────── */
        var activeFlash = null;   // { el, svg, timer }

        function clearActiveFlash() {
            if (!activeFlash) return;
            clearTimeout(activeFlash.timer);
            activeFlash.el.classList.remove('copy-flash');
            if (activeFlash.svg) activeFlash.svg.innerHTML = copySvg;
            activeFlash = null;
        }

        function flashElement(el) {
            // Cancel any existing flash first
            clearActiveFlash();

            el.classList.add('copy-flash');
            var svg = el.querySelector('.icon-btn svg') || el.querySelector('svg');
            if (svg) svg.innerHTML = checkSvg;

            var timer = setTimeout(function() {
                el.classList.remove('copy-flash');
                if (svg) svg.innerHTML = copySvg;
                activeFlash = null;
            }, 1200);

            activeFlash = { el: el, svg: svg, timer: timer };
        }

        /* ── Copy ── */
        function handleCopy(id, el) {
            handleCopyText(document.getElementById(id).innerText, el);
        }

        function handleCopyText(text, el) {
            if (navigator.clipboard && navigator.clipboard.writeText) {
                navigator.clipboard.writeText(text).then(function() {
                    showToast();
                    flashElement(el);
                }).catch(function() { fallbackCopy(text, el); });
            } else {
                fallbackCopy(text, el);
            }
        }

        function fallbackCopy(text, el) {
            var ta = document.createElement('textarea');
            ta.value = text;
            ta.style.cssText = 'position:fixed;left:-9999px';
            document.body.appendChild(ta);
            ta.focus();
            ta.select();
            try {
                if (document.execCommand('copy')) {
                    showToast();
                    flashElement(el);
                }
            } catch(e) { console.error(e); }
            document.body.removeChild(ta);
        }

        /* ── History ── */
        function saveHistory(data) {
            localStorage.setItem('dfx_history_v3', JSON.stringify(data));
            renderHistory(data);
        }

        function loadHistory() {
            var saved = localStorage.getItem('dfx_history_v3');
            if (saved) renderHistory(JSON.parse(saved));
            else renderHistory([]);
        }

        function renderHistory(items) {
            var list = document.getElementById('historyList');
            var h = '';
            items.forEach(function(item) {
                h += '<div class="history-row" onclick="handleCopyText(\'' + item.uuid + '\', this)">' +
                     '<div class="h-info"><span class="h-ts">' + item.timestamp + '</span>' +
                     '<span class="h-id">' + item.uuid + '</span></div>' +
                     '<div class="icon-btn"><svg viewBox="0 0 24 24">' + copySvg + '</svg></div></div>';
            });
            list.innerHTML = h;
        }

        /* ── Generate ── */
        async function generate() {
            var res = await fetch('/api/generate');
            var data = await res.json();
            document.getElementById('qrImg').src = 'data:image/png;base64,' + data.qr_data;
            document.getElementById('currentId').innerText = data.uuid;
            saveHistory(data.history);
        }

        /* ── Fullscreen QR ── */
        function openFullscreen() {
            document.getElementById('fsImg').src = document.getElementById('qrImg').src;
            document.getElementById('modal').classList.add('show');
        }

        function closeFullscreen() {
            document.getElementById('modal').classList.remove('show');
        }

        window.onload = loadHistory;
    </script>
</body>
</html>`
