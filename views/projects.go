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


func tagColor(tag string) lipgloss.Color {
	langs := map[string]bool{
		"Python": true, "Go": true, "JavaScript": true, "TypeScript": true,
		"C++": true, "Rust": true, "SQL": true,
	}
	ml := map[string]bool{
		"PyTorch": true, "TensorFlow": true, "LoRA": true, "GGUF": true,
		"YOLOv7": true, "TensorRT": true, "OpenCV": true, "CNNs": true,
		"Scikit-learn": true, "Librosa": true, "BlenderProc": true, "CUDA": true,
	}
	cloud := map[string]bool{
		"AWS": true, "Firebase": true, "Docker": true, "Railway": true,
		"Vercel": true, "Oracle Cloud": true, "Fly.io": true, "Alpaca API": true,
	}
	db := map[string]bool{"MongoDB": true, "PostgreSQL": true, "SQLite": true, "Redis": true}

	if langs[tag] {
		return lipgloss.Color("#00DFDF") // cyan — languages
	}
	if ml[tag] {
		return lipgloss.Color("#FF6AC1") // magenta — ML/AI
	}
	if cloud[tag] {
		return lipgloss.Color("#50FA7B") // green — cloud/infra
	}
	if db[tag] {
		return lipgloss.Color("#FFB86C") // orange — databases
	}
	return lipgloss.Color("#BD93F9") // purple — frameworks/other
}

func statusStyle(r *lipgloss.Renderer, status string) string {
	switch status {
	case "Live":
		return r.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true).Render("[Live]")
	case "WIP":
		return r.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true).Render("[WIP]")
	case "Research":
		return r.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Bold(true).Render("[Research]")
	default:
		return ""
	}
}

func RenderProjects(r *lipgloss.Renderer, width, height, cursor, scroll, projectsReveal, tagPopReveal int, livePulse bool) string {
	goldStyle        := r.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	descStyle        := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	greenStyle       := r.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	cyanStyle        := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle         := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	orangeStyle      := r.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	selectedBorder   := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	unselectedBorder := r.NewStyle().Foreground(lipgloss.Color("#333333"))

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ Projects") + "\n")
	b.WriteString(" " + dimStyle.Render("Things I've built") + "\n\n")

	contentWidth := width - 8
	if contentWidth < 50 {
		contentWidth = 50
	}
	if contentWidth > 90 {
		contentWidth = 90
	}

	for i, p := range AllProjects {
		// Cascade: only show up to projectsReveal
		if i >= projectsReveal {
			b.WriteString("\n")
			continue
		}
		isSelected := i == cursor
		num := fmt.Sprintf("%02d", i+1)

		var prefix string
		if isSelected {
			prefix = selectedBorder.Render("▎ ")
		} else {
			prefix = unselectedBorder.Render("  ")
		}

		// Title line
		titleLine := prefix + dimStyle.Render(num+". ") + goldStyle.Render(p.Title)

		// Status badge with [Live] pulse
		if p.Status == "Live" {
			if livePulse {
				titleLine += "  " + r.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true).Render("[Live]")
			} else {
				titleLine += "  " + r.NewStyle().Foreground(lipgloss.Color("#2A7A4A")).Bold(true).Render("[Live]")
			}
		} else {
			titleLine += "  " + statusStyle(r, p.Status)
		}
		if p.Highlight != "" {
			titleLine += "  " + greenStyle.Render("⚡ "+p.Highlight)
		}
		b.WriteString(" " + titleLine + "\n")

		if isSelected {
			desc := wrapText(p.Description, contentWidth-6)
			for _, line := range strings.Split(desc, "\n") {
				b.WriteString("     " + descStyle.Render(line) + "\n")
			}
			if p.GitHubURL != "" {
				b.WriteString("     " + orangeStyle.Render("→ "+p.GitHubURL) + "\n")
			} else {
				b.WriteString("     " + dimStyle.Render("→ Private / Institutional repo") + "\n")
			}
		}

		// Tags — pop in one by one for selected, all at once for others
		var tags []string
		showCount := len(p.Tags)
		if isSelected {
			showCount = tagPopReveal
			if showCount > len(p.Tags) { showCount = len(p.Tags) }
		}
		for _, t := range p.Tags[:showCount] {
			tagS := r.NewStyle().Foreground(tagColor(t))
			tags = append(tags, tagS.Render("["+t+"]"))
		}
		if len(tags) > 0 {
			b.WriteString("     " + strings.Join(tags, " ") + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(" " + dimStyle.Render("[↑↓ browse · enter selects · esc back]") + "\n")
	return b.String()
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
