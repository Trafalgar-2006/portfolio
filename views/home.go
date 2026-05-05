package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII art portrait — custom user provided art
var portraitArt = []string{
	" ⠄⠠⡀⠄⠠⠀⠄⠠⠀⠄⢠⠀⠄⡠⠀⡄⠠⠀⠀⠠⠀⠄⠠⠀⠄⠀⠀⠀⠀⠀⠠⢀⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⡀",
	" ⢈⠐⠠⢈⠡⠈⠄⠡⡈⢐⠠⢈⠐⠠⠁⠀⠀⠀⠠⠁⡈⠄⡑⠈⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠂⠁⠈⠄⠡⠈⠄⠡⠈⠄⠡⠈⠄⡁⢂⠐",
	" ⠠⠌⡐⠂⠄⠡⢈⠐⡀⠢⠐⠂⢉⠐⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠁⠌⠠⠁⠌⠠⢁⠂⡐⠠⢈",
	" ⡐⠠⠄⠡⢈⡐⢈⡐⠠⢁⡘⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠈⠄⡁⢂⠐⠠⢁⠂",
	" ⠠⠁⠌⡐⠠⠐⠠⢀⠃⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠂⠐⡀⠂⠌⡐⠠⢈",
	" ⠡⠘⠠⠄⡑⢈⡁⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠀⠀⠄⡁⠂⠄⡁⢂",
	" ⠂⡡⢁⢂⠰⢀⠰⢁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠡⢀⠡⢂⠐⠠",
	" ⠡⠐⣀⠢⠐⢂⠰⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⢁⠂⠤⠘⠠",
	" ⡈⠔⡀⢂⠁⡂⠄⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠤⡈⠔⠠⠈⠀⠑⠢⢄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠠⠈⠄⡑⠈",
	" ⡁⠆⡈⠔⢂⠡⠌⠀⠀⠀⠀⠀⠀⠀⠀⣀⠂⡌⢀⠈⠀⡁⠀⠀⠀⠀⠀⣉⠒⠠⠄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠂⠠⠈⠄⡡",
	" ⠐⡠⠁⢌⠠⠒⡈⠄⠀⠀⠀⠀⠀⠀⡐⠄⢮⣐⠣⣄⠡⠀⠄⠀⠀⢠⠰⢁⡀⠰⠁⠀⠑⠂⠤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡁⠆⡐",
	" ⠊⠄⡑⠠⢂⠁⠆⢀⠠⠒⠤⢀⠀⡰⢈⡜⢮⣇⡛⢤⢋⡔⡂⢆⣡⢎⠶⣡⠀⠄⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠰⠀⠒⠠",
	" ⡡⠘⠠⢁⠂⣉⠰⠀⠆⢠⡐⠡⢂⠡⠒⢌⠲⠩⠜⡀⠂⠈⠰⣎⡵⣮⢳⡡⠊⢄⡀⠂⡀⠀⠀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠄⡈⠤⠁",
	" ⠡⠌⡁⠂⢌⠀⠆⢡⠚⡌⢡⠃⢌⠂⡉⠄⠃⠁⠀⠀⢤⡱⢃⡈⠀⠁⠃⠱⢁⠦⡐⢆⠠⢄⠠⡀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠐⠠⢁⢂⠡",
	" ⡐⡈⠤⠑⡀⠊⠌⡄⠁⠈⠤⣉⠰⡈⢐⠈⠄⠀⠠⠙⢢⠣⠅⡀⠀⠀⠈⠀⠀⠂⠈⢍⠺⣌⠑⠠⠀⠀⠀⠀⠀⠀⠀⠀⠡⠈⠔⠂⠄⠊",
	" ⠒⡀⠆⠡⠐⣉⠰⠠⠀⠀⠒⡠⠑⡄⠃⠌⠀⠀⠀⠤⣀⡤⣐⣀⠂⠘⢀⠂⡄⠃⠀⢀⠃⠤⠉⠄⠀⠀⠀⠀⠀⡀⢂⠡⠌⠡⢈⡐⠉⡄",
	" ⠆⠡⠌⢂⠡⠄⢂⠡⠃⠀⢁⠢⢑⠠⢉⠂⠠⡄⠀⠀⠀⠀⠁⠙⠚⠣⢄⠀⠀⠌⠀⠀⡈⠄⡑⠀⠀⠀⠄⢂⠡⠐⢂⡐⠈⠤⠁⡄⠡⢀",
	" ⡐⠡⡈⢂⠔⠨⠄⢂⠀⠀⡀⠆⣁⠢⠁⡌⢰⢩⠖⡔⢶⡠⠀⠀⠀⠀⠀⠀⠀⠀⡀⠐⢀⠐⠀⠀⠠⢉⠰⠈⠄⠃⠤⢀⡉⠄⡡⠠⢁⠂",
	" ⡠⠡⠐⠡⣈⠂⢅⠂⠀⠀⡐⠄⡀⠂⢅⢢⡁⢎⡘⠌⠠⠙⠹⠒⠧⢀⠠⠘⠀⠠⠀⡁⠠⢈⠂⠌⡐⢁⠂⡉⠤⢉⠐⡠⠐⢂⠁⢂⠡⠈",
	" ⢁⢂⠉⠔⠠⡈⠄⢊⢀⠰⡐⠠⢀⠉⠄⢂⠹⣄⠲⣈⢄⡐⠀⠁⡐⠠⢂⠡⠈⢀⠐⠀⡰⠀⠌⡐⠠⠂⢄⢁⠂⡂⠡⠄⡑⢈⠰⠈⠄⡡",
	" ⠂⡄⠊⠌⡐⠠⠑⣈⠐⢢⠡⠘⠠⢈⠐⠠⠁⠌⡱⠌⠦⡑⢎⡰⢀⠡⠂⡐⠀⠀⠀⢠⠁⠌⡐⠠⠡⠌⡀⠆⢂⡁⠒⣈⠐⠌⡠⠑⡈⠄",
	" ⡡⠄⠉⠒⠈⠁⠈⠀⠈⢆⠡⢃⠁⠂⠌⡐⠀⢀⠀⠌⠐⠁⠎⠐⠁⠂⠁⠀⠀⠰⢈⠄⠌⡐⠠⠑⠠⢂⠡⠈⡄⠰⠁⡄⠌⢂⠄⡁⠆⡈",
	" ⠀⠀⠀⠀⠀⠀⠀⠀⢌⠢⡑⠌⡌⠌⡐⠀⠄⠂⠀⠀⡀⠀⠀⠀⠀⠀⠀⠄⡈⠐⠂⠌⡐⠠⠑⡈⢁⠂⠄⠃⠄⡡⠒⠠⠘⣀⠢⠐⢂⠡",
	" ⠀⠀⠀⠀⠀⠀⢰⢁⠢⢡⠘⡐⢌⠢⡁⠌⠀⠀⠀⠄⠀⠀⠀⠈⢀⠠⢁⢂⠐⠀⠈⠀⠐⠡⡁⠌⠠⠌⡀⠃⠌⠠⢁⠊⠡⢀⠂⠱⢀⠂",
}

// Name banner — figlet-style ASCII art for "Mohith Akshay"
var nameBanner = []string{
	"███╗   ███╗   ██████╗   ██╗  ██╗  ██╗  ████████╗  ██╗  ██╗",
	"████╗ ████║  ██╔═══██╗  ██║  ██║  ██║  ╚══██╔══╝  ██║  ██║",
	"██╔████╔██║  ██║   ██║  ███████║  ██║     ██║     ███████║",
	"██║╚██╔╝██║  ██║   ██║  ██╔══██║  ██║     ██║     ██╔══██║",
	"██║ ╚═╝ ██║  ╚██████╔╝  ██║  ██║  ██║     ██║     ██║  ██║",
	"╚═╝     ╚═╝   ╚═════╝   ╚═╝  ╚═╝  ╚═╝     ╚═╝     ╚═╝  ╚═╝",
	"",
	" █████╗   ██╗  ██╗  ███████╗  ██╗  ██╗   █████╗   ██╗   ██╗",
	"██╔══██╗  ██║ ██╔╝  ██╔════╝  ██║  ██║  ██╔══██╗  ╚██╗ ██╔╝",
	"███████║  █████╔╝   ███████╗  ███████║  ███████║   ╚████╔╝ ",
	"██╔══██║  ██╔═██╗   ╚════██║  ██╔══██║  ██╔══██║    ╚██╔╝  ",
	"██║  ██║  ██║  ██╗  ███████║  ██║  ██║  ██║  ██║     ██║   ",
	"╚═╝  ╚═╝  ╚═╝  ╚═╝  ╚══════╝  ╚═╝  ╚═╝  ╚═╝  ╚═╝     ╚═╝   ",
}

// Star field decorations
var starPositions = []struct {
	x, y int
	char string
}{
	{2, 0, "✧"},
	{15, 1, "*"},
	{8, 2, "·"},
	{20, 0, "✦"},
	{5, 3, "*"},
	{22, 2, "✧"},
	{12, 4, "·"},
	{1, 5, "✦"},
}

// TaglineText is the full tagline revealed by the typewriter animation.
const TaglineText = "is an engineer, builder & creator who turns ideas into products."

// BannerLines returns the number of lines in the name banner for the animation ticker.
func BannerLines() int {
	return len(nameBanner)
}

func RenderHome(r *lipgloss.Renderer, width, height, revealIdx int, blink bool, taglineIdx int, taglineDone bool, cursorBlink bool) string {
	cyanStyle     := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))
	dimStyle      := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	whiteStyle    := r.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	magentaStyle  := r.NewStyle().Foreground(lipgloss.Color("#FF6AC1"))
	starStyle     := r.NewStyle().Foreground(lipgloss.Color("#888888"))
	brightStarStyle := r.NewStyle().Foreground(lipgloss.Color("#00DFDF"))

	// Calculate available space
	maxContentWidth := width - 4
	if maxContentWidth < 60 {
		maxContentWidth = 60
	}

	// Left column: portrait art — reveals at 2× banner speed
	portraitWidth := 50
	portReveal := revealIdx * 2
	if portReveal > len(portraitArt) {
		portReveal = len(portraitArt)
	}
	var leftCol strings.Builder
	for i, line := range portraitArt {
		if i < portReveal {
			leftCol.WriteString(cyanStyle.Render(line))
		} else {
			// Hold space so layout doesn't jump
			visLen := len([]rune(stripAnsi(line)))
			if visLen < 1 { visLen = 1 }
			leftCol.WriteString(strings.Repeat(" ", visLen))
		}
		leftCol.WriteString("\n")
	}

	// Right column: name + bio info
	var rightCol strings.Builder

	// Stars decoration at top — blink on slow tick
	if blink {
		rightCol.WriteString(starStyle.Render("    ✧") + brightStarStyle.Render("*") + starStyle.Render("·✦") + "\n")
	} else {
		rightCol.WriteString(starStyle.Render("    · ") + brightStarStyle.Render("✦") + starStyle.Render("·· ") + "\n")
	}
	rightCol.WriteString(starStyle.Render(" +") + "     " + starStyle.Render("*") + "\n")
	rightCol.WriteString("\n")

	// Name banner — reveal one line per tick
	visibleBanner := revealIdx
	if visibleBanner > len(nameBanner) {
		visibleBanner = len(nameBanner)
	}
	for i, line := range nameBanner {
		if i < visibleBanner {
			rightCol.WriteString(cyanStyle.Render(line) + "\n")
		} else {
			rightCol.WriteString("\n") // hold space
		}
	}
	// DUGGIRALA subtitle — only after banner fully revealed
	if revealIdx >= len(nameBanner) {
		rightCol.WriteString(magentaStyle.Render("  ·  D U G G I R A L A  ·") + "\n")
	} else {
		rightCol.WriteString("\n")
	}
	rightCol.WriteString("\n")

	// Stars
	rightCol.WriteString("  " + brightStarStyle.Render("✦") + "  " + starStyle.Render("·") + "\n")
	rightCol.WriteString("\n")

	// Bio text — right side
	// Tagline (typewriter)
	if revealIdx >= len(nameBanner) {
		visible := []rune(TaglineText)[:taglineIdx]
		cursor := ""
		if taglineIdx < len([]rune(TaglineText)) {
			cursor = "█"
		} else if cursorBlink {
			cursor = "█"
		}
		rightCol.WriteString("  " + whiteStyle.Render(string(visible)+cursor) + "\n")
	} else {
		rightCol.WriteString("\n")
	}

	bioLines := []string{
		whiteStyle.Render("is an engineer, builder &"),
		whiteStyle.Render("creator who turns ideas"),
		whiteStyle.Render("into products."),
		"",
		dimStyle.Render("Founder & Lead Designer of"),
		magentaStyle.Render("Webcraft Studios") + dimStyle.Render(","),
		dimStyle.Render("building full-stack apps,"),
		dimStyle.Render("Computer Vision pipelines,"),
		dimStyle.Render("and scalable AI perception systems."),
		"",
		dimStyle.Render("B.Tech in Electronics & Computer"),
		dimStyle.Render("Engineering @ Manipal Institute"),
		dimStyle.Render("of Technology (Aug 2023–Jul 2027),"),
		dimStyle.Render("Bengaluru."),
		"",
		dimStyle.Render("President of ") + cyanStyle.Render("MBOSC"),
		dimStyle.Render("(Manipal Bengaluru Open Source"),
		dimStyle.Render("Community)"),
		"",
		dimStyle.Render("Former CV Research Intern at"),
		cyanStyle.Render("ISRO – LEOS") + dimStyle.Render("."),
	}

	for _, line := range bioLines {
		rightCol.WriteString("  " + line + "\n")
	}

	leftContent := leftCol.String()
	rightContent := rightCol.String()

	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	// Pad to equal height
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	for len(leftLines) < maxLines {
		leftLines = append(leftLines, strings.Repeat(" ", portraitWidth))
	}
	for len(rightLines) < maxLines {
		rightLines = append(rightLines, "")
	}

	// Join columns side by side
	var combined strings.Builder
	combined.WriteString("\n")
	gap := "   "

	availHeight := height - 6 // leave room for tabs + hint
	if availHeight < 10 {
		availHeight = maxLines
	}

	renderLines := maxLines
	if renderLines > availHeight {
		renderLines = availHeight
	}

	for i := 0; i < renderLines; i++ {
		left := leftLines[i]
		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}
		// Pad left column to consistent width
		leftVisible := stripAnsi(left)
		padding := portraitWidth - len([]rune(leftVisible))
		if padding < 0 {
			padding = 0
		}
		combined.WriteString(" " + left + strings.Repeat(" ", padding) + gap + right + "\n")
	}

	return combined.String()
}

// stripAnsi removes ANSI escape codes for width calculation
func stripAnsi(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
