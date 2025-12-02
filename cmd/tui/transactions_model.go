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
	CategoryID    uuid.UUID       `json:"category_id"`
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

	confirmCursor int
}

var accountNames = []string{
	// api to get list of users accounts' names
}

var categoryNames = []string{
	// api to get list of users categories' names
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
				// case "enter":
				// 	m.mode = transactionsModeDetails
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
	}

	return "Unknown transaction mode"
}
