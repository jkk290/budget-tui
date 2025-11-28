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

type section int

const (
	sectionBudget section = iota
	sectionCategories
	sectionAccounts
	sectionTransactions
)

type focus int

const (
	focusNav focus = iota
	focusMain
	focusDetail
)

type model struct {
	navItems  []string
	navCursor int

	currentSection section

	// budgetModel budgetModel
	// categoriesModel categoriesModel
	accountsModel accountsModel
	// transactionsModel transactionsModel

	focus focus
}

func initialModel() model {
	return model{
		navItems:       []string{"Budget", "Categories", "Accounts", "Transactions"},
		navCursor:      0,
		currentSection: sectionBudget,
		accountsModel:  initialAccountModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}

		if key == "tab" {
			if m.focus == focusNav {
				m.focus = focusMain
			} else {
				m.focus = focusNav
			}
			return m, nil
		}

		switch m.focus {
		case focusNav:
			switch msg.String() {
			case "up", "k":
				if m.navCursor > 0 {
					m.navCursor--
				}
			case "down", "j":
				if m.navCursor < len(m.navItems)-1 {
					m.navCursor++
				}
			case "enter":
				m.currentSection = section(m.navCursor)
				m.focus = focusMain
			}
		case focusMain:
			switch m.currentSection {
			case sectionAccounts:
				var cmd tea.Cmd
				m.accountsModel, cmd = m.accountsModel.Update(msg)
				return m, cmd
			case sectionBudget:
			case sectionCategories:
			case sectionTransactions:
			default:
				return m, nil
			}
		case focusDetail:
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	sidebar := m.navView()

	var main string

	switch m.currentSection {
	case sectionBudget:
		main = "Budget View\n"
	case sectionCategories:
		main = "Categories View\n"
	case sectionAccounts:
		main = m.accountsModel.View()
	case sectionTransactions:
		main = "Transactions View\n"
	}

	return sidebar + "\n---\n\n" + main
}

func (m model) navView() string {
	s := "BudgeTUI\n\n"
	s += fmt.Sprintf("\nCurrent View: %v\n\n", m.currentSection)

	for i, navItem := range m.navItems {
		cursor := " "
		if m.navCursor == i {
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
