package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Account struct {
	name        string
	accountType string
}
type accountViewModel struct {
	accounts        []Account
	cursor          int
	selectedAccount string
}

func initialAccountModel() accountViewModel {
	return accountViewModel{
		accounts: []Account{
			{
				name:        "Account 1",
				accountType: "Checking",
			},
			{
				name:        "Account 2",
				accountType: "Credit Card",
			},
		},
	}
}

func (m accountViewModel) Init() tea.Cmd {
	return nil
}

func (m accountViewModel) Update(msg tea.Msg) (accountViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.accounts)-1 {
				m.cursor++
			}
		case "enter":
			m.selectedAccount = m.accounts[m.cursor].name
		}
	}
	return m, nil
}

func (m accountViewModel) View() string {
	s := fmt.Sprintf("\nSelected Account: %s\n\n", m.selectedAccount)

	for i, account := range m.accounts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, account.name)
	}

	return s
}
