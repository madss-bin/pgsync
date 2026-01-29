package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorNeonBlue   = lipgloss.Color("#00f2ff")
	colorNeonPink   = lipgloss.Color("#ff00d4")
	colorNeonGreen  = lipgloss.Color("#00ff41")
	colorNeonPurple = lipgloss.Color("#bd00ff")
	colorDarkGray   = lipgloss.Color("#1a1b26")
	colorLightGray  = lipgloss.Color("#a9b1d6")
	TitleStyle      = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorNeonBlue).
			Padding(0, 1).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(colorNeonPink).
			MarginBottom(1)
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(colorNeonPurple).
			Bold(true).
			MarginBottom(1)
	PromptStyle = lipgloss.NewStyle().
			Foreground(colorNeonGreen).
			Bold(true)
	InputStyle = lipgloss.NewStyle().
			Foreground(colorNeonBlue).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorNeonPurple).
			Padding(0, 1)
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(colorDarkGray).
				Background(colorNeonBlue).
				Bold(true).
				Padding(0, 1).
				MarginLeft(1)
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(colorLightGray).
				Padding(0, 1).
				MarginLeft(1)
	ProgressStartColor = "#00f2ff"
	ProgressEndColor   = "#bd00ff"
	PlaceholderStyle   = lipgloss.NewStyle().
				Foreground(colorLightGray)
	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(colorNeonPink).
				Bold(true).
				MarginTop(1)
	HintStyle = lipgloss.NewStyle().
			Foreground(colorLightGray).
			Italic(true)
	ProgressBarStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)
	ProgressTextStyle = lipgloss.NewStyle().
				Foreground(colorNeonBlue).
				MarginBottom(1)
	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorNeonBlue).
			Padding(1, 2)
	HelpStyle = lipgloss.NewStyle().
			Foreground(colorDarkGray).
			Background(colorNeonGreen).
			Padding(0, 1).
			MarginTop(1)
	WarningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorNeonPurple)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorNeonGreen)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorNeonPink)
)
