package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII braille portrait — trimmed to 24 visible chars so it fits on 80-col terminals
var portraitArt = []string{
	" ⠄⠠⡀⠄⠠⠀⠄⠠⠀⠄⢠⠀⠄⡠⠀⡄⠠⠀⠀⠠⠀⠄",
	" ⢈⠐⠠⢈⠡⠈⠄⠡⡈⢐⠠⢈⠐⠠⠁⠀⠀⠀⠠⠁⡈⠄",
	" ⠠⠌⡐⠂⠄⠡⢈⠐⡀⠢⠐⠂⢉⠐⡀⠀⠀⠀⠀⠀⠀⠀",
	" ⡐⠠⠄⠡⢈⡐⢈⡐⠠⢁⡘⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠠⠁⠌⡐⠠⠐⠠⢀⠃⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠡⠘⠠⠄⡑⢈⡁⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠂⡡⢁⢂⠰⢀⠰⢁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	" ⠡⠐⣀⠢⠐⢂⠰⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀",
	" ⡈⠔⡀⢂⠁⡂⠄⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠤⡈⠔⠠",
	" ⡁⠆⡈⠔⢂⠡⠌⠀⠀⠀⠀⠀⠀⠀⠀⣀⠂⡌⢀⠈⠀⡁",
	" ⠐⡠⠁⢌⠠⠒⡈⠄⠀⠀⠀⠀⠀⠀⡐⠄⢮⣐⠣⣄⠡⠀",
	" ⠊⠄⡑⠠⢂⠁⠆⢀⠠⠒⠤⢀⠀⡰⢈⡜⢮⣇⡛⢤⢋⡔",
	" ⡡⠘⠠⢁⠂⣉⠰⠀⠆⢠⡐⠡⢂⠡⠒⢌⠲⠩⠜⡀⠂⠈",
	" ⠡⠌⡁⠂⢌⠀⠆⢡⠚⡌⢡⠃⢌⠂⡉⠄⠃⠁⠀⠀⢤⡱",
	" ⡐⡈⠤⠑⡀⠊⠌⡄⠁⠈⠤⣉⠰⡈⢐⠈⠄⠀⠠⠙⢢⠣",
	" ⠒⡀⠆⠡⠐⣉⠰⠠⠀⠀⠒⡠⠑⡄⠃⠌⠀⠀⠀⠤⣀⡤",
	" ⠆⠡⠌⢂⠡⠄⢂⠡⠃⠀⢁⠢⢑⠠⢉⠂⠠⡄⠀⠀⠀⠀",
	" ⡐⠡⡈⢂⠔⠨⠄⢂⠀⠀⡀⠆⣁⠢⠁⡌⢰⢩⠖⡔⢶⡠",
	" ⡠⠡⠐⠡⣈⠂⢅⠂⠀⠀⡐⠄⡀⠂⢅⢢⡁⢎⡘⠌⠠⠙",
	" ⢁⢂⠉⠔⠠⡈⠄⢊⢀⠰⡐⠠⢀⠉⠄⢂⠹⣄⠲⣈⢄⡐",
	" ⠂⡄⠊⠌⡐⠠⠑⣈⠐⢢⠡⠘⠠⢈⠐⠠⠁⠌⡱⠌⠦⡑",
	" ⡡⠄⠉⠒⠈⠁⠈⠀⠈⢆⠡⢃⠁⠂⠌⡐⠀⢀⠀⠌⠐⠁",
	" ⠀⠀⠀⠀⠀⠀⠀⠀⢌⠢⡑⠌⡌⠌⡐⠀⠄⠂⠀⠀⡀⠀",
	" ⠀⠀⠀⠀⠀⠀⢰⢁⠢⢡⠘⡐⢌⠢⡁⠌⠀⠀⠀⠄⠀⠀",
}

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

const portraitVisibleWidth = 22
const TaglineText = "is an engineer, builder & creator who turns ideas into products."

func BannerLines() int       { return len(nameBanner) }
func NameBannerLines() []string { return nameBanner }

func RenderHome(r *lipgloss.Renderer, width, height, revealIdx int, starBright []bool, taglineIdx int, taglineDone bool, cursorBlink bool, glitchFrames int, glitchRunes [][]rune, lastCommit string, sessionID string, connectedSecs int, buildInfo string, scanlineY int, idleGlitch bool, theme Theme) string {
	cyanStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	dimStyle        := r.NewStyle().Foreground(lipgloss.Color(theme.Dim))
	whiteStyle      := r.NewStyle().Foreground(lipgloss.Color(theme.Text))
	magentaStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.Secondary))
	starStyle       := r.NewStyle().Foreground(lipgloss.Color(theme.StarDim))
	brightStarStyle := r.NewStyle().Foreground(lipgloss.Color(theme.StarBright))

	// Responsive: show portrait only when terminal wide enough
	showPortrait := width >= 90

	// ── LEFT COLUMN: portrait ─────────────────────────────────────────────
	var leftCol strings.Builder
	portReveal := revealIdx * 2
	if portReveal > len(portraitArt) {
		portReveal = len(portraitArt)
	}

	if showPortrait {
		for i, line := range portraitArt {
			// Trim to portraitVisibleWidth to avoid overflow
			runes := []rune(line)
			if len(runes) > portraitVisibleWidth {
				runes = runes[:portraitVisibleWidth]
			}
			trimmed := string(runes)
			if i < portReveal {
				leftCol.WriteString(cyanStyle.Render(trimmed))
			} else {
				leftCol.WriteString(strings.Repeat(" ", portraitVisibleWidth))
			}
			leftCol.WriteString("\n")
		}
	}

	// ── RIGHT COLUMN: name + bio ───────────────────────────────────────────
	var rightCol strings.Builder

	// Stars — independent twinkle
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

	// Subtitle
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

	// Bio lines — updated with SenseOps
	bioLines := []string{
		dimStyle.Render("Founder & Lead Engineer of"),
		magentaStyle.Render("Webcraft Studios") + dimStyle.Render(","),
		dimStyle.Render("building full-stack apps,"),
		dimStyle.Render("Computer Vision pipelines,"),
		dimStyle.Render("and scalable AI perception systems."),
		"",
		dimStyle.Render("B.Tech ECE @ Manipal Institute of Technology"),
		dimStyle.Render("Bengaluru  ·  Aug 2023 – Jul 2027"),
		"",
		dimStyle.Render("President, ") + cyanStyle.Render("MBOSC"),
		"",
		dimStyle.Render("Software Engineering Intern @ ") + cyanStyle.Render("SenseOps"),
		dimStyle.Render("CV Research Intern @ ") + cyanStyle.Render("ISRO – LEOS"),
	}
	for _, line := range bioLines {
		rightCol.WriteString("  " + line + "\n")
	}

	// ── COMBINE columns ───────────────────────────────────────────────────
	leftContent  := leftCol.String()
	rightContent := rightCol.String()

	leftLines  := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	maxLines := len(rightLines)
	if showPortrait && len(leftLines) > maxLines {
		maxLines = len(leftLines)
	}
	for len(leftLines)  < maxLines { leftLines  = append(leftLines, "") }
	for len(rightLines) < maxLines { rightLines = append(rightLines, "") }

	var combined strings.Builder
	combined.WriteString("\n")
	gap := "  "

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
			// Pad left column to consistent width
			lVis := len([]rune(stripAnsi(left)))
			if lVis < portraitVisibleWidth {
				left += strings.Repeat(" ", portraitVisibleWidth-lVis)
			}
			combined.WriteString(" " + left + gap + right + "\n")
		} else {
			combined.WriteString(" " + right + "\n")
		}
	}

	// Session info line
	if sessionID != "" {
		mins := connectedSecs / 60
		secs := connectedSecs % 60
		metaStyle    := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim))
		sessionStyle := r.NewStyle().Foreground(lipgloss.Color(theme.FooterText))
		combined.WriteString("\n " + sessionStyle.Render(fmt.Sprintf("session: %s  connected: %02d:%02d", sessionID, mins, secs)) + "  " + metaStyle.Render(buildInfo) + "\n")
	}
	if lastCommit != "" {
		commitStyle := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim)).Italic(true)
		combined.WriteString(" " + commitStyle.Render("last pushed: "+lastCommit) + "\n")
	}

	result := combined.String()

	// CRT scanline overlay — use theme color
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
		lines := strings.Split(result, "\n")
		glitchS := r.NewStyle().Foreground(lipgloss.Color(theme.VeryDim))
		for i, line := range lines {
			visual := []rune(stripAnsi(line))
			for j := range visual {
				if visual[j] != ' ' && len(visual) > 0 {
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
