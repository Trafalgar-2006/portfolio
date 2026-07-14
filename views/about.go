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
	dimDarkStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	whiteStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	magentaStyle := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	goldStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Accent))
	greenStyle   := r.NewStyle().Foreground(lipgloss.Color(theme.Success))
	purpleStyle  := r.NewStyle().Foreground(lipgloss.Color(theme.Purple))
	boxStyle     := r.NewStyle().Foreground(lipgloss.Color(theme.BoxBorder))
	divider      := dimDarkStyle.Render("  " + strings.Repeat("─", 50))

	boxW := 56
	if width < boxW+6 { boxW = width - 6 }
	secTop := func(label string) string {
		inner := " " + label + " "
		pad   := boxW - len(inner) - 2
		if pad < 0 { pad = 0 }
		return "  " + boxStyle.Render("╭"+inner+strings.Repeat("─", pad)+"╮")
	}
	secBot := "  " + boxStyle.Render("╰"+strings.Repeat("─", boxW-2)+"╯")
	secRow := func(s string) string {
		vis := lipgloss.Width(s)
		pad := ""
		if boxW-2-vis > 0 { pad = strings.Repeat(" ", boxW-2-vis) }
		return "  " + boxStyle.Render("│") + s + pad + boxStyle.Render("│")
	}
	secBlank := func() string { return secRow("") }

	var b strings.Builder
	b.WriteString("\n")

	// Self-aware time-based greeting (IST)
	ist  := time.FixedZone("IST", 5*60*60+30*60)
	hour := time.Now().In(ist).Hour()
	var greeting string
	switch {
	case hour >= 2 && hour < 5:
		greeting = "you're up late. so am I, probably."
	case hour >= 5 && hour < 9:
		greeting = "early start. respect."
	case hour >= 9 && hour < 18:
		greeting = "currently probably in class or debugging something."
	case hour >= 18 && hour < 22:
		greeting = "golden hours. this is when the best code gets written."
	default:
		greeting = "late night build session energy in here."
	}
	b.WriteString("  " + dimStyle.Italic(true).Render(greeting) + "\n\n")

	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ About") + "\n")
	b.WriteString(divider + "\n\n")

	// Name & tagline
	b.WriteString("  " + goldStyle.Bold(true).Render("Mohith Akshay Duggirala") + "\n")
	b.WriteString("  " + magentaStyle.Italic(true).Render("ML Engineer · Full-Stack Developer · Builder") + "\n")
	b.WriteString("  " + dimStyle.Render("Bengaluru, Karnataka, India") + "\n\n")

	// Education
	b.WriteString(secTop(cyanStyle.Bold(true).Render("◆ Education")) + "\n")
	b.WriteString(secBlank() + "\n")
	b.WriteString(secRow("  "+whiteStyle.Bold(true).Render("Manipal Institute of Technology, Bengaluru")) + "\n")
	b.WriteString(secRow("  "+whiteStyle.Render("B.Tech — Electronics & Computer Engineering")) + "\n")
	b.WriteString(secRow("  "+dimStyle.Render("Aug 2023 – Jul 2027  ·  GPA: In Progress")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+whiteStyle.Render("President, ")+purpleStyle.Render("MBOSC")+dimStyle.Render(" (Manipal Bengaluru Open Source Community)")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("Coursework: DSA, Algorithms, Network Protocols, Embedded Systems")) + "\n")
	b.WriteString(secBlank() + "\n")
	b.WriteString(secBot + "\n\n")

	// Experience
	b.WriteString(secTop(cyanStyle.Bold(true).Render("◆ Experience")) + "\n")
	b.WriteString(secBlank() + "\n")

	// ISRO
	b.WriteString(secRow("  "+goldStyle.Bold(true).Render("Computer Vision Intern")) + "\n")
	b.WriteString(secRow("  "+magentaStyle.Render("ISRO – LEOS")+"  "+dimStyle.Render("Dec 2025 – Jan 2026  ·  Bengaluru")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("BlenderProc pipeline → 6,000+ COCO/YOLO annotated images")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("YOLOv7 on Jetson Xavier via TensorRT — 22 FPS real-time")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("GAN domain adaptation to close sim-to-real gap")) + "\n")
	b.WriteString(secBlank() + "\n")

	// SenseOps
	b.WriteString(secRow("  "+goldStyle.Bold(true).Render("Software Engineering Intern")) + "\n")
	b.WriteString(secRow("  "+magentaStyle.Render("SenseOps Tech Solutions")+"  "+dimStyle.Render("May – Jul 2025  ·  Bengaluru")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("Rebuilt portfolio: mobile-first, Lighthouse gains")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("Custom UDP/TCP packet analyser for vibration sensor network")) + "\n")
	b.WriteString(secBlank() + "\n")

	// Webcraft
	b.WriteString(secRow("  "+goldStyle.Bold(true).Render("Founder & Lead Engineer")) + "\n")
	b.WriteString(secRow("  "+magentaStyle.Render("Webcraft Studios")+"  "+dimStyle.Render("Jan 2025 – Present  ·  Bengaluru")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("Digital agency, multiple intl. SaaS clients in year one")) + "\n")
	b.WriteString(secRow("  "+greenStyle.Render("▸ ")+dimStyle.Render("React / Node.js / MongoDB / Stripe — subscription billing")) + "\n")
	b.WriteString(secBlank() + "\n")
	b.WriteString(secBot + "\n\n")

	// Skills
	b.WriteString(secTop(cyanStyle.Bold(true).Render("◆ Skills")) + "\n")
	b.WriteString(secBlank() + "\n")

	type skillBar struct {
		name    string
		pct     int
		colorFn func(string) string
	}

	langColor  := func(s string) string { return r.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Render(s) }
	mlColor    := func(s string) string { return r.NewStyle().Foreground(lipgloss.Color(theme.Secondary)).Render(s) }
	frameColor := func(s string) string { return r.NewStyle().Foreground(lipgloss.Color(theme.Purple)).Render(s) }
	infraColor := func(s string) string { return r.NewStyle().Foreground(lipgloss.Color(theme.Success)).Render(s) }

	skills := []skillBar{
		{"Python",      95, langColor},
		{"Go",          82, langColor},
		{"TypeScript",  75, langColor},
		{"C / C++",     65, langColor},
		{"PyTorch",     88, mlColor},
		{"TensorFlow",  78, mlColor},
		{"YOLOv7",      85, mlColor},
		{"TensorRT",    72, mlColor},
		{"React",       80, frameColor},
		{"Node.js",     77, frameColor},
		{"Docker",      83, infraColor},
		{"AWS",         68, infraColor},
	}

	barTotal := 18
	for _, s := range skills {
		filled := (s.pct * barTotal) / 100
		empty  := barTotal - filled
		bar := s.colorFn(strings.Repeat("█", filled)) + dimDarkStyle.Render(strings.Repeat("░", empty))
		pctStr  := fmt.Sprintf("%3d%%", s.pct)
		namePad := fmt.Sprintf("%-12s", s.name)
		rowContent := " " + dimStyle.Render(namePad) + "  " + bar + "  " + dimDarkStyle.Render(pctStr)
		b.WriteString(secRow(rowContent) + "\n")
	}
	b.WriteString(secBlank() + "\n")
	b.WriteString(secBot + "\n\n")

	// What drives me
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ What drives me") + "\n")
	b.WriteString("  " + dimStyle.Render("Shipping things that actually work in production —") + "\n")
	b.WriteString("  " + dimStyle.Render("from satellite CV at ISRO, to a live autonomous") + "\n")
	b.WriteString("  " + dimStyle.Render("trading system, to this SSH portfolio you're in now.") + "\n\n")

	b.WriteString("  " + dimDarkStyle.Render("Resume →") + "\n")
	b.WriteString("  " + r.NewStyle().Foreground(lipgloss.Color(theme.Warning)).Render("github.com/Trafalgar-2006/portflio/raw/master/Mohith_Akshay_Duggirala_Resume.pdf") + "\n\n")

	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimDarkStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
