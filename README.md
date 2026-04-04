# Asset Tag Generator

A lightweight tool for generating unique asset identifiers and QR codes for OT edge device management. Built in Go, containerised, and deployed on RKE2.

![Light Mode](light-mode-screenshot.png)
![Dark Mode](dark-mode-screenshot.png)

## Features

**Tag Generation** — creates `dfx-` prefixed UUIDs for tracking OT hardware and edge devices. Each tag gets a matching QR code rendered server-side as a Base64 PNG with no external API calls.

**Manual Editing** — tap the pencil icon to edit a tag directly. Validates against the `dfx-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` format, regenerates the QR to match, and pushes the edit into history. Invalid input shows an inline error with a one-click regenerate button.

**Fullscreen QR** — tap any QR code to open a fullscreen viewer for easy scanning.

**History** — recent tags are persisted in localStorage. Each history entry has a QR button, a save button, and a copy button.

**Copy Behaviour** — exclusive flash logic ensures only one row highlights at a time. Rapid successive copies keep the toast visible for a consistent 2.5s without flickering.

**Dark Mode** — auto-detects system `prefers-color-scheme` on first load. Manual toggle (animated sun/moon icon) overrides and persists via localStorage.

**Responsive** — three breakpoints for desktop, tablet, and mobile.

### My Devices (requires account)

**Accounts** — registration and login restricted to `@deltaflare.com` email addresses. Sessions persist for 30 days via a secure HttpOnly cookie. Passwords are hashed with bcrypt (cost 12).

**Save Tags** — any generated or historical tag can be saved to My Devices with the following fields:

| Field | Details |
|---|---|
| Device Name | Free text |
| Environment | Dev / Test / Preproduction / Production / Staging / Cadent / SGN / Other |
| Serial Number | Optional — `E` + 6-digit format (e.g. `E604930`) |
| Location | Dev Rack / Preproduction Rack / Production Rack / OT Test Rack / Other |
| Visibility | Private (owner only) or Shared (visible and editable by all users) |

**My Devices tab** — searchable list of your saved devices and all shared entries. Each row shows the environment badge, location, tag ID, and actions to view QR, edit, or delete.

## Tech Stack

| Layer | Detail |
|---|---|
| Backend | Go 1.25 — `net/http`, `google/uuid`, `skip2/go-qrcode`, `golang.org/x/crypto` |
| Database | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |
| Frontend | Vanilla JS, CSS3, DM Sans + JetBrains Mono. Zero frameworks |
| Container | Multi-stage Docker build (`golang:1.25-alpine` → `alpine:3.20`) |
| Orchestration | RKE2 (Kubernetes) — Deployment + NodePort Service |
| OS | Rocky Linux |
| Access | Cloudflare Tunnel → `https://getdfx.uk` |

## Architecture

```
Internet → Cloudflare Tunnel → NodePort :30092 → K8s Service → DFX Pod (:9092)
```

The app runs as a single-replica Deployment in the `dfx` namespace. The SQLite database is mounted at `/data/dfx.db` via a persistent volume. The Cloudflare Tunnel provides TLS termination and secure context (required for the Clipboard API) without opening inbound firewall ports.

## API

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/` | GET | — | Serves the UI |
| `/api/generate` | GET | — | Generates a new tag — returns UUID, QR base64, and history |
| `/api/qr?text=` | GET | — | Returns a QR code PNG (base64) for any given text |
| `/api/auth/register` | POST | — | Create a `@deltaflare.com` account |
| `/api/auth/login` | POST | — | Login — sets a 30-day session cookie |
| `/api/auth/logout` | POST | — | Clears the session |
| `/api/auth/me` | GET | — | Returns the current user or null |
| `/api/devices` | GET | ✓ | List visible devices (own + shared) |
| `/api/devices` | POST | ✓ | Save a new device entry |
| `/api/devices/{id}` | PUT | ✓ | Update a device (owner or shared) |
| `/api/devices/{id}` | DELETE | ✓ | Delete a device (owner or shared) |

## Configuration

| Environment Variable | Default | Description |
|---|---|---|
| `DB_PATH` | `./dfx.db` | Path to the SQLite database file |
| `ADMIN_PASSWORD` | `DFXadmin1!` | Password for the seeded `admin@deltaflare.com` account |
