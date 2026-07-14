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
	{Icon: "(@)", Label: "Email",    Value: "d.mohithakshay@gmail.com"},
	{Icon: "(~)", Label: "Website",  Value: "webcraftstudios.co.in"},
	{Icon: "(in)", Label: "LinkedIn", Value: "linkedin.com/in/dmohithakshay"},
	{Icon: "(gh)", Label: "GitHub",  Value: "github.com/trafalgar-2006"},
}

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

// padToWidth pads a styled string (with ANSI codes) to a target visual width using spaces.
func padToWidth(content string, targetWidth int) string {
	vis := lipgloss.Width(content)
	if vis >= targetWidth {
		return content
	}
	return content + strings.Repeat(" ", targetWidth-vis)
}

func RenderContacts(r *lipgloss.Renderer, width, height, contactsReveal, sshFlash int, theme Theme) string {
	cyanStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	dimMid     := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))
	whiteStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	goldStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Accent)).Bold(true)
	linkStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Underline(true)
	boxStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.BoxBorder))
	divider    := dimStyle.Render("  " + strings.Repeat("─", 46))

	// Card dimensions — inner content width (between │ and │)
	const cardInner = 44 // visual chars between the two border pipes
	top    := boxStyle.Render("╭" + strings.Repeat("─", cardInner) + "╮")
	bottom := boxStyle.Render("╰" + strings.Repeat("─", cardInner) + "╯")
	lBdr   := boxStyle.Render("│")
	rBdr   := boxStyle.Render("│")

	// Build a card row: border + content padded to cardInner + border
	cardRow := func(content string) string {
		return "  " + lBdr + padToWidth(" "+content, cardInner) + rBdr
	}

	iconColors := map[string]lipgloss.Color{
		"(@)":  "#FF6AC1",
		"(~)":  "#50FA7B",
		"(in)": "#0088CC",
		"(gh)": "#E0E0E0",
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ Contacts") + "\n")
	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimMid.Render("Let's connect! Feel free to reach out.") + "\n\n")

	for i, c := range AllContacts {
		if i >= contactsReveal {
			b.WriteString("\n")
			continue
		}

		iconColor := iconColors[c.Icon]
		if iconColor == "" {
			iconColor = lipgloss.Color(theme.DimMid)
		}
		iconStyle := r.NewStyle().Foreground(iconColor).Bold(true)

		labelContent := iconStyle.Render(c.Icon) + "  " + goldStyle.Render(c.Label)
		valueContent := linkStyle.Render(c.Value)

		b.WriteString("  " + top + "\n")
		b.WriteString(cardRow(labelContent) + "\n")
		b.WriteString(cardRow(valueContent) + "\n")
		b.WriteString("  " + bottom + "\n\n")
	}

	b.WriteString(divider + "\n\n")

	sshLine := whiteStyle.Render("You're viewing this over ") + cyanStyle.Bold(true).Render("SSH") + whiteStyle.Render("!")
	if sshFlash > 0 && sshFlash%2 == 0 {
		sshLine = r.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("You're viewing this over SSH! ✦")
	}
	b.WriteString("  " + sshLine + "\n")
	b.WriteString("  " + dimStyle.Render("Built with Go + Bubbletea + Wish  ·  github.com/trafalgar-2006/ssh-portfolio") + "\n\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
