package views

import (
	"fmt"
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

// NameBannerLines returns the raw banner strings (used by glitch effect)
func NameBannerLines() []string {
	return nameBanner
}

func RenderHome(r *lipgloss.Renderer, width, height, revealIdx int, starBright []bool, taglineIdx int, taglineDone bool, cursorBlink bool, glitchFrames int, glitchRunes [][]rune, lastCommit string, sessionID string, connectedSecs int, buildInfo string, scanlineY int, idleGlitch bool, theme Theme) string {
	cyanStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle        := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	whiteStyle      := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	magentaStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	starStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.StarDim))
	brightStarStyle := r.NewStyle().Foreground(lipgloss.Color(theme.StarBright))

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

	// Stars decoration at top — each star blinks independently
	var starChars = []string{"✧", "*", "·", "✦", "*", "✧", "·", "✦"}
	var starRow1, starRow2 strings.Builder
	for i, ch := range starChars {
		bright := i < len(starBright) && starBright[i]
		var s string
		if bright {
			s = brightStarStyle.Render(ch)
		} else {
			s = starStyle.Render(ch)
		}
		if i < 4 {
			starRow1.WriteString(s)
			if i < 3 {
				starRow1.WriteString(" ")
			}
		} else {
			starRow2.WriteString(s)
			if i < 7 {
				starRow2.WriteString(" ")
			}
		}
	}
	rightCol.WriteString("    " + starRow1.String() + "\n")
	rightCol.WriteString(" + " + starRow2.String() + "\n")
	rightCol.WriteString("\n")

	// Name banner — reveal one line per tick, glitch on completion
	visibleBanner := revealIdx
	if visibleBanner > len(nameBanner) {
		visibleBanner = len(nameBanner)
	}
	for i, line := range nameBanner {
		if i < visibleBanner {
			// If glitching, use corrupted runes
			if glitchFrames > 0 && glitchRunes != nil && i < len(glitchRunes) {
				rightCol.WriteString(r.NewStyle().Foreground(lipgloss.Color("#FF6AC1")).Render(string(glitchRunes[i])) + "\n")
			} else {
				rightCol.WriteString(cyanStyle.Render(line) + "\n")
			}
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

	// Session info + last GitHub commit — dim lines at the bottom
	var bottomLines []string
	if sessionID != "" {
		mins := connectedSecs / 60
		secs := connectedSecs % 60
		metaStyle := r.NewStyle().Foreground(lipgloss.Color("#333333"))
		sessionStyle := r.NewStyle().Foreground(lipgloss.Color("#2A5A5A"))
		bottomLines = append(bottomLines, " "+sessionStyle.Render(fmt.Sprintf("session: %s  connected: %02d:%02d", sessionID, mins, secs))+"  "+metaStyle.Render(buildInfo))
	}
	if lastCommit != "" {
		commitStyle := r.NewStyle().Foreground(lipgloss.Color("#444444")).Italic(true)
		bottomLines = append(bottomLines, " "+commitStyle.Render("last pushed: "+lastCommit))
	}
	for _, bl := range bottomLines {
		combined.WriteString("\n" + bl + "\n")
	}

	result := combined.String()

	// CRT scanline overlay — a faint horizontal line sweeping down
	if scanlineY >= 0 {
		lines := strings.Split(result, "\n")
		if scanlineY < len(lines) {
			scanS := r.NewStyle().Foreground(lipgloss.Color("#0A2A2A")).Faint(true)
			// Overlay the scanline as a dim tinted version of that line
			visual := stripAnsi(lines[scanlineY])
			if len(visual) > 0 {
				lines[scanlineY] = scanS.Render(visual)
			}
			result = strings.Join(lines, "\n")
		}
	}

	// Idle ambient glitch — one-frame corruption across the whole screen
	if idleGlitch {
		glitchChars := []rune{'▓', '░', '▒', '▌', '▐', '╬', '╫', '╪'}
		lines := strings.Split(result, "\n")
		glitchS := r.NewStyle().Foreground(lipgloss.Color("#1A1A3A"))
		for i, line := range lines {
			visual := []rune(stripAnsi(line))
			for j := range visual {
				if visual[j] != ' ' && len(visual) > 0 {
					// ~8% corruption per non-space char
					if (i*len(visual)+j)%13 == 0 {
						visual[j] = glitchChars[(i*7+j*3)%len(glitchChars)]
					}
				}
			}
			lines[i] = glitchS.Render(string(visual))
		}
		result = strings.Join(lines, "\n")
	}

	return result
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
