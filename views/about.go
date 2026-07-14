package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func RenderAbout(r *lipgloss.Renderer, width, height int, theme Theme) string {
	cyanStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle     := r.NewStyle().Foreground(lipgloss.Color(theme.DimMid))
	whiteStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	magentaStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	goldStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Accent))
	greenStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Success))
	purpleStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Purple))
	orangeStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Warning))
	divider      := r.NewStyle().Foreground(lipgloss.Color(theme.Dim)).Render("  ─────────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")

	// Self-aware time-based greeting (IST)
	ist := time.FixedZone("IST", 5*60*60+30*60)
	hour := time.Now().In(ist).Hour()
	var timeGreeting string
	switch {
	case hour >= 2 && hour < 5:
		timeGreeting = "you're up late. so am I, probably."
	case hour >= 5 && hour < 9:
		timeGreeting = "early start. respect."
	case hour >= 9 && hour < 18:
		timeGreeting = "currently probably in class or debugging something."
	case hour >= 18 && hour < 22:
		timeGreeting = "golden hours. this is when the best code gets written."
	default:
		timeGreeting = "late night build session energy in here."
	}
	b.WriteString("  " + dimStyle.Italic(true).Render(timeGreeting) + "\n\n")

	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ About") + "\n")
	b.WriteString(divider + "\n\n")

	// Name & tagline
	b.WriteString("  " + goldStyle.Bold(true).Render("Mohith Akshay Duggirala") + "\n")
	b.WriteString("  " + magentaStyle.Italic(true).Render("AI/ML Engineer · Full-Stack Developer · Builder") + "\n")
	b.WriteString("  " + dimStyle.Render("Bengaluru, Karnataka, India") + "\n\n")

	// Education
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Education") + "\n")
	b.WriteString("  " + whiteStyle.Bold(true).Render("Manipal Institute of Technology, Bengaluru") + "\n")
	b.WriteString("  " + whiteStyle.Render("B.Tech — Electronics & Computer Engineering") + "\n")
	b.WriteString("  " + dimStyle.Render("Aug 2023 – Jul 2027") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("President, ") + purpleStyle.Render("MBOSC") + dimStyle.Render(" (Manipal Bengaluru Open Source Community)") + "\n\n")

	// Experience timeline (absorbed from removed Resume tab)
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Experience") + "\n\n")

	b.WriteString("  " + goldStyle.Bold(true).Render("Computer Vision Research Intern") + "\n")
	b.WriteString("  " + magentaStyle.Render("ISRO – LEOS") + "  " + dimStyle.Render("Dec 2025 – Jan 2026") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Synthetic data generation with BlenderProc (6,000+ images)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("YOLOv7+Segmentation on NVIDIA Jetson via TensorRT @ 22 FPS") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("GAN-based sim-to-real style transfer (Unreal Engine 5)") + "\n\n")

	b.WriteString("  " + goldStyle.Bold(true).Render("Founder & Lead Developer") + "\n")
	b.WriteString("  " + magentaStyle.Render("Webcraft Studios") + "  " + dimStyle.Render("2025 – Present") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Full-stack apps, CV pipelines, scalable AI products") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Product design, engineering & client delivery") + "\n\n")

	b.WriteString("  " + goldStyle.Bold(true).Render("Full Stack Developer Intern") + "\n")
	b.WriteString("  " + magentaStyle.Render("SnuqSq Tech Solutions") + "  " + dimStyle.Render("Mar 2025 – Jul 2025") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Built and shipped production features") + "\n\n")

	// Skills with bar meters
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Skills") + "\n\n")

	type skillBar struct{ cat, bar string; pct int }
	skillBars := []skillBar{
		{"Python",     "#00DFDF", 95},
		{"Go",         "#00DFDF", 82},
		{"TypeScript", "#00DFDF", 75},
		{"C / C++",    "#00DFDF", 65},
		{"PyTorch",    "#FF6AC1", 88},
		{"TensorFlow", "#FF6AC1", 78},
		{"YOLOv7",     "#FF6AC1", 85},
		{"TensorRT",   "#FF6AC1", 72},
		{"React",      "#BD93F9", 80},
		{"Node.js",    "#BD93F9", 77},
		{"Docker",     "#50FA7B", 83},
		{"AWS",        "#50FA7B", 68},
	}
	barTotal := 16
	for _, s := range skillBars {
		filled := (s.pct * barTotal) / 100
		empty  := barTotal - filled
		bar := r.NewStyle().Foreground(lipgloss.Color(s.bar)).Render(strings.Repeat("█", filled)) +
			dimStyle.Render(strings.Repeat("░", empty))
		pctStr := fmt.Sprintf("%3d%%", s.pct)
		namePad := fmt.Sprintf("%-12s", s.cat)
		b.WriteString("  " + orangeStyle.Render(namePad) + "  " + bar + "  " + dimStyle.Render(pctStr) + "\n")
	}
	b.WriteString("\n")

	// Interests
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ What drives me") + "\n")
	b.WriteString("  " + dimStyle.Render("Shipping things that actually work in production —") + "\n")
	b.WriteString("  " + dimStyle.Render("from satellite CV at ISRO, to a live autonomous") + "\n")
	b.WriteString("  " + dimStyle.Render("trading system, to this SSH portfolio you're in now.") + "\n\n")

	b.WriteString("  " + dimStyle.Render("Resume PDF →") + "\n")
	b.WriteString("  " + orangeStyle.Render("github.com/Trafalgar-2006/portflio/raw/master/Mohith_Akshay_Duggirala_Resume.pdf") + "\n\n")

	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
