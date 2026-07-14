package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/trafalgar-2006/ssh-portfolio/config"
)

type Project struct {
	Title       string
	Description string
	Tags        []string
	Highlight   string
	Status      string // "Live", "Research", "WIP"
	GitHubURL   string
}

var AllProjects = []Project{
	{
		Title:       "EmbedGen — Embedded Systems LLM",
		Description: "Custom LLM fine-tuned on embedded-systems codebases using LoRA adapters. Quantized to GGUF for deployment on edge devices. Trained on RTX A6000 at Manipal research lab.",
		Tags:        []string{"Python", "PyTorch", "LoRA", "GGUF", "CUDA"},
		Status:      "WIP",
		GitHubURL:   "github.com/trafalgar-2006/EmbedGen",
	},
	{
		Title:       "ISRO LEOS — Sim-to-Real Vision",
		Description: "Synthetic data generation using BlenderProc for 6,000+ satellite images with COCO annotations. YOLOv7+Segmentation on NVIDIA Jetson via TensorRT. GAN-based style transfer with Unreal Engine 5.",
		Tags:        []string{"BlenderProc", "YOLOv7", "TensorRT", "Python", "OpenCV"},
		Highlight:   "22 FPS real-time",
		Status:      "Research",
		GitHubURL:   "",
	},
	{
		Title:       "Autonomous Trading System",
		Description: "AI-driven swing trader with market-wide signal generation, paper trading via Alpaca API, risk management, and automated order execution. Deployed on Oracle Cloud Always Free tier.",
		Tags:        []string{"Python", "SQLite", "Alpaca API", "Oracle Cloud"},
		Status:      "Live",
		GitHubURL:   "github.com/trafalgar-2006/trading-agent",
	},
	{
		Title:       "Webcraft Studios Platform",
		Description: "Full digital agency platform with SEO-driven architecture, scalable AI pipelines, and performance-optimized delivery for global client base.",
		Tags:        []string{"React", "Node.js", "MongoDB", "AWS", "Python"},
		Status:      "Live",
		GitHubURL:   "",
	},
	{
		Title:       "Bird Sound Species Detector",
		Description: "CNN-based audio classifier for bird species identification using spectrogram analysis. Trained on a custom dataset.",
		Tags:        []string{"Python", "TensorFlow", "CNNs"},
		Highlight:   "90%+ accuracy",
		Status:      "Research",
		GitHubURL:   "",
	},
	{
		Title:       "Frog Call Classifier",
		Description: "ML pipeline for frog species identification across 803+ species using MFCC-based feature extraction from environmental audio recordings.",
		Tags:        []string{"Python", "Librosa", "Scikit-learn"},
		Status:      "Research",
		GitHubURL:   "",
	},
	{
		Title:       "SSH Portfolio (this one)",
		Description: "This portfolio — an interactive TUI accessible over SSH from anywhere in the world. Auto-deploys from GitHub via Railway.",
		Tags:        []string{"Go", "Docker", "Railway"},
		Status:      "Live",
		GitHubURL:   "github.com/trafalgar-2006/portflio",
	},
}

// LoadFromConfig overwrites AllProjects and AllContacts from the loaded YAML config.
// Called at startup if content.yaml was found. Falls back to hardcoded data if not.
func LoadFromConfig() {
	if config.Loaded == nil {
		return
	}
	var projects []Project
	for _, p := range config.Loaded.Projects {
		projects = append(projects, Project{
			Title:       p.Title,
			Description: p.Description,
			Tags:        p.Tags,
			Status:      p.Status,
			GitHubURL:   p.GitHubURL,
			Highlight:   p.Highlight,
		})
	}
	if len(projects) > 0 {
		AllProjects = projects
	}
	loadContactsFromConfig()
}


func tagColor(tag string, theme Theme) lipgloss.Color {
	langs := map[string]bool{
		"Python": true, "Go": true, "JavaScript": true, "TypeScript": true,
		"C++": true, "Rust": true, "SQL": true, "HTML": true,
	}
	ml := map[string]bool{
		"PyTorch": true, "TensorFlow": true, "LoRA": true, "GGUF": true, "ONNX": true,
		"YOLOv7": true, "TensorRT": true, "OpenCV": true, "CNNs": true, "GANs": true,
		"Scikit-learn": true, "Librosa": true, "BlenderProc": true, "CUDA": true,
		"ResNet": true, "Transformers": true, "NumPy": true, "Pandas": true,
	}
	cloud := map[string]bool{
		"AWS": true, "Firebase": true, "Docker": true, "Railway": true,
		"Vercel": true, "Oracle Cloud": true, "Fly.io": true, "Alpaca API": true,
		"Wish": true, "SSH": true,
	}
	db := map[string]bool{"MongoDB": true, "PostgreSQL": true, "SQLite": true, "Redis": true}

	if langs[tag] { return theme.Primary }
	if ml[tag]    { return theme.Secondary }
	if cloud[tag] { return theme.Success }
	if db[tag]    { return theme.Warning }
	return theme.Purple
}

func statusDot(status string, theme Theme) string {
	switch status {
	case "Live":     return string(theme.Success)
	case "WIP":      return string(theme.Warning)
	case "Research": return string(theme.Purple)
	default:         return string(theme.Dim)
	}
}

func statusBadge(r *lipgloss.Renderer, status string, theme Theme) string {
	switch status {
	case "Live":
		return r.NewStyle().Foreground(theme.Success).Bold(true).Render("[Live]")
	case "WIP":
		return r.NewStyle().Foreground(theme.Warning).Bold(true).Render("[WIP]")
	case "Research":
		return r.NewStyle().Foreground(theme.Purple).Bold(true).Render("[Research]")
	default:
		return ""
	}
}

func RenderProjects(r *lipgloss.Renderer, width, height, cursor, scroll, projectsReveal, tagPopReveal int, livePulse bool, highlightY float64, decryptIdx int, decryptRunes []rune, theme Theme) string {
	goldStyle       := r.NewStyle().Foreground(theme.Accent).Bold(true)
	cyanStyle       := r.NewStyle().Foreground(theme.Primary)
	dimStyle        := r.NewStyle().Foreground(theme.Dim)
	dimMidStyle     := r.NewStyle().Foreground(theme.DimMid)
	orangeStyle     := r.NewStyle().Foreground(theme.Warning)
	greenStyle      := r.NewStyle().Foreground(theme.Success).Bold(true)
	selectedBg      := r.NewStyle().Background(theme.BoxBorder).Foreground(theme.Primary).Bold(true)
	unselectedStyle := r.NewStyle().Foreground(theme.DimMid)
	boxStyle        := r.NewStyle().Foreground(theme.BoxBorder)
	decryptStyle    := r.NewStyle().Foreground(theme.Text)
	hintStyle       := r.NewStyle().Foreground(theme.VeryDim).Italic(true)

	// Visual cursor row from lerp (rounded)
	visualCursor := int(highlightY + 0.5)
	if visualCursor < 0                  { visualCursor = 0 }
	if visualCursor >= len(AllProjects)  { visualCursor = len(AllProjects) - 1 }

	// ── Layout math ───────────────────────────────────────────────
	totalW  := width - 4
	leftW   := 32
	dividerW := 1
	rightW  := totalW - leftW - dividerW - 2
	if rightW < 30 { rightW = 30 }

	p := AllProjects[cursor]

	// ── LEFT PANEL — project list ──────────────────────────────────
	var left strings.Builder
	left.WriteString(cyanStyle.Bold(true).Render(" ✦ Projects") + "\n")
	left.WriteString(dimStyle.Render(" "+fmt.Sprintf("%d projects", len(AllProjects))) + "\n")
	left.WriteString(boxStyle.Render(" "+strings.Repeat("─", leftW-2)) + "\n")

	for i, proj := range AllProjects {
		if i >= projectsReveal {
			if i == projectsReveal {
				scanChar := []string{"▶", "▷", " "}[(projectsReveal/2)%3]
				left.WriteString(dimStyle.Render(" "+scanChar) + "\n")
			} else {
				left.WriteString("\n")
			}
			continue
		}

		isSelected := i == cursor
		isHover    := i == visualCursor && visualCursor != cursor
		num        := fmt.Sprintf("%d", i+1)

		title := proj.Title
		maxTitleW := leftW - 5
		if len([]rune(title)) > maxTitleW {
			runes := []rune(title)
			title = string(runes[:maxTitleW-1]) + "…"
		}

		var line string
		if isSelected {
			padding := leftW - len(num) - len([]rune(title)) - 3
			if padding < 0 { padding = 0 }
			line = selectedBg.Render(" "+num+" "+title+strings.Repeat(" ", padding))
		} else if isHover {
			line = r.NewStyle().Foreground(theme.Primary).Render(" "+num+" ") + dimMidStyle.Render(title)
		} else {
			line = unselectedStyle.Render(" "+num+" ") + dimStyle.Render(title)
		}

		// Status dot
		var dot string
		switch proj.Status {
		case "Live":
			if livePulse { dot = r.NewStyle().Foreground(theme.Success).Render("●") } else { dot = dimStyle.Render("○") }
		case "WIP":
			dot = r.NewStyle().Foreground(theme.Warning).Render("◐")
		case "Research":
			dot = r.NewStyle().Foreground(theme.Purple).Render("◇")
		default:
			dot = " "
		}
		left.WriteString(line + dot + "\n")
	}

	left.WriteString(boxStyle.Render(" "+strings.Repeat("─", leftW-2)) + "\n")
	left.WriteString(hintStyle.Render(" jk/↑↓  gg G  Nj  t theme") + "\n")

	// ── RIGHT PANEL — selected project detail with box border ──────
	var rightLines []string

	// Box top
	boxTop    := boxStyle.Render("╭"+strings.Repeat("─", rightW-2)+"╮")
	boxBottom := boxStyle.Render("╰"+strings.Repeat("─", rightW-2)+"╯")
	borderL   := boxStyle.Render("│")

	// Title + status badge
	titleLine := goldStyle.Render(p.Title)
	if p.Status == "Live" {
		if livePulse {
			titleLine += "  " + greenStyle.Render("● Live")
		} else {
			titleLine += "  " + dimStyle.Render("○ Live")
		}
	} else {
		titleLine += "  " + statusBadge(r, p.Status, theme)
	}
	if p.Highlight != "" {
		titleLine += "  " + greenStyle.Render("⚡ "+p.Highlight)
	}

	rightLines = append(rightLines, boxTop)
	for _, line := range strings.Split(wrapText(titleLine, rightW-4), "\n") {
		rightLines = append(rightLines, borderL+" "+line)
	}
	rightLines = append(rightLines, borderL+" "+boxStyle.Render(strings.Repeat("─", rightW-4)))
	rightLines = append(rightLines, borderL)

	// Description — decrypt reveal
	var descText string
	if len(decryptRunes) > 0 && decryptIdx < len(decryptRunes) {
		desc   := []rune(p.Description)
		merged := make([]rune, len(desc))
		for j := range desc {
			if j < decryptIdx {
				merged[j] = desc[j]
			} else if j < len(decryptRunes) {
				merged[j] = decryptRunes[j]
			} else {
				merged[j] = desc[j]
			}
		}
		descText = wrapText(string(merged), rightW-4)
	} else {
		descText = wrapText(p.Description, rightW-4)
	}
	for _, line := range strings.Split(descText, "\n") {
		rightLines = append(rightLines, borderL+" "+decryptStyle.Render(line))
	}
	rightLines = append(rightLines, borderL)

	// GitHub link
	rightLines = append(rightLines, borderL+" "+boxStyle.Render(strings.Repeat("─", rightW-4)))
	if p.GitHubURL != "" {
		rightLines = append(rightLines, borderL+" "+orangeStyle.Render("→ "+p.GitHubURL))
	} else {
		rightLines = append(rightLines, borderL+" "+dimStyle.Render("→ Private / Institutional repo"))
	}

	// Tags — pop in
	showCount := tagPopReveal
	if showCount > len(p.Tags) { showCount = len(p.Tags) }
	var tagParts []string
	for _, t := range p.Tags[:showCount] {
		tagS := r.NewStyle().Foreground(tagColor(t, theme))
		tagParts = append(tagParts, tagS.Render("["+t+"]"))
	}
	if len(tagParts) > 0 {
		rightLines = append(rightLines, borderL+" "+strings.Join(tagParts, " "))
	}

	// Other projects list
	rightLines = append(rightLines, borderL)
	rightLines = append(rightLines, borderL+" "+boxStyle.Render(strings.Repeat("─", rightW-4)))
	rightLines = append(rightLines, borderL+" "+dimStyle.Render("other projects"))
	for i, op := range AllProjects {
		if i == cursor || i >= projectsReveal { continue }
		dot := dimStyle.Render("·")
		rightLines = append(rightLines, borderL+" "+dot+" "+dimStyle.Render(fmt.Sprintf("%02d", i+1)+". ")+dimMidStyle.Render(op.Title))
	}
	rightLines = append(rightLines, boxBottom)

	// ── COMBINE — zip left and right lines side by side ────────────
	leftLines := strings.Split(left.String(), "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines { maxLines = len(rightLines) }
	for len(leftLines)  < maxLines { leftLines  = append(leftLines, "") }
	for len(rightLines) < maxLines { rightLines = append(rightLines, "") }

	sep := boxStyle.Render("│")
	var out strings.Builder
	out.WriteString("\n")
	for i := range leftLines {
		lLine := leftLines[i]
		rLine := rightLines[i]

		lVis := len([]rune(stripAnsi(lLine)))
		if lVis < leftW {
			lLine += strings.Repeat(" ", leftW-lVis)
		}
		out.WriteString(" "+lLine+" "+sep+" "+rLine+"\n")
	}
	return out.String()
}


func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	words := strings.Fields(text)
	var lines []string
	var currentLine strings.Builder
	for _, word := range words {
		if currentLine.Len()+len(word)+1 > maxWidth {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	return strings.Join(lines, "\n")
}
