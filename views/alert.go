package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderAlert renders the "Unauthorized Access" fake alert (phase 0)
// or the "just kidding" reveal (phase 1)
func RenderAlert(r *lipgloss.Renderer, width, height, phase int) string {
	var b strings.Builder

	if phase == 0 {
		redBg    := r.NewStyle().Background(lipgloss.Color("#FF0000")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
		redStyle := r.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
		dimStyle := r.NewStyle().Foreground(lipgloss.Color("#888888"))

		boxW := 55
		if width < boxW+4 {
			boxW = width - 4
		}
		inner := boxW - 2

		top    := "┌" + strings.Repeat("─", inner) + "┐"
		mid    := "│" + strings.Repeat(" ", inner) + "│"
		bottom := "└" + strings.Repeat("─", inner) + "┘"

		pad := (height - 12) / 2
		if pad < 2 { pad = 2 }
		for i := 0; i < pad; i++ {
			b.WriteString("\n")
		}

		b.WriteString("  " + redStyle.Render(top) + "\n")
		b.WriteString("  " + redStyle.Render(mid) + "\n")
		b.WriteString("  " + redStyle.Render("│") + "  " + redBg.Render("  ⚠  UNAUTHORIZED ACCESS DETECTED  ") + strings.Repeat(" ", inner-39) + redStyle.Render("│") + "\n")
		b.WriteString("  " + redStyle.Render(mid) + "\n")
		b.WriteString("  " + redStyle.Render("│") + "  " + dimStyle.Render("Tracing connection origin...         ") + strings.Repeat(" ", inner-37) + redStyle.Render("│") + "\n")
		b.WriteString("  " + redStyle.Render("│") + "  " + dimStyle.Render("Logging session metadata...          ") + strings.Repeat(" ", inner-37) + redStyle.Render("│") + "\n")
		b.WriteString("  " + redStyle.Render("│") + "  " + dimStyle.Render("Notifying system administrator...    ") + strings.Repeat(" ", inner-37) + redStyle.Render("│") + "\n")
		b.WriteString("  " + redStyle.Render(mid) + "\n")
		b.WriteString("  " + redStyle.Render(bottom) + "\n")

	} else {
		// Phase 1: just kidding
		dimStyle  := r.NewStyle().Foreground(lipgloss.Color("#888888"))
		cyanStyle := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))

		pad := height / 2
		for i := 0; i < pad; i++ {
			b.WriteString("\n")
		}
		b.WriteString("  " + cyanStyle.Render("just kidding.") + "\n")
		b.WriteString("  " + dimStyle.Render("welcome. :)") + "\n")
	}

	return b.String()
}
