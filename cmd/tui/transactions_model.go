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
				m.mode = transactionsModeFormNew
				return m, func() tea.Msg {
					return transactionsNewRequestedMsg{}
				}
			case "enter":
				m.mode = transactionsModeDetails
			}
		case transactionsModeDetails:
			switch key {
			case "esc":
				m.mode = transactionsModeList
			case "e":
			case "d":
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

			s += fmt.Sprintf("%s %v | %s | %s\n", cursor, transaction.TxDate, transaction.TxDescription, transaction.Amount)
		}
		s += "\n(Use 'j'/'k' to move, 'enter' to view details, 'n' to create a new transaction, 'd' to delete transaction)\n"
		return s
	case transactionsModeDetails:
		tx := m.transactions[m.cursor]
		s := "Transaction Details\n\n"
		s += fmt.Sprintf("Amount: %s\n", tx.Amount.String())
		s += fmt.Sprintf("Description: %s\n", tx.TxDescription)
		s += fmt.Sprintf("Date: %s\n", tx.TxDate.String())
		s += fmt.Sprintf("Posted: %v\n", tx.Posted)
		s += fmt.Sprintf("Account: %s\n", tx.AccountName)
		s += fmt.Sprintf("Category: %s\n\n", tx.CategoryName)

		s += "(Press 'esc' to go back, 'e' to edit, 'd' to delete)\n"

		return s
	case transactionsModeFormNew:
		s := "New Category\n\n"
	
		s += m.errorView()
	
		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}
	
		s += fmt.Sprintf("%s Amount: %s\n", currentRow(txFormFieldAmount), m.amountInput.View())
		s += fmt.Sprintf("%s Description: %s\n", currentRow(txFormFieldDescription), m.descriptionInput.View())
		
		s += fmt.Sprintf("%s Account ('h'/'l' to change): %s\n", currentRow(txFormFieldAccount), m.accountOptions[m.formAccountIndex])
		s += fmt.Sprintf("%s Category ('h'/'l' to change): %s\n", currentRow(txFormFieldCategory), m.categoryOptions[m.formCategoryIndex])
		s += fmt.Sprintf("%s Initial Balance: %s\n", currentRow(formFieldBalance), m.balanceInput.View())
		s += "\n"
		s += fmt.Sprintf("%s [ Save ]\n", currentRow(formFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"
	
		return s
	}

	return "Unknown transaction mode"
}
