# ssh-portfolio

> An interactive TUI portfolio served over raw SSH. No browser. No JavaScript. No loading spinners.

## 👾 Connect

```bash
# One-time install — drops a `mohith` command on your machine
curl -fsSL https://raw.githubusercontent.com/Trafalgar-2006/portfolio/master/install.sh | bash

# Then just run anytime
mohith
```

```bash
# Or connect directly (no install needed)
ssh trolley.proxy.rlwy.net -p 41074
```

<!-- TODO: switch above to `ssh mohith.is-a.dev -p 41074` once is-a-dev PR is merged -->

---

## What visitors see

1. **Matrix rain** — katakana + digits fall across the full terminal
2. **Name solidification** — ASCII name crystallises out of the chaos, cell by cell
3. **Boot sequence** — fake SSH handshake, corrupted loading line, signal-lost recovery
4. **"Unauthorized Access" alert** — red full-screen warning → "just kidding. welcome. :)"
5. **Home splash** — name glitch effect, typewriter tagline, live GitHub commit at the bottom
6. **Projects** — cascade drop-in, [Live] badge pulse, tag pop in detail, colour-coded stack tags
7. **About** — IST time-aware greeting that changes based on when you're reading it
8. **Contacts** — stagger reveal, SSH line flash

---

## Stack

| Layer | Tech |
|---|---|
| TUI framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| SSH server | [Wish](https://github.com/charmbracelet/wish) + [charmbracelet/ssh](https://github.com/charmbracelet/ssh) |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Deployment | [Railway](https://railway.app) (Docker) |
| Custom domain | `mohith.is-a.dev` (via [is-a.dev](https://is-a.dev)) |
| Keep-alive | UptimeRobot TCP monitor |

---

## Deployment (Railway)

This project is deployed on Railway with a persistent SSH host key so the fingerprint never changes between redeploys.

**Environment variables required:**

| Variable | Value | Purpose |
|---|---|---|
| `SSH_ENABLED` | `true` | Enables SSH server mode |
| `SSH_PORT` | `23234` | Internal SSH port (decoupled from Railway's `PORT`) |
| `SSH_HOST_KEY` | base64-encoded private key | Persistent host key — prevents "REMOTE HOST CHANGED" warnings |
| `COLORTERM` | `truecolor` | Forces 24-bit colour in SSH sessions |

**Generate a persistent host key:**
```bash
ssh-keygen -t ed25519 -f host_key
[Convert]::ToBase64String([IO.File]::ReadAllBytes("host_key")) | clip
# Paste the clipboard value into Railway as SSH_HOST_KEY
```

**Networking:**
- HTTP health check → port `8080` (Railway's injected `PORT`)
- SSH TCP proxy → internal `23234`, external `41074`

---

## Customization

**Update content without redeploying:** edit `content.yaml`

```yaml
projects:
  - title: "Your Project"
    description: "What it does and why it matters."
    tags: [Go, Python, Docker]
    status: Live        # Live | WIP | Research
    github: "github.com/you/repo"
    highlight: "optional metric"  # shown as ⚡ highlight

contacts:
  - icon: "(@)"
    label: "Email"
    value: "you@example.com"
```

Push to `master` → Railway rebuilds automatically.

**Change bio/about text:** edit `views/about.go`
**Change ASCII name banner or portrait:** edit `views/home.go`

---

## Controls

| Key | Action |
|---|---|
| `← →` or `h l` | Switch tabs on home screen |
| `Enter` | Open selected tab |
| `↑ ↓` or `j k` | Browse projects |
| `Esc` or `q` | Go back / quit |

---

## Architecture

```
main.go          — dual-port setup: SSH on SSH_PORT, HTTP health on 8080
model.go         — single global 50ms ticker drives all animations
views/
  matrix.go      — matrix rain renderer (grouped ANSI segments for SSH efficiency)
  boot.go        — SSH handshake + corrupted recovery boot sequence
  alert.go       — "unauthorized access" fake alert
  home.go        — splash screen, typewriter, name glitch, banner reveal
  projects.go    — cascade, tag pop, live pulse, YAML loader
  about.go       — IST time-based greeting, experience timeline
  contacts.go    — stagger reveal, SSH line flash
config/
  loader.go      — YAML parser for content.yaml
content.yaml     — all projects and contacts (edit this, not Go files)
entrypoint.sh    — decodes SSH_HOST_KEY env var → persistent key at startup
```

---

## License

MIT
