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

type View int

const (
	ViewBoot View = iota
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

type Model struct {
	renderer  *lipgloss.Renderer
	width     int
	height    int
	quitting  bool
	tickCount int

	// Boot animation
	bootLines    []views.BootLine
	bootSchedule []int // tick number when each line should appear
	bootVisible  int   // how many lines currently shown

	// Alert animation
	alertPhase     int // 0=warning, 1=just kidding
	alertPhaseTick int // tickCount when phase started

	// Home / Splash
	currentView View
	activeTab   int
	revealIdx   int  // banner reveal
	glitchFrames int  // >0 = banner is glitching
	glitchRunes  [][]rune // corrupted banner lines during glitch
	taglineIdx  int  // typewriter char index
	taglineDone bool
	cursorBlink bool
	cursorLeft  int // blink cycles remaining after typewrite done
	blink        bool
	lastCommit   string // cached GitHub commit message

	// Wipe transition
	wipePhase   int  // 0=none, 1=wipe-out, 2=wipe-in
	wipeLines   int  // lines exposed so far
	pendingView View // view to switch to after wipe-out
	pendingTab  int  // tab to activate after wipe-out

	// Projects
	projectCursor  int
	projectScroll  int
	projectsReveal int // cascade drop-in
	tagPopReveal   int // tag pop in detail
	lastCursor     int
	livePulse      bool

	// Contacts
	contactsReveal int
	sshFlash       int // counts down from 4
}

func NewModel(r *lipgloss.Renderer) Model {
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}
	lines := views.NewBootSequence()
	// Pre-compute the tick at which each line should appear
	schedule := make([]int, len(lines))
	t := 0
	for i, l := range lines {
		schedule[i] = t
		t += l.DelayTicks
	}
	return Model{
		renderer:     r,
		width:        80,
		height:       40,
		currentView:  ViewBoot,
		bootLines:    lines,
		bootSchedule: schedule,
		bootVisible:  0,
		cursorLeft:   6, // blink 3 times (on+off = 2 per cycle * 3)
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
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

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
				prev := m.projectCursor
				m.projectCursor--
				if m.projectCursor < 0 {
					m.projectCursor = len(views.AllProjects) - 1
				}
				if m.projectCursor != prev {
					m.tagPopReveal = 0
				}
			}
			return m, nil

		case "down", "j":
			if m.currentView == ViewProjects {
				prev := m.projectCursor
				m.projectCursor++
				if m.projectCursor >= len(views.AllProjects) {
					m.projectCursor = 0
				}
				if m.projectCursor != prev {
					m.tagPopReveal = 0
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

		// ── Boot sequence ──────────────────────────────────────────
		if m.currentView == ViewBoot {
			// Reveal lines based on schedule
			for i, scheduledTick := range m.bootSchedule {
				if m.tickCount >= scheduledTick && i >= m.bootVisible {
					m.bootVisible = i + 1
				}
			}
			// All lines shown + extra pause (10 ticks = 500ms) → go to alert
			lastTick := m.bootSchedule[len(m.bootSchedule)-1]
			if m.bootVisible >= len(m.bootLines) && m.tickCount >= lastTick+10 {
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

		// ── Home / Splash animations ───────────────────────────────
		m.blink = !m.blink

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

		// ── Live pulse (always, every 10 ticks = 500ms) ───────────
		if m.tickCount%10 == 0 {
			m.livePulse = !m.livePulse
		}

		// ── Project cascade + tag pop ─────────────────────────────
		if m.currentView == ViewProjects {
			if m.tickCount%2 == 0 && m.projectsReveal < len(views.AllProjects) {
				m.projectsReveal++
			}
			if m.tagPopReveal < len(views.AllProjects[m.projectCursor].Tags) {
				m.tagPopReveal++
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
		return "\n  Thanks for visiting! ✦\n\n"
	}

	var content string

	switch m.currentView {
	case ViewBoot:
		return views.RenderBoot(m.renderer, m.width, m.height, m.bootVisible, m.bootLines)
	case ViewAlert:
		return views.RenderAlert(m.renderer, m.width, m.height, m.alertPhase)
	case ViewHome:
		content = views.RenderHome(m.renderer, m.width, m.height, m.revealIdx, m.blink, m.taglineIdx, m.taglineDone, m.cursorBlink, m.glitchFrames, m.glitchRunes, m.lastCommit)
		content += m.renderTabBar()
	case ViewProjects:
		content = views.RenderProjects(m.renderer, m.width, m.height, m.projectCursor, m.projectScroll, m.projectsReveal, m.tagPopReveal, m.livePulse)
	case ViewAbout:
		content = views.RenderAbout(m.renderer, m.width, m.height)
	case ViewContacts:
		content = views.RenderContacts(m.renderer, m.width, m.height, m.contactsReveal, m.sshFlash)
	default:
		return ""
	}

	// Apply wipe mask
	if m.wipePhase != 0 {
		lines := strings.Split(content, "\n")
		for i := range lines {
			var blank bool
			if m.wipePhase == 1 {
				// Wipe-out: blank lines from top downward
				blank = i < m.wipeLines
			} else {
				// Wipe-in: reveal lines from top downward
				blank = i >= m.wipeLines
			}
			if blank {
				lines[i] = ""
			}
		}
		content = strings.Join(lines, "\n")
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

func (m Model) renderTabBar() string {
	r := m.renderer
	activeStyle   := r.NewStyle().Bold(true).Background(lipgloss.Color("#00DFDF")).Foreground(lipgloss.Color("#0A0A0A")).Padding(0, 1)
	inactiveStyle := r.NewStyle().Foreground(lipgloss.Color("#555555"))
	hintStyle     := r.NewStyle().Foreground(lipgloss.Color("#444444")).Italic(true)
	sep           := "  "

	var tabs []string
	for i, name := range tabNames {
		if i == m.activeTab {
			tabs = append(tabs, activeStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveStyle.Render(name))
		}
	}

	tabBar := "\n " + joinStrings(tabs, sep) + "\n"
	tabBar += "\n " + hintStyle.Render("[← → tabs · enter open · ↑↓ browse · q quit]") + "\n"
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
