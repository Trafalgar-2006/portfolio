# SSH TUI Portfolio — Improvement Tracker
> Last updated: 2026-05-05

---

## ✅ Completed (26 items)

| # | Item | Notes |
|---|------|-------|
| 1 | SSH host key fix | Generated at Docker build time |
| 2 | Truecolor colors | SSH session renderer per connection |
| 3 | Resume tab removed | Decision made — absorbed into About |
| 4 | Experience timeline in About | ISRO, Webcraft Studios, SnuqSq with full bullets |
| 5 | Skills in About | Re-added after Resume tab removal |
| 6 | PDF download link | About page footer → GitHub raw URL |
| 7 | EmbedGen as Project #1 | PyTorch · LoRA · GGUF · CUDA — `[WIP]` |
| 8 | Trading System added | Alpaca API · SQLite · Oracle Cloud — `[Live]` |
| 9 | ARAK removed | Decided not needed |
| 10 | SSH Portfolio added as project | Go · Docker · Railway — `[Live]` |
| 11 | Status badges | `[Live]` `[WIP]` `[Research]` per project |
| 12 | Color-coded tags | Cyan=lang, Magenta=ML, Green=cloud, Orange=DB, Purple=fw |
| 13 | GitHub links in projects | Shown in expanded project view |
| 14 | "Things I've built" | Subtitle changed from overused "shipped" framing |
| 15 | Banner "D" fix | Removed broken glyph → `· D U G G I R A L A ·` subtitle |
| 16 | Tab highlight | Active tab: cyan background + dark text |
| 17 | ASCII contacts | `(@)` `(~)` `(in)` `(gh)` — cross-platform |
| 18 | Specific interests | About mentions ISRO, trading system, SSH portfolio |
| 19 | Health check | HTTP server on `:8080` → `/health` returns 200 OK |
| 20 | Duplicate skills resolved | Only in About now |
| 21 | Typewriter animation | Banner reveals line-by-line at 55ms |
| 22 | Portrait animation | Reveals at 2× banner speed |
| 23 | Blinking stars | Decorative stars blink at 600ms after reveal |
| 24 | SSH_PORT fix | Decoupled from Railway's `PORT` env var |
| 25 | UptimeRobot keep-alive | TCP monitor on `trolley.proxy.rlwy.net:41074` @ 5 min |
| 26 | Railway networking fix | HTTP domain → 8080 · TCP proxy → 23234 (ext: 41074) |

---

## 🔧 To Do

### 🔴 SSH known_hosts annoyance
Every redeploy regenerates the host key → users get the scary "REMOTE HOST IDENTIFICATION HAS CHANGED" warning and have to clear it manually.

**Fix options:**
- A) Add `StrictHostKeyChecking no` + `UserKnownHostsFile /dev/null` to user's `~/.ssh/config` for this host (client-side, instant fix)
- B) Use a **persistent SSH host key** stored in Railway environment variable instead of regenerating each build (server-side, permanent fix — recommended)

```bash
# Option A — add to ~/.ssh/config:
Host trolley.proxy.rlwy.net
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null

# Option B — generate key once, store in RAILWAY env var:
# 1. ssh-keygen -t ed25519 -f host_key -N ""
# 2. Add HOST_KEY_B64 = base64(host_key) to Railway vars
# 3. Dockerfile: decode HOST_KEY_B64 → .ssh/id_ed25519 at startup
```

---

### 🟡 SQLite Visitor Analytics
Log every session: timestamp, IP, sections visited, duration.

```go
func analyticsMiddleware(next ssh.Handler) ssh.Handler {
    return func(s ssh.Session) {
        start := time.Now()
        next(s)
        logSession(s.RemoteAddr(), time.Since(start))
    }
}
```
**Effort:** ~2 hrs · Great interview talking point ("X visitors in Y weeks")

---

### 🟡 SSH Key Admin View
Your own SSH key → secret admin panel with visitor stats.

```go
wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
    return wish.KeysEqual(key, adminKey)
}) → route to AdminModel
```
**Effort:** ~4 hrs · Very impressive for interviews

---

### 🟡 YAML Content Config
Projects and bio hardcoded — every update = rebuild + redeploy.

```yaml
projects:
  - title: EmbedGen
    status: WIP
    github: github.com/trafalgar-2006/EmbedGen
    tags: [Python, PyTorch, LoRA]
```
**Effort:** ~2 hrs · Makes updates instant without touching Go code

---

### 🟢 Visitor Counter Easter Egg
"You are visitor #N" on the splash screen. Easy once analytics is in place.
**Effort:** ~30 min (after SQLite middleware)

---

### 💡 Guestbook Tab
Visitors leave a one-liner. Stored in SQLite. Visible to next visitor.
Viral potential — people share things with their name in them.
**Effort:** ~3 hrs

---

## 📊 Current Status

```
Live URL:       ssh trolley.proxy.rlwy.net -p 41074
GitHub:         github.com/Trafalgar-2006/portflio
Auto-deploy:    ✅ Railway (push to master → rebuild)
Health check:   ✅ http://portflio-production.up.railway.app/health
Keep-alive:     ✅ UptimeRobot TCP @ 5 min (port 41074)
Colors:         ✅ SSH session renderer
Animation:      ✅ Typewriter + blinking stars
Projects:       7 items (EmbedGen #1)
Tabs:           Projects · About · Contacts

Known annoyance: Host key regenerates on every redeploy
Recommended next: Fix persistent host key (Option B above)
```
