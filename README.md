DFX Tag Generator

A lightweight, standalone Go tool for generating unique identifiers and QR codes for OT edge device asset management.

Key Features
Asset Tagging: Generates custom dfx- prefixed UUIDs specifically for tracking OT hardware and edge devices.

On-the-Fly QR Encoding: Renders high-quality QR codes as Base64 PNGs on the fly—no external APIs or data leaks.

Exclusive Copy Logic: Refined UI behavior—clicking a new ID instantly clears the highlight and checkmark of any previous ones to avoid confusion.

Manual Overrides: Tap the pencil icon to manually edit a tag. It validates against the expected UUID format and refreshes the QR code automatically.

Interface: Modern, mobile-responsive design with a dedicated Dark Mode
Local History: Keeps the last 5 tags synced in the browser so you can verify recent deployments even after a refresh.

Tech Stack
Backend: Go (using net/http and skip2/go-qrcode).

Frontend: Vanilla JS and CSS3. Zero frameworks.

Infrastructure: Rocky Linux + systemd.

Access: Cloudflare Tunnel for secure, encrypted traffic to the edge without opening inbound firewall ports.

📂 Structure
Plaintext
.
├── main.go           # Application logic + HTML/CSS template
├── go.mod            # Dependencies
└── go.sum            # Checksums

🌐 Networking (Cloudflare)
The app is served via https://getdfx.uk.

Using a Cloudflare Tunnel ensures a "Secure Context," which is required for the browser to allow the Copy to Clipboard feature to work natively