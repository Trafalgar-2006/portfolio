package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Full-width braille portrait (50 chars per line).
// Show only when terminal is wide enough (>= 115 cols).
var portraitArt = []string{
	" ⠄⠠⡀⠄⠠⠀⠄⠠⠀⠄⢠⠀⠄⡠⠀⡄⠠⠀⠀⠠⠀⠄⠠⠀⠄⠀⠀⠀⠀⠀⠠⢀⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⠠⠀⠄⡀",
	" ⢈⠐⠠⢈⠡⠈⠄⠡⡈⢐⠠⢈⠐⠠⠁⠀⠀⠀⠠⠁⡈⠄⠀⠀⠠⠁⡈⠄⠀⠀⠠⠁⡈⠄⠀⠀⠠⠁⡈⠄⠀⠠⠁⡈⠄⠀⠠⠁⡈⠄",
	" ⠠⠌⡐⠂⠄⠡⢈⠐⡀⠢⠐⠂⢉⠐⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡐⠠⠄⠡⢈⡐⢈⡐⠠⢁⡘⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠠⠁⠌⡐⠠⠐⠠⢀⠃⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠡⠘⠠⠄⡑⢈⡁⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠂⡡⢁⢂⠰⢀⠰⢁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠡⠐⣀⠢⠐⢂⠰⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡈⠔⡀⢂⠁⡂⠄⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠤⡈⠔⠠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡁⠆⡈⠔⢂⠡⠌⠀⠀⠀⠀⠀⠀⠀⠀⣀⠂⡌⢀⠈⠀⡁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡡",
	" ⠐⡠⠁⢌⠠⠒⡈⠄⠀⠀⠀⠀⠀⠀⡐⠄⢮⣐⠣⣄⠡⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠊⠄⡑⠠⢂⠁⠆⢀⠠⠒⠤⢀⠀⡰⢈⡜⢮⣇⡛⢤⢋⡔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡡⠘⠠⢁⠂⣉⠰⠀⠆⢠⡐⠡⢂⠡⠒⢌⠲⠩⠜⡀⠂⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡡",
	" ⠡⠌⡁⠂⢌⠀⠆⢡⠚⡌⢡⠃⢌⠂⡉⠄⠃⠁⠀⠀⢤⡱⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡐⡈⠤⠑⡀⠊⠌⡄⠁⠈⠤⣉⠰⡈⢐⠈⠄⠀⠠⠙⢢⠣⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡐",
	" ⠒⡀⠆⠡⠐⣉⠰⠠⠀⠀⠒⡠⠑⡄⠃⠌⠀⠀⠀⠤⣀⡤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠆⠡⠌⢂⠡⠄⢂⠡⠃⠀⢁⠢⢑⠠⢉⠂⠠⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⡐⠡⡈⢂⠔⠨⠄⢂⠀⠀⡀⠆⣁⠢⠁⡌⢰⢩⠖⡔⢶⡠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡐",
	" ⡠⠡⠐⠡⣈⠂⢅⠂⠀⠀⡐⠄⡀⠂⢅⢢⡁⢎⡘⠌⠠⠙⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠",
	" ⢁⢂⠉⠔⠠⡈⠄⢊⢀⠰⡐⠠⢀⠉⠄⢂⠹⣄⠲⣈⢄⡐⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢁",
	" ⠂⡄⠊⠌⡐⠠⠑⣈⠐⢢⠡⠘⠠⢈⠐⠠⠁⠌⡱⠌⠦⡑⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠂",
	" ⡡⠄⠉⠒⠈⠁⠈⠀⠈⢆⠡⢃⠁⠂⠌⡐⠀⢀⠀⠌⠐⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡡",
	" ⠀⠀⠀⠀⠀⠀⠀⠀⢌⠢⡑⠌⡌⠌⡐⠀⠄⠂⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠀⠀⠀⠀⠀⠀⢰⢁⠢⢡⠘⡐⢌⠢⡁⠌⠀⠀⠀⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
}

// portraitWidth is the visual character width of the portrait art lines above.
const portraitWidth = 50

// Name banner — figlet-style block letters for MOHITH / AKSHAY
var nameBanner = []string{
	"███╗   ███╗  ██████╗  ██╗  ██╗  ██╗  ████████╗  ██╗  ██╗",
	"████╗ ████║ ██╔═══██╗ ██║  ██║  ██║  ╚══██╔══╝  ██║  ██║",
	"██╔████╔██║ ██║   ██║ ███████║  ██║     ██║     ███████║",
	"██║╚██╔╝██║ ██║   ██║ ██╔══██║  ██║     ██║     ██╔══██║",
	"██║ ╚═╝ ██║ ╚██████╔╝ ██║  ██║  ██║     ██║     ██║  ██║",
	"╚═╝     ╚═╝  ╚═════╝  ╚═╝  ╚═╝  ╚═╝     ╚═╝     ╚═╝  ╚═╝",
	"",
	" █████╗  ██╗ ██╗  ███████╗ ██╗  ██╗  █████╗   ██╗   ██╗",
	"██╔══██╗ ██║ ██╔╝ ██╔════╝ ██║  ██║ ██╔══██╗  ╚██╗ ██╔╝",
	"███████║ █████╔╝  ███████╗ ███████║ ███████║   ╚████╔╝ ",
	"██╔══██║ ██╔═██╗  ╚════██║ ██╔══██║ ██╔══██║    ╚██╔╝  ",
	"██║  ██║ ██║  ██╗ ███████║ ██║  ██║ ██║  ██║     ██║   ",
	"╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝     ╚═╝   ",
}

const TaglineText = "is an engineer, builder & creator who turns ideas into products."

func BannerLines() int         { return len(nameBanner) }
func NameBannerLines() []string { return nameBanner }

func RenderHome(r *lipgloss.Renderer, width, height, revealIdx int, starBright []bool, taglineIdx int, taglineDone bool, cursorBlink bool, glitchFrames int, glitchRunes [][]rune, lastCommit string, sessionID string, connectedSecs int, buildInfo string, scanlineY int, idleGlitch bool, theme Theme) string {
	cyanStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle        := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	whiteStyle      := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	magentaStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	starStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.StarDim))
	brightStarStyle := r.NewStyle().Foreground(lipgloss.Color(theme.StarBright))

	// Show portrait only on wide-enough terminals (portrait=50 + gap=4 + banner~62 = 116)
	showPortrait := width >= 115

	// ── LEFT COLUMN: portrait ──────────────────────────────────────────────
	var leftCol strings.Builder
	portReveal := revealIdx * 2
	if portReveal > len(portraitArt) {
		portReveal = len(portraitArt)
	}

	if showPortrait {
		for i, line := range portraitArt {
			if i < portReveal {
				leftCol.WriteString(cyanStyle.Render(line))
			} else {
				leftCol.WriteString(strings.Repeat(" ", portraitWidth))
			}
			leftCol.WriteString("\n")
		}
	}

	// ── RIGHT COLUMN: name + bio ───────────────────────────────────────────
	var rightCol strings.Builder

	// Stars — independent twinkle
	starChars := []string{"✧", "*", "·", "✦", "*", "✧", "·", "✦"}
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
			if i < 3 { starRow1.WriteString(" ") }
		} else {
			starRow2.WriteString(s)
			if i < 7 { starRow2.WriteString(" ") }
		}
	}
	rightCol.WriteString("    " + starRow1.String() + "\n")
	rightCol.WriteString(" +  " + starRow2.String() + "\n")
	rightCol.WriteString("\n")

	// Name banner — reveal one line per tick
	visibleBanner := revealIdx
	if visibleBanner > len(nameBanner) {
		visibleBanner = len(nameBanner)
	}
	for i, line := range nameBanner {
		if i < visibleBanner {
			if glitchFrames > 0 && glitchRunes != nil && i < len(glitchRunes) {
				rightCol.WriteString(r.NewStyle().Foreground(lipgloss.Color(theme.Secondary)).Render(string(glitchRunes[i])) + "\n")
			} else {
				rightCol.WriteString(cyanStyle.Render(line) + "\n")
			}
		} else {
			rightCol.WriteString("\n")
		}
	}

	// Surname subtitle
	if revealIdx >= len(nameBanner) {
		rightCol.WriteString(magentaStyle.Render("  ·  D U G G I R A L A  ·") + "\n")
	} else {
		rightCol.WriteString("\n")
	}
	rightCol.WriteString("\n")
	rightCol.WriteString("  " + brightStarStyle.Render("✦") + "  " + starStyle.Render("·") + "\n")
	rightCol.WriteString("\n")

	// Typewriter tagline
	if revealIdx >= len(nameBanner) {
		runes := []rune(TaglineText)
		end := taglineIdx
		if end > len(runes) { end = len(runes) }
		cursor := ""
		if taglineIdx < len(runes) {
			cursor = "█"
		} else if cursorBlink {
			cursor = "█"
		}
		rightCol.WriteString("  " + whiteStyle.Render(string(runes[:end])+cursor) + "\n")
	} else {
		rightCol.WriteString("\n")
	}

	// Bio lines
	bioLines := []string{
		dimStyle.Render("Founder & Lead Engineer of") + " " + magentaStyle.Render("Webcraft Studios") + dimStyle.Render(","),
		dimStyle.Render("building full-stack apps, CV pipelines, and AI systems."),
		"",
		dimStyle.Render("B.Tech ECE  ·  Manipal Institute of Technology, Bengaluru"),
		dimStyle.Render("Aug 2023 – Jul 2027"),
		"",
		dimStyle.Render("President, ") + cyanStyle.Render("MBOSC") + dimStyle.Render("  ·  Manipal Bengaluru Open Source Community"),
		"",
		dimStyle.Render("SWE Intern @ ") + cyanStyle.Render("SenseOps Tech Solutions") + dimStyle.Render("  (May–Jul 2025)"),
		dimStyle.Render("CV Research Intern @ ") + cyanStyle.Render("ISRO – LEOS") + dimStyle.Render("  (Dec 2025–Jan 2026)"),
	}
	for _, line := range bioLines {
		rightCol.WriteString("  " + line + "\n")
	}

	// ── COMBINE columns ────────────────────────────────────────────────────
	leftLines  := strings.Split(leftCol.String(), "\n")
	rightLines := strings.Split(rightCol.String(), "\n")

	maxLines := len(rightLines)
	if showPortrait && len(leftLines) > maxLines {
		maxLines = len(leftLines)
	}
	for len(leftLines)  < maxLines { leftLines  = append(leftLines, "") }
	for len(rightLines) < maxLines { rightLines = append(rightLines, "") }

	var combined strings.Builder
	combined.WriteString("\n")
	gap := "    " // gap between portrait and name banner

	availHeight := height - 6
	if availHeight < 10 { availHeight = maxLines }
	renderLines := maxLines
	if renderLines > availHeight { renderLines = availHeight }

	for i := 0; i < renderLines; i++ {
		left  := ""
		right := ""
		if i < len(rightLines) { right = rightLines[i] }

		if showPortrait {
			if i < len(leftLines) { left = leftLines[i] }
			// Pad left column to consistent width so right column stays aligned
			lVis := len([]rune(stripAnsi(left)))
			if lVis < portraitWidth {
				left += strings.Repeat(" ", portraitWidth-lVis)
			}
			combined.WriteString(" " + left + gap + right + "\n")
		} else {
			combined.WriteString(" " + right + "\n")
		}
	}

	// Session / commit info
	if sessionID != "" {
		mins := connectedSecs / 60
		secs := connectedSecs % 60
		metaStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim))
		sessionStyle := r.NewStyle().Foreground(lipgloss.Color(theme.FooterText))
		combined.WriteString("\n " + sessionStyle.Render(fmt.Sprintf("session: %s  connected: %02d:%02d", sessionID, mins, secs)) +
			"  " + metaStyle.Render(buildInfo) + "\n")
	}
	if lastCommit != "" {
		commitStyle := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim)).Italic(true)
		combined.WriteString(" " + commitStyle.Render("last pushed: "+lastCommit) + "\n")
	}

	result := combined.String()

	// CRT scanline overlay
	if scanlineY >= 0 {
		lines := strings.Split(result, "\n")
		if scanlineY < len(lines) {
			scanS  := r.NewStyle().Foreground(lipgloss.Color(theme.ScanlineColor)).Faint(true)
			visual := stripAnsi(lines[scanlineY])
			if len(visual) > 0 {
				lines[scanlineY] = scanS.Render(visual)
			}
			result = strings.Join(lines, "\n")
		}
	}

	// Idle ambient glitch
	if idleGlitch {
		glitchChars := []rune{'▓', '░', '▒', '▌', '▐', '╬', '╫', '╪'}
		lines  := strings.Split(result, "\n")
		glitchS := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim))
		for i, line := range lines {
			visual := []rune(stripAnsi(line))
			for j := range visual {
				if visual[j] != ' ' && (i*len(visual)+j)%13 == 0 {
					visual[j] = glitchChars[(i*7+j*3)%len(glitchChars)]
				}
			}
			lines[i] = glitchS.Render(string(visual))
		}
		result = strings.Join(lines, "\n")
	}

	return result
}

// stripAnsi removes ANSI escape codes for visual-width calculation.
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
