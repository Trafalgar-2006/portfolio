package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/trafalgar-2006/ssh-portfolio/config"
)

type Contact struct {
	Icon  string
	Label string
	Value string
}

var AllContacts = []Contact{
	{Icon: "(@)", Label: "Email", Value: "d.mohithakshay@gmail.com"},
	{Icon: "(~)", Label: "Website", Value: "webcraftstudios.co.in"},
	{Icon: "(in)", Label: "LinkedIn", Value: "linkedin.com/in/dmohithakshay"},
	{Icon: "(gh)", Label: "GitHub", Value: "github.com/trafalgar-2006"},
}

// loadContactsFromConfig overwrites AllContacts from the YAML config.
// Called by LoadFromConfig() in projects.go after config is loaded.
func loadContactsFromConfig() {
	if config.Loaded == nil || len(config.Loaded.Contacts) == 0 {
		return
	}
	var contacts []Contact
	for _, c := range config.Loaded.Contacts {
		contacts = append(contacts, Contact{
			Icon:  c.Icon,
			Label: c.Label,
			Value: c.Value,
		})
	}
	AllContacts = contacts
}

func RenderContacts(r *lipgloss.Renderer, width, height, contactsReveal, sshFlash int, theme Theme) string {
	cyanStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	dimMid     := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))
	whiteStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	goldStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Accent))
	linkStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Underline(true)
	boxStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.BoxBorder))
	divider    := dimStyle.Render("  ─────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ Contacts") + "\n")
	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimMid.Render("Let's connect! Feel free to reach out.") + "\n\n")

	iconColors := map[string]lipgloss.Color{
		"(@)":  "#FF6AC1",
		"(~)":  "#50FA7B",
		"(in)": "#0088CC",
		"(gh)": "#E0E0E0",
	}

	for i, c := range AllContacts {
		if i >= contactsReveal {
			b.WriteString("\n")
			continue
		}

		iconColor := iconColors[c.Icon]
		if iconColor == "" {
			iconColor = "#888888"
		}
		iconStyle := r.NewStyle().Foreground(iconColor).Bold(true)

		// Box card per contact
		b.WriteString("  " + boxStyle.Render("╭──────────────────────────────────────────╮") + "\n")
		b.WriteString("  " + boxStyle.Render("│") + "  " +
			iconStyle.Render(c.Icon) + "  " +
			goldStyle.Bold(true).Render(c.Label) +
			strings.Repeat(" ", 36-len(c.Icon)-len(c.Label)) +
			boxStyle.Render("│") + "\n")
		b.WriteString("  " + boxStyle.Render("│") + "     " +
			linkStyle.Render(c.Value) +
			strings.Repeat(" ", 36-len(c.Value)) +
			boxStyle.Render("│") + "\n")
		b.WriteString("  " + boxStyle.Render("╰──────────────────────────────────────────╯") + "\n\n")
	}

	// SSH fun note
	b.WriteString(divider + "\n\n")
	sshLine := whiteStyle.Render("You're viewing this over ") + cyanStyle.Bold(true).Render("SSH") + whiteStyle.Render("!")
	if sshFlash > 0 && sshFlash%2 == 0 {
		sshLine = r.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("You're viewing this over SSH! ✦")
	}
	b.WriteString("  " + sshLine + "\n")
	b.WriteString("  " + dimStyle.Render("Built with Go + Bubbletea + Wish · github.com/trafalgar-2006/ssh-portfolio") + "\n\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
