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
	divider      := dimStyle.Render("  ─────────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")
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

	// Skills (re-added since Resume tab is removed)
	b.WriteString("  " + cyanStyle.Bold(true).Render("◆ Skills") + "\n\n")

	skills := []struct{ cat, items string }{
		{"AI / ML",       "PyTorch · TensorFlow · LoRA · GGUF · YOLOv7 · TensorRT · OpenCV"},
		{"Languages",     "Python · Go · TypeScript · JavaScript · C/C++"},
		{"Web / Backend", "React · Next.js · Node.js · FastAPI · Firebase"},
		{"Tools",         "Docker · Git · CUDA · NVIDIA Jetson · Linux · AWS"},
	}
	for _, s := range skills {
		b.WriteString("  " + orangeStyle.Bold(true).Render(s.cat) + "\n")
		b.WriteString("    " + dimStyle.Render(s.items) + "\n\n")
	}

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
