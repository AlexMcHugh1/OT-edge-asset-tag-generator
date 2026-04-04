package main

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
            --error-text: #dc2626;
            --error-bg: #fef2f2;
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
            --error-text: #f87171;
            --error-bg: rgba(220,38,38,0.1);
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
            position: sticky; top: 0; z-index: 100;
            display: flex; align-items: center; justify-content: space-between;
            padding: 0 1.25rem; height: 56px;
            background: var(--bg-header);
            backdrop-filter: blur(20px) saturate(1.4);
            -webkit-backdrop-filter: blur(20px) saturate(1.4);
            border-bottom: 1px solid var(--border);
            transition: background var(--tr) ease, border-color var(--tr) ease;
        }

        .header-left { display: flex; align-items: center; gap: 0.6rem; }
        .header-right { display: flex; align-items: center; gap: 0.5rem; }
        .dfx-logo { width: 34px; height: 34px; flex-shrink: 0; }
        .dfx-logo svg { width: 100%; height: 100%; display: block; }

        .app-title {
            font-weight: 700; font-size: 0.88rem; color: var(--accent);
            letter-spacing: -0.01em; transition: color var(--tr) ease;
        }

        /* ── Auth Area ── */
        .auth-btn {
            height: 32px; padding: 0 0.9rem;
            background: var(--accent); color: #fff; border: none;
            border-radius: 7px; font-family: inherit; font-size: 0.75rem;
            font-weight: 700; cursor: pointer; white-space: nowrap;
            transition: background 0.2s ease;
        }
        .auth-btn:hover { background: var(--accent-hover); }

        .user-pill {
            display: flex; align-items: center; gap: 0.35rem;
            background: var(--accent-light); border: 1px solid var(--accent);
            border-radius: 8px; padding: 0.2rem 0.4rem 0.2rem 0.65rem;
            color: var(--accent); font-size: 0.7rem; font-weight: 600;
        }
        .user-email {
            max-width: 130px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
        }
        .signout-btn {
            background: none; border: none; cursor: pointer; padding: 3px;
            color: var(--accent); display: flex; align-items: center;
            border-radius: 4px; transition: background 0.2s ease; flex-shrink: 0;
        }
        .signout-btn:hover { background: rgba(0,0,0,0.06); }
        .signout-btn svg { width: 14px; height: 14px; stroke: currentColor; stroke-width: 2.5; fill: none; display: block; }

        /* ── Theme Toggle ── */
        .theme-btn {
            width: 36px; height: 36px; border: none; background: transparent;
            cursor: pointer; padding: 0; display: flex; align-items: center;
            justify-content: center; border-radius: 8px; color: var(--text-muted);
            position: relative; transition: color 0.25s ease, background 0.25s ease;
        }
        .theme-btn:hover { color: var(--accent); background: var(--accent-light); }
        .theme-btn svg {
            width: 20px; height: 20px; position: absolute;
            transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease;
        }
        .theme-btn .icon-sun { opacity: 1; transform: rotate(0deg) scale(1); }
        .theme-btn .icon-moon { opacity: 0; transform: rotate(-60deg) scale(0.5); }
        [data-theme="dark"] .theme-btn .icon-sun { opacity: 0; transform: rotate(90deg) scale(0.5); }
        [data-theme="dark"] .theme-btn .icon-moon { opacity: 1; transform: rotate(0deg) scale(1); }

        /* ── Content ── */
        .content {
            flex: 1; max-width: 560px; width: 100%;
            margin: 0 auto; padding: 2.5rem 1rem 3rem; text-align: center;
        }

        .page-title {
            color: var(--accent); font-size: 0.72rem; font-weight: 800;
            text-transform: uppercase; letter-spacing: 0.15em; margin-bottom: 1.75rem;
            transition: color var(--tr) ease;
        }

        .card {
            background: var(--bg-card); border: 1px solid var(--border);
            border-radius: 14px; padding: 1.75rem 1.5rem; margin-bottom: 2rem;
            box-shadow: var(--shadow-card);
            transition: background var(--tr) ease, border-color var(--tr) ease, box-shadow var(--tr) ease;
        }

        .qr-frame {
            border: 1.5px solid var(--qr-border); padding: 1rem; border-radius: 10px;
            display: inline-block; margin-bottom: 1.25rem; background: var(--qr-bg);
            position: relative; cursor: pointer;
            transition: border-color 0.25s ease, box-shadow 0.25s ease;
        }
        .qr-frame:hover { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-light); }
        .qr-frame img { width: 180px; height: 180px; display: block; border-radius: 4px; }

        .fs-btn {
            position: absolute; top: 8px; right: 8px;
            background: rgba(255,255,255,0.85); backdrop-filter: blur(6px);
            border: 1px solid rgba(0,0,0,0.06); border-radius: 6px;
            padding: 5px; cursor: pointer; transition: background 0.2s ease;
        }
        .fs-btn:hover { background: rgba(255,255,255,0.95); }
        .fs-btn svg { width: 14px; height: 14px; stroke: #666; stroke-width: 2.5; fill: none; display: block; }

        /* ── ID Display ── */
        .id-display {
            display: flex; align-items: center;
            background: var(--bg-input); border: 1.5px solid var(--border);
            border-radius: 10px; padding: 0.7rem 0.9rem; margin-bottom: 1.25rem;
            transition: background 0.25s ease, border-color 0.25s ease; gap: 0.5rem;
        }
        .id-display:hover { border-color: var(--border-hover); }
        .id-display.copy-flash { background: var(--copy-flash-bg) !important; border-color: var(--copy-flash-border) !important; }
        .id-display code {
            flex: 1; font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.8rem; text-align: left; overflow: hidden;
            text-overflow: ellipsis; white-space: nowrap; color: var(--text-sec);
            cursor: pointer; transition: color var(--tr) ease;
        }
        .edit-input {
            flex: 1; font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.8rem; color: var(--text); background: transparent;
            border: none; outline: none; padding: 0; margin: 0; min-width: 0;
        }
        .id-display.editing { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-light); }
        .id-display.editing .edit-only,
        .id-display.editing .copy-only { display: none; }
        .id-display.edit-error-state { border-color: var(--error-text); box-shadow: 0 0 0 3px rgba(220,38,38,0.1); }

        /* ── Error bar ── */
        .edit-error {
            display: none; align-items: center; justify-content: space-between;
            background: var(--error-bg); border-radius: 8px;
            padding: 0.5rem 0.75rem; margin-top: -0.85rem; margin-bottom: 1rem; gap: 0.5rem;
        }
        .edit-error.show { display: flex; }
        .edit-error-text { font-size: 0.72rem; font-weight: 600; color: var(--error-text); }

        .regen-btn {
            display: flex; align-items: center; gap: 0.3rem;
            background: none; border: 1.5px solid var(--error-text); color: var(--error-text);
            border-radius: 6px; padding: 0.3rem 0.6rem; cursor: pointer;
            font-family: inherit; font-size: 0.65rem; font-weight: 700;
            text-transform: uppercase; letter-spacing: 0.04em;
            transition: background 0.2s ease, color 0.2s ease;
            white-space: nowrap; flex-shrink: 0;
        }
        .regen-btn:hover { background: var(--error-text); color: #fff; }
        .regen-btn svg { width: 12px; height: 12px; flex-shrink: 0; }

        /* ── Buttons ── */
        .btn-main {
            background: var(--accent); color: #fff; border: none;
            padding: 0.85rem; border-radius: 10px; font-family: inherit;
            font-weight: 700; text-transform: uppercase; letter-spacing: 0.06em;
            cursor: pointer; width: 100%; font-size: 0.8rem;
            transition: background 0.2s ease, transform 0.1s ease, box-shadow 0.2s ease;
        }
        .btn-main:hover { background: var(--accent-hover); box-shadow: 0 4px 16px rgba(232,93,42,0.25); }
        .btn-main:active { transform: scale(0.98); }
        .btn-main:disabled { opacity: 0.6; cursor: not-allowed; transform: none; }

        .btn-secondary {
            background: transparent; color: var(--accent); border: 1.5px solid var(--accent);
            padding: 0.7rem; border-radius: 10px; font-family: inherit;
            font-weight: 700; text-transform: uppercase; letter-spacing: 0.06em;
            cursor: pointer; width: 100%; font-size: 0.8rem; margin-top: 0.6rem;
            transition: background 0.2s ease, color 0.2s ease, transform 0.1s ease;
        }
        .btn-secondary:hover { background: var(--accent-light); }
        .btn-secondary:active { transform: scale(0.98); }

        /* ── Tabs ── */
        .tab-bar {
            display: flex; border-bottom: 1px solid var(--border);
            margin-bottom: 1.5rem;
            transition: border-color var(--tr) ease;
        }
        .tab-btn {
            background: none; border: none; cursor: pointer;
            font-family: inherit; font-size: 0.68rem; font-weight: 700;
            color: var(--text-muted); padding: 0.5rem 0.85rem;
            text-transform: uppercase; letter-spacing: 0.1em;
            border-bottom: 2px solid transparent; margin-bottom: -1px;
            transition: color 0.2s ease, border-color 0.2s ease;
        }
        .tab-btn.active { color: var(--accent); border-bottom-color: var(--accent); }
        .tab-btn:hover:not(.active) { color: var(--text-sec); }

        /* ── History ── */
        .history { text-align: left; transition: border-color var(--tr) ease; }

        .history-label {
            font-size: 0.62rem; font-weight: 700; color: var(--text-muted);
            text-transform: uppercase; letter-spacing: 0.12em; margin-bottom: 0.75rem;
            display: block; text-align: center; transition: color var(--tr) ease;
        }
        .history-row {
            display: flex; align-items: center; justify-content: space-between;
            padding: 0.75rem 0.5rem; border-bottom: 1px solid var(--hist-border);
            gap: 0.5rem; cursor: pointer; border-radius: 8px; margin: 0 -0.5rem;
            transition: background 0.2s ease, border-color var(--tr) ease;
        }
        .history-row:hover { background: var(--hist-hover); }
        .history-row:last-child { border-bottom: none; }
        .history-row.copy-flash { background: var(--copy-flash-bg); }
        .h-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; flex: 1; cursor: pointer; }
        .h-id {
            font-family: 'JetBrains Mono', ui-monospace, monospace; font-size: 0.75rem;
            color: var(--text-sec); overflow: hidden; text-overflow: ellipsis;
            white-space: nowrap; transition: color var(--tr) ease; pointer-events: none;
        }
        .h-ts { font-size: 0.62rem; color: var(--text-muted); font-weight: 600; letter-spacing: 0.02em; transition: color var(--tr) ease; pointer-events: none; }
        .history-actions { display: flex; align-items: center; gap: 0.15rem; flex-shrink: 0; }
        .icon-btn {
            background: none; border: none; cursor: pointer; padding: 6px;
            color: var(--text-muted); border-radius: 6px; flex-shrink: 0;
            transition: color 0.2s ease, background 0.2s ease;
        }
        .icon-btn:hover { color: var(--accent); background: var(--accent-light); }
        .icon-btn svg { width: 16px; height: 16px; fill: none; stroke: currentColor; stroke-width: 2; display: block; }
        .icon-btn.danger:hover { color: var(--error-text); background: var(--error-bg); }

        /* ── Devices Section ── */
        .devices-section { text-align: left; }
        .devices-guest {
            text-align: center; padding: 2.5rem 1rem; color: var(--text-muted);
        }
        .devices-guest p { font-size: 0.85rem; margin-bottom: 1.25rem; line-height: 1.5; }
        .devices-guest .btn-main { max-width: 180px; margin: 0 auto; }
        .device-search {
            width: 100%; background: var(--bg-input); border: 1.5px solid var(--border);
            border-radius: 8px; padding: 0.5rem 0.75rem; font-family: inherit;
            font-size: 0.8rem; color: var(--text); margin-bottom: 0.75rem;
            outline: none; transition: border-color 0.2s ease, box-shadow 0.2s ease;
        }
        .device-search:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-light); }
        .device-search::placeholder { color: var(--text-muted); }
        .device-row {
            display: flex; align-items: flex-start; justify-content: space-between;
            padding: 0.7rem 0.5rem; border-bottom: 1px solid var(--hist-border);
            gap: 0.5rem; border-radius: 8px; margin: 0 -0.5rem;
            transition: background 0.2s ease;
        }
        .device-row:hover { background: var(--hist-hover); }
        .device-row:last-child { border-bottom: none; }
        .device-row.copy-flash { background: var(--copy-flash-bg) !important; }
        .device-info { display: flex; flex-direction: column; gap: 3px; min-width: 0; flex: 1; cursor: pointer; }
        .device-name {
            font-size: 0.82rem; font-weight: 600; color: var(--text);
            white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        }
        .device-meta { display: flex; align-items: center; gap: 0.35rem; flex-wrap: wrap; margin-top: 1px; }
        .device-location { font-size: 0.72rem; color: var(--text-sec); }
        .device-tag {
            font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.65rem; color: var(--text-muted);
            white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 220px;
        }
        .device-by { font-size: 0.65rem; color: var(--text-muted); font-style: italic; }
        .device-actions { display: flex; align-items: center; gap: 0.1rem; flex-shrink: 0; padding-top: 1px; }
        .devices-empty { text-align: center; padding: 2rem 0; color: var(--text-muted); font-size: 0.8rem; }

        .env-badge {
            display: inline-flex; align-items: center; padding: 2px 7px;
            border-radius: 4px; font-size: 0.6rem; font-weight: 700;
            text-transform: uppercase; letter-spacing: 0.05em; flex-shrink: 0;
        }
        .env-dev             { background: rgba(59,130,246,0.1);  color: #3b82f6; }
        .env-test            { background: rgba(139,92,246,0.1);  color: #8b5cf6; }
        .env-preprod         { background: rgba(245,158,11,0.1);  color: #d97706; }
        .env-prod            { background: rgba(239,68,68,0.1);   color: #dc2626; }
        .env-staging         { background: rgba(234,179,8,0.1);   color: #ca8a04; }
        .env-cadent          { background: rgba(16,185,129,0.1);  color: #059669; }
        .env-sgn             { background: rgba(20,184,166,0.1);  color: #0d9488; }
        .env-custom          { background: rgba(107,114,128,0.1); color: #6b7280; }
        .env-shared          { background: var(--accent-light); color: var(--accent); }
        [data-theme="dark"] .env-dev     { background: rgba(59,130,246,0.15); }
        [data-theme="dark"] .env-test    { background: rgba(139,92,246,0.15); }
        [data-theme="dark"] .env-preprod { background: rgba(245,158,11,0.15); }
        [data-theme="dark"] .env-prod    { background: rgba(239,68,68,0.15);  }
        [data-theme="dark"] .env-staging { background: rgba(234,179,8,0.15);  }
        [data-theme="dark"] .env-cadent  { background: rgba(16,185,129,0.15); }
        [data-theme="dark"] .env-sgn     { background: rgba(20,184,166,0.15); }
        [data-theme="dark"] .env-custom  { background: rgba(107,114,128,0.15);}

        /* ── Dialog Overlay ── */
        .dialog-overlay {
            display: none; position: fixed; inset: 0;
            background: rgba(0,0,0,0.4); backdrop-filter: blur(6px);
            -webkit-backdrop-filter: blur(6px); z-index: 1500;
            align-items: center; justify-content: center; padding: 1rem;
        }
        .dialog-overlay.show { display: flex; }
        .dialog {
            background: var(--bg-card); border: 1px solid var(--border);
            border-radius: 16px; padding: 1.5rem; width: 100%; max-width: 400px;
            box-shadow: 0 24px 64px rgba(0,0,0,0.18); position: relative;
            max-height: 90vh; overflow-y: auto;
            transition: background var(--tr) ease, border-color var(--tr) ease;
        }
        .dialog-close {
            position: absolute; top: 1rem; right: 1rem; background: none; border: none;
            cursor: pointer; color: var(--text-muted); font-size: 1.2rem; line-height: 1;
            padding: 4px 6px; border-radius: 6px;
            transition: color 0.2s ease, background 0.2s ease;
        }
        .dialog-close:hover { color: var(--text); background: var(--hist-hover); }
        .dialog-title {
            font-size: 0.95rem; font-weight: 700; color: var(--text);
            margin-bottom: 1.25rem; padding-right: 2rem;
        }

        /* ── Auth Tabs ── */
        .auth-tabs {
            display: flex; border-bottom: 1px solid var(--border); margin-bottom: 1.25rem;
            transition: border-color var(--tr) ease;
        }
        .auth-tab {
            flex: 1; background: none; border: none; cursor: pointer; font-family: inherit;
            font-size: 0.72rem; font-weight: 700; color: var(--text-muted);
            padding: 0.5rem; text-transform: uppercase; letter-spacing: 0.08em;
            border-bottom: 2px solid transparent; margin-bottom: -1px;
            transition: color 0.2s ease, border-color 0.2s ease;
        }
        .auth-tab.active { color: var(--accent); border-bottom-color: var(--accent); }

        /* ── Form Elements ── */
        .form-group { margin-bottom: 0.9rem; }
        .form-label {
            display: block; font-size: 0.68rem; font-weight: 700;
            text-transform: uppercase; letter-spacing: 0.08em;
            color: var(--text-sec); margin-bottom: 0.3rem;
            transition: color var(--tr) ease;
        }
        .form-input, .form-select {
            width: 100%; background: var(--bg-input); border: 1.5px solid var(--border);
            border-radius: 8px; padding: 0.6rem 0.75rem; font-family: inherit;
            font-size: 0.82rem; color: var(--text); outline: none;
            transition: border-color 0.2s ease, box-shadow 0.2s ease, background var(--tr) ease;
            -webkit-appearance: none; appearance: none;
        }
        .form-input:focus, .form-select:focus {
            border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-light);
        }
        .form-input::placeholder { color: var(--text-muted); }
        .form-select {
            cursor: pointer; padding-right: 2rem;
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%23999' stroke-width='2.5'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
            background-repeat: no-repeat; background-position: right 0.75rem center;
        }
        .form-input-mt { margin-top: 0.4rem; }
        .form-error {
            display: none; background: var(--error-bg); color: var(--error-text);
            border-radius: 7px; padding: 0.5rem 0.75rem; font-size: 0.75rem;
            font-weight: 600; margin-bottom: 0.5rem;
            transition: background var(--tr) ease, color var(--tr) ease;
        }
        .form-error.show { display: block; }
        .checkbox-row {
            display: flex; align-items: flex-start; gap: 0.5rem;
            font-size: 0.78rem; color: var(--text-sec); cursor: pointer;
            user-select: none; line-height: 1.4;
        }
        .checkbox-row input[type="checkbox"] {
            margin-top: 2px; accent-color: var(--accent); flex-shrink: 0; cursor: pointer;
        }
        .dialog-actions { display: flex; gap: 0.5rem; margin-top: 1.1rem; }
        .btn-cancel {
            flex: 1; background: var(--bg-input); color: var(--text-sec);
            border: 1.5px solid var(--border); padding: 0.65rem; border-radius: 8px;
            font-family: inherit; font-weight: 700; font-size: 0.78rem;
            text-transform: uppercase; letter-spacing: 0.05em; cursor: pointer;
            transition: background 0.2s ease, border-color 0.2s ease;
        }
        .btn-cancel:hover { background: var(--hist-hover); border-color: var(--border-hover); }
        .dialog-actions .btn-main { flex: 2; padding: 0.65rem; font-size: 0.78rem; }
        .device-tag-preview {
            background: var(--bg-input); border: 1px solid var(--border);
            border-radius: 8px; padding: 0.45rem 0.75rem; margin-bottom: 1rem;
            font-family: 'JetBrains Mono', ui-monospace, monospace;
            font-size: 0.72rem; color: var(--text-sec); word-break: break-all;
            transition: background var(--tr) ease, border-color var(--tr) ease;
        }

        /* ── QR Fullscreen Modal ── */
        #modal {
            display: none; position: fixed; inset: 0; background: var(--modal-bg);
            backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);
            z-index: 2000; flex-direction: column; align-items: center;
            justify-content: center; padding: 20px; transition: background var(--tr) ease;
        }
        #modal.show { display: flex; }
        #modal img {
            width: min(95vmin, 800px); max-width: 800px; max-height: 800px;
            object-fit: contain; aspect-ratio: 1/1; image-rendering: pixelated;
            border-radius: 12px; box-shadow: 0 8px 40px rgba(0,0,0,0.1);
        }
        .close-fs {
            position: absolute; top: 1rem; right: 1.25rem; width: 40px; height: 40px;
            border-radius: 10px; display: flex; align-items: center; justify-content: center;
            background: var(--bg-card); border: 1px solid var(--border); cursor: pointer;
            font-size: 1.25rem; color: var(--text-muted); transition: all 0.2s ease;
        }
        .close-fs:hover { color: var(--text); border-color: var(--border-hover); }

        /* ── Toast ── */
        #toast {
            position: fixed; bottom: 1.5rem; left: 50%;
            transform: translateX(-50%) translateY(80px);
            background: var(--toast-bg); color: var(--toast-text);
            padding: 0.65rem 1.25rem; border-radius: 10px; font-size: 0.8rem;
            font-weight: 600; opacity: 0;
            transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1),
                        opacity 0.35s ease, background var(--tr) ease, color var(--tr) ease;
            box-shadow: 0 4px 20px rgba(0,0,0,0.12); z-index: 3000;
            white-space: nowrap; pointer-events: none;
        }
        #toast.show { transform: translateX(-50%) translateY(0); opacity: 1; }

        /* ── Animations ── */
        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(14px); }
            to   { opacity: 1; transform: translateY(0); }
        }
        .page-title { animation: fadeInUp 0.4s ease both; }
        .card       { animation: fadeInUp 0.45s ease 0.05s both; }
        .tab-bar, .history, .devices-section { animation: fadeInUp 0.45s ease 0.1s both; }

        @media (max-width: 768px) { .content { padding: 2rem 1rem 2.5rem; } }

        @media (max-width: 480px) {
            .header { padding: 0 0.75rem; height: 50px; }
            .dfx-logo { width: 28px; height: 28px; }
            .app-title { font-size: 0.82rem; }
            .theme-btn { width: 32px; height: 32px; }
            .theme-btn svg { width: 18px; height: 18px; }
            .auth-btn { font-size: 0.7rem; padding: 0 0.65rem; height: 28px; }
            .user-email { max-width: 80px; }
            .content { padding: 1.5rem 0.85rem 2rem; }
            .card { padding: 1.25rem 1rem; border-radius: 12px; }
            .qr-frame { padding: 0.75rem; }
            .qr-frame img { width: 150px; height: 150px; }
            .id-display { padding: 0.6rem 0.75rem; }
            .id-display code { font-size: 0.7rem; }
            .edit-input { font-size: 0.7rem; }
            .btn-main { padding: 0.75rem; font-size: 0.75rem; border-radius: 8px; }
            .h-id { font-size: 0.68rem; }
            .h-ts { font-size: 0.58rem; }
            .page-title { font-size: 0.65rem; margin-bottom: 1.25rem; }
            .dialog { padding: 1.25rem; border-radius: 14px; }
        }

        @media (max-width: 340px) {
            .qr-frame img { width: 130px; height: 130px; }
            .id-display code { font-size: 0.62rem; }
            .edit-input { font-size: 0.62rem; }
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
        <div class="header-right">
            <div id="authArea"></div>
            <button class="theme-btn" id="themeToggle" aria-label="Toggle dark mode">
                <svg class="icon-sun" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="5"/>
                    <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
                    <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
                    <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
                    <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
                </svg>
                <svg class="icon-moon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
                </svg>
            </button>
        </div>
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
            <div class="id-display" id="mainIdBox">
                <code id="currentId" onclick="handleCopy('currentId', document.getElementById('mainIdBox'))">{{.UUID}}</code>
                <input type="text" id="editInput" class="edit-input" style="display:none" autocomplete="off" spellcheck="false">
                <div class="icon-btn edit-only" id="editBtn" onclick="startEdit()" title="Edit tag">
                    <svg viewBox="0 0 24 24"><path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                </div>
                <div class="icon-btn copy-only" id="copyBtn" onclick="handleCopy('currentId', document.getElementById('mainIdBox'))">
                    <svg viewBox="0 0 24 24"><path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"/><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"/></svg>
                </div>
            </div>
            <div class="edit-error" id="editError">
                <span class="edit-error-text">Invalid tag</span>
                <button class="regen-btn" onclick="regenFromError()" type="button">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2v6h-6"/><path d="M3 12a9 9 0 0 1 15-6.7L21 8"/><path d="M3 22v-6h6"/><path d="M21 12a9 9 0 0 1-15 6.7L3 16"/></svg>
                    Regenerate
                </button>
            </div>
            <button class="btn-main" onclick="generate()">Generate Tag</button>
            <button class="btn-secondary" onclick="onSaveTag()">Save Tag</button>
        </div>

        <div class="tab-bar">
            <button class="tab-btn active" id="tabBtnHistory" onclick="switchTab('history')">Recent</button>
            <button class="tab-btn" id="tabBtnDevices" onclick="switchTab('devices')">My Devices</button>
        </div>

        <div class="history" id="historySection">
            <span class="history-label">Recent History</span>
            <div id="historyList"></div>
        </div>

        <div class="devices-section" id="devicesSection" style="display:none">
            <div id="devicesGuest" style="display:none">
                <div class="devices-guest">
                    <p>Sign in to save tags and keep track of your enrolled devices.</p>
                    <button class="btn-main" onclick="openAuthModal('login')">Sign In</button>
                </div>
            </div>
            <div id="devicesContent" style="display:none">
                <input type="text" class="device-search" id="deviceSearch"
                       placeholder="Search by name, tag, environment or location..."
                       oninput="filterDevices()">
                <div id="devicesRows"></div>
            </div>
        </div>
    </div>

    <div class="dialog-overlay" id="authOverlay">
        <div class="dialog" id="authDialog">
            <button class="dialog-close" onclick="closeAuthModal()" aria-label="Close">&times;</button>
            <div class="auth-tabs">
                <button class="auth-tab active" id="authTabLogin" onclick="switchAuthTab('login')">Sign In</button>
                <button class="auth-tab" id="authTabRegister" onclick="switchAuthTab('register')">Register</button>
            </div>
            <form id="authForm" onsubmit="submitAuth(event)">
                <div class="form-group">
                    <label class="form-label" for="authEmail">Email</label>
                    <input class="form-input" type="email" id="authEmail" name="email"
                           placeholder="you@deltaflare.com" required autocomplete="email">
                </div>
                <div class="form-group">
                    <label class="form-label" for="authPassword">Password</label>
                    <input class="form-input" type="password" id="authPassword" name="password"
                           placeholder="••••••••" required autocomplete="current-password">
                </div>
                <div class="form-error" id="authError"></div>
                <button type="submit" class="btn-main" id="authSubmitBtn" style="margin-top:0.4rem">Sign In</button>
            </form>
        </div>
    </div>

    <div class="dialog-overlay" id="deviceOverlay">
        <div class="dialog" id="deviceDialog">
            <button class="dialog-close" onclick="closeDeviceModal()" aria-label="Close">&times;</button>
            <div class="dialog-title" id="deviceModalTitle">Save Tag</div>
            <div class="device-tag-preview" id="deviceTagDisplay"></div>
            <form id="deviceForm" onsubmit="submitDevice(event)">
                <div class="form-group">
                    <label class="form-label" for="dName">Device Name</label>
                    <input class="form-input" type="text" id="dName" maxlength="100" required
                           placeholder="PP1">
                </div>
                <div class="form-group">
                    <label class="form-label" for="dEnv">Environment</label>
                    <select class="form-select" id="dEnv" onchange="onEnvChange()" required>
                        <option value="">Select environment</option>
                        <option value="dev">Dev</option>
                        <option value="test">Test</option>
                        <option value="preproduction">Preproduction</option>
                        <option value="production">Production</option>
                        <option value="staging">Staging</option>
                        <option value="cadent">Cadent</option>
                        <option value="sgn">SGN</option>
                        <option value="other">Other</option>
                    </select>
                    <input class="form-input form-input-mt" type="text" id="dEnvCustom"
                           maxlength="100" placeholder="e.g. Preproduction"
                           style="display:none">
                </div>
                <div class="form-group">
                    <label class="form-label" for="dSerial">Serial Number <span style="font-weight:400;text-transform:none;letter-spacing:0;color:var(--text-muted)">(optional)</span></label>
                    <input class="form-input" type="text" id="dSerial" maxlength="100"
                           placeholder="e.g. E604930">
                </div>
                <div class="form-group">
                    <label class="form-label" for="dLocation">Location</label>
                    <select class="form-select" id="dLocation" onchange="onLocationChange()" required>
                        <option value="">Select location</option>
                        <option value="dev rack">Dev Rack</option>
                        <option value="preproduction rack">Preproduction Rack</option>
                        <option value="production rack">Production Rack</option>
                        <option value="ot test rack">OT Test Rack</option>
                        <option value="other">Other</option>
                    </select>
                    <input class="form-input form-input-mt" type="text" id="dLocationCustom"
                           maxlength="100" placeholder="e.g. Preproduction Rack"
                           style="display:none">
                </div>
                <div class="form-group">
                    <label class="checkbox-row">
                        <input type="checkbox" id="dGlobal">
                        <span>Share with all users</span>
                    </label>
                </div>
                <div class="form-error" id="deviceError"></div>
                <div class="dialog-actions">
                    <button type="button" class="btn-cancel" onclick="closeDeviceModal()">Cancel</button>
                    <button type="submit" class="btn-main" id="deviceSubmitBtn">Save</button>
                </div>
            </form>
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
            if (manual && saved) { setTheme(saved); }
            else if (window.matchMedia('(prefers-color-scheme: dark)').matches) { setTheme('dark'); }
            else { setTheme('light'); }
        })();

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
            if (!localStorage.getItem('dfx-theme-manual')) setTheme(e.matches ? 'dark' : 'light');
        });

        toggleBtn.addEventListener('click', function() {
            setTheme(html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark');
            localStorage.setItem('dfx-theme-manual', 'true');
        });

        /* ── Utils ── */
        function escHtml(s) {
            if (!s) return '';
            return String(s)
                .replace(/&/g, '&amp;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;')
                .replace(/"/g, '&quot;')
                .replace(/'/g, '&#39;');
        }

        /* ── SVGs ── */
        var checkSvg = '<path d="M20 6L9 17L4 12" stroke-linecap="round" stroke-linejoin="round"/>';
        var copySvg = '<path d="M8 4v12a2 2 0 002 2h8a2 2 0 002-2V7.242a2 2 0 00-.586-1.414l-3.242-3.242A2 2 0 0014.758 2H10a2 2 0 00-2 2z"/><path d="M16 18v2a2 2 0 01-2 2H6a2 2 0 01-2-2V9a2 2 0 012-2h2"/>';
        var qrIconSvg = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="8" height="8" rx="1"/><rect x="14" y="2" width="8" height="8" rx="1"/><rect x="2" y="14" width="8" height="8" rx="1"/><rect x="5" y="5" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="17" y="5" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="5" y="17" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="14" y="14" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="18" y="14" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="14" y="18" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/><rect x="18" y="18" width="2" height="2" rx="0.3" fill="currentColor" stroke="none"/></svg>';

        /* ── Toast ── */
        var toastTimer = null;
        var toastEl = document.getElementById('toast');

        function showToast(msg) {
            if (toastTimer) clearTimeout(toastTimer);
            toastEl.textContent = msg || 'Copied to clipboard';
            toastEl.classList.add('show');
            toastTimer = setTimeout(function() {
                toastEl.classList.remove('show');
                toastTimer = null;
            }, 2500);
        }

        /* ── Copy flash ── */
        var activeFlash = null;

        function clearActiveFlash() {
            if (!activeFlash) return;
            clearTimeout(activeFlash.timer);
            activeFlash.el.classList.remove('copy-flash');
            if (activeFlash.svg) activeFlash.svg.innerHTML = copySvg;
            activeFlash = null;
        }

        function flashElement(el) {
            clearActiveFlash();
            el.classList.add('copy-flash');
            var svg = el.querySelector('.copy-only svg');
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
            if (isEditing) return;
            handleCopyText(document.getElementById(id).innerText, el);
        }

        function handleCopyText(text, el) {
            if (navigator.clipboard && navigator.clipboard.writeText) {
                navigator.clipboard.writeText(text).then(function() {
                    showToast(); flashElement(el);
                }).catch(function() { fallbackCopy(text, el); });
            } else { fallbackCopy(text, el); }
        }

        function fallbackCopy(text, el) {
            var ta = document.createElement('textarea');
            ta.value = text; ta.style.cssText = 'position:fixed;left:-9999px';
            document.body.appendChild(ta); ta.focus(); ta.select();
            try { if (document.execCommand('copy')) { showToast(); flashElement(el); } }
            catch(e) { console.error(e); }
            document.body.removeChild(ta);
        }

        /* ── Edit tag ── */
        var isEditing = false;
        var dfxPattern = /^dfx-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

        function clearEditState() {
            isEditing = false;
            var box   = document.getElementById('mainIdBox');
            var code  = document.getElementById('currentId');
            var input = document.getElementById('editInput');
            var err   = document.getElementById('editError');
            input.style.display = 'none';
            code.style.display  = '';
            box.classList.remove('editing', 'edit-error-state');
            err.classList.remove('show');
        }

        function startEdit() {
            if (isEditing) return;
            isEditing = true;
            var box   = document.getElementById('mainIdBox');
            var code  = document.getElementById('currentId');
            var input = document.getElementById('editInput');
            var err   = document.getElementById('editError');
            err.classList.remove('show');
            box.classList.remove('edit-error-state');
            input.value = code.innerText;
            code.style.display  = 'none';
            input.style.display = 'block';
            box.classList.add('editing');
            input.focus();
            input.setSelectionRange(input.value.length, input.value.length);
        }

        async function finishEdit() {
            if (!isEditing) return;
            var box   = document.getElementById('mainIdBox');
            var code  = document.getElementById('currentId');
            var input = document.getElementById('editInput');
            var err   = document.getElementById('editError');
            var val   = input.value.trim().toLowerCase();
            if (val === code.innerText) { clearEditState(); return; }
            if (!val || !dfxPattern.test(val)) {
                box.classList.add('edit-error-state');
                err.classList.add('show');
                return;
            }
            err.classList.remove('show');
            box.classList.remove('edit-error-state');
            isEditing = false;
            code.innerText = val;
            input.style.display = 'none';
            code.style.display  = '';
            box.classList.remove('editing');
            try {
                var res  = await fetch('/api/qr?text=' + encodeURIComponent(val));
                var data = await res.json();
                document.getElementById('qrImg').src = 'data:image/png;base64,' + data.qr_data;
            } catch(e) { console.error('QR fetch failed:', e); }
            var ts = new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
            var saved = localStorage.getItem('dfx_history_v3');
            var hist  = saved ? JSON.parse(saved) : [];
            hist.unshift({ uuid: val, timestamp: ts });
            if (hist.length > 5) hist = hist.slice(0, 5);
            saveHistory(hist);
        }

        function cancelEdit() { clearEditState(); }

        async function regenFromError() { clearEditState(); await generate(); }

        document.getElementById('editInput').addEventListener('keydown', function(e) {
            if (e.key === 'Enter')  { e.preventDefault(); finishEdit(); }
            if (e.key === 'Escape') { e.preventDefault(); cancelEdit(); }
        });

        document.getElementById('editInput').addEventListener('input', function() {
            if (dfxPattern.test(this.value.trim())) {
                document.getElementById('editError').classList.remove('show');
                document.getElementById('mainIdBox').classList.remove('edit-error-state');
            }
        });

        document.getElementById('editInput').addEventListener('blur', function() {
            setTimeout(function() {
                if (!isEditing) return;
                var val = document.getElementById('editInput').value.trim();
                if (dfxPattern.test(val)) { finishEdit(); } else { cancelEdit(); }
            }, 200);
        });

        /* ── History ── */
        function saveHistory(data) {
            localStorage.setItem('dfx_history_v3', JSON.stringify(data));
            renderHistory(data);
        }

        function loadHistory() {
            var saved = localStorage.getItem('dfx_history_v3');
            renderHistory(saved ? JSON.parse(saved) : []);
        }

        function renderHistory(items) {
            var list = document.getElementById('historyList');
            var h = '';
            items.forEach(function(item) {
                h += '<div class="history-row">' +
                    '<div class="h-info" onclick="handleCopyText(\'' + escHtml(item.uuid) + '\', this.parentElement)">' +
                    '<span class="h-ts">' + escHtml(item.timestamp) + '</span>' +
                    '<span class="h-id">' + escHtml(item.uuid) + '</span></div>' +
                    '<div class="history-actions">' +
                    '<div class="icon-btn" onclick="event.stopPropagation();showHistoryQR(\'' + escHtml(item.uuid) + '\')" title="Show QR">' + qrIconSvg + '</div>' +
                    '<div class="icon-btn" onclick="event.stopPropagation();saveFromHistory(\'' + escHtml(item.uuid) + '\')" title="Save to My Devices">' +
                    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/></svg></div>' +
                    '<div class="icon-btn copy-only" onclick="event.stopPropagation();handleCopyText(\'' + escHtml(item.uuid) + '\', this.closest(\'.history-row\'))" title="Copy">' +
                    '<svg viewBox="0 0 24 24">' + copySvg + '</svg></div>' +
                    '</div></div>';
            });
            list.innerHTML = h;
        }

        /* ── Device Saving Logic ── */
        function saveFromHistory(tag) {
            if (!currentUser) {
                pendingAction = { type: 'saveTag', tag: tag };
                openAuthModal('login');
                return;
            }
            openSaveDeviceModal(tag);
        }

        function onSaveTag() {
            saveFromHistory(null);
        }

        function openSaveDeviceModal(tag) {
            editingDeviceId = null;
            document.getElementById('deviceModalTitle').textContent = 'Save Tag';
            document.getElementById('deviceTagDisplay').textContent =
                tag || document.getElementById('currentId').innerText;
            document.getElementById('dName').value       = '';
            document.getElementById('dSerial').value     = '';
            document.getElementById('dEnv').value        = '';
            document.getElementById('dEnvCustom').style.display = 'none';
            document.getElementById('dEnvCustom').value  = '';
            document.getElementById('dLocation').value   = '';
            document.getElementById('dLocationCustom').style.display = 'none';
            document.getElementById('dLocationCustom').value = '';
            document.getElementById('dGlobal').checked   = false;
            document.getElementById('deviceError').classList.remove('show');
            document.getElementById('deviceSubmitBtn').textContent = 'Save';
            document.getElementById('deviceOverlay').classList.add('show');
            setTimeout(function() { document.getElementById('dName').focus(); }, 60);
        }

        /* ── Generate ── */
        async function generate() {
            clearEditState();
            clearActiveFlash();
            var res  = await fetch('/api/generate');
            var data = await res.json();
            document.getElementById('qrImg').src   = 'data:image/png;base64,' + data.qr_data;
            document.getElementById('currentId').innerText = data.uuid;
            var saved = localStorage.getItem('dfx_history_v3');
            var local = saved ? JSON.parse(saved) : [];
            local.unshift({ uuid: data.uuid, timestamp: data.history[0].timestamp });
            if (local.length > 5) local = local.slice(0, 5);
            saveHistory(local);
        }

        /* ── Fullscreen QR ── */
        function openFullscreen() {
            document.getElementById('fsImg').src = document.getElementById('qrImg').src;
            document.getElementById('modal').classList.add('show');
        }

        function openFullscreenSrc(src) {
            document.getElementById('fsImg').src = src;
            document.getElementById('modal').classList.add('show');
        }

        function closeFullscreen() { document.getElementById('modal').classList.remove('show'); }

        async function showHistoryQR(uuid) {
            try {
                var res  = await fetch('/api/qr?text=' + encodeURIComponent(uuid));
                var data = await res.json();
                openFullscreenSrc('data:image/png;base64,' + data.qr_data);
            } catch(e) { console.error('QR fetch failed:', e); }
        }

        /* ═══════════════════════════════════════════
           Auth
        ═══════════════════════════════════════════ */
        var currentUser   = null;
        var authMode      = 'login';
        var pendingAction = null;

        function initAuth() {
            fetch('/api/auth/me')
                .then(function(r) { return r.json(); })
                .then(function(data) {
                    currentUser = data.user;
                    renderAuthArea();
                    if (currentUser) loadDevices();
                })
                .catch(function() { renderAuthArea(); });
        }

        function renderAuthArea() {
            var area = document.getElementById('authArea');
            if (!currentUser) {
                area.innerHTML = '<button class="auth-btn" onclick="openAuthModal(\'login\')">Sign In</button>';
            } else {
                var displayName = currentUser.email.replace('@deltaflare.com', '');
                area.innerHTML =
                    '<div class="user-pill">' +
                    '<span class="user-email">' + escHtml(displayName) + '</span>' +
                    '<button class="signout-btn" onclick="doLogout()" title="Sign Out">' +
                    '<svg viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">' +
                    '<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>' +
                    '<polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/>' +
                    '</svg></button></div>';
            }
        }

        function openAuthModal(mode) {
            authMode = mode || 'login';
            switchAuthTab(authMode);
            document.getElementById('authEmail').value    = '';
            document.getElementById('authPassword').value = '';
            document.getElementById('authError').classList.remove('show');
            document.getElementById('authOverlay').classList.add('show');
            setTimeout(function() { document.getElementById('authEmail').focus(); }, 60);
        }

        function closeAuthModal() {
            document.getElementById('authOverlay').classList.remove('show');
        }

        document.getElementById('authOverlay').addEventListener('click', function(e) {
            if (e.target === this) closeAuthModal();
        });

        function switchAuthTab(mode) {
            authMode = mode;
            document.getElementById('authTabLogin').classList.toggle('active', mode === 'login');
            document.getElementById('authTabRegister').classList.toggle('active', mode === 'register');
            document.getElementById('authSubmitBtn').textContent =
                mode === 'login' ? 'Sign In' : 'Create Account';
            document.getElementById('authPassword').autocomplete =
                mode === 'login' ? 'current-password' : 'new-password';
            document.getElementById('authError').classList.remove('show');
        }

        async function submitAuth(e) {
            e.preventDefault();
            var email    = document.getElementById('authEmail').value.trim();
            var password = document.getElementById('authPassword').value;
            var errEl    = document.getElementById('authError');
            var btn      = document.getElementById('authSubmitBtn');
            var label    = authMode === 'login' ? 'Sign In' : 'Create Account';

            btn.disabled    = true;
            btn.textContent = '...';
            errEl.classList.remove('show');

            try {
                var url = authMode === 'login' ? '/api/auth/login' : '/api/auth/register';
                var res = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: email, password: password })
                });
                var data = await res.json();
                if (!res.ok) {
                    errEl.textContent = data.error || 'Something went wrong';
                    errEl.classList.add('show');
                    return;
                }
                currentUser = data.user;
                renderAuthArea();
                closeAuthModal();
                await loadDevices();
                if (pendingAction && pendingAction.type === 'saveTag') {
                    var tag = pendingAction.tag;
                    pendingAction = null;
                    openSaveDeviceModal(tag);
                } else {
                    switchTab('devices');
                }
            } catch(err) {
                errEl.textContent = 'Connection error. Please try again.';
                errEl.classList.add('show');
            } finally {
                btn.disabled    = false;
                btn.textContent = label;
            }
        }

        async function doLogout() {
            try { await fetch('/api/auth/logout', { method: 'POST' }); } catch(e) {}
            currentUser = null;
            renderAuthArea();
            allDevices = [];
            if (document.getElementById('devicesSection').style.display !== 'none') {
                renderDevicesArea();
                renderDevices([]);
            }
        }

        /* ═══════════════════════════════════════════
           Tabs
        ═══════════════════════════════════════════ */
        function switchTab(tab) {
            var isHistory = (tab === 'history');
            document.getElementById('historySection').style.display  = isHistory ? '' : 'none';
            document.getElementById('devicesSection').style.display  = isHistory ? 'none' : '';
            document.getElementById('tabBtnHistory').classList.toggle('active', isHistory);
            document.getElementById('tabBtnDevices').classList.toggle('active', !isHistory);
            if (!isHistory) { renderDevicesArea(); }
            localStorage.setItem('dfx-active-tab', tab);
        }

        /* ═══════════════════════════════════════════
           Devices
        ═══════════════════════════════════════════ */
        var allDevices = [];

        function renderDevicesArea() {
            var guest   = document.getElementById('devicesGuest');
            var content = document.getElementById('devicesContent');
            if (!currentUser) {
                guest.style.display   = '';
                content.style.display = 'none';
            } else {
                guest.style.display   = 'none';
                content.style.display = '';
            }
        }

        async function loadDevices() {
            if (!currentUser) return;
            try {
                var res = await fetch('/api/devices');
                allDevices = await res.json();
                renderDevicesArea();
                renderDevices(allDevices);
            } catch(e) { console.error('Failed to load devices', e); }
        }

        function filterDevices() {
            var q = document.getElementById('deviceSearch').value.toLowerCase();
            if (!q) { renderDevices(allDevices); return; }
            renderDevices(allDevices.filter(function(d) {
                return d.device_name.toLowerCase().indexOf(q) !== -1 ||
                       d.tag.toLowerCase().indexOf(q) !== -1 ||
                       d.environment.toLowerCase().indexOf(q) !== -1 ||
                       d.location.toLowerCase().indexOf(q) !== -1;
            }));
        }

        var envClassMap = {
            'dev': 'env-dev', 'test': 'env-test',
            'preproduction': 'env-preprod', 'production': 'env-prod',
            'staging': 'env-staging', 'cadent': 'env-cadent', 'sgn': 'env-sgn'
        };

        function envBadge(env) {
            var cls = envClassMap[env.toLowerCase()] || 'env-custom';
            return '<span class="env-badge ' + cls + '">' + escHtml(env) + '</span>';
        }

        function renderDevices(devices) {
            var rows = document.getElementById('devicesRows');
            if (!devices || devices.length === 0) {
                rows.innerHTML = '<div class="devices-empty">No saved devices yet.</div>';
                return;
            }
            var h = '';
            devices.forEach(function(d) {
                var isMine      = currentUser && d.user_email === currentUser.email;
                var sharedBadge = d.is_global ? '<span class="env-badge env-shared">Shared</span>' : '';
                var byLine      = (!isMine && d.is_global)
                    ? '<span class="device-by">by ' + escHtml(d.user_email.split('@')[0]) + '</span>'
                    : '';
                h += '<div class="device-row" data-id="' + d.id + '" data-tag="' + escHtml(d.tag) + '">' +
                    '<div class="device-info" onclick="handleCopyText(\'' + escHtml(d.tag) + '\', this.parentElement)">' +
                    '<span class="device-name">' + escHtml(d.device_name) + '</span>' +
                    (d.serial_number ? '<span style="font-size:0.65rem; color:var(--text-muted);">S/N: ' + escHtml(d.serial_number) + '</span>' : '') +
                    '<div class="device-meta">' + envBadge(d.environment) + sharedBadge + '</div>' +
                    '<span class="device-location">' + escHtml(d.location) + '</span>' +
                    '<span class="device-tag">' + escHtml(d.tag) + '</span>' +
                    byLine +
                    '</div>' +
                    '<div class="device-actions">' +
                    '<div class="icon-btn" onclick="showDeviceQR(this.closest(\'.device-row\').dataset.tag)" title="Show QR">' + qrIconSvg + '</div>' +
                    '<div class="icon-btn copy-only" onclick="event.stopPropagation();handleCopyText(\'' + escHtml(d.tag) + '\', this.closest(\'.device-row\'))" title="Copy">' +
                    '<svg viewBox="0 0 24 24">' + copySvg + '</svg></div>' +
                    '<div class="icon-btn" onclick="openEditDevice(+this.closest(\'.device-row\').dataset.id)" title="Edit">' +
                    '<svg viewBox="0 0 24 24"><path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>' +
                    '</div>' +
                    '<div class="icon-btn danger" onclick="confirmDeleteDevice(+this.closest(\'.device-row\').dataset.id)" title="Delete">' +
                    '<svg viewBox="0 0 24 24"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/></svg>' +
                    '</div>' +
                    '</div>' +
                    '</div>';
            });
            rows.innerHTML = h;
        }

        async function showDeviceQR(tag) {
            try {
                var res  = await fetch('/api/qr?text=' + encodeURIComponent(tag));
                var data = await res.json();
                openFullscreenSrc('data:image/png;base64,' + data.qr_data);
            } catch(e) { console.error(e); }
        }

        var editingDeviceId = null;

        function openEditDevice(id) {
            var device = allDevices.find(function(d) { return d.id === id; });
            if (!device) return;
            editingDeviceId = id;
            document.getElementById('deviceModalTitle').textContent = 'Edit Device';
            document.getElementById('deviceTagDisplay').textContent = device.tag;

            document.getElementById('dName').value       = device.device_name;
            document.getElementById('dSerial').value     = device.serial_number || '';
            document.getElementById('dGlobal').checked   = device.is_global;

            var knownLocs = ['dev rack','preproduction rack','production rack','ot test rack'];
            var locLower  = device.location.toLowerCase();
            if (knownLocs.indexOf(locLower) !== -1) {
                document.getElementById('dLocation').value = locLower;
                document.getElementById('dLocationCustom').style.display = 'none';
                document.getElementById('dLocationCustom').value = '';
            } else {
                document.getElementById('dLocation').value = 'other';
                document.getElementById('dLocationCustom').style.display = '';
                document.getElementById('dLocationCustom').value = device.location;
            }

            var knownEnvs = ['dev','test','preproduction','production','staging','cadent','sgn'];
            var envLower  = device.environment.toLowerCase();
            if (knownEnvs.indexOf(envLower) !== -1) {
                document.getElementById('dEnv').value = envLower;
                document.getElementById('dEnvCustom').style.display = 'none';
                document.getElementById('dEnvCustom').value = '';
            } else {
                document.getElementById('dEnv').value = 'other';
                document.getElementById('dEnvCustom').style.display = '';
                document.getElementById('dEnvCustom').value = device.environment;
            }

            document.getElementById('deviceError').classList.remove('show');
            document.getElementById('deviceSubmitBtn').textContent = 'Update';
            document.getElementById('deviceOverlay').classList.add('show');
            setTimeout(function() { document.getElementById('dName').focus(); }, 60);
        }

        function closeDeviceModal() {
            document.getElementById('deviceOverlay').classList.remove('show');
        }

        document.getElementById('deviceOverlay').addEventListener('click', function(e) {
            if (e.target === this) closeDeviceModal();
        });

        function onEnvChange() {
            var val    = document.getElementById('dEnv').value;
            var custom = document.getElementById('dEnvCustom');
            custom.style.display = (val === 'other') ? '' : 'none';
            if (val !== 'other') custom.value = '';
        }

        function onLocationChange() {
            var val    = document.getElementById('dLocation').value;
            var custom = document.getElementById('dLocationCustom');
            custom.style.display = (val === 'other') ? '' : 'none';
            if (val !== 'other') custom.value = '';
        }

        async function submitDevice(e) {
            e.preventDefault();
            var name      = document.getElementById('dName').value.trim();
            var serial    = document.getElementById('dSerial').value.trim();
            var envSel    = document.getElementById('dEnv').value;
            var envCustom = document.getElementById('dEnvCustom').value.trim();
            var env       = (envSel === 'other') ? envCustom : envSel;
            var locSel    = document.getElementById('dLocation').value;
            var locCustom = document.getElementById('dLocationCustom').value.trim();
            var location  = (locSel === 'other') ? locCustom : locSel;
            var isGlobal  = document.getElementById('dGlobal').checked;
            var errEl     = document.getElementById('deviceError');
            var btn       = document.getElementById('deviceSubmitBtn');

            if (!name || !env || !location) {
                errEl.textContent = 'All fields except serial number are required.';
                errEl.classList.add('show');
                return;
            }

            btn.disabled = true;
            errEl.classList.remove('show');

            try {
                var url, method, body;
                if (editingDeviceId === null) {
                    url    = '/api/devices';
                    method = 'POST';
                    body   = {
                        tag:           document.getElementById('deviceTagDisplay').textContent,
                        device_name:   name,
                        serial_number: serial,
                        environment:   env,
                        location:      location,
                        is_global:     isGlobal
                    };
                } else {
                    url    = '/api/devices/' + editingDeviceId;
                    method = 'PUT';
                    body   = { device_name: name, serial_number: serial, environment: env, location: location, is_global: isGlobal };
                }
                var res  = await fetch(url, {
                    method: method,
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(body)
                });
                var data = await res.json();
                if (!res.ok) {
                    errEl.textContent = data.error || 'Something went wrong';
                    errEl.classList.add('show');
                    return;
                }
                closeDeviceModal();
                await loadDevices();
                switchTab('devices');
                showToast(editingDeviceId === null ? 'Tag saved' : 'Device updated');
            } catch(err) {
                errEl.textContent = 'Connection error. Please try again.';
                errEl.classList.add('show');
            } finally {
                btn.disabled = false;
            }
        }

        async function confirmDeleteDevice(id) {
            if (!confirm('Delete this device entry? This cannot be undone.')) return;
            try {
                var res = await fetch('/api/devices/' + id, { method: 'DELETE' });
                if (res.ok) {
                    await loadDevices();
                    showToast('Device deleted');
                }
            } catch(e) { console.error(e); }
        }

        /* ── Init ── */
        window.onload = function() {
            loadHistory();
            initAuth();
            const lastTab = localStorage.getItem('dfx-active-tab');
            if (lastTab) switchTab(lastTab);
        };
    </script>
</body>
</html>`
