package main

import "github.com/charmbracelet/lipgloss"

var (
	colorBg        = lipgloss.Color("#0f0f10")
	colorSidebarBg = lipgloss.Color("#1b1b1d")
	colorAccent    = lipgloss.Color("#8be9fd")
	colorText      = lipgloss.Color("#f8f8f2")
	colorMuted     = lipgloss.Color("#666875")
	colorDanger    = lipgloss.Color("#ff5555")

	sidebarWidth = 22
)

var (
	appBGStyle = lipgloss.NewStyle().
			Background(colorBg)

	navBaseStyle = baseContainerStyle().
			Width(sidebarWidth).
			Background(colorSidebarBg)

	navFocusedStyle = withFocus(navBaseStyle)

	loginStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorText).
			Padding(1, 2)

	loginTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	sidebarTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent)

	sidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorMuted)

	sidebarItemActiveStyle = sidebarItemStyle.
				Foreground(colorAccent).
				Bold(true)

	mainBaseStyle = baseContainerStyle()

	mainFocusedStyle = withFocus(mainBaseStyle)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger)
)

func baseContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(colorBg).
		Foreground(colorText).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorMuted)
}

func withFocus(style lipgloss.Style) lipgloss.Style {
	return style.BorderForeground(colorAccent)
}
