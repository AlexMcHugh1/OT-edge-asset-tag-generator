package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// privacyHandler serves the UK GDPR privacy notice.
func privacyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("privacy").Parse(privacyTmpl))
	t.Execute(w, nil)
}

// erasureRequestHandler accepts a right-to-erasure request without exposing
// any contact details publicly. Requests are stored for manual review.
func erasureRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Reason   string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	username := strings.TrimSpace(strings.ToLower(req.Username))
	if idx := strings.Index(username, "@"); idx != -1 {
		username = username[:idx]
	}
	if username == "" {
		jsonError(w, "Username is required", 400)
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if len(reason) > 1000 {
		reason = reason[:1000]
	}
	if _, err := db.Exec(
		"INSERT INTO erasure_requests (username, reason) VALUES (?, ?)",
		username, reason,
	); err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// deleteAccountHandler permanently removes the authenticated user's account,
// all their sessions, and all their device records. Password confirmation is
// required as an additional safeguard.
func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromRequest(r)
	if user == nil {
		jsonError(w, "Authentication required", 401)
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", 400)
		return
	}
	if req.Password == "" {
		jsonError(w, "Password confirmation is required", 400)
		return
	}
	var storedHash string
	if err := db.QueryRow("SELECT password_hash FROM users WHERE id = ?", user.ID).
		Scan(&storedHash); err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		jsonError(w, "Incorrect password", 403)
		return
	}
	// Deleting the user cascades to sessions and devices via FK constraints.
	if _, err := db.Exec("DELETE FROM users WHERE id = ?", user.ID); err != nil {
		jsonError(w, "Server error", 500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "dfx_session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// startInactiveAccountCleanup launches a background goroutine that runs once
// per day at 03:00 and permanently deletes any user account that has been
// inactive for more than 12 months, in accordance with the retention policy
// stated in the privacy notice.
func startInactiveAccountCleanup() {
	go func() {
		for {
			deleteInactiveUsers()
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 3, 0, 0, 0, now.Location())
			time.Sleep(time.Until(next))
		}
	}()
}

const privacyTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Privacy Notice — DFX Tag Generator</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700&display=swap" rel="stylesheet">
    <style>
        *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
        :root {
            --bg: #fafafa; --bg-card: #ffffff; --text: #1a1a1a; --text-sec: #555;
            --text-muted: #888; --border: #e5e7eb; --accent: #e85d2a;
            --accent-light: rgba(232,93,42,0.08); --error-bg: #fef2f2; --error-text: #dc2626;
            --success-bg: #f0fdf4; --success-text: #166534;
            color-scheme: light;
        }
        @media (prefers-color-scheme: dark) {
            :root {
                --bg: #1a1a1e; --bg-card: #242428; --text: #e8e8ea; --text-sec: #a0a0a5;
                --text-muted: #666669; --border: #323236; --accent: #f06b35;
                --accent-light: rgba(240,107,53,0.1); --error-bg: rgba(220,38,38,0.1);
                --error-text: #f87171; --success-bg: rgba(22,163,74,0.1);
                --success-text: #4ade80;
                color-scheme: dark;
            }
        }
        body {
            font-family: 'DM Sans', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg); color: var(--text);
            line-height: 1.7; -webkit-font-smoothing: antialiased;
        }
        .page { max-width: 740px; margin: 0 auto; padding: 3rem 1.5rem 5rem; }
        .back {
            display: inline-flex; align-items: center; gap: 0.4rem;
            color: var(--accent); text-decoration: none; font-size: 0.85rem;
            font-weight: 600; margin-bottom: 2.5rem;
        }
        .back:hover { opacity: 0.8; }
        h1 { font-size: 1.8rem; font-weight: 800; color: var(--accent); margin-bottom: 0.3rem; }
        .meta { font-size: 0.8rem; color: var(--text-muted); margin-bottom: 2.5rem; }
        h2 { font-size: 1.05rem; font-weight: 700; margin: 2rem 0 0.5rem; }
        p { color: var(--text-sec); margin-bottom: 0.75rem; font-size: 0.93rem; }
        ul { color: var(--text-sec); font-size: 0.93rem; padding-left: 1.4rem; margin-bottom: 0.75rem; }
        li { margin-bottom: 0.3rem; }
        .card {
            background: var(--bg-card); border: 1px solid var(--border);
            border-radius: 12px; padding: 1.75rem; margin-top: 3rem;
        }
        .card h2 { margin-top: 0; }
        .form-group { margin-bottom: 1rem; }
        label { display: block; font-size: 0.82rem; font-weight: 600; margin-bottom: 0.35rem; }
        input, textarea {
            width: 100%; padding: 0.6rem 0.8rem; border: 1px solid var(--border);
            border-radius: 8px; background: var(--bg); color: var(--text);
            font-family: inherit; font-size: 0.9rem; outline: none;
            transition: border-color 0.2s ease;
        }
        input:focus, textarea:focus { border-color: var(--accent); }
        textarea { resize: vertical; min-height: 90px; }
        .btn {
            padding: 0.65rem 1.5rem; background: var(--accent); color: #fff;
            border: none; border-radius: 8px; font-family: inherit; font-size: 0.88rem;
            font-weight: 700; cursor: pointer; transition: opacity 0.2s ease;
        }
        .btn:hover { opacity: 0.88; }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; }
        .msg { padding: 0.75rem 1rem; border-radius: 8px; font-size: 0.88rem; margin-top: 0.75rem; display: none; }
        .msg.error { background: var(--error-bg); color: var(--error-text); display: block; }
        .msg.success { background: var(--success-bg); color: var(--success-text); display: block; }
        hr { border: none; border-top: 1px solid var(--border); margin: 2rem 0; }
        a { color: var(--accent); }
    </style>
</head>
<body>
<div class="page">
    <a class="back" href="/">&#8592; Back to DFX Tag Generator</a>

    <h1>Privacy Notice</h1>
    <p class="meta">DFX Tag Generator &mdash; Last updated: April 2026</p>

    <h2>1. About this tool</h2>
    <p>DFX Tag Generator is an internal tool for creating and tracking asset tags for OT/IT devices. It is not affiliated with or officially operated by any organisation. For data queries, use the contact form at the bottom of this page.</p>

    <h2>2. What data we store</h2>
    <ul>
        <li><strong>Username</strong> — a short identifier you choose at registration. No email addresses are stored.</li>
        <li><strong>Password</strong> — stored as a bcrypt hash (cost&nbsp;12). Your plaintext password is never stored or logged.</li>
        <li><strong>Session token</strong> — a cryptographically random token stored in an HttpOnly cookie for up to 30 days.</li>
        <li><strong>Device records</strong> — device name, serial number, environment, location, and the generated asset tag UUID you choose to save.</li>
        <li><strong>Account timestamps</strong> — when your account was created and when it was last active.</li>
    </ul>

    <h2>3. What we do not collect</h2>
    <ul>
        <li>No email addresses.</li>
        <li>No tracking cookies or third-party analytics.</li>
        <li>No personal data is shared with any third party for marketing purposes.</li>
    </ul>

    <h2>4. Data retention</h2>
    <p>Your account and all associated data (devices, sessions) are permanently deleted if your account has been <strong>inactive for 12 consecutive months</strong>. This runs automatically. You may also delete your account at any time from within the application.</p>

    <h2>5. Infrastructure</h2>
    <p>Network traffic is proxied through <strong>Cloudflare</strong>, which may process IP addresses and request metadata to provide DDoS protection and TLS termination. All application data is stored on self-hosted infrastructure.</p>

    <h2>6. Your rights</h2>
    <p>You can delete your account and all associated data at any time using the "Delete Account" option in the app. To request manual deletion or a copy of your data, use the form below.</p>

    <h2>7. Security</h2>
    <p>Passwords are hashed with bcrypt (cost 12). Session tokens use a cryptographically secure RNG. Cookies are HttpOnly and SameSite=Lax. All connections are TLS-encrypted via Cloudflare.</p>

    <hr>

    <div class="card">
        <h2>Data Request</h2>
        <p style="margin-bottom:1.25rem">Use this form to request deletion of your data or a copy of what is stored against your account.</p>
        <form id="erasureForm" onsubmit="submitErasure(event)">
            <div class="form-group">
                <label for="erUsername">Your username</label>
                <input type="text" id="erUsername" name="username" placeholder="your username" required autocomplete="username">
            </div>
            <div class="form-group">
                <label for="erReason">Request details</label>
                <textarea id="erReason" name="reason" placeholder="e.g. Please delete all data associated with my account." maxlength="1000"></textarea>
            </div>
            <button type="submit" class="btn" id="erBtn">Submit Request</button>
            <div class="msg" id="erMsg"></div>
        </form>
    </div>
</div>
<script>
    async function submitErasure(e) {
        e.preventDefault();
        var username = document.getElementById('erUsername').value.trim();
        var reason   = document.getElementById('erReason').value.trim();
        var btn      = document.getElementById('erBtn');
        var msg      = document.getElementById('erMsg');
        msg.className = 'msg';
        msg.textContent = '';
        btn.disabled = true;
        btn.textContent = '...';
        try {
            var res  = await fetch('/api/privacy/erasure-request', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: username, reason: reason })
            });
            var data = await res.json();
            if (!res.ok) {
                msg.textContent = data.error || 'Something went wrong. Please try again.';
                msg.className = 'msg error';
            } else {
                msg.textContent = 'Request received.';
                msg.className = 'msg success';
                document.getElementById('erasureForm').reset();
            }
        } catch(err) {
            msg.textContent = 'Connection error. Please try again.';
            msg.className = 'msg error';
        } finally {
            btn.disabled = false;
            btn.textContent = 'Submit Request';
        }
    }
</script>
</body>
</html>`
