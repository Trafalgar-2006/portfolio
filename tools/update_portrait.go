//go:build ignore

// update_portrait writes the user-supplied braille art directly into views/home.go.
// Usage: go run tools/update_portrait.go
package main

import (
	"fmt"
	"os"
	"strings"
)

// newPortraitSection is the exact replacement for everything from the
// "// Braille portrait" comment down to (and including) the
// portraitWidth const line.
const newPortraitSection = `// Braille portrait — hand-crafted 68×24 char art, face region.
// Clipped to portraitWidth cols (face content lives in cols ~25–47).
var portraitArt = []string{
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⡀⡀⡀⠐⠐⡀⡀⡀⢀⡀⠠⡀⡀⡀⡂⡀⡀⡀⡀⡀⠠⡀⡀⡀⠐⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⠄⠅⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⠠⡀⡀⠂⡀⠂⠂⡀⡀⠄⠐⡀⠄⡀⢄⠄⡀⡀⠂⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⠂⡀⡀⡀⡀⡀⡀⢀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⠤⡀⡀⡀⠂⡀⡀⠄⠄⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⢀⡀⡀⡀⡀⡀⡀⠠⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⡄⡀⢀⢀⢀⡀⡀⠤⠒⡀⡀⡀⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢠⢞⣡⢠⡀⠈⣴⣯⣛⣉⡻⢀⣩⡥⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⠂⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⡀⡀⡀⡀⢀⡀⣿⣾⣍⣷⣶⣶⣷⣵⣛⣿⣿⡿⢶⣭⢷⢶⣦⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⠠⡀⠐⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⣹⡶⣫⣿⠿⡵⡾⢫⣽⣷⣿⣿⣿⣿⣿⣿⣭⣍⡗⣿⠢⠂⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⣴⣭⠻⣵⣾⣗⣿⣽⣿⣿⣿⣿⣿⣿⣷⣾⣿⣿⣿⣿⣿⡚⠻⢳⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢿⣏⣻⣿⣶⣿⣿⣿⣿⣿⣿⢿⣃⣞⣫⣿⣻⢝⣿⣿⣿⣿⣝⣏⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⠠⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⡀⡀⣻⣿⣿⣿⣿⣷⣠⠭⢛⠒⠿⠿⣿⣿⣿⣿⣿⣴⣿⣯⢿⣟⢌⣹⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢕⣿⣿⣿⣿⡿⠋⢁⣀⣁⠛⣿⣦⡍⠻⢿⢿⢟⣿⣿⣿⣟⣿⣽⣿⠄⡦⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢸⣿⡟⣿⡿⠐⡀⡀⠉⡓⠳⣿⡟⠁⢭⢡⠦⣙⠻⣻⣯⣿⣿⣿⣽⣿⣶⡀⡀⡀⡀⠂⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡞⠋⣈⢻⠁⡀⡀⡀⡀⡀⠈⡀⠁⡀⢸⣞⣗⢍⠙⢲⣍⠯⣿⣿⢽⣩⠁⡄⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢸⠚⡀⡀⡀⢀⡀⡀⢤⢚⠃⡀⡀⡀⠈⠈⢥⣙⣻⡄⣷⢣⣿⣿⣟⡍⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⠄⡀⡀⡀⠒⣾⠂⡀⡀⣩⣇⣴⣤⣠⡀⡀⡀⢸⢣⣿⣿⣿⣧⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⣿⡀⡀⡀⡀⢸⣿⠓⠒⠒⢬⡍⢍⠡⢺⠂⡀⡀⣶⢿⣻⡟⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⣼⡄⡀⡀⡀⡟⠙⣿⢿⣿⣶⣤⠉⡧⣼⠡⢠⢠⣿⠉⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⣿⡀⡀⡀⡀⡀⡀⡀⡀⠈⠛⣿⠿⣻⠁⡀⢀⢋⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢿⡀⡀⡀⡄⡀⡀⠨⢍⣺⠗⠐⡀⢁⠖⡡⠹⠄⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⠰⠂⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢂⣀⡀⡀⡀⠈⠁⣁⠣⣒⠃⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⢀⣀⣀⣴⣶⣷⣿⣿⡀⡀⡀⡀⠈⡆⠰⣨⣒⢲⡖⡞⣊⠿⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⡀⣿⣃⣲⣸⠬⣏⡦⣭⣿⣯⣾⣿⣾⣿⡀⡀⡀⡀⡀⡀⢸⣎⡖⡒⢳⣚⡞⠈⡀⣵⣄⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⡀⡀⡀⡀⢠⣘⣉⣩⣟⣟⣿⣿⣟⣿⣟⡿⠋⡀⠁⡀⡀⡀⡀⡀⢀⢰⠥⠢⣬⠞⠁⡀⡀⠘⣿⣷⣏⡛⠦⣀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀",
	"⡀⡀⡀⡀⠉⡀⠈⠉⠉⠉⠉⠉⡀⡀⡀⠈⠉⠉⡀⡀⡀⡀⡀⡀⡀⡀⡀⡀⠁⠁⠈⠁⡀⡀⡀⡀⡀⡀⠉⠉⠉⠉⠉⡀⡀⠈⠁⡀⡀⡀⡀⡀⡀⡀⡀⡀",
}

// portraitWidth is the display columns rendered — lines are 68 chars, face is cols ~25-47.
const portraitWidth = 56
`

func main() {
	const target = "views/home.go"

	b, err := os.ReadFile(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read:", err)
		os.Exit(1)
	}
	src := string(b)

	// Find the start of the portrait comment block
	startMarker := "// Braille portrait"
	startIdx := strings.Index(src, startMarker)
	if startIdx < 0 {
		fmt.Fprintln(os.Stderr, "could not find portrait comment block")
		os.Exit(1)
	}

	// Find the end: the line containing "const portraitWidth" + newline
	endMarker := "const portraitWidth"
	endIdx := strings.Index(src[startIdx:], endMarker)
	if endIdx < 0 {
		fmt.Fprintln(os.Stderr, "could not find portraitWidth const")
		os.Exit(1)
	}
	endIdx += startIdx

	// Skip to end of that line
	nlIdx := strings.Index(src[endIdx:], "\n")
	if nlIdx < 0 {
		fmt.Fprintln(os.Stderr, "no newline after portraitWidth")
		os.Exit(1)
	}
	endIdx += nlIdx + 1

	newSrc := src[:startIdx] + newPortraitSection + src[endIdx:]

	if err := os.WriteFile(target, []byte(newSrc), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
	fmt.Println("portrait updated ✓")
}
