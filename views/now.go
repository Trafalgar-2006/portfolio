package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderNow(r *lipgloss.Renderer, width, height int, theme Theme) string {
	cyanStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	goldStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Accent)).Bold(true)
	dimStyle     := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	dimMidStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))
	magentaStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	greenStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Success))
	boxStyle     := r.NewStyle().Foreground(lipgloss.Color(theme.BoxBorder))
	hintStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim)).Italic(true)
	divider      := dimStyle.Render("  " + strings.Repeat("─", 50))

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ /now") + "\n")
	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimMidStyle.Italic(true).Render("A live snapshot of what I'm building and focused on.") + "\n\n")

	// Currently building
	boxW := 52
	if width < boxW+6 { boxW = width - 6 }
	bTop := "  " + boxStyle.Render("╭"+strings.Repeat("─", boxW)+"╮")
	bBot := "  " + boxStyle.Render("╰"+strings.Repeat("─", boxW)+"╯")
	bRow := func(s string) string {
		vis := lipgloss.Width(s)
		pad := ""
		if boxW-1-vis > 0 { pad = strings.Repeat(" ", boxW-1-vis) }
		return "  " + boxStyle.Render("│") + " " + s + pad + boxStyle.Render("│")
	}

	b.WriteString("  " + goldStyle.Render("◆ Currently Building") + "\n")
	b.WriteString(bTop + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("EmbedGen — LLM fine-tuning on embedded-systems code")) + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("Autonomous trading agent (paper trading, live)")) + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("This SSH portfolio — always iterating on it")) + "\n")
	b.WriteString(bBot + "\n\n")

	// Currently learning
	b.WriteString("  " + goldStyle.Render("◆ Currently Learning") + "\n")
	b.WriteString(bTop + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("Distributed systems & consensus algorithms (Raft)")) + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("Rust for embedded targets")) + "\n")
	b.WriteString(bRow(greenStyle.Render("▸ ") + dimStyle.Render("Advanced quantization — AWQ, GPTQ")) + "\n")
	b.WriteString(bBot + "\n\n")

	// Reading
	b.WriteString("  " + goldStyle.Render("◆ Reading") + "\n")
	b.WriteString("  " + dimStyle.Render("  \"Designing Data-Intensive Applications\" — Kleppmann") + "\n")
	b.WriteString("  " + dimStyle.Render("  \"The Pragmatic Programmer\" — Hunt & Thomas") + "\n\n")

	// Status / availability
	b.WriteString("  " + goldStyle.Render("◆ Status") + "\n")
	b.WriteString("  " + magentaStyle.Render("B.Tech ECE") + dimStyle.Render(" @ Manipal Institute of Technology, Bengaluru (2023–2027)") + "\n")
	b.WriteString("  " + greenStyle.Render("Open to") + dimStyle.Render(" SWE / ML internships and research roles — from Jul 2025") + "\n\n")

	b.WriteString("  " + boxStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString("  " + dimMidStyle.Italic(true).Render("last updated: July 2025") + "\n\n")
	b.WriteString("  " + hintStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
