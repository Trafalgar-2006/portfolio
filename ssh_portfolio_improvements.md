# SSH TUI Portfolio — Improvement Tracker
> Last updated: 2026-05-05

---

## ✅ Completed

| # | Item | Notes |
|---|------|-------|
| 1 | **SSH host key fix** | Key generated at Docker build time — no more connection reset |
| 2 | **Truecolor fix** | SSH session renderer (`bubbletea.MakeRenderer`) passed to all views — colors now work |
| 3 | **Resume tab removed** | Decision made — content absorbed into About |
| 4 | **Experience timeline in About** | ISRO, Webcraft Studios, SnuqSq now with full bullet detail |
| 5 | **Skills in About** | Re-added after Resume tab removal (AI/ML, Languages, Web, Tools) |
| 6 | **PDF download link** | In About page footer → GitHub raw URL |
| 7 | **EmbedGen added as Project #1** | PyTorch · LoRA · GGUF · CUDA — `[WIP]` badge |
| 8 | **Trading System added** | Alpaca API · SQLite · Oracle Cloud — `[Live]` badge |
| 9 | **ARAK added** | React · Firebase · TypeScript — `[Live]` badge |
| 10 | **SSH Portfolio added** | Go · Docker · Railway — `[Live]` badge |
| 11 | **Status badges** | `[Live]` `[WIP]` `[Research]` per project |
| 12 | **Color-coded tags** | Cyan=language, Magenta=ML/AI, Green=cloud, Orange=DB, Purple=framework |
| 13 | **GitHub links in Projects** | Shown in expanded project view |
| 14 | **"Things I've built"** | Subtitle changed from "shipped" to honest framing |
| 15 | **Banner "D" fix** | Removed lone broken glyph → `· D U G G I R A L A ·` subtitle |
| 16 | **Tab highlight** | Active tab: cyan background + dark text (inverted, clearly visible) |
| 17 | **ASCII contacts** | `(@)` `(~)` `(in)` `(gh)` — cross-platform safe |
| 18 | **Specific interests** | About page mentions ISRO, trading system, SSH portfolio by name |
| 19 | **Health check** | HTTP server on `:8080` → `/health` returns `200 OK` for Railway |
| 20 | **Duplicate skills** | Resolved — only in About now |
| 21 | **GitHub push** | All code live at `github.com/Trafalgar-2006/portflio` |
| 22 | **Railway deployment** | Auto-deploys on every push to master |

---

## 🔧 In Progress / Pending

---

### 🔴 UptimeRobot Keep-Alive (DO THIS NOW)
Railway Hobby plan sleeps containers after inactivity. Without this, first-time visitors wait 10+ seconds.

**Setup (5 min, free):**
1. Go to **https://uptimerobot.com** → Sign up free
2. Click **Add New Monitor**
3. Monitor Type: **TCP Port**
4. Friendly Name: `SSH Portfolio`
5. Hostname: `trolley.proxy.rlwy.net`
6. Port: `23115`
7. Monitoring Interval: **5 minutes**
8. Click **Create Monitor** ✅

> Also add an HTTP monitor for the health check:
> - Type: HTTP(s)
> - URL: `http://portflio-production.up.railway.app/health`
> - Interval: 5 min

---

### 🟡 Typewriter Animation on Intro Screen
Reveal the name banner line-by-line when a visitor first connects. Makes the first impression memorable.

**Implementation plan:**
```go
// In model.go — add to Model struct:
revealIdx int  // how many banner lines to show
animDone  bool // true once fully revealed

// Init() → start ticker:
func (m Model) Init() tea.Cmd {
    return tea.Tick(time.Millisecond*55, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// Update() → handle tick:
case tickMsg:
    if m.revealIdx < totalBannerLines {
        m.revealIdx++
        return m, tea.Tick(55ms, ...)
    }
    m.animDone = true
    return m, nil

// RenderHome() → pass revealIdx:
for i, line := range nameBanner {
    if i < revealIdx {
        rightCol.WriteString(cyanStyle.Render(line) + "\n")
    } else {
        rightCol.WriteString("\n") // hold space
    }
}
```

**Estimated effort:** ~1 hour · High visual impact

---

### 🟡 SSH Key Admin View
If your own SSH public key connects, route to a secret admin panel showing session stats.

```go
wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
    return wish.KeysEqual(key, adminKey)
}) → route to AdminModel
```

**Estimated effort:** Half a day · Very impressive for interviews

---

### 🟡 SQLite Visitor Analytics
Log every session: timestamp, IP, sections visited, duration. Zero infra cost on Railway volume.

```go
// wish middleware
func analyticsMiddleware(next ssh.Handler) ssh.Handler {
    return func(s ssh.Session) {
        start := time.Now()
        next(s)
        logSession(s.RemoteAddr(), time.Since(start))
    }
}
```
Output: `"500 unique visitors in 2 weeks"` → great talking point.

**Estimated effort:** ~2 hours

---

### 🟡 YAML Content Config
Projects and bio are hardcoded — every update = rebuild + redeploy. Move to `content.yaml`.

```yaml
projects:
  - title: EmbedGen
    status: WIP
    github: github.com/trafalgar-2006/EmbedGen
    tags: [Python, PyTorch, LoRA]
```

**Estimated effort:** ~2 hours · Makes the portfolio self-serve

---

### 🟢 Visitor Counter Easter Egg
"You are visitor #N" on the splash screen. Stored in a flat file or SQLite.
Easy win once analytics middleware is in place.

---

### 💡 Guestbook Tab
Let visitors leave a one-liner. Stored in SQLite, visible to the next visitor.
Viral potential — people share things that have their name in them.

---

## 📊 Status Summary

```
Total items:     22 done · 7 remaining
Colors working:  ✅ (SSH session renderer)
Live URL:        ssh trolley.proxy.rlwy.net -p 23115
GitHub:          github.com/Trafalgar-2006/portflio
Auto-deploy:     ✅ Railway (pushes to master trigger rebuild)
Health check:    ✅ http://portflio-production.up.railway.app/health
Keep-alive:      ⏳ UptimeRobot not set up yet
Animation:       ⏳ Not implemented yet
```
