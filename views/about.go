package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderAbout(width, height int) string {
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	magentaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6AC1"))
	goldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	purpleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	orangeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	divider := dimStyle.Render("  ─────────────────────────────────────────")

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

	// Technical Skills
	b.WriteString("  " + cyanStyle.Bold(true).Render("⚙ Technical Skills") + "\n")
	b.WriteString("  " + orangeStyle.Render("Languages  ") + dimStyle.Render("JavaScript, TypeScript, Python, C++, HTML/CSS, SQL, Go") + "\n")
	b.WriteString("  " + orangeStyle.Render("Frameworks ") + dimStyle.Render("React, Next.js, Node.js, YOLOv7, BlenderProc, PyTorch,") + "\n")
	b.WriteString("             " + dimStyle.Render("TensorRT, OpenCV, React Native, Tailwind CSS, Svelte") + "\n")
	b.WriteString("  " + orangeStyle.Render("Tools      ") + dimStyle.Render("NVIDIA Jetson, Unreal Engine 5, Max-Q (SPICE), AWS,") + "\n")
	b.WriteString("             " + dimStyle.Render("MongoDB, PostgreSQL, Firebase, Git, Docker") + "\n\n")

	// Interests
	b.WriteString("  " + cyanStyle.Bold(true).Render("🌟 Interests") + "\n")
	b.WriteString("  " + dimStyle.Render("Building products that bridge the gap between") + "\n")
	b.WriteString("  " + dimStyle.Render("research and real-world deployment. Passionate") + "\n")
	b.WriteString("  " + dimStyle.Render("about open source, developer tools, and making") + "\n")
	b.WriteString("  " + dimStyle.Render("AI systems that actually work in production.") + "\n\n")

	b.WriteString(divider + "\n")
	b.WriteString("\n")
	b.WriteString("  " + dimStyle.Render("[esc to go back]") + "\n")

	return b.String()
}
