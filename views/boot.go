package views

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BootLine is a single line in the boot sequence
type BootLine struct {
	Text      string
	Color     string // "normal", "red", "green", "dim"
	DelayTicks int   // ticks before next line appears (at 50ms/tick)
}

// NewBootSequence builds the boot lines with a unique random fingerprint
func NewBootSequence() []BootLine {
	fp := randomFingerprint()
	return []BootLine{
		{"> connecting to ssh-portfolio:22...", "dim", 6},
		{"  SSH-2.0-OpenSSH_9.6p1 portfolio", "dim", 4},
		{"> negotiating: chacha20-poly1305@openssh.com", "dim", 5},
		{fmt.Sprintf("> host key: SHA256:%s", fp), "dim", 5},
		{"> løâdïñg modules...", "dim", 3},
		{"> ERROR: signal lost — attempting recovery...", "red", 8},
		{"> .", "dim", 4},
		{"> ..", "dim", 4},
		{"> ...", "dim", 4},
		{"> recovery successful.", "green", 6},
		{"> decrypting portfolio...", "dim", 4},
		{"> ✓ Identity verified. Welcome, visitor.", "green", 10},
	}
}

func randomFingerprint() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789+/"
	b := make([]byte, 28)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func RenderBoot(r *lipgloss.Renderer, width, height, visibleLines int, lines []BootLine) string {
	dimStyle   := r.NewStyle().Foreground(lipgloss.Color("#666666"))
	redStyle   := r.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	greenStyle := r.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	cyanStyle  := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))

	var b strings.Builder

	// Center vertically
	topPad := (height - len(lines) - 4) / 2
	if topPad < 2 {
		topPad = 2
	}
	for i := 0; i < topPad; i++ {
		b.WriteString("\n")
	}

	b.WriteString("  " + cyanStyle.Render("ssh-portfolio") + "\n\n")

	show := visibleLines
	if show > len(lines) {
		show = len(lines)
	}
	for i := 0; i < show; i++ {
		line := lines[i]
		var rendered string
		switch line.Color {
		case "red":
			rendered = redStyle.Render(line.Text)
		case "green":
			rendered = greenStyle.Render(line.Text)
		default:
			rendered = dimStyle.Render(line.Text)
		}
		b.WriteString("  " + rendered + "\n")
	}

	return b.String()
}
