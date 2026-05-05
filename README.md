# SSH Portfolio — Mohith Akshay Duggirala

A beautiful SSH-based TUI portfolio built with Go + [Charm](https://charm.sh/) ecosystem.

## 🌐 Live Demo

```bash
ssh trolley.proxy.rlwy.net -p 23115
```

## ✨ Features

- **Interactive TUI** — Navigate with keyboard (arrow keys, enter, esc)
- **SSH Accessible** — Anyone can view your portfolio via `ssh`
- **Beautiful Design** — Braille art, color palette, ASCII name banner
- **Docker Ready** — One-command deployment

## 🚀 Quick Start

### Run Locally (no SSH)

```bash
go run .
```

### Run as SSH Server

```bash
# Copy and edit the env file
cp .env.example .env

# Set SSH_ENABLED=true in .env, then:
go run .

# In another terminal:
ssh localhost -p 23234
```

### Deploy with Docker

```bash
cp .env.example .env
# Edit .env as needed
docker-compose up --build
```

## 🌐 Deployment Options

Since this is an SSH server, it needs a **persistent host** (not GitHub Pages/Vercel).

### Fly.io (Recommended — Free Tier)
```bash
flyctl launch
flyctl deploy
```

### Any VPS (DigitalOcean, Linode, Hetzner)
```bash
# On your server:
git clone https://github.com/trafalgar-2006/ssh-portfolio.git
cd ssh-portfolio
docker-compose up -d --build
```

## 🎨 Customization

Edit the view files in `views/` to change content:
- `views/home.go` — Home screen, ASCII art, bio
- `views/projects.go` — Your projects
- `views/about.go` — About section
- `views/contacts.go` — Contact links

Edit `styles.go` to change colors and styling.

## 🛠️ Tech Stack

- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** — TUI framework
- **[Wish](https://github.com/charmbracelet/wish)** — SSH server middleware
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** — Terminal styling

## 📋 Controls

| Key | Action |
|---|---|
| `← →` | Switch tabs |
| `Enter` | Open selected tab |
| `↑ ↓` | Browse projects |
| `Esc` | Go back |
| `q` | Quit |

## 📄 License

MIT
