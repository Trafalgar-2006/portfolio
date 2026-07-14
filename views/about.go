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
	divider      := dimDarkStyle.Render("  " + strings.Repeat("─", 50))

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
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Education") + "\n")
	b.WriteString("  " + whiteStyle.Bold(true).Render("Manipal Institute of Technology, Bengaluru") + "\n")
	b.WriteString("  " + whiteStyle.Render("B.Tech — Electronics & Computer Engineering") + "\n")
	b.WriteString("  " + dimStyle.Render("Aug 2023 – Jul 2027  ·  GPA: In Progress") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + whiteStyle.Render("President, ") + purpleStyle.Render("MBOSC") + dimStyle.Render(" (Manipal Bengaluru Open Source Community)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Coursework: Data Structures, Algorithms, Network Protocols, Embedded Systems, Probability & Statistics") + "\n\n")

	// Experience
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Experience") + "\n\n")

	// ISRO
	b.WriteString("  " + goldStyle.Bold(true).Render("Computer Vision Intern") + "\n")
	b.WriteString("  " + magentaStyle.Render("ISRO – LEOS") + "  " + dimStyle.Render("Dec 2025 – Jan 2026  ·  Bengaluru") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Architected BlenderProc synthetic data pipeline → 6,000+ annotated images (COCO/YOLO)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Implemented SOG2 rotations + spherical camera trajectories for 360° angular coverage") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Reduced data generation time 70% by benchmarking path-scaling parameters") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Deployed YOLOv7 on NVIDIA Jetson Xavier via TensorRT — 22 FPS real-time inference") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Built GAN-based domain adaptation to close sim-to-real gap for offline missions") + "\n\n")

	// SenseOps
	b.WriteString("  " + goldStyle.Bold(true).Render("Software Engineering Intern") + "\n")
	b.WriteString("  " + magentaStyle.Render("SenseOps Tech Solutions Pvt. Ltd.") + "  " + dimStyle.Render("May 2025 – Jul 2025  ·  Bengaluru") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Rebuilt company portfolio site with mobile-first responsive design, measurable Lighthouse gains") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Designed & implemented custom UDP/TCP packet analyser for distributed vibration sensor network") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Improved real-time data transmission reliability and reduced packet loss across 10+ hardware nodes") + "\n\n")

	// Webcraft
	b.WriteString("  " + goldStyle.Bold(true).Render("Founder & Lead Engineer") + "\n")
	b.WriteString("  " + magentaStyle.Render("Webcraft Studios") + "  " + dimStyle.Render("Jan 2025 – Present  ·  Bengaluru") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Founded digital agency, scaled to multiple international SaaS clients within first year") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Engineered full-stack platforms: React, Node.js, MongoDB, Stripe — subscription billing + dashboards") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Built delivery pipeline (Discovery→Design→Dev→Launch) cutting cycle time by 30%+") + "\n\n")

	// Skills
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Skills") + "\n\n")

	type skillBar struct {
		name    string
		pct     int
		colorFn func(string) string // returns colored bar segment
	}

	// Color functions using theme
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
		b.WriteString("  " + dimStyle.Render(namePad) + "  " + bar + "  " + dimDarkStyle.Render(pctStr) + "\n")
	}
	b.WriteString("\n")

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
