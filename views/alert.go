package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderAlert renders the "Unauthorized Access" fake alert (phase 0)
// or the "just kidding" reveal (phase 1).
func RenderAlert(r *lipgloss.Renderer, width, height, phase int, theme Theme) string {
	var b strings.Builder

	if phase == 0 {
		redBg    := r.NewStyle().Background(lipgloss.Color("#CC0000")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
		redStyle := r.NewStyle().Foreground(lipgloss.Color("#FF3333")).Bold(true)
		dimStyle := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))

		boxW := 55
		if width < boxW+4 {
			boxW = width - 4
		}
		inner := boxW - 2 // chars between the two border pipes

		// Horizontal centering: pad so box sits in the middle of the terminal
		hPad := (width - boxW) / 2
		if hPad < 0 { hPad = 0 }
		prefix := strings.Repeat(" ", hPad)

		// Vertical centering: pad top so box sits in the middle
		boxHeight := 8 // number of box lines we'll draw
		vPad := (height - boxHeight) / 2
		if vPad < 2 { vPad = 2 }
		for i := 0; i < vPad; i++ {
			b.WriteString("\n")
		}

		// Helper: build one box row with left border, content padded to inner, right border
		row := func(content string) string {
			vis := lipgloss.Width(content)
			pad := inner - vis
			if pad < 0 { pad = 0 }
			return prefix + redStyle.Render("│") + content + strings.Repeat(" ", pad) + redStyle.Render("│")
		}

		// Headline content (fixed text, measure to center inside box)
		headline := "  " + redBg.Render("  ⚠  UNAUTHORIZED ACCESS DETECTED  ") + "  "

		b.WriteString(prefix + redStyle.Render("┌"+strings.Repeat("─", inner)+"┐") + "\n")
		b.WriteString(prefix + redStyle.Render("│"+strings.Repeat(" ", inner)+"│") + "\n")
		b.WriteString(row(headline) + "\n")
		b.WriteString(prefix + redStyle.Render("│"+strings.Repeat(" ", inner)+"│") + "\n")
		b.WriteString(row("  " + dimStyle.Render("Tracing connection origin...")) + "\n")
		b.WriteString(row("  " + dimStyle.Render("Logging session metadata...")) + "\n")
		b.WriteString(row("  " + dimStyle.Render("Notifying system administrator...")) + "\n")
		b.WriteString(prefix + redStyle.Render("│"+strings.Repeat(" ", inner)+"│") + "\n")
		b.WriteString(prefix + redStyle.Render("└"+strings.Repeat("─", inner)+"┘") + "\n")

	} else {
		// Phase 1: just kidding — centered
		dimStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))
		cyanStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))

		vPad := (height / 2) - 1
		if vPad < 0 { vPad = 0 }
		hPad := (width / 2) - 7
		if hPad < 0 { hPad = 0 }
		prefix := strings.Repeat(" ", hPad)

		for i := 0; i < vPad; i++ {
			b.WriteString("\n")
		}
		b.WriteString(prefix + cyanStyle.Render("just kidding.") + "\n")
		b.WriteString(prefix + dimStyle.Render("welcome. :)") + "\n")
	}

	return b.String()
}
