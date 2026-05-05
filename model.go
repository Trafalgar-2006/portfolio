package main

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trafalgar-2006/ssh-portfolio/views"
)

// View represents the current screen
type View int

const (
	ViewHome View = iota
	ViewProjects
	ViewAbout
	ViewContacts
)

// Tab names for navigation
var tabNames = []string{"Projects", "About", "Contacts"}

// Model is the main Bubbletea model
type Model struct {
	renderer      *lipgloss.Renderer
	width         int
	height        int
	currentView   View
	activeTab     int // 0=Projects, 1=About, 2=Contacts, 3=Resume
	projectCursor int
	projectScroll int
	quitting      bool
}

// NewModel creates a new model. Pass nil to use the default renderer (local TUI mode).
func NewModel(r *lipgloss.Renderer) Model {
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}
	return Model{
		renderer:    r,
		width:       80,
		height:      40,
		currentView: ViewHome,
		activeTab:   0,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
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
			// Go back to home
			m.currentView = ViewHome
			return m, nil

		case "esc":
			if m.currentView != ViewHome {
				m.currentView = ViewHome
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

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
				m.projectCursor--
				if m.projectCursor < 0 {
					m.projectCursor = len(views.AllProjects) - 1
				}
			}
			return m, nil

		case "down", "j":
			if m.currentView == ViewProjects {
				m.projectCursor++
				if m.projectCursor >= len(views.AllProjects) {
					m.projectCursor = 0
				}
			}
			return m, nil

		case "enter":
			if m.currentView == ViewHome {
				switch m.activeTab {
				case 0:
					m.currentView = ViewProjects
					m.projectCursor = 0
				case 1:
					m.currentView = ViewAbout
				case 2:
					m.currentView = ViewContacts
				}
			}
			return m, nil
		}
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "\n  Thanks for visiting! ✦\n\n"
	}

	var content string

	switch m.currentView {
	case ViewHome:
		content = views.RenderHome(m.renderer, m.width, m.height)
		content += m.renderTabBar()
	case ViewProjects:
		content = views.RenderProjects(m.renderer, m.width, m.height, m.projectCursor, m.projectScroll)
	case ViewAbout:
		content = views.RenderAbout(m.renderer, m.width, m.height)
	case ViewContacts:
		content = views.RenderContacts(m.renderer, m.width, m.height)
	}

	return content
}

// renderTabBar renders the bottom navigation tabs using the session renderer
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
