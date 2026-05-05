package main

import (
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
	taglineIdx  int  // typewriter char index
	taglineDone bool
	cursorBlink bool
	cursorLeft  int // blink cycles remaining after typewrite done
	blink       bool

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
	return tickCmd()
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
			m.currentView = ViewHome
			return m, nil

		case "esc":
			if m.currentView != ViewHome && m.currentView != ViewBoot && m.currentView != ViewAlert {
				m.currentView = ViewHome
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
					m.currentView = ViewProjects
					m.projectCursor = 0
					m.projectsReveal = 0
					m.tagPopReveal = 0
				case 1:
					m.currentView = ViewAbout
				case 2:
					m.currentView = ViewContacts
					m.contactsReveal = 0
					m.sshFlash = 4
				}
			}
			return m, nil
		}

	case tickMsg:
		m.tickCount++

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
			} else {
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
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "\n  Thanks for visiting! ✦\n\n"
	}

	switch m.currentView {
	case ViewBoot:
		return views.RenderBoot(m.renderer, m.width, m.height, m.bootVisible, m.bootLines)

	case ViewAlert:
		return views.RenderAlert(m.renderer, m.width, m.height, m.alertPhase)

	case ViewHome:
		content := views.RenderHome(m.renderer, m.width, m.height, m.revealIdx, m.blink, m.taglineIdx, m.taglineDone, m.cursorBlink)
		content += m.renderTabBar()
		return content

	case ViewProjects:
		return views.RenderProjects(m.renderer, m.width, m.height, m.projectCursor, m.projectScroll, m.projectsReveal, m.tagPopReveal, m.livePulse)

	case ViewAbout:
		return views.RenderAbout(m.renderer, m.width, m.height)

	case ViewContacts:
		return views.RenderContacts(m.renderer, m.width, m.height, m.contactsReveal, m.sshFlash)
	}
	return ""
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
