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
	appStyle = lipgloss.NewStyle().Background(colorBg)

	appBGStyle = lipgloss.NewStyle().
			Background(colorBg)

	loginStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorText).
			Padding(1, 2)

	loginTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	sidebarStyle = lipgloss.NewStyle().
			Background(colorSidebarBg).
			Foreground(colorText).
			Padding(1, 2).
			Width(sidebarWidth)

	sidebarTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent)

	sidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorMuted)

	sidebarItemActiveStyle = sidebarItemStyle.
				Foreground(colorAccent).
				Bold(true)

	mainStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(colorText)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger)
)
