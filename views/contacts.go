package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Contact struct {
	Icon  string
	Label string
	Value string
}

var AllContacts = []Contact{
	{Icon: "📧", Label: "Email", Value: "d.mohithakshay@gmail.com"},
	{Icon: "🌐", Label: "Website", Value: "webcraftstudios.co.in"},
	{Icon: "💼", Label: "LinkedIn", Value: "linkedin.com/in/dmohithakshay"},
	{Icon: "🐙", Label: "GitHub", Value: "github.com/trafalgar-2006"},
}

func RenderContacts(r *lipgloss.Renderer, width, height int) string {
	cyanStyle  := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle   := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	whiteStyle := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	goldStyle  := r.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	linkStyle  := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0")).Underline(true)
	divider    := dimStyle.Render("  ─────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ Contacts") + "\n")
	b.WriteString(divider + "\n\n")

	b.WriteString("  " + dimStyle.Render("Let's connect! Feel free to reach out.") + "\n\n")

	for _, c := range AllContacts {
		b.WriteString("  " + c.Icon + "  " + goldStyle.Bold(true).Render(c.Label) + "\n")
		b.WriteString("     " + linkStyle.Render(c.Value) + "\n\n")
	}

	// Fun SSH note
	b.WriteString(divider + "\n\n")
	b.WriteString("  " + whiteStyle.Render("You're viewing this over ") + cyanStyle.Bold(true).Render("SSH") + whiteStyle.Render("!") + "\n")
	b.WriteString("  " + dimStyle.Render("Built with Go + Bubbletea + Wish") + "\n")
	b.WriteString("  " + dimStyle.Render("github.com/trafalgar-2006/ssh-portfolio") + "\n\n")

	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
