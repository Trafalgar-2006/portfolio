package views

import "github.com/charmbracelet/lipgloss"

// Theme holds all color tokens for the portfolio.
// All render functions receive a *Theme so every screen reacts to theme switches.
type Theme struct {
	Name string

	// Core palette
	Primary   lipgloss.Color // main accent (cyan in Dracula, green in Matrix…)
	Secondary lipgloss.Color // second accent (magenta, etc.)
	Accent    lipgloss.Color // gold/highlight
	Success   lipgloss.Color // green / positive
	Warning   lipgloss.Color // orange / WIP
	Purple    lipgloss.Color // purple / framework tags

	// Text levels
	Text    lipgloss.Color // bright readable text
	DimMid  lipgloss.Color // medium dim (#888)
	Dim     lipgloss.Color // darkest readable (#555)
	VeryDim lipgloss.Color // nearly invisible (#333)

	// UI chrome
	BoxBorder  lipgloss.Color // box drawing lines
	FooterBg   lipgloss.Color // footer bar background
	FooterText lipgloss.Color // footer bar text
	TabActive  lipgloss.Color // active tab bg

	// Matrix-specific
	MatrixHead   lipgloss.Color
	MatrixBright lipgloss.Color
	MatrixMid    lipgloss.Color
	MatrixDim    lipgloss.Color
	MatrixLocked lipgloss.Color // name crystallisation colour

	// Star colours
	StarBright lipgloss.Color
	StarDim    lipgloss.Color

	// Scanline / glitch overlay
	ScanlineColor lipgloss.Color
}

// ──────────────────────────────────────────────────────────────
// Theme presets
// ──────────────────────────────────────────────────────────────

var Themes = []Theme{
	ThemeDracula,
	ThemeMatrix,
	ThemeAmber,
	ThemeNord,
}

// ThemeDracula — the original. Cyan + magenta on deep dark.
var ThemeDracula = Theme{
	Name:      "dracula",
	Primary:   "#00DFDF",
	Secondary: "#FF6AC1",
	Accent:    "#FFD700",
	Success:   "#50FA7B",
	Warning:   "#FFB86C",
	Purple:    "#BD93F9",

	Text:    "#E0E0E0",
	DimMid:  "#888888",
	Dim:     "#555555",
	VeryDim: "#333333",

	BoxBorder:  "#1A4040",
	FooterBg:   "#0A1A1A",
	FooterText: "#1A5050",
	TabActive:  "#00DFDF",

	MatrixHead:   "#FFFFFF",
	MatrixBright: "#00FF41",
	MatrixMid:    "#00AA00",
	MatrixDim:    "#004400",
	MatrixLocked: "#00DFDF",

	StarBright:    "#00DFDF",
	StarDim:       "#444444",
	ScanlineColor: "#0A2A2A",
}

// ThemeMatrix — pure green-on-black hacker. No distractions.
var ThemeMatrix = Theme{
	Name:      "matrix",
	Primary:   "#00FF41",
	Secondary: "#00CC33",
	Accent:    "#00FF41",
	Success:   "#00FF41",
	Warning:   "#88FF44",
	Purple:    "#00BB33",

	Text:    "#00DD33",
	DimMid:  "#007722",
	Dim:     "#004411",
	VeryDim: "#002208",

	BoxBorder:  "#003310",
	FooterBg:   "#000A00",
	FooterText: "#005522",
	TabActive:  "#00FF41",

	MatrixHead:   "#FFFFFF",
	MatrixBright: "#00FF41",
	MatrixMid:    "#00AA00",
	MatrixDim:    "#003300",
	MatrixLocked: "#AAFFAA",

	StarBright:    "#00FF41",
	StarDim:       "#003311",
	ScanlineColor: "#001100",
}

// ThemeAmber — warm retro CRT. Phosphor amber glow on near-black.
var ThemeAmber = Theme{
	Name:      "amber",
	Primary:   "#FFB300",
	Secondary: "#FF8C00",
	Accent:    "#FFD54F",
	Success:   "#CDDC39",
	Warning:   "#FF6F00",
	Purple:    "#FFA040",

	Text:    "#FFE082",
	DimMid:  "#8B6914",
	Dim:     "#4A3800",
	VeryDim: "#2A2000",

	BoxBorder:  "#3A2800",
	FooterBg:   "#0A0800",
	FooterText: "#5A3A00",
	TabActive:  "#FFB300",

	MatrixHead:   "#FFFF88",
	MatrixBright: "#FFB300",
	MatrixMid:    "#996600",
	MatrixDim:    "#442200",
	MatrixLocked: "#FFD700",

	StarBright:    "#FFD700",
	StarDim:       "#4A3000",
	ScanlineColor: "#1A0E00",
}

// ThemeNord — cool arctic palette. Blues and muted cyans.
var ThemeNord = Theme{
	Name:      "nord",
	Primary:   "#88C0D0",
	Secondary: "#B48EAD",
	Accent:    "#EBCB8B",
	Success:   "#A3BE8C",
	Warning:   "#D08770",
	Purple:    "#B48EAD",

	Text:    "#D8DEE9",
	DimMid:  "#6C7A8A",
	Dim:     "#3B4252",
	VeryDim: "#2E3440",

	BoxBorder:  "#2E4060",
	FooterBg:   "#1A2030",
	FooterText: "#3C6080",
	TabActive:  "#88C0D0",

	MatrixHead:   "#ECEFF4",
	MatrixBright: "#88C0D0",
	MatrixMid:    "#5E81AC",
	MatrixDim:    "#2E3440",
	MatrixLocked: "#88C0D0",

	StarBright:    "#88C0D0",
	StarDim:       "#2E3440",
	ScanlineColor: "#1A2030",
}
