package main

import "github.com/charmbracelet/lipgloss"

var cursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#daba38")).Foreground(lipgloss.Color("#000000")).Underline(true).Bold(true)

func formatCursor(s string) string {
	return cursorStyle.Render(s)
}

var extraRuneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#733235"))

func formatExtraRune(r rune) string {
	return extraRuneStyle.Render(string(r))
}

var notTypedRuneStyle = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("9"))

func formatNotTyped(r interface{}) string {
	var s string
	switch v := r.(type) {
	case string:
		s = v
	case rune:
		s = string(v)
	case byte:
		s = string(v)
	}
	return notTypedRuneStyle.Render(s)
}

var correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

func formatCorrect(r interface{}) string {
	var s string
	switch v := r.(type) {
	case string:
		s = v
	case rune:
		s = string(v)
	case byte:
		s = string(v)
	}

	return correctStyle.Render(s)
}

var incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func formatIncorrect(r interface{}) string {
	var s string
	switch v := r.(type) {
	case string:
		s = v
	case rune:
		s = string(v)
	case byte:
		s = string(v)
	}

	return incorrectStyle.Render(s)
}
