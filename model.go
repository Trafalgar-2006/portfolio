package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trafalgar-2006/ssh-portfolio/views"
)

// Build info — injected via ldflags at build time
var (
	BuildVersion = "dev"
	BuildCommit  = "unknown"
	BuildDate    = "unknown"
)

type View int

const (
	ViewMatrix View = iota // NEW: matrix rain pre-splash
	ViewBoot
	ViewAlert
	ViewHome
	ViewProjects
	ViewAbout
	ViewContacts
)

var tabNames = []string{"Projects", "About", "Contacts"}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const numStars = 8

// StarState tracks per-star independent twinkle timing
type StarState struct {
	Bright   bool
	FlipAt   int // tickCount when this star should next flip
}

type Model struct {
	renderer  *lipgloss.Renderer
	width     int
	height    int
	quitting  bool
	tickCount int

	// Matrix rain animation (pre-boot)
	matrixCols    []views.MatrixColumn
	matrixLocked  map[[2]int]rune
	matrixPending [][2]int
	matrixPhase   int
	matrixNameX   int
	matrixNameY   int

	// Boot animation
	bootLines     []views.BootLine
	bootSchedule  []int
	bootVisible   int
	bootStartTick int

	// Alert animation
	alertPhase     int
	alertPhaseTick int

	// Home / Splash
	currentView  View
	activeTab    int
	revealIdx    int
	glitchFrames int
	glitchRunes  [][]rune
	taglineIdx   int
	taglineDone  bool
	cursorBlink  bool
	cursorLeft   int
	lastCommit   string

	// Independent star twinkle
	stars [numStars]StarState

	// CRT scanline sweep
	scanlineY int // -1 = disabled, 0..height = current line

	// Idle ambient glitch
	idleTicks  int  // resets on any keypress
	idleGlitch bool // true for exactly 1 tick (one-frame flicker)

	// Session info
	sessionID    string
	sessionStart time.Time

	// Theme
	themeIdx   int // index into views.Themes
	themeFlash int // counts down from 4, triggers flash overlay

	// Wipe transition
	wipePhase   int
	wipeLines   int
	pendingView View
	pendingTab  int

	// Projects
	projectCursor  int
	projectScroll  int
	projectsReveal int
	tagPopReveal   int
	lastCursor     int
	livePulse      bool

	// Projects — smooth highlight + momentum
	highlightY float64 // visual lerp position of the selection bar
	velocity   float64 // scroll momentum (decays each tick)

	// Projects — decrypt reveal on open
	decryptIdx   int      // chars revealed so far in description
	decryptRunes []rune   // scrambled desc runes, resolved left-to-right

	// Vim number prefix buffer (e.g. "3" before j)
	numBuf string

	// Contacts
	contactsReveal int
	sshFlash       int
}

func NewModel(r *lipgloss.Renderer) Model {
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}
	lines, sid := views.NewBootSequence()
	schedule := make([]int, len(lines))
	t := 0
	for i, l := range lines {
		schedule[i] = t
		t += l.DelayTicks
	}

	// Matrix: initialise with default terminal size; resized on first WindowSizeMsg
	w, h := 80, 24
	nameX, nameY := views.MatrixNameOrigin(w, h)
	allCells := views.ComputeNameCells(nameX, nameY, views.NameBannerLines())
	pending := make([][2]int, 0, len(allCells))
	for pos := range allCells {
		pending = append(pending, pos)
	}
	// Shuffle pending so cells lock in random order
	rand.Shuffle(len(pending), func(i, j int) { pending[i], pending[j] = pending[j], pending[i] })

	// Initialise stars with random independent flip times
	var stars [numStars]StarState
	for i := range stars {
		stars[i] = StarState{
			Bright: rand.Intn(2) == 0,
			FlipAt: rand.Intn(30) + 10, // 0.5s–1.5s first flip
		}
	}

	return Model{
		renderer:     r,
		width:        w,
		height:       h,
		currentView:  ViewMatrix,
		matrixCols:   views.NewMatrixColumns(w, h),
		matrixLocked: make(map[[2]int]rune),
		matrixPending: pending,
		matrixNameX:  nameX,
		matrixNameY:  nameY,
		bootLines:    lines,
		bootSchedule: schedule,
		bootVisible:  0,
		cursorLeft:   6,
		sessionID:    sid,
		sessionStart: time.Now(),
		stars:        stars,
		scanlineY:    -1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), fetchCommit())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Reset idle counter on any keypress
		m.idleTicks = 0
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "t":
			// Cycle theme with flash
			m.themeIdx = (m.themeIdx + 1) % len(views.Themes)
			m.themeFlash = 4
			return m, nil

		case "q":
			if m.currentView == ViewHome {
				m.quitting = true
				return m, tea.Quit
			}
			m.startWipe(ViewHome, m.activeTab)
			return m, nil

		case "esc":
			if m.currentView != ViewHome && m.currentView != ViewBoot && m.currentView != ViewAlert {
				m.startWipe(ViewHome, m.activeTab)
				return m, nil
			}
			if m.currentView == ViewHome {
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil

		case "left", "h":
			if m.currentView == ViewHome {
				m.activeTab--
				if m.activeTab < 0 {
					m.activeTab = len(tabNames) - 1
				}
			}
			return m, nil

		case "right", "l":
			if m.currentView == ViewHome {
				m.activeTab++
				if m.activeTab >= len(tabNames) {
					m.activeTab = 0
				}
			}
			return m, nil

		case "up", "k":
			if m.currentView == ViewProjects {
				steps := m.consumeNum(1)
				prev := m.projectCursor
				m.projectCursor -= steps
				if m.projectCursor < 0 {
					m.projectCursor = 0
				}
				if m.projectCursor != prev {
					m.tagPopReveal = 0
					m.velocity -= float64(steps) * 0.4
					m.startDecrypt()
				}
			}
			return m, nil

		case "down", "j":
			if m.currentView == ViewProjects {
				steps := m.consumeNum(1)
				prev := m.projectCursor
				m.projectCursor += steps
				if m.projectCursor >= len(views.AllProjects) {
					m.projectCursor = len(views.AllProjects) - 1
				}
				if m.projectCursor != prev {
					m.tagPopReveal = 0
					m.velocity += float64(steps) * 0.4
					m.startDecrypt()
				}
			}
			return m, nil

		case "g":
			// gg — handled as two consecutive g presses
			if m.currentView == ViewProjects {
				if m.numBuf == "g" {
					m.numBuf = ""
					prev := m.projectCursor
					m.projectCursor = 0
					if m.projectCursor != prev { m.tagPopReveal = 0; m.startDecrypt() }
				} else {
					m.numBuf = "g"
				}
			}
			return m, nil

		case "G":
			if m.currentView == ViewProjects {
				prev := m.projectCursor
				m.projectCursor = len(views.AllProjects) - 1
				if m.projectCursor != prev { m.tagPopReveal = 0; m.startDecrypt() }
				m.numBuf = ""
			}
			return m, nil

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.currentView == ViewProjects {
				if m.numBuf != "g" {
					m.numBuf += msg.String()
				}
			}
			return m, nil

		case "enter":
			if m.currentView == ViewHome {
				switch m.activeTab {
				case 0:
					m.startWipe(ViewProjects, 0)
				case 1:
					m.startWipe(ViewAbout, 1)
				case 2:
					m.startWipe(ViewContacts, 2)
				}
			}
			return m, nil
		}

	case tickMsg:
		m.tickCount++

		// ── Wipe transition ─────────────────────────────────────────────
		if m.wipePhase != 0 {
			step := m.height / 4
			if step < 4 { step = 4 }
			m.wipeLines += step
			if m.wipePhase == 1 && m.wipeLines >= m.height {
				// Switch to pending view
				m.currentView = m.pendingView
				if m.pendingView == ViewProjects { m.projectCursor = 0; m.projectsReveal = 0; m.tagPopReveal = 0 }
				if m.pendingView == ViewContacts { m.contactsReveal = 0; m.sshFlash = 4 }
				m.wipePhase = 2
				m.wipeLines = 0
			} else if m.wipePhase == 2 && m.wipeLines >= m.height {
				m.wipePhase = 0
				m.wipeLines = 0
			}
			return m, tickCmd()
		}

		// ── Matrix rain ────────────────────────────────────────────
		if m.currentView == ViewMatrix {
			// Tick matrix every 2 ticks (~10fps) for SSH efficiency
			if m.tickCount%2 == 0 {
				m.matrixCols = views.TickMatrixColumns(m.matrixCols, m.height)
			}
			// Phase 0 → 1 after 40 ticks (2s)
			if m.matrixPhase == 0 && m.tickCount >= 40 {
				m.matrixPhase = 1
			}
			// Phase 1: lock 6 random cells per tick
			if m.matrixPhase == 1 {
				allCells := views.ComputeNameCells(m.matrixNameX, m.matrixNameY, views.NameBannerLines())
				for i := 0; i < 6 && len(m.matrixPending) > 0; i++ {
					pos := m.matrixPending[0]
					m.matrixPending = m.matrixPending[1:]
					m.matrixLocked[pos] = allCells[pos]
				}
				// All locked → phase 2 (fade)
				if len(m.matrixPending) == 0 {
					m.matrixPhase = 2
				}
			}
			// Phase 2 → boot after 10 ticks (500ms)
			if m.matrixPhase == 2 && m.tickCount >= 40+len(views.NameBannerLines())*6/6+10 {
				m.currentView = ViewBoot
				m.bootStartTick = m.tickCount // record when boot screen starts
			}
			return m, tickCmd()
		}

		// ── Boot sequence ──────────────────────────────────────────
		if m.currentView == ViewBoot {
			// Use elapsed ticks since boot screen appeared — NOT absolute tickCount
			// (tickCount is already ~70+ when boot starts, so all schedule entries
			//  would fire instantly if we compared against raw tickCount)
			elapsedBoot := m.tickCount - m.bootStartTick
			for i, scheduledTick := range m.bootSchedule {
				if elapsedBoot >= scheduledTick && i >= m.bootVisible {
					m.bootVisible = i + 1
				}
			}
			// All lines shown + extra pause (10 ticks = 500ms) → go to alert
			lastTick := m.bootSchedule[len(m.bootSchedule)-1]
			if m.bootVisible >= len(m.bootLines) && elapsedBoot >= lastTick+10 {
				m.currentView = ViewAlert
				m.alertPhase = 0
				m.alertPhaseTick = m.tickCount
			}
			return m, tickCmd()
		}

		// ── Alert ─────────────────────────────────────────────────
		if m.currentView == ViewAlert {
			elapsed := m.tickCount - m.alertPhaseTick
			if m.alertPhase == 0 && elapsed >= 30 { // 1500ms
				m.alertPhase = 1
				m.alertPhaseTick = m.tickCount
			} else if m.alertPhase == 1 && elapsed >= 16 { // 800ms
				m.currentView = ViewHome
			}
			return m, tickCmd()
		}

		// ── Independent star twinkle ──────────────────────────────
		for i := range m.stars {
			if m.tickCount >= m.stars[i].FlipAt {
				m.stars[i].Bright = !m.stars[i].Bright
				// Next flip in 10–40 ticks (500ms–2000ms)
				m.stars[i].FlipAt = m.tickCount + 10 + rand.Intn(30)
			}
		}

		// ── CRT scanline sweep (every tick, wraps around) ─────────
		if m.currentView == ViewHome {
			m.scanlineY = (m.tickCount / 2) % (m.height + 5)
			if m.scanlineY >= m.height {
				m.scanlineY = -1 // hide during reset gap
			}
		} else {
			m.scanlineY = -1
		}

		// ── Idle ambient glitch ────────────────────────────────────
		m.idleGlitch = false
		if m.currentView == ViewHome {
			m.idleTicks++
			// 400 ticks @ 50ms = 20s idle threshold
			if m.idleTicks >= 400 && m.tickCount%400 == 0 {
				m.idleGlitch = true
			}
		}

		// ── Home / Splash animations ───────────────────────────────


		if m.currentView == ViewHome {
			if m.revealIdx < views.BannerLines() {
				m.revealIdx++
				// Trigger glitch when banner just completed
				if m.revealIdx == views.BannerLines() && m.glitchFrames == 0 {
					m.glitchFrames = 3
					m.glitchRunes = glitchBanner()
				}
			} else {
				// Advance glitch frames every 2 ticks (100ms per frame)
				if m.glitchFrames > 0 && m.tickCount%2 == 0 {
					m.glitchFrames--
					if m.glitchFrames > 0 {
						m.glitchRunes = glitchBanner()
					}
				}
				tagline := views.TaglineText
				if m.taglineIdx < len([]rune(tagline)) {
					m.taglineIdx++
				} else if m.cursorLeft > 0 {
					m.cursorBlink = !m.cursorBlink
					if !m.cursorBlink {
						m.cursorLeft--
					}
				} else {
					m.taglineDone = true
					m.cursorBlink = false
				}
			}
		}

		// ── Theme flash decay ───────────────────────────────────────
		if m.themeFlash > 0 {
			m.themeFlash--
		}

		// ── Live pulse (always, every 10 ticks = 500ms) ───────────
		if m.tickCount%10 == 0 {
			m.livePulse = !m.livePulse
		}

		// ── Project cascade + tag pop + lerp + decrypt ────────────
		if m.currentView == ViewProjects {
			if m.tickCount%2 == 0 && m.projectsReveal < len(views.AllProjects) {
				m.projectsReveal++
			}
			if m.tagPopReveal < len(views.AllProjects[m.projectCursor].Tags) {
				m.tagPopReveal++
			}

			// Lerp highlightY toward projectCursor (smooth slide)
			target := float64(m.projectCursor)
			m.highlightY += (target - m.highlightY) * 0.25
			if abs64(m.highlightY-target) < 0.01 {
				m.highlightY = target
			}

			// Momentum decay (rubber band feel)
			if abs64(m.velocity) > 0.01 {
				m.velocity *= 0.78
			} else {
				m.velocity = 0
			}

			// Decrypt reveal — advance one char per tick
			if m.decryptIdx < len(m.decryptRunes) {
				m.decryptIdx++
			}
		}

		// ── Contacts stagger + SSH flash ──────────────────────────
		if m.currentView == ViewContacts {
			if m.tickCount%3 == 0 && m.contactsReveal < len(views.AllContacts) {
				m.contactsReveal++
			}
			if m.sshFlash > 0 && m.tickCount%4 == 0 {
				m.sshFlash--
			}
		}

		return m, tickCmd()

	// ── GitHub commit fetch result ──────────────────────────────
	case commitMsg:
		m.lastCommit = string(msg)
		return m, nil

	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		r := m.renderer
		cyanS := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
		dimS  := r.NewStyle().Foreground(lipgloss.Color("#555555"))
		return "\n  " + cyanS.Render("✶ Thanks for visiting!") + "\n  " + dimS.Render("github.com/trafalgar-2006/portflio") + "\n  " + dimS.Render("ssh.mohith.is-a.dev") + "\n\n"
	}

	var content string
	theme := views.Themes[m.themeIdx]

	switch m.currentView {
	case ViewMatrix:
		return views.RenderMatrix(m.renderer, m.width, m.height, m.matrixCols, m.matrixLocked, m.matrixPhase == 2, theme)
	case ViewBoot:
		return views.RenderBoot(m.renderer, m.width, m.height, m.bootVisible, m.bootLines, theme)
	case ViewAlert:
		return views.RenderAlert(m.renderer, m.width, m.height, m.alertPhase, theme)
	case ViewHome:
		starBright := make([]bool, numStars)
		for i, s := range m.stars {
			starBright[i] = s.Bright
		}
		buildInfo := fmt.Sprintf("build %s · go1.22", BuildCommit)
		connectedSecs := int(time.Since(m.sessionStart).Seconds())
		content = views.RenderHome(m.renderer, m.width, m.height, m.revealIdx, starBright, m.taglineIdx, m.taglineDone, m.cursorBlink, m.glitchFrames, m.glitchRunes, m.lastCommit, m.sessionID, connectedSecs, buildInfo, m.scanlineY, m.idleGlitch, theme)
		content += m.renderTabBar(theme)
	case ViewProjects:
		content = views.RenderProjects(m.renderer, m.width, m.height, m.projectCursor, m.projectScroll, m.projectsReveal, m.tagPopReveal, m.livePulse, m.highlightY, m.decryptIdx, m.decryptRunes, theme)
	case ViewAbout:
		content = views.RenderAbout(m.renderer, m.width, m.height, theme)
	case ViewContacts:
		content = views.RenderContacts(m.renderer, m.width, m.height, m.contactsReveal, m.sshFlash, theme)
	default:
		return ""
	}

	// Apply wipe mask
	if m.wipePhase != 0 {
		lines := strings.Split(content, "\n")
		for i := range lines {
			var blank bool
			if m.wipePhase == 1 {
				blank = i < m.wipeLines
			} else {
				blank = i >= m.wipeLines
			}
			if blank {
				lines[i] = ""
			}
		}
		content = strings.Join(lines, "\n")
	}

	// Global status footer bar — shown on all non-intro screens
	if m.currentView >= ViewHome {
		content += m.renderFooterBar()
	}

	// Theme flash overlay — a brief bright flicker on theme switch
	if m.themeFlash > 0 {
		r := m.renderer
		flashAlpha := float64(m.themeFlash) / 4.0
		_ = flashAlpha
		theme := views.Themes[m.themeIdx]
		flashS := r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Bold(true)
		name   := "  ✨ theme: " + theme.Name
		content = flashS.Render(name) + "\n" + content
	}

	return content
}

// startWipe begins a wipe-out → switch → wipe-in transition
func (m *Model) startWipe(target View, tab int) {
	m.pendingView = target
	m.pendingTab = tab
	m.wipePhase = 1
	m.wipeLines = 0
}

func (m Model) renderTabBar(theme views.Theme) string {
	r := m.renderer
	activeStyle   := r.NewStyle().Bold(true).Background(theme.TabActive).Foreground(lipgloss.Color("#0A0A0A")).Padding(0, 1)
	inactiveStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	hintStyle     := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim)).Italic(true)
	sep           := "  "

	var tabs []string
	for i, name := range tabNames {
		if i == m.activeTab {
			tabs = append(tabs, activeStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveStyle.Render(name))
		}
	}

	themeName := r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Italic(true).Render("[t] "+theme.Name)
	tabBar := "\n " + joinStrings(tabs, sep) + "   " + themeName + "\n"
	tabBar += "\n " + hintStyle.Render("[← → tabs · enter open · ↑↓ browse · t theme · q quit]") + "\n"
	return tabBar
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func (m Model) renderFooterBar() string {
	r := m.renderer
	theme  := views.Themes[m.themeIdx]
	barBg    := r.NewStyle().Foreground(lipgloss.Color(theme.FooterText)).Background(lipgloss.Color(theme.FooterBg))
	sidStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Background(lipgloss.Color(theme.FooterBg))
	sepStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Dim)).Background(lipgloss.Color(theme.FooterBg))
	qStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim)).Background(lipgloss.Color(theme.FooterBg)).Italic(true)

	secs := int(time.Since(m.sessionStart).Seconds())
	mins := secs / 60
	s    := secs % 60

	viewName := "home"
	switch m.currentView {
	case ViewProjects: viewName = "projects"
	case ViewAbout:    viewName = "about"
	case ViewContacts: viewName = "contacts"
	}

	sid := m.sessionID
	if sid == "" { sid = "--------" }

	bar := sidStyle.Render("  "+sid) +
		sepStyle.Render("  ·  ") +
		barBg.Render(fmt.Sprintf("connected: %02d:%02d", mins, s)) +
		sepStyle.Render("  ·  ") +
		barBg.Render(viewName) +
		sepStyle.Render("  ·  ") +
		barBg.Render(theme.Name) +
		sepStyle.Render("  ·  ") +
		qStyle.Render("[q] quit  ")

	return "\n" + bar + "\n"
}

// commitMsg carries the last GitHub commit line back to the model
type commitMsg string

// fetchCommit fetches the latest commit from the GitHub API asynchronously
func fetchCommit() tea.Cmd {
	return func() tea.Msg {
		type ghCommit struct {
			Commit struct {
				Message string `json:"message"`
			} `json:"commit"`
			CommitDate string // parsed below
		}
		client := &http.Client{Timeout: 4 * time.Second}
		resp, err := client.Get("https://api.github.com/repos/Trafalgar-2006/portflio/commits?per_page=1")
		if err != nil {
			return commitMsg("")
		}
		defer resp.Body.Close()
		var commits []struct {
			Commit struct {
				Message   string `json:"message"`
				Author struct {
					Date string `json:"date"`
				} `json:"author"`
			} `json:"commit"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil || len(commits) == 0 {
			return commitMsg("")
		}
		msg := commits[0].Commit.Message
		// Trim to first line only
		if idx := strings.Index(msg, "\n"); idx != -1 {
			msg = msg[:idx]
		}
		if len(msg) > 52 {
			msg = msg[:52] + "..."
		}
		// Relative time
		pushedAt, err := time.Parse(time.RFC3339, commits[0].Commit.Author.Date)
		ago := ""
		if err == nil {
			diff := time.Since(pushedAt)
			switch {
			case diff < time.Hour:
				ago = fmt.Sprintf("%dm ago", int(diff.Minutes()))
			case diff < 24*time.Hour:
				ago = fmt.Sprintf("%dh ago", int(diff.Hours()))
			default:
				ago = fmt.Sprintf("%dd ago", int(diff.Hours()/24))
			}
		}
		return commitMsg(fmt.Sprintf("\"%s\" · %s", msg, ago))
	}
}

// glitchBanner returns a 2D slice of runes with ~20%% of non-space chars corrupted
func glitchBanner() [][]rune {
	glyphChars := []rune{'#', '@', '%', '$', '!', '?', '█', '▓', '&', '*'}
	result := make([][]rune, len(views.NameBannerLines()))
	for i, line := range views.NameBannerLines() {
		runes := []rune(line)
		for j, r := range runes {
			if r != ' ' && rand.Float32() < 0.20 {
				runes[j] = glyphChars[rand.Intn(len(glyphChars))]
			}
		}
		result[i] = runes
	}
	return result
}

// consumeNum reads m.numBuf as an integer (defaulting to def if empty), then clears it.
func (m *Model) consumeNum(def int) int {
	if m.numBuf == "" || m.numBuf == "g" {
		m.numBuf = ""
		return def
	}
	n := 0
	for _, ch := range m.numBuf {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		}
	}
	m.numBuf = ""
	if n == 0 {
		return def
	}
	return n
}

// startDecrypt seeds a fresh decrypt animation for the current project's description.
func (m *Model) startDecrypt() {
	if m.projectCursor >= len(views.AllProjects) {
		return
	}
	desc := []rune(views.AllProjects[m.projectCursor].Description)
	scramble := []rune("!@#$%^&*<>?/\\|~`[]{}ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	runes := make([]rune, len(desc))
	for i, r := range desc {
		if r == ' ' {
			runes[i] = ' '
		} else {
			runes[i] = scramble[rand.Intn(len(scramble))]
		}
	}
	m.decryptRunes = runes
	m.decryptIdx = 0
}

// abs64 returns the absolute value of a float64.
func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
