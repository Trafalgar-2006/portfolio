package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	Title       string
	Description string
	Tags        []string
	Highlight   string // e.g. "85% accuracy", "22 FPS"
}

var AllProjects = []Project{
	{
		Title:       "Webcraft Studios Platform",
		Description: "Built comprehensive digital agency platform with modern web technologies and clean UI/UX design. Implemented SEO optimization and performance-driven architecture for global client base.",
		Tags:        []string{"React", "Node.js", "MongoDB", "AWS", "Python"},
	},
	{
		Title:       "Myntra Minis Integration",
		Description: "Developed custom e-commerce solutions and integration for Myntra's mini-app ecosystem. Delivered scalable design system resulting in improved conversion rates.",
		Tags:        []string{"JavaScript", "React Native", "REST APIs", "Firebase"},
	},
	{
		Title:       "Str8Fire.io Platform",
		Description: "Built complete SaaS platform with modern web stack, user administration, and dashboard analytics for content creators.",
		Tags:        []string{"Next.js", "TypeScript", "PostgreSQL", "Vercel", "Stripe API"},
	},
	{
		Title:       "Frog Call Classifier",
		Description: "Created a machine learning model to identify frog species with over 803 species using MFCC-based feature extraction. Implemented preprocessing pipeline for environmental audio classification.",
		Tags:        []string{"Python", "Librosa", "Scikit-learn", "MFCC Analysis"},
	},
	{
		Title:       "Bird Sound Species Detector",
		Description: "Built an AI system to classify bird species from audio clips using spectrogram analysis and CNNs.",
		Tags:        []string{"Python", "TensorFlow", "CNNs", "Spectrogram Analysis"},
		Highlight:   "90%+ accuracy",
	},
	{
		Title:       "ISRO LEOS — Sim-to-Real Vision",
		Description: "Synthetic data generation using BlenderProc for 6,000+ high-fidelity satellite images with COCO/RILE annotations. Trained YOLOv7+Segmentation model and deployed on NVIDIA Jetson using TensorRT, achieving 22 FPS. Implemented GAN-based style transfer using Unreal Engine 5.",
		Tags:        []string{"BlenderProc", "YOLOv7", "TensorRT", "NVIDIA Jetson", "Unreal Engine 5"},
		Highlight:   "22 FPS real-time",
	},
}

func RenderProjects(width, height, cursor, scroll int) string {
	goldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	selectedBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	unselectedBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(" " + cyanStyle.Bold(true).Render("✦ Projects") + "\n")
	b.WriteString(" " + dimStyle.Render("Things I've built and shipped") + "\n\n")

	contentWidth := width - 8
	if contentWidth < 50 {
		contentWidth = 50
	}
	if contentWidth > 80 {
		contentWidth = 80
	}

	for i, p := range AllProjects {
		isSelected := i == cursor

		// Project number
		num := fmt.Sprintf("%02d", i+1)

		// Border character
		var prefix string
		if isSelected {
			prefix = selectedBorder.Render("▎ ")
		} else {
			prefix = unselectedBorder.Render("  ")
		}

		// Title line
		titleLine := prefix + dimStyle.Render(num+". ") + goldStyle.Render(p.Title)
		if p.Highlight != "" {
			titleLine += "  " + greenStyle.Render("⚡ "+p.Highlight)
		}
		b.WriteString(" " + titleLine + "\n")

		// Description (only for selected)
		if isSelected {
			// Wrap description to fit
			desc := wrapText(p.Description, contentWidth-6)
			for _, line := range strings.Split(desc, "\n") {
				b.WriteString("     " + descStyle.Render(line) + "\n")
			}
		}

		// Tags
		var tags []string
		for _, t := range p.Tags {
			tags = append(tags, tagStyle.Render("["+t+"]"))
		}
		b.WriteString("     " + strings.Join(tags, " ") + "\n")
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(" " + dimStyle.Render("[↑↓ to browse · esc to go back]") + "\n")

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
