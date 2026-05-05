# SSH Portfolio — Complete Animation Plan

> A full reference for every animation: what it looks like, how it works, implementation details, and priority order.

---

## Table of Contents

1. [Priority 1 — Ship These First](#priority-1--ship-these-first)
2. [Priority 2 — Polish Layer](#priority-2--polish-layer)
3. [Priority 3 — Showstoppers](#priority-3--showstoppers)
4. [Priority 4 — Nuclear Options](#priority-4--nuclear-options)
5. [Implementation Order Summary](#implementation-order-summary)
6. [General Bubbletea Tips](#general-bubbletea-tips)

---

## Priority 1 — Ship These First

> High impact, low-to-medium effort. These alone will make people screenshot the TUI.

---

### 1. Boot Sequence (Pre-Splash)

**Where:** Runs once, before the splash screen appears.

**What it looks like:**
```
> connecting to portfolio.mohithakshay.dev:22...
> initializing ssh-portfolio...
> authenticating identity...
> loading modules...
> decrypting portfolio...
> done.
```
Each line prints one at a time with a ~300–400ms delay. After the last line, a brief pause (500ms), then the splash screen snaps in.

**How it works:**
- On app start, enter a `BootModel` state instead of going straight to splash
- Use `tea.Tick(350ms)` to reveal one line per tick
- Keep a `lineIndex int` in the model — each tick increments it and appends the next line to a `[]string` buffer
- Once all lines are shown, wait one tick then send a `BootDoneMsg` to transition to the splash

**Feel:** Like a real system booting. Hacker aesthetic. Every visitor will read it slowly the first time.

---

### 2. Typewriter on Tagline

**Where:** Splash screen — the tagline under the name.

**What it looks like:**
The tagline types itself out character by character with a blinking cursor `█` at the end. Once complete, the cursor blinks a few times then disappears.

```
is an engineer, builder & creator█
is an engineer, builder & creator who turns ideas into products.█
is an engineer, builder & creator who turns ideas into products.
```

**How it works:**
- Store the full tagline string and a `charIndex int` in the splash model
- Use `tea.Tick(40ms)` — each tick reveals one more character: `tagline[:charIndex]`
- Append `█` to the visible string while typing is in progress
- After the last character, switch to a `cursorBlink bool` that toggles every 500ms for ~3 blinks then stops

**Feel:** Like the TUI is introducing you personally. Slow enough to read, fast enough not to be annoying.

---

### 3. Project List Cascade Drop-In

**Where:** Projects tab — the list of 7 projects.

**What it looks like:**
When the Projects tab loads, projects don't all appear at once. They drop in one by one:
```
  01. EmbedGen ...          ← appears first
    02. ISRO LEOS ...       ← 80ms later
      03. Autonomous ...    ← 80ms later
```
Each row also slides in subtly from the left (starts with a leading space that shrinks to zero).

**How it works:**
- Track `visibleCount int` in the projects model
- On tab enter, start a `tea.Tick(80ms)` — each tick increments `visibleCount` by 1
- Only render projects where `index < visibleCount`
- For the slide effect: on the tick the row first appears, prefix it with 2 spaces; next tick, 1 space; tick after, 0 spaces (fully settled)

**Feel:** Like the projects are loading in from a server. Gives the list a sense of weight and importance.

---

### 4. `[Live]` Badge Pulse

**Where:** Projects list — next to projects marked `[Live]`.

**What it looks like:**
The `[Live]` badge alternates between full brightness cyan and a dimmer cyan, like a heartbeat. Slow and subtle — not flashy. About 1 pulse per second.

```
03. Autonomous Trading System  [Live]   ← full cyan
03. Autonomous Trading System  [Live]   ← dim cyan (500ms later)
```

**How it works:**
- Global ticker at `500ms` interval
- Toggle a `livePulse bool` on each tick
- When rendering, if `livePulse == true` render `[Live]` in bright cyan; if false, render in dim cyan
- Only applies to `[Live]` badges — `[WIP]` and `[Research]` stay static

**Feel:** Like the project is actually running right now. Proves it's real and deployed.

---

### 5. SSH Handshake Fake Negotiation

**Where:** Boot sequence — replaces or extends the simple `> connecting...` lines.

**What it looks like:**
```
Connecting to portfolio.mohithakshay.dev:22...
SSH-2.0-OpenSSH_9.0
Negotiating encryption: chacha20-poly1305@openssh.com
Host key fingerprint: SHA256:xK9mPqR3vL8nT2wZ...
✓ Identity verified. Welcome, visitor.
```

**How it works:**
- Same `BootModel` tick pattern as animation #1
- Each "line" is hardcoded and reveals on its own tick interval
- The fingerprint hash is randomly generated on each session startup for authenticity
- Add a deliberate 600ms pause before the final `✓ Identity verified.` line — the tension builds

**Feel:** People who know SSH will lose their minds. People who don't will think you're a wizard. Everyone screenshots this.

---

## Priority 2 — Polish Layer

> Medium effort, great polish. These make the TUI feel premium.

---

### 6. Twinkling Star Symbols

**Where:** Splash screen — the scattered `✦` `·` `+` `*` symbols around the ASCII name.

**What it looks like:**
The symbols randomly flicker between visible and dim. Not all at once — each symbol has its own random timing, so they feel organic, like stars in a night sky.

**How it works:**
- Store each symbol as a struct: `{char string, x int, y int, bright bool, nextFlip time.Time}`
- Global ticker at `120ms`
- On each tick, loop through symbols — if `time.Now().After(nextFlip)`, toggle `bright` and set `nextFlip = now + random(500ms, 2000ms)`
- Render bright symbols in full white/cyan, dim ones in dark gray

**Feel:** The splash feels alive and breathing even when you're just sitting on it. Subtle but hypnotic.

---

### 7. Smooth Row Highlight Slide

**Where:** Projects list — the selected row indicator when browsing up/down.

**What it looks like:**
When you press ↑ or ↓, the highlight bar doesn't jump instantly. It slides smoothly from the current row to the next one over ~3 frames (about 60ms total).

**How it works:**
- Track `currentRow int` and `highlightY float64` (the visual position)
- On keypress, update `currentRow` immediately (for logic) but let `highlightY` animate toward it
- Use `tea.Tick(20ms)` — each tick move `highlightY` 33% closer to `currentRow` (lerp)
- Render the highlight bar at `int(math.Round(highlightY))`

**Feel:** Buttery smooth. Makes the list feel like a real UI, not a terminal. People will scroll up and down just to watch it.

---

### 8. Scroll Momentum (Rubber Band)

**Where:** Projects list — when you press ↑/↓ rapidly.

**What it looks like:**
When you scroll fast, the list slightly overshoots by one row, then snaps back — like iOS rubber banding. Only noticeable if you're moving quickly; feels natural otherwise.

**How it works:**
- Track a `float64 velocity` that increases with each rapid keypress and decays over time
- On each tick, apply `velocity` to `highlightY` and multiply `velocity` by 0.85 (friction)
- Clamp `highlightY` to valid row bounds; if it hits a bound, reverse a small fraction of velocity for the bounce feel
- Works together with animation #7's lerp

**Feel:** Tiny detail. Huge payoff. Makes the list feel physical.

---

### 9. Column Wipe Tab Transition

**Where:** Between all tabs — Projects, About, Contacts.

**What it looks like:**
When switching tabs, the current content "wipes" off screen column by column from left to right, then the new content wipes in from left to right. Each wipe takes ~200ms total.

```
Before:  [Projects content fully visible]
Frame 1: [██ rojects content visible    ]
Frame 2: [████ jects content visible    ]
Frame 3: [██████ ects content visible   ]
...
After:   [New tab content fully visible ]
```

**How it works:**
- Track `wipeProgress float64` from 0.0 to 1.0 during transition
- On each tick (20ms), increment `wipeProgress` by 0.1
- When rendering, mask out columns to the left of `wipeProgress * termWidth` with `█` or spaces
- At `wipeProgress >= 1.0`, snap to new content and start the wipe-in phase

**Feel:** Each tab feels like its own world. The transition gives a sense of physical space between sections.

---

### 10. Contacts Stagger Reveal

**Where:** Contacts page — the `(@)` `(~)` `(in)` `(gh)` entries.

**What it looks like:**
When you enter the Contacts tab, each entry appears top to bottom with a ~120ms delay between them.

```
(@)  Email          ← appears first
                    ← 120ms gap
(~)  Website        ← appears second
                    ← 120ms gap
(in) LinkedIn       ← appears third
```

**How it works:**
- Same `visibleCount` pattern as the project cascade
- `tea.Tick(120ms)` — each tick reveals one more contact entry
- No slide effect — just clean appear (the contacts page is minimal, keep the animation minimal too)

**Feel:** Clean and intentional. Doesn't overstay its welcome.

---

### 11. "You're viewing this over SSH!" Flash

**Where:** Contacts page — the meta line at the bottom.

**What it looks like:**
When the Contacts page first loads, the line flashes bright cyan twice, then settles into its normal color. Like a notification ping.

```
You're viewing this over SSH!   ← bright cyan flash
You're viewing this over SSH!   ← normal color (200ms later)
You're viewing this over SSH!   ← bright cyan again
You're viewing this over SSH!   ← settles to normal
```

**How it works:**
- On Contacts page enter, set `sshFlash int = 4` (4 frames of flash)
- `tea.Tick(200ms)` — each tick decrements `sshFlash`
- While `sshFlash > 0` and `sshFlash % 2 == 0`, render the line in bright cyan; otherwise normal color
- After `sshFlash == 0`, stop the ticker for this element

**Feel:** Draws attention to the coolest thing about the portfolio — that you're reading this over SSH. A little wink to the visitor.

---

### 12. Tech Stack Tag Pop

**Where:** Project detail view — the `[Python]` `[PyTorch]` `[LoRA]` tags.

**What it looks like:**
When you open a project, the tech stack tags don't all appear at once. They pop in left to right, one by one, like cards being dealt:
```
[Python]
[Python] [PyTorch]
[Python] [PyTorch] [LoRA]
[Python] [PyTorch] [LoRA] [GGUF]
```

**How it works:**
- Track `visibleTags int` per project detail view
- On view enter, start `tea.Tick(35ms)` — each tick increments `visibleTags`
- Render only the first `visibleTags` tags from the tag slice
- Total duration is `numTags * 35ms` — fast enough to feel snappy, slow enough to register

**Feel:** Adds energy to the detail view. Makes the stack feel like a flex being laid out one piece at a time.

---

## Priority 3 — Showstoppers

> Hardest to implement. Save for when everything else is done. These are the ones that go viral.

---

### 13. Matrix Rain → Name Solidification

**Where:** Pre-splash — replaces the simple boot sequence entirely.

**What it looks like:**
Green falling characters rain down the screen for ~2 seconds. Then, one by one, the characters at positions that belong to the final "MOHITH AKSHAY DUGGIRALA" ASCII art freeze and turn cyan — like your name is being assembled from chaos. The rain continues in the background until all name positions are locked, then fades out.

```
Phase 1 (0–2s):    Falling green chars everywhere
Phase 2 (2–3s):    Chars at name positions start "locking" to the correct glyph in cyan
Phase 3 (3–3.5s):  Rain fades out, locked name stays in cyan
Phase 4 (3.5s+):   Full splash screen loads
```

**How it works:**
- Pre-compute a set of `(x, y)` positions that belong to each ASCII art character of your name
- Render the "rain": each column has a `[]rune` buffer that scrolls down every 50ms, filled with random chars from `[]rune{'ｦ','ｱ','ｳ','ｴ','ｵ','ｶ','ｷ',...}` (katakana for authenticity)
- Maintain a `lockedCells map[Point]rune` — the final correct character at each name position
- After 2s, start a `lockTicker` at 20ms — each tick randomly picks 5–10 unlocked name-position cells and "locks" them (they stop raining and render in cyan with the correct final char)
- Once all name cells are locked (~1s), fade out the remaining rain by reducing brightness over 500ms
- Transition to full splash

**Feel:** The most dramatic entrance possible in a terminal. Nobody has seen their name crystallize out of falling characters. This is the screenshot moment. Every person who sees this will show it to someone else.

---

### 14. Name Glitch Effect

**Where:** Splash screen — "MOHITH AKSHAY DUGGIRALA" ASCII art, on first load.

**What it looks like:**
After the splash loads (or after the matrix rain), the ASCII name briefly corrupts — random characters replace real ones for 2–3 frames — then snaps back to the real name cleanly. Like a signal locking in. Happens only once per session.

```
Frame 1: M#H!TH  AK$H%Y  DU&&IR!L@   ← corrupted
Frame 2: MO#ITH  AK#HAY  DU#GIRA#A   ← partially recovered
Frame 3: MOHITH  AKSHAY  DUGGIRALA   ← clean
```

**How it works:**
- Store the ASCII name as a `[][]rune` (2D grid of characters)
- On splash enter, start a `glitchFrames int = 3` counter
- Each frame (tick at `80ms`): replace ~20% of non-space runes with random chars from `[]rune{'#','@','%','$','!','?','█','▓'}`
- After 3 frames, render the original clean name and stop the glitch ticker
- Only triggers once — set a `glitchDone bool` flag

**Feel:** Like your identity is being decoded. Very hacker aesthetic. Pairs perfectly with the Matrix rain if both are implemented.

---

### 15. Radar Sweep Project Reveal

**Where:** Projects tab — the list of 7 projects.

**What it looks like:**
A single `│` sweep line travels left to right across the terminal. Projects to the left of the line are visible; projects to the right are rendered as `░░░░░░░░░` placeholder blocks. As the sweep line passes each project row, it "decrypts" and reveals the real content — like a sonar ping.

```
Frame 5:   01. EmbedGen...      │  ░░░░░░░░░░░░░░░
Frame 10:  01. EmbedGen...         02. ISRO...    │  ░░░░░░░░
Frame 15:  01. EmbedGen...         02. ISRO...       03. Autonomous...
```

**How it works:**
- Track `sweepX float64` — the current X position of the sweep line
- `tea.Tick(16ms)` — each tick advance `sweepX` by `termWidth / 50` (full sweep in ~800ms)
- For each project row, compare its content width to `sweepX`; render real content if `sweepX` has passed, render `░` blocks otherwise
- The `│` sweep line is rendered as a vertical bar at `int(sweepX)` across all rows
- After sweep completes, remove the bar and show all content clean

**Feel:** One of the most original animations on this list. Looks like actual radar or sonar. Nobody does this in a terminal portfolio.

---

### 16. Particle Drift on Dot Portrait

**Where:** Splash screen — the large dot particle portrait on the left side.

**What it looks like:**
The thousands of cyan dots that form the portrait slowly drift — each dot moves 1–2 characters in a random direction over several seconds, then drifts back. The overall shape stays recognizable but feels like it's breathing or vibrating with energy.

**How it works:**
- Parse the dot art into a list of `{x, y, offsetX, offsetY, targetX, targetY}` structs
- Global ticker at `100ms`
- Each tick, move each dot 10% closer to its target offset (lerp)
- Every 3–5 seconds, assign each dot a new random target offset (max ±2 chars)
- Render dots at their current offset positions

**Important:** Cap the number of animated dots at ~200 (sample from the full set) to avoid terminal performance issues. The effect reads the same with a subset.

**Feel:** The portrait feels alive. Like a living watermark. This is the animation that will get recorded and posted online.

---

## Priority 4 — Nuclear Options

> Low implementation effort but massive psychological impact. These make visitors physically react.

---

### 17. "Unauthorized Access" Fake Alert

**Where:** 500ms after the TUI fully loads — a full-screen interrupt.

**What it looks like:**
The entire terminal goes red. A large ASCII warning fills the screen:

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│   ⚠  UNAUTHORIZED ACCESS DETECTED                  │
│                                                     │
│   Tracing connection origin...                      │
│   Logging session metadata...                       │
│   Notifying system administrator...                 │
│                                                     │
└─────────────────────────────────────────────────────┘
```

Then after 1.5 seconds, the entire thing cuts to black, and in small monospace text at the center:

```
  just kidding. welcome. :)
```

Then the normal splash loads.

**How it works:**
- After `BootDoneMsg`, instead of going straight to splash, enter `FakeAlertModel` for 2.5s
- Render a full-screen red box using `lipgloss` with the warning text
- After 1.5s, switch to the `justKidding` state and render the small reassurance text
- After another 800ms, transition to normal splash

**Feel:** People will literally yell out loud. Then immediately laugh. Then screenshot it. Then show their friends. This is the most shareable moment in the whole portfolio.

---

### 18. Self-Aware Time-Based About Text

**Where:** About section — a single dynamic line that changes based on time of day (IST).

**What it looks like:**
The About section has one line that reacts to when the visitor is reading it:

```
2:00 AM – 5:00 AM:   "you're up late. so am I, probably."
5:00 AM – 9:00 AM:   "early start. respect."
9:00 AM – 6:00 PM:   "currently probably in class or debugging something."
6:00 PM – 10:00 PM:  "golden hours. this is when the best code gets written."
10:00 PM – 2:00 AM:  "late night build session energy in here."
```

**How it works:**
- On About view render, call `time.Now().In(ist)` to get current IST time
- Switch on `hour` to select the appropriate string
- No animation needed — just a static string swap. Zero engineering effort.

**Feel:** Costs nothing to build. People will feel seen. They'll check it at different times of day just to see what it says. Wildly personal for something that's technically just a `switch` statement.

---

### 19. Live GitHub Commit on Splash

**Where:** Splash screen — small text line at the very bottom.

**What it looks like:**
```
last pushed: "fix: typewriter cursor timing" · 3h ago
```
Updates every time the portfolio is deployed. Shows you're actively maintaining it.

**How it works:**
- On startup, hit `https://api.github.com/repos/trafalgar-2006/ssh-portfolio/commits?per_page=1`
- Parse the commit message and timestamp from the JSON response
- Format as `last pushed: "{message}" · {relative time}`
- Cache the result for the session so it only fetches once
- Render in dim gray at the bottom of the splash — subtle, not the focus

**Feel:** Proves the portfolio is alive and actively maintained. Shows discipline. The commit message itself becomes part of the presentation — you'll start writing commit messages with visitors in mind.

---

### 20. Corrupted Boot Recovery

**Where:** Mid boot sequence — a deliberate "glitch" before recovery.

**What it looks like:**
```
> initializing ssh-portfolio...
> authenticating identity...
> l̷o̸a̵d̷i̵n̷g̸ ̸m̸o̷d̷u̴l̸e̵s̸.̵.̵.̸       ← text glitches mid-line
> ERROR: signal lost — attempting recovery...
> .
> ..
> ...
> recovery successful. resuming.
> decrypting portfolio...
> done.
```

**How it works:**
- The third boot line renders with ~15% of characters replaced by zalgo-style combining diacritics
- After it renders, wait 300ms then flash a red `> ERROR: signal lost — attempting recovery...` line
- Show the three `> .` `> ..` `> ...` lines with 400ms between them (tension build)
- Then `> recovery successful.` in bright green, and continue normally

**Feel:** Makes the connection feel genuinely unstable. The recovery makes visitors feel like they just witnessed something real. Pairs perfectly with the SSH handshake negotiation.

---

## Implementation Order Summary

| # | Animation | Location | Effort | Mind-Blow |
|---|-----------|----------|--------|-----------|
| 1 | Boot sequence | Pre-splash | Low | 🔥🔥🔥 |
| 2 | Typewriter tagline | Splash | Low | 🔥🔥🔥 |
| 3 | Project cascade drop-in | Projects | Low | 🔥🔥🔥 |
| 4 | `[Live]` badge pulse | Projects | Low | 🔥🔥 |
| 5 | SSH handshake fake negotiation | Pre-splash | Low | 🔥🔥🔥🔥 |
| 6 | Twinkling star symbols | Splash | Medium | 🔥🔥 |
| 7 | Smooth row highlight slide | Projects | Medium | 🔥🔥 |
| 8 | Scroll momentum rubber band | Projects | Medium | 🔥🔥🔥 |
| 9 | Column wipe tab transition | All tabs | Medium | 🔥🔥🔥 |
| 10 | Contacts stagger reveal | Contacts | Low | 🔥 |
| 11 | SSH line flash | Contacts | Low | 🔥 |
| 12 | Tech stack tag pop | Project detail | Low | 🔥🔥 |
| 13 | Matrix rain → name solidify | Pre-splash | High | 🤯🤯🤯🤯🤯 |
| 14 | Name glitch effect | Splash | High | 🤯🤯🤯 |
| 15 | Radar sweep project reveal | Projects | High | 🤯🤯🤯🤯 |
| 16 | Particle drift dot portrait | Splash | High | 🤯🤯🤯🤯 |
| 17 | "Unauthorized access" fake alert | Pre-splash | Low | 🤯🤯🤯🤯🤯 |
| 18 | Self-aware time-based text | About | Trivial | 🤯🤯🤯 |
| 19 | Live GitHub commit on splash | Splash | Medium | 🤯🤯🤯🤯 |
| 20 | Corrupted boot recovery | Pre-splash | Low | 🤯🤯🤯 |

---

## General Bubbletea Tips

- **Never run multiple independent tickers** if you can share one global tick command — one global `tea.Tick(16ms)` that all animations respond to is cleaner than 8 separate tickers
- **Always gate animations behind a `done bool` flag** so they don't re-trigger on window resize
- **Test on Windows Terminal, macOS Terminal, and a Linux TTY** — animation timing feels different across them; what looks smooth on Windows may strobe on a slow TTY
- **Keep total tick rate above 60ms for SSH connections** — anything faster may overwhelm high-latency connections and cause visual stuttering; use a `isSlowConnection bool` flag to disable the heaviest animations
- **For the matrix rain**, limit the render area to what's actually visible — rendering offscreen characters wastes CPU
- **For particle drift**, sample ~200 dots from the full set — the effect reads identically and saves significant render time
- **Use `lipgloss.NewStyle().Foreground()`** for color toggling in the pulse/flash effects — it's faster than manual ANSI escape codes
- **The `Unauthorized Access` alert** needs a `once.Do()` guard so it only fires on first session, not on every window resize event
- **Commit messages for the live GitHub feed** — start writing them as if visitors will read them. `fix: typewriter cursor was 5ms too fast` is way cooler than `minor fixes`

---

*Built with Go + Bubbletea. Deployed on Railway. Accessible via `ssh mohithakshay.dev`*
