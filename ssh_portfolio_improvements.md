# SSH TUI Portfolio — Improvement Tracker

**4** Critical fixes · **6** UI polish · **5** Content gaps · **6** Backend ideas

---

## Critical Fixes

### 🔴 Splash name truncated — "D" looks broken
The ASCII banner shows "MOHITH AKSHAY D" but the last initial renders as a lone, half-formed glyph. Either include your full surname "DUGGIRALA" (sized down) or drop the initial entirely. A clean two-line splash is better than a cut-off three-liner.

```
Option: bigletters.SetFont(...) with smaller font size for full surname
```

---

### ~~🔴 Typo in resume PDF URL — "portflio"~~
> **Resolved** — URL matches actual repo name `portflio`. No action needed.

---

### 🔴 No interactive project drill-down
Pressing Enter on a project in the list does nothing (or isn't shown). This is the most expected interaction — recruiter selects a project, expects to see a detail view with GitHub link, live demo URL, longer description. Right now it's just a list.

```
Add ProjectDetail view: description, links, stack, status (Live / WIP)
```

---

### ~~🔴 Resume and About sections duplicate skills~~
> **Resolved** — Resume tab removed entirely. Skills now live only in About.

---

### ~~🔴 No loading/error state for SSH connection drops~~
> **Deferred** — Railway keep-alive via UptimeRobot reduces cold-starts significantly.

---

## UI / UX

### 🟡 Add a subtle animated intro sequence
Right now the splash loads instantly. A typewriter effect on your name (character-by-character) or a fade-in on the ASCII art would make the first impression far more memorable. Bubbletea supports tickers and viewport updates — use them.

```
tea.Tick(50ms) → reveal rune-by-rune on the banner
```

---

### 🟡 Navigation tab highlight is too subtle
The active tab uses a diamond prefix. It's easy to miss at a glance. Consider inverting the background — highlight color background with dark text — so the active tab is immediately obvious, especially on terminals with non-default color schemes.

---

### 🟡 Skill tags on projects feel flat — add color coding
All stack tags like `[React]` `[Python]` `[Go]` render in the same color. Color-coding by category (language = cyan, framework = yellow, cloud = green) would make the skill spread scannable at a glance without reading each tag.

---

### 🟢 Contacts page emoji icons render inconsistently
The emoji icons (📧 🌐 💼 🔔) before each contact entry look different on Windows Terminal vs macOS Terminal vs Linux. Use pure-ASCII alternatives (`@` for email, `>` for links) or Nerd Font glyphs with a fallback for consistent cross-platform rendering.

---

### 🟢 Footer "esc to go back" hint is inconsistent
Some views show it, some don't. Also on the main nav, the hint says "q to quit" but nowhere explains that arrow keys also change tabs. Standardize hints across all views and add a persistent one-liner at the very bottom of every screen.

---

### 🟢 Add a visitor counter easter egg
A small "You are visitor #N" counter displayed on the splash or contacts page would be a fun nerd detail that also proves the app is live and seeing real traffic. Store count in a simple file or Redis on Railway.

---

## Content

### ~~🔴 Resume tab removed — decided~~
> **Decision** — Resume tab dropped. Experience timeline and skills absorbed into About. PDF download link added to About page footer. EmbedGen remains prominently in Projects as #1.

---

### ~~🔴 EmbedGen project missing from the list~~
> **Resolved** — EmbedGen is now Project #1 in the Projects tab.

---

### 🟡 "Shipped" framing is overused and weakened
"Things I've built and shipped" sets a high bar. Frog Call Classifier and some others are research/course projects, not deployed products. Either change the section subtitle to "Things I've built" or add a `[Live]` / `[Research]` / `[WIP]` status badge per project so the distinction is honest and clear.

---

### 🟡 No GitHub links on any project
Every project should have a clickable GitHub URL in its detail view. Recruiters expect this. It's also how they verify the claims you make about accuracy metrics and real-time performance figures.

---

### 🟢 About page "Interests" blurb is generic
"Building products that bridge the gap between research and real-world deployment" reads like every other ML engineer's bio. Make it specific: mention ISRO, the trading system, the SSH portfolio itself as proof of those interests rather than just stating them.

---

## Backend & Infra

### 🔴 Railway free tier will sleep — add a keep-alive ping
Railway's Hobby plan puts containers to sleep after inactivity. Add a simple cron job (or UptimeRobot free tier) that pings your SSH port every 5 minutes to prevent cold starts. Nobody wants to wait 10 seconds for the TUI to wake up.

```
UptimeRobot → TCP monitor → trolley.proxy.rlwy.net:23115
```

---

### 🔴 No persistent visitor/analytics log
You have zero visibility into who's connecting, from where, how long they stay, which sections they visit. Add a lightweight SQLite-backed logger inside the Go app — log timestamp, section visited, session duration. Gives you real interview talking points ("500 unique visitors in 2 weeks").

```
charm/wish middleware → log session events to SQLite on Railway volume
```

---

### 🟡 Add SSH key-based auth for a secret admin view
Public visitors get the portfolio. If your own SSH key is detected, you get an extra admin panel showing visitor stats, session logs, and a live "update content" mode. This is a genuinely impressive feature that shows you understand SSH auth primitives.

```
wish.WithPublicKeyAuth() → check key fingerprint → route to AdminModel
```

---

### 🟡 Content is hardcoded — move to a config file
Projects, bio, contacts are all compiled into the binary. Every update requires a rebuild and redeploy. Move content to a YAML or JSON file loaded at startup — Railway volumes can persist it. Now you can update your portfolio without touching Go code.

```
embed config.yaml → viper or plain encoding/yaml → hot-reload on SIGHUP
```

---

### 🟢 No health check endpoint for Railway
Railway needs an HTTP health check to know the container is healthy. Spin up a minimal `net/http` server on a second port (e.g. 8080) that just returns `200 OK`. Without it, Railway can't distinguish a healthy app from a crashed one.

```go
go func() { http.ListenAndServe(":8080", healthHandler) }()
```

---

### 💡 Add a live "guestbook" section
Let visitors leave a short note (name + one line) that's visible to the next visitor. Stored in SQLite. Shows up as a scrollable list in a new Guestbook tab. It's the kind of delightful detail that makes people share the project link — perfect for virality on Twitter/X and LinkedIn.
