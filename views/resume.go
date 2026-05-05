package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderResume(r *lipgloss.Renderer, width, height int) string {
	cyanStyle    := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	goldStyle    := r.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	magentaStyle := r.NewStyle().Foreground(lipgloss.Color("#FF6AC1"))
	greenStyle   := r.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	purpleStyle  := r.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	dimStyle     := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	whiteStyle   := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	orangeStyle  := r.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	divider      := dimStyle.Render("  ──────────────────────────────────────────────────")

	var b strings.Builder
	b.WriteString("\n")

	// ── Header ──────────────────────────────────────────────────────────────
	b.WriteString("  " + cyanStyle.Bold(true).Render("MOHITH AKSHAY DUGGIRALA") + "\n")
	b.WriteString("  " + magentaStyle.Render("AI/ML Engineer · Full-Stack Developer · Computer Vision") + "\n")
	b.WriteString("  " + dimStyle.Render("d.mohithakshay@gmail.com") +
		dimStyle.Render("  ·  ") +
		dimStyle.Render("linkedin.com/in/dmohithakshay") +
		dimStyle.Render("  ·  ") +
		dimStyle.Render("github.com/trafalgar-2006") + "\n")
	b.WriteString(divider + "\n\n")

	// ── Education ────────────────────────────────────────────────────────────
	b.WriteString("  " + goldStyle.Bold(true).Render("◆ EDUCATION") + "\n\n")

	b.WriteString("  " + whiteStyle.Bold(true).Render("Manipal Institute of Technology, Bengaluru") + "\n")
	b.WriteString("  " + cyanStyle.Render("B.Tech — Electronics & Computer Engineering") + "\n")
	b.WriteString("  " + dimStyle.Render("Aug 2023 – Jul 2027") + "\n\n")

	b.WriteString(divider + "\n\n")

	// ── Experience ───────────────────────────────────────────────────────────
	b.WriteString("  " + goldStyle.Bold(true).Render("◆ EXPERIENCE") + "\n\n")

	b.WriteString("  " + whiteStyle.Bold(true).Render("Computer Vision Research Intern") + "\n")
	b.WriteString("  " + cyanStyle.Render("ISRO – LEOS") + "  " + dimStyle.Render("(Indian Space Research Organisation)") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Developed CV pipelines for satellite image analysis") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Worked on perception systems for space applications") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Applied deep learning models to remote sensing data") + "\n\n")

	b.WriteString("  " + whiteStyle.Bold(true).Render("Founder & Lead Developer") + "\n")
	b.WriteString("  " + cyanStyle.Render("Webcraft Studios") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Building full-stack web apps & AI-powered products") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Leading product design and engineering") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Scalable AI perception & Computer Vision systems") + "\n\n")

	b.WriteString("  " + whiteStyle.Bold(true).Render("President") + "\n")
	b.WriteString("  " + cyanStyle.Render("MBOSC") + " " + dimStyle.Render("— Manipal Bengaluru Open Source Community") + "\n")
	b.WriteString("  " + greenStyle.Render("▸ ") + dimStyle.Render("Leading the open source developer community on campus") + "\n\n")

	b.WriteString(divider + "\n\n")

	// ── Skills ───────────────────────────────────────────────────────────────
	b.WriteString("  " + goldStyle.Bold(true).Render("◆ SKILLS") + "\n\n")

	skills := []struct{ cat, items string }{
		{"AI / ML",       "PyTorch · TensorFlow · Transformers · LoRA · GGUF"},
		{"Languages",     "Python · Go · TypeScript · C/C++ · Rust"},
		{"CV / Vision",   "OpenCV · YOLO · Segmentation · Remote Sensing"},
		{"Web / Backend", "React · Next.js · Node.js · FastAPI · Docker"},
		{"Tools",         "Git · Linux · CUDA · Firebase · Railway · Fly.io"},
	}
	for _, s := range skills {
		b.WriteString("  " + purpleStyle.Bold(true).Render(s.cat) + "\n")
		b.WriteString("  " + dimStyle.Render("  "+s.items) + "\n\n")
	}

	b.WriteString(divider + "\n\n")

	// ── Key Projects ──────────────────────────────────────────────────────────
	b.WriteString("  " + goldStyle.Bold(true).Render("◆ KEY PROJECTS") + "\n\n")

	projects := []struct{ name, tech, desc string }{
		{
			"EmbedGen",
			"Python · PyTorch · LoRA · GGUF",
			"Custom LLM fine-tuned on embedded-systems codebases. Quantized\n" +
				"  " + "  " + "  for deployment on edge devices.",
		},
		{
			"SSH Portfolio",
			"Go · Bubbletea · Wish · Docker",
			"This portfolio — interactive TUI accessible over SSH from anywhere.",
		},
		{
			"ISRO CV Pipeline",
			"Python · OpenCV · Deep Learning",
			"Computer vision pipeline for satellite imagery analysis at ISRO LEOS.",
		},
	}
	for _, p := range projects {
		b.WriteString("  " + whiteStyle.Bold(true).Render(p.name) +
			"  " + dimStyle.Render("("+p.tech+")") + "\n")
		b.WriteString("  " + "  " + dimStyle.Render(p.desc) + "\n\n")
	}

	b.WriteString(divider + "\n\n")

	// ── Download ─────────────────────────────────────────────────────────────
	b.WriteString("  " + goldStyle.Bold(true).Render("◆ DOWNLOAD RESUME (PDF)") + "\n\n")
	b.WriteString("  " + orangeStyle.Render("https://github.com/Trafalgar-2006/portflio/raw/master/Mohith_Akshay_Duggirala_Resume.pdf") + "\n\n")

	b.WriteString(divider + "\n\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
