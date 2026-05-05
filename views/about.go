package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderAbout(r *lipgloss.Renderer, width, height int) string {
	cyanStyle    := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle     := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	whiteStyle   := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	magentaStyle := r.NewStyle().Foreground(lipgloss.Color("#FF6AC1"))
	goldStyle    := r.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	greenStyle   := r.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	purpleStyle  := r.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	orangeStyle  := r.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	divider      := dimStyle.Render("  ─────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ About") + "\n")
	b.WriteString(divider + "\n\n")

	// Name & tagline
	b.WriteString("  " + goldStyle.Bold(true).Render("Mohith Akshay Duggirala") + "\n")
	b.WriteString("  " + magentaStyle.Italic(true).Render("Engineer · Builder · Creator") + "\n")
	b.WriteString("  " + dimStyle.Render("Bengaluru, Karnataka, India") + "\n\n")

	// Education
	b.WriteString("  " + cyanStyle.Bold(true).Render("🎓 Education") + "\n")
	b.WriteString("  " + whiteStyle.Render("B.Tech in Electronics & Computer Engineering") + "\n")
	b.WriteString("  " + whiteStyle.Render("Manipal Institute of Technology, Bengaluru") + "\n")
	b.WriteString("  " + dimStyle.Render("Aug 2023 – Jul 2027") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Active member and President of ") + purpleStyle.Render("MBOSC") + "\n")
	b.WriteString("  " + dimStyle.Render("  Coursework: Computer Architecture, Data Structures,") + "\n")
	b.WriteString("  " + dimStyle.Render("  Network Protocols, Embedded Systems") + "\n\n")

	// Experience
	b.WriteString("  " + cyanStyle.Bold(true).Render("💼 Experience") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Computer Vision Research Intern, ") + magentaStyle.Render("ISRO – LEOS") + "\n")
	b.WriteString("    " + dimStyle.Render("Dec 2025 – Jan 2026") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Founder & Lead Designer, ") + magentaStyle.Render("Webcraft Studios") + "\n")
	b.WriteString("    " + dimStyle.Render("2025 – Present  ·  Bengaluru, IN") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Full Stack Developer Intern, ") + magentaStyle.Render("SnuqSq Tech Solutions") + "\n")
	b.WriteString("    " + dimStyle.Render("Mar 2025 – Jul 2025") + "\n\n")

	// Specializations
	b.WriteString("  " + cyanStyle.Bold(true).Render("🔬 Specializations") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Computer Vision (Sim-to-Real Transfer)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Full-Stack Development & Scalable AI Pipelines") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("Edge Deployment (NVIDIA Jetson, TensorRT)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("SEO Optimization & Embedded Systems") + "\n\n")

	// Technical Skills — removed (see Resume tab)

	// Interests
	b.WriteString("  " + cyanStyle.Bold(true).Render("🌟 Interests") + "\n")
	b.WriteString("  " + dimStyle.Render("Shipping things that actually work in production —") + "\n")
	b.WriteString("  " + dimStyle.Render("from satellite CV pipelines at ISRO, to a live") + "\n")
	b.WriteString("  " + dimStyle.Render("autonomous trading system on Oracle Cloud, to this") + "\n")
	b.WriteString("  " + dimStyle.Render("very SSH portfolio you're viewing right now.") + "\n\n")

	b.WriteString(divider + "\n")
	b.WriteString("\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
