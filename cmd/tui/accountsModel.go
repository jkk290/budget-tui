package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type accountsMode int

const (
	accountsModeList accountsMode = iota
	accountsModeDetails
	accountsModeFormNew
	accountsModeFormEdit
)

type Account struct {
	name        string
	accountType string
}
type accountsModel struct {
	mode     accountsMode
	accounts []Account
	cursor   int
	// fields for selected account / form inputs
}

func initialAccountModel() accountsModel {
	return accountsModel{
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

func (m accountsModel) Update(msg tea.Msg) (accountsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch m.mode {
		case accountsModeList:
			switch key {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.accounts)-1 {
					m.cursor++
				}
			case "enter":
				m.mode = accountsModeDetails
			}
		case accountsModeDetails:
			switch key {
			case "esc":
				m.mode = accountsModeList
			}
		case accountsModeFormNew, accountsModeFormEdit:
			// form input handling
		}

	}

	return m, nil
}

func (m accountsModel) View() string {
	switch m.mode {
	case accountsModeList:
		s := "Accounts\n\n"
		for i, account := range m.accounts {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, account.name)
		}
		s += "\n(Use j/k to move, Enter to view details)\n"
		return s
	case accountsModeDetails:
		acc := m.accounts[m.cursor]
		return fmt.Sprintf(
			"Account Details\n\nName: %s\nType: %s\n\n(Press Esc to go back)\n",
			acc.name,
			acc.accountType,
		)
	}

	return "Unknown accounts mode\n"
}
