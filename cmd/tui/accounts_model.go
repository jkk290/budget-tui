package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type accountsMode int

const (
	accountsModeList accountsMode = iota
	accountsModeDetails
	accountsModeFormNew
	accountsModeFormEdit
	accountsModeDelete
)

const (
	formFieldName = iota
	formFieldType
	formFieldBalance
	formFieldSave
)

const (
	confirmYes = iota
	confirmCancel
)

type Account struct {
	ID             uuid.UUID       `json:"id"`
	AccountName    string          `json:"account_name"`
	AccountType    string          `json:"account_type"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	UserID         uuid.UUID       `json:"user_id"`
	AccountBalance decimal.Decimal `json:"account_balance"`
}

type accountsModel struct {
	mode       accountsMode
	accounts   []Account
	accountTxs []Transaction
	cursor     int

	formEditing     bool
	formFieldCursor int
	nameInput       textinput.Model
	balanceInput    textinput.Model
	formTypeIndex   int

	confirmCursor int

	errorMsg string
}

var accountTypes = []string{
	"Checking",
	"Savings",
	"Credit Card",
	"Investing",
}

func initialAccountModel() accountsModel {
	name := textinput.New()
	name.CharLimit = 64
	name.Blur()

	balance := textinput.New()
	balance.CharLimit = 16
	balance.Blur()

	return accountsModel{
		accounts:        []Account{},
		cursor:          0,
		formEditing:     false,
		formFieldCursor: formFieldName,
		nameInput:       name,
		balanceInput:    balance,
		formTypeIndex:   0,
		confirmCursor:   confirmCancel,
		errorMsg:        "",
	}
}

func (m accountsModel) Update(msg tea.Msg) (accountsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case accountCreatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.mode = accountsModeList
		return m, func() tea.Msg {
			return accountsReloadRequestedMsg{}
		}

	case loadAccountTxsMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.accountTxs = msg.accountTxs
		m.mode = accountsModeDetails
		return m, nil

	case accountUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.mode = accountsModeList
		return m, func() tea.Msg {
			return accountsReloadRequestedMsg{}
		}

	case accountDeletedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}

		filtered := m.accounts[:0]
		for _, account := range m.accounts {
			if account.ID != msg.accountID {
				filtered = append(filtered, account)
			}
		}
		m.accounts = filtered

		if len(m.accounts) == 1 {
			m.cursor = 0
			m.mode = accountsModeList
		} else {
			prevCursor := m.cursor
			if prevCursor >= len(m.accounts) {
				m.cursor = len(m.accounts) - 1
			} else {
				m.cursor = prevCursor
			}
			m.mode = accountsModeList
		}
		return m, nil

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
				// m.mode = accountsModeDetails
				return m, submitLoadAccountTxsMsg(m.accounts[m.cursor].ID)

			case "n":
				m.mode = accountsModeFormNew

				m.formFieldCursor = formFieldName
				m.formEditing = false

				m.nameInput.SetValue("")
				m.nameInput.Blur()

				m.balanceInput.SetValue("")
				m.balanceInput.Blur()

				m.formTypeIndex = 0
				m.errorMsg = ""
			case "d":
				if len(m.accounts) > 0 {
					m.mode = accountsModeDelete
					m.confirmCursor = confirmCancel
				}
			}
		case accountsModeDetails:
			switch key {
			case "esc":
				m.mode = accountsModeList
			case "d":
				m.mode = accountsModeDelete
				m.confirmCursor = confirmCancel
			case "e":
				m.mode = accountsModeFormEdit
				m.formFieldCursor = formFieldName
				m.formEditing = false

				m.nameInput.SetValue(m.accounts[m.cursor].AccountName)
				m.nameInput.Blur()
				m.formTypeIndex = slices.Index(accountTypes, (m.accounts[m.cursor].AccountType))
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
				m.errorMsg = ""
				return m, nil
			case "up", "k":
				if m.formFieldCursor > formFieldName {
					m.formFieldCursor--
				}
				if m.mode == accountsModeFormEdit && m.formFieldCursor == formFieldBalance {
					m.formFieldCursor = formFieldType
				}
			case "down", "j":
				if m.formFieldCursor < formFieldSave {
					m.formFieldCursor++
				}
				if m.mode == accountsModeFormEdit && m.formFieldCursor == formFieldBalance {
					m.formFieldCursor = formFieldSave
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
						if m.mode == accountsModeFormNew {
							m.balanceInput.Focus()
						}
					}
				case formFieldSave:
					switch m.mode {
					case accountsModeFormNew:
						name := m.nameInput.Value()
						accountType := accountTypes[m.formTypeIndex]
						balance := m.balanceInput.Value()
						return m, submitCreateAccountMsg(name, accountType, balance)
					case accountsModeFormEdit:
						id := m.accounts[m.cursor].ID
						name := m.nameInput.Value()
						return m, submitUpdateAccountMsg(id, name)
					}
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
					if m.mode == accountsModeFormNew {
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
		case accountsModeDelete:
			switch key {
			case "esc":
				m.mode = accountsModeList
			case "up", "k":
				if m.confirmCursor > confirmYes {
					m.confirmCursor--
				}
			case "down", "j":
				if m.confirmCursor < confirmCancel {
					m.confirmCursor++
				}
			case "enter":
				switch m.confirmCursor {
				case confirmYes:
					return m, submitDeleteAccountMsg(m.accounts[m.cursor].ID)
				case confirmCancel:
					m.mode = accountsModeList
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

		s += m.errorView()

		for i, account := range m.accounts {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n  - $%s\n\n", cursor, account.AccountName, account.AccountBalance.String())
		}
		s += "\n(Use 'j'/'k' to move, 'enter' to view details, 'n' to create new account, 'd' to delete account)\n"
		return s
	case accountsModeDetails:
		acc := m.accounts[m.cursor]
		s := "Account Details\n\n"
		s += m.errorView()
		s += fmt.Sprintf("Name: %s\nType: %s\nBalance: $%s\n\n", acc.AccountName,
			acc.AccountType,
			acc.AccountBalance.String())

		s += "Transactions\n\n"
		for _, transaction := range m.accountTxs {

			dateStr := transaction.TxDate.Format("2006-01-02")
			s += fmt.Sprintf("%s | %s | %s\n", dateStr, transaction.TxDescription, transaction.Amount)
		}

		s += "\n(Press 'esc' to go back, 'e' to edit, 'd' to delete)\n"

		return s

	case accountsModeFormNew:
		s := "New Account\n\n"

		s += m.errorView()

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Name: %s\n", currentRow(formFieldName), m.nameInput.View())
		s += fmt.Sprintf("%s Type ('h'/'l' to change): %s\n", currentRow(formFieldType), accountTypes[m.formTypeIndex])
		s += fmt.Sprintf("%s Initial Balance: %s\n", currentRow(formFieldBalance), m.balanceInput.View())
		s += "\n"
		s += fmt.Sprintf("%s [ Save ]\n", currentRow(formFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"

		return s
	case accountsModeFormEdit:
		s := "Edit Account\n\n"

		s += m.errorView()

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Name: %s\n", currentRow(formFieldName), m.nameInput.View())
		s += fmt.Sprintf("%s Type: %s\n", currentRow(formFieldType), accountTypes[m.formTypeIndex])
		s += "\n"
		s += fmt.Sprintf("%s [ Save ]\n", currentRow(formFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"

		return s
	case accountsModeDelete:
		s := "Delete Account\n\n"
		s += fmt.Sprintf("Are you sure you want to delete Account '%s'?\n", m.accounts[m.cursor].AccountName)
		s += "This will also delete all of the account's transactions.\n\n"

		currentRow := func(field int) string {
			if m.confirmCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s [ Yes ]\n", currentRow(confirmYes))
		s += fmt.Sprintf("%s [ Cancel ]\n", currentRow(confirmCancel))

		s += "\n(Use 'j'/'k' to move, 'enter' to select, 'esc' to cancel)"

		return s
	}

	return "Unknown accounts mode\n"
}

func (m accountsModel) errorView() string {
	if m.errorMsg == "" {
		return ""
	}
	return fmt.Sprintf("Error: %s\n\n", m.errorMsg)
}

func (m accountsModel) IsEditing() bool {
	return m.formEditing
}
