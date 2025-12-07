package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	tea "github.com/charmbracelet/bubbletea"
)

type transactionsMode int

const (
	transactionsModeList transactionsMode = iota
	transactionsModeDetails
	transactionsModeFormNew
	transactionsModeFormEdit
	transactionsModeDelete
)

const (
	txFormFieldAmount = iota
	txFormFieldDescription
	txFormFieldDate
	txFormFieldPosted
	txFormFieldAccount
	txFormFieldCategory
	txFormFieldSave
)

const (
	txConfirmYes = iota
	txConfirmCancel
)

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	Amount        decimal.Decimal `json:"amount"`
	TxDescription string          `json:"tx_description"`
	TxDate        time.Time       `json:"tx_date"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Posted        bool            `json:"posted"`
	AccountID     uuid.UUID       `json:"account_id"`
	AccountName   string          `json:"account_name"`
	CategoryID    uuid.UUID       `json:"category_id"`
	CategoryName  string          `json:"category_name"`
}

type txAccountOption struct {
	ID   uuid.UUID
	Name string
}

type txCategoryOption struct {
	ID   uuid.UUID
	Name string
}

var postedValues = []bool{
	true,
	false,
}

type transactionsModel struct {
	mode         transactionsMode
	transactions []Transaction
	cursor       int

	formEditing       bool
	formFieldCursor   int
	amountInput       textinput.Model
	descriptionInput  textinput.Model
	dateInput         textinput.Model
	formPostedIndex   int
	formAccountIndex  int
	formCategoryIndex int
	accountOptions    []txAccountOption
	categoryOptions   []txCategoryOption

	confirmCursor int
	errorMsg      string
}

func initialTransactionsModel() transactionsModel {
	txAmount := textinput.New()
	txAmount.CharLimit = 64
	txAmount.Blur()

	txDescription := textinput.New()
	txDescription.CharLimit = 64
	txDescription.Blur()

	txDate := textinput.New()
	txDate.CharLimit = 64
	txDate.Blur()

	return transactionsModel{
		transactions:      []Transaction{},
		cursor:            0,
		formEditing:       false,
		formFieldCursor:   txFormFieldAmount,
		amountInput:       txAmount,
		descriptionInput:  txDescription,
		dateInput:         txDate,
		formPostedIndex:   0,
		formAccountIndex:  0,
		formCategoryIndex: 0,
		accountOptions:    []txAccountOption{},
		categoryOptions:   []txCategoryOption{},
		confirmCursor:     txConfirmCancel,
	}
}

func (m transactionsModel) Update(msg tea.Msg) (transactionsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case transactionCreatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.mode = transactionsModeList
		return m, func() tea.Msg {
			return transactionsReloadRequestedMsg{}
		}

	case transactionUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.mode = transactionsModeList
		return m, func() tea.Msg {
			return transactionsReloadRequestedMsg{}
		}

	case transactionDeletedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}

		filtered := m.transactions[:0]
		for _, transaction := range m.transactions {
			if transaction.ID != msg.transactionID {
				filtered = append(filtered, transaction)
			}
		}
		m.transactions = filtered

		if len(m.transactions) == 1 {
			m.cursor = 0
			m.mode = transactionsModeList
		} else {
			prevCursor := m.cursor
			if prevCursor >= len(m.transactions) {
				m.cursor = len(m.transactions) - 1
			} else {
				m.cursor = prevCursor
			}
			m.mode = transactionsModeList
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		switch m.mode {
		case transactionsModeList:
			switch key {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.transactions)-1 {
					m.cursor++
				}
			case "n":
				return m, func() tea.Msg {
					return transactionsNewRequestedMsg{}
				}
			case "d":
				if len(m.transactions) > 0 {
					m.mode = transactionsModeDelete
					m.confirmCursor = txConfirmCancel
				}
			case "enter":
				m.mode = transactionsModeDetails
			}
		case transactionsModeDetails:
			switch key {
			case "esc":
				m.mode = transactionsModeList
			case "e":
				return m, func() tea.Msg {
					return transactionsEditRequestedMsg{}
				}
			case "d":
				m.mode = transactionsModeDelete
				m.confirmCursor = txConfirmCancel
			}
		case transactionsModeFormNew, transactionsModeFormEdit:
			if m.formEditing {
				if key == "esc" {
					m.formEditing = false
					m.amountInput.Blur()
					m.descriptionInput.Blur()
					m.dateInput.Blur()
					return m, nil
				}

				switch m.formFieldCursor {
				case txFormFieldAmount:
					var cmd tea.Cmd
					m.amountInput, cmd = m.amountInput.Update(msg)
					return m, cmd
				case txFormFieldDescription:
					var cmd tea.Cmd
					m.descriptionInput, cmd = m.descriptionInput.Update(msg)
					return m, cmd
				case txFormFieldDate:
					var cmd tea.Cmd
					m.dateInput, cmd = m.dateInput.Update(msg)
					return m, cmd
				}
			}

			switch key {
			case "esc":
				m.mode = transactionsModeList
				m.errorMsg = ""
				return m, nil
			case "up", "k":
				if m.formFieldCursor > txFormFieldAmount {
					m.formFieldCursor--
				}
			case "down", "j":
				if m.formFieldCursor < txFormFieldSave {
					m.formFieldCursor++
				}
			case "enter":
				switch m.formFieldCursor {
				case txFormFieldAmount, txFormFieldDescription, txFormFieldDate:
					m.formEditing = true
					m.amountInput.Blur()
					m.descriptionInput.Blur()
					m.dateInput.Blur()

					switch m.formFieldCursor {
					case txFormFieldAmount:
						m.amountInput.Focus()
					case txFormFieldDescription:
						m.descriptionInput.Focus()
					case txFormFieldDate:
						m.dateInput.Focus()
					}
				case txFormFieldSave:
					switch m.mode {
					case transactionsModeFormNew:
						amount := m.amountInput.Value()
						description := m.descriptionInput.Value()
						date := m.dateInput.Value()
						account := m.accountOptions[m.formAccountIndex].ID
						category := m.categoryOptions[m.formCategoryIndex].ID
						posted := postedValues[m.formPostedIndex]
						return m, submitCreateTransactionMsg(amount, description, date, account.String(), category.String(), posted)
					case transactionsModeFormEdit:
						amount := m.amountInput.Value()
						description := m.descriptionInput.Value()
						date := m.dateInput.Value()
						account := m.accountOptions[m.formAccountIndex].ID
						category := m.categoryOptions[m.formCategoryIndex].ID
						posted := postedValues[m.formPostedIndex]
						return m, submitUpdateTransactionMsg(m.transactions[m.cursor].ID, amount, description, date, posted, account, category)
					}
				}
			default:
				switch m.formFieldCursor {
				case txFormFieldAmount:
					var cmd tea.Cmd
					m.amountInput, cmd = m.amountInput.Update(msg)
					return m, cmd
				case txFormFieldDescription:
					var cmd tea.Cmd
					m.descriptionInput, cmd = m.descriptionInput.Update(msg)
					return m, cmd
				case txFormFieldDate:
					var cmd tea.Cmd
					m.dateInput, cmd = m.dateInput.Update(msg)
					return m, cmd
				case txFormFieldPosted:
					if key == "left" || key == "h" {
						if m.formPostedIndex > 0 {
							m.formPostedIndex--
						}
					}
					if key == "right" || key == "l" {
						if m.formPostedIndex < len(postedValues)-1 {
							m.formPostedIndex++
						}
					}
				case txFormFieldAccount:
					if key == "left" || key == "h" {
						if m.formAccountIndex > 0 {
							m.formAccountIndex--
						}
					}
					if key == "right" || key == "l" {
						if m.formAccountIndex < len(m.accountOptions)-1 {
							m.formAccountIndex++
						}
					}
				case txFormFieldCategory:
					if key == "left" || key == "h" {
						if m.formCategoryIndex > 0 {
							m.formCategoryIndex--
						}
					}
					if key == "right" || key == "l" {
						if m.formCategoryIndex < len(m.categoryOptions)-1 {
							m.formCategoryIndex++
						}
					}
				}
			}
		case transactionsModeDelete:
			switch key {
			case "esc":
				m.mode = transactionsModeList
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
					return m, submitDeleteTransactionMsg(m.transactions[m.cursor].ID)
				case confirmCancel:
					m.mode = transactionsModeList
				}
			}
		}
	}
	return m, nil
}

func (m transactionsModel) View() string {
	switch m.mode {
	case transactionsModeList:
		s := "Transactions\n\n"
		for i, transaction := range m.transactions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			dateStr := transaction.TxDate.Format("2006-01-02")
			s += fmt.Sprintf("%s %s | %s | %s\n", cursor, dateStr, transaction.TxDescription, transaction.Amount)
		}
		s += "\n(Use 'j'/'k' to move, 'enter' to view details, 'n' to create a new transaction, 'd' to delete transaction)\n"
		return s
	case transactionsModeDetails:
		tx := m.transactions[m.cursor]
		s := "Transaction Details\n\n"
		s += fmt.Sprintf("Amount: %s\n", tx.Amount.String())
		s += fmt.Sprintf("Description: %s\n", tx.TxDescription)
		s += fmt.Sprintf("Date: %s\n", tx.TxDate.Format("2006-01-02"))
		s += fmt.Sprintf("Posted: %v\n", tx.Posted)
		s += fmt.Sprintf("Account: %s\n", tx.AccountName)
		s += fmt.Sprintf("Category: %s\n\n", tx.CategoryName)

		s += "(Press 'esc' to go back, 'e' to edit, 'd' to delete)\n"

		return s
	case transactionsModeFormNew, transactionsModeFormEdit:
		s := "New Transaction\n\n"
		if m.mode == transactionsModeFormEdit {
			s = "Edit Transaction\n\n"
		}

		s += m.errorView()

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Amount: %s\n", currentRow(txFormFieldAmount), m.amountInput.View())
		s += fmt.Sprintf("%s Description: %s\n", currentRow(txFormFieldDescription), m.descriptionInput.View())
		s += fmt.Sprintf("%s Date(YYYY-MM-DD): %s\n", currentRow(txFormFieldDate), m.dateInput.View())
		s += fmt.Sprintf("%s Posted ('h'/'l' to change): %v\n", currentRow(txFormFieldPosted), postedValues[m.formPostedIndex])
		s += fmt.Sprintf("%s Account ('h'/'l' to change): %s\n", currentRow(txFormFieldAccount), m.accountOptions[m.formAccountIndex].Name)
		s += fmt.Sprintf("%s Category ('h'/'l' to change): %s\n\n", currentRow(txFormFieldCategory), m.categoryOptions[m.formCategoryIndex].Name)

		s += fmt.Sprintf("%s [ Save ]\n", currentRow(txFormFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"

		return s
	case transactionsModeDelete:
		s := "Delete Transaction\n\n"
		s += fmt.Sprintf("Are you sure you want to delete transaction '%s'?\n", m.transactions[m.cursor].TxDescription)

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

	return "Unknown transaction mode"
}

func (m transactionsModel) errorView() string {
	if m.errorMsg == "" {
		return ""
	}
	return fmt.Sprintf("Error: %s\n\n", m.errorMsg)
}

func (m transactionsModel) IsEditing() bool {
	return m.formEditing
}
