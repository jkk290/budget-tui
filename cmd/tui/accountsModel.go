package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
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

	formEditing     bool
	formFieldCursor int
	nameInput       textinput.Model
	balanceInput    textinput.Model
	formTypeIndex   int
}

var accountTypes = []string{
	"Checking",
	"Savings",
	"Credit Card",
	"Investing",
}

func initialAccountModel() accountsModel {
	name := textinput.New()
	name.Placeholder = "Account Name"
	name.CharLimit = 64
	name.Focus()

	balance := textinput.New()
	balance.Placeholder = "0.00"
	balance.CharLimit = 16

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
		nameInput:     name,
		balanceInput:  balance,
		formTypeIndex: 0,
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
				m.formEditing = false

				m.nameInput.SetValue("")
				m.nameInput.Blur()

				m.balanceInput.SetValue("")
				m.balanceInput.Blur()

				m.formTypeIndex = 0
			}
		case accountsModeDetails:
			switch key {
			case "esc":
				m.mode = accountsModeList
			}
		case accountsModeFormNew, accountsModeFormEdit:
			if m.formEditing {
				if key == "esc" {
					m.formEditing = false
					m.nameInput.Blur()
					m.balanceInput.Blur()
					return m, nil
				}

				switch m.formFieldCursor {
				case formFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				case formFieldBalance:
					var cmd tea.Cmd
					m.balanceInput, cmd = m.balanceInput.Update(msg)
					return m, cmd
				}
			}

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
				case formFieldName, formFieldBalance:
					m.formEditing = true
					m.nameInput.Blur()
					m.balanceInput.Blur()

					switch m.formFieldCursor {
					case formFieldName:
						m.nameInput.Focus()
					case formFieldBalance:
						m.balanceInput.Focus()
					}
				case formFieldSave:
					newAccount := Account{
						name:        m.nameInput.Value(),
						accountType: accountTypes[m.formTypeIndex],
					}
					m.accounts = append(m.accounts, newAccount)
					m.mode = accountsModeList
				}
			default:
				switch m.formFieldCursor {
				case formFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				case formFieldBalance:
					var cmd tea.Cmd
					m.balanceInput, cmd = m.balanceInput.Update(msg)
					return m, cmd
				case formFieldType:
					if key == "left" || key == "h" {
						if m.formTypeIndex > 0 {
							m.formTypeIndex--
						}
					}
					if key == "right" || key == "l" {
						if m.formTypeIndex < len(accountTypes)-1 {
							m.formTypeIndex++
						}
					}
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

		s += fmt.Sprintf("%s Name: %s\n", currentRow(formFieldName), m.nameInput.View())
		s += fmt.Sprintf("%s Type (h/l to change): %s\n", currentRow(formFieldType), accountTypes[m.formTypeIndex])
		s += fmt.Sprintf("%s Initial Balance: %s\n", currentRow(formFieldBalance), m.balanceInput.View())
		s += "\n"
		s += fmt.Sprintf("%s [ Save ]\n", currentRow(formFieldSave))
		s += "\n(Use j/k to move, Enter to edit field, Esc to stop editing, Esc again to cancel)\n"

		return s
	}

	return "Unknown accounts mode\n"
}
