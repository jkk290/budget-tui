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

const (
	formFieldName = iota
	formFieldType
	formFieldBalance
	formFieldSave
)

type Account struct {
	name        string
	accountType string
}
type accountsModel struct {
	mode     accountsMode
	accounts []Account
	cursor   int

	formFieldCursor int
	formName        string
	formTypeIndex   int
	formBalance     string
}

var accountTypes = []string{
	"Checking",
	"Savings",
	"Credit Card",
	"Investing",
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
			case "n":
				m.mode = accountsModeFormNew

				m.formFieldCursor = formFieldName
				m.formName = ""
				m.formTypeIndex = 0
				m.formBalance = ""
			}
		case accountsModeDetails:
			switch key {
			case "esc":
				m.mode = accountsModeList
			}
		case accountsModeFormNew, accountsModeFormEdit:
			switch key {
			case "esc":
				m.mode = accountsModeList
				return m, nil
			case "up", "k":
				if m.formFieldCursor > formFieldName {
					m.formFieldCursor--
				}
			case "down", "j":
				if m.formFieldCursor < formFieldSave {
					m.formFieldCursor++
				}
			case "enter":
				switch m.formFieldCursor {
				case formFieldSave:
					newAccount := Account{
						name:        m.formName,
						accountType: accountTypes[m.formTypeIndex],
					}
					m.accounts = append(m.accounts, newAccount)
					m.mode = accountsModeList
				}
			default:
				switch m.formFieldCursor {
				case formFieldName:
					// create form input
				case formFieldBalance:
					// create form input
				case formFieldType:
					// implement dropdown or type selection
				}
			}
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
		s += "\n(Use j/k to move, Enter to view details, n to create new account)\n"
		return s
	case accountsModeDetails:
		acc := m.accounts[m.cursor]
		return fmt.Sprintf(
			"Account Details\n\nName: %s\nType: %s\n\n(Press Esc to go back)\n",
			acc.name,
			acc.accountType,
		)
	case accountsModeFormNew:
		s := "New Account\n\n"

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Name: %s\n", currentRow(formFieldName), m.formName)
		s += fmt.Sprintf("%s Type: %s\n", currentRow(formFieldType), accountTypes[m.formTypeIndex])
		s += fmt.Sprintf("%s Initial Balance: %s\n", currentRow(formFieldBalance), m.formBalance)
		s += "\n"
		s += fmt.Sprintf("%s [ Save ]\n", currentRow(formFieldSave))
		s += "\n(Use j/k to move, Esc to cancel)\n"

		return s
	}

	return "Unknown accounts mode\n"
}
