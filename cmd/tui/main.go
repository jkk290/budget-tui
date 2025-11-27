package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type navigationItem int

const (
	navBudget navigationItem = iota
	navCategories
	navAccounts
	navTransactions
)

type focusLevel int

const (
	focusNavbar focusLevel = iota
	focusList
)

type model struct {
	navItems     []string
	cursor       int
	currentView  navigationItem
	accountsView accountViewModel
	focus        focusLevel
}

func initialModel() model {
	return model{
		navItems:     []string{"Budget", "Categories", "Accounts", "Transactions"},
		cursor:       0,
		currentView:  navBudget,
		accountsView: initialAccountModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.navItems)-1 {
				m.cursor++
			}
		case "enter":
			m.currentView = navigationItem(m.cursor)
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "BudgeTUI\n\n"
	s += fmt.Sprintf("\nCurrent View: %v\n\n", m.currentView)

	for i, navItem := range m.navItems {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, navItem)
	}

	s += "\nPress 'q' to quit\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
