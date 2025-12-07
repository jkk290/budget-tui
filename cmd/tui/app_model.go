package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type screen int

const (
	screenLogin screen = iota
	screenMain
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
	jwt           string
	screen        screen
	loginErr      string
	loginUsername textinput.Model
	loginPassword textinput.Model

	client          *Client
	accountsAPI     AccountsAPI
	transactionsAPI TransactionsAPI
	categoriesAPI   CategoriesAPI

	navItems  []string
	navCursor int

	currentSection section

	// budgetModel budgetModel
	categoriesModel   categoriesModel
	accountsModel     accountsModel
	transactionsModel transactionsModel

	focus focus
}

func initialModel(client *Client) model {
	username := textinput.New()
	username.Focus()
	username.Prompt = "Username: "

	password := textinput.New()
	password.EchoMode = textinput.EchoPassword
	password.Prompt = "Password: "

	return model{
		client:        client,
		screen:        screenLogin,
		jwt:           "",
		loginUsername: username,
		loginPassword: password,

		navItems:          []string{"Budget", "Categories", "Accounts", "Transactions"},
		navCursor:         0,
		currentSection:    sectionBudget,
		accountsModel:     initialAccountModel(),
		accountsAPI:       client.Accounts(),
		transactionsModel: initialTransactionsModel(),
		transactionsAPI:   client.Transactions(),
		categoriesModel:   initialCategoriesModel(),
		categoriesAPI:     client.Categories(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenLogin:
		return m.updateLogin(msg)
	case screenMain:
		return m.updateMain(msg)
	default:
		return m, nil
	}

}

func (m model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Accounts
	case accountsReloadRequestedMsg:
		cmd := loadAccountsCmd(m.accountsAPI)
		return m, cmd

	case accountsLoadedMsg:
		if msg.err != nil {
			var cmd tea.Cmd
			m.accountsModel, cmd = m.accountsModel.Update(accountsLoadedMsg{
				err: msg.err,
			})
			return m, cmd
		}

		m.accountsModel.accounts = msg.accounts
		m.accountsModel.cursor = 0
		return m, nil

	case loadAccountTxsSubmittedMsg:
		return m, loadAccountTxsCmd(m.accountsAPI, msg.accountID)

	case loadAccountTxsMsg:
		var cmd tea.Cmd
		m.accountsModel, cmd = m.accountsModel.Update(msg)
		return m, cmd

	case accountCreateSubmittedMsg:
		balanceDecimal, err := decimal.NewFromString(msg.BalanceText)
		if err != nil {
			var cmd tea.Cmd
			m.accountsModel, cmd = m.accountsModel.Update(accountCreatedMsg{
				err: fmt.Errorf("invalid balance: %w", err),
			})
			return m, cmd
		}

		req := CreateAccountRequest{
			AccountName:    msg.Name,
			AccountType:    msg.Type,
			InitialBalance: balanceDecimal,
		}

		return m, createAccountCmd(m.accountsAPI, req)

	case accountCreatedMsg:
		var cmd tea.Cmd
		m.accountsModel, cmd = m.accountsModel.Update(msg)
		return m, cmd

	case accountUpdateSubmittedMsg:
		req := UpdateAccountRequest{
			AccountName: msg.Name,
		}

		return m, updateAccountCmd(m.accountsAPI, msg.AccountID, req)

	case accountUpdatedMsg:
		var cmd tea.Cmd
		m.accountsModel, cmd = m.accountsModel.Update(msg)
		return m, cmd

	case accountDeleteSubmittedMsg:
		return m, deleteAccountCmd(m.accountsAPI, msg.AccountID)

	case accountDeletedMsg:
		var cmd tea.Cmd
		m.accountsModel, cmd = m.accountsModel.Update(msg)
		return m, cmd

	// Categories
	case categoriesLoadedMsg:
		if msg.err != nil {
			var cmd tea.Cmd
			m.categoriesModel, cmd = m.categoriesModel.Update(categoriesLoadedMsg{
				err: msg.err,
			})
			return m, cmd
		}

		m.categoriesModel.categories = msg.categories
		return m, nil

	// Transactions
	case transactionsReloadRequestedMsg:
		cmd := loadTransactionsCmd(m.transactionsAPI)
		return m, cmd

	case transactionsLoadedMsg:
		if msg.err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionsLoadedMsg{
				err: msg.err,
			})
			return m, cmd
		}

		m.transactionsModel.transactions = msg.transactions
		m.transactionsModel.cursor = 0
		return m, nil

	case transactionCreateSubmittedMsg:
		amountDecimal, err := decimal.NewFromString(msg.AmountText)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionCreatedMsg{
				err: fmt.Errorf("invalid amount: %w", err),
			})
			return m, cmd
		}
		txDateTime, err := time.Parse(time.DateOnly, msg.TxDate)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionCreatedMsg{
				err: fmt.Errorf("invalid date: %w", err),
			})
			return m, cmd
		}
		accountID, err := uuid.Parse(msg.AccountID)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionCreatedMsg{
				err: fmt.Errorf("invalid account ID: %w", err),
			})
			return m, cmd
		}
		categoryID, err := uuid.Parse(msg.CategoryID)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionCreatedMsg{
				err: fmt.Errorf("invalid category ID: %w", err),
			})
			return m, cmd
		}

		req := CreateTransactionRequest{
			Amount:        amountDecimal,
			TxDescription: msg.TxDescription,
			TxDate:        txDateTime,
			Posted:        msg.Posted,
			AccountID:     accountID,
			CategoryID:    categoryID,
		}
		return m, createTransactionCmd(m.transactionsAPI, req)

	case transactionCreatedMsg:
		var cmd tea.Cmd
		m.transactionsModel, cmd = m.transactionsModel.Update(msg)
		return m, cmd

	case transactionUpdateSubmittedMsg:
		amountDecimal, err := decimal.NewFromString(msg.Amount)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionUpdatedMsg{
				err: fmt.Errorf("invalid amount: %w", err),
			})
			return m, cmd
		}
		txDateTime, err := time.Parse(time.DateOnly, msg.Date)
		if err != nil {
			var cmd tea.Cmd
			m.transactionsModel, cmd = m.transactionsModel.Update(transactionUpdatedMsg{
				err: fmt.Errorf("invalid date: %w", err),
			})
			return m, cmd
		}

		req := UpdateTransactionRequest{
			Amount:        amountDecimal,
			TxDescription: msg.Description,
			TxDate:        txDateTime,
			Posted:        msg.Posted,
			AccountID:     msg.AccountID,
			CategoryID:    msg.CategoryID,
		}
		return m, updateTransactionCmd(m.transactionsAPI, msg.TransactionID, req)

	case transactionUpdatedMsg:
		var cmd tea.Cmd
		m.transactionsModel, cmd = m.transactionsModel.Update(msg)
		return m, cmd

	case transactionDeleteSubmittedMsg:
		return m, deleteTransactionCmd(m.transactionsAPI, msg.transactionID)

	case transactionDeletedMsg:
		var cmd tea.Cmd
		m.transactionsModel, cmd = m.transactionsModel.Update(msg)
		return m, cmd

	case transactionsNewRequestedMsg:
		accountOptions := make([]txAccountOption, len(m.accountsModel.accounts))
		for i, account := range m.accountsModel.accounts {
			accountOptions[i] = txAccountOption{
				ID:   account.ID,
				Name: account.AccountName,
			}
		}

		categoryOptions := make([]txCategoryOption, len(m.categoriesModel.categories))
		for i, category := range m.categoriesModel.categories {
			categoryOptions[i] = txCategoryOption{
				ID:   category.ID,
				Name: category.CategoryName,
			}
		}

		tm := m.transactionsModel
		tm.accountOptions = accountOptions
		tm.categoryOptions = categoryOptions
		tm.mode = transactionsModeFormNew
		tm.formFieldCursor = txFormFieldAmount
		tm.formEditing = false
		tm.amountInput.SetValue("")
		tm.amountInput.Blur()
		tm.descriptionInput.SetValue("")
		tm.descriptionInput.Blur()
		tm.dateInput.SetValue("")
		tm.dateInput.Blur()
		tm.formPostedIndex = 0
		tm.formAccountIndex = 0
		tm.formCategoryIndex = 0
		tm.errorMsg = ""
		m.transactionsModel = tm
		return m, nil

	case transactionsEditRequestedMsg:
		accountOptions := make([]txAccountOption, len(m.accountsModel.accounts))
		for i, account := range m.accountsModel.accounts {
			accountOptions[i] = txAccountOption{
				ID:   account.ID,
				Name: account.AccountName,
			}
		}

		categoryOptions := make([]txCategoryOption, len(m.categoriesModel.categories))
		for i, category := range m.categoriesModel.categories {
			categoryOptions[i] = txCategoryOption{
				ID:   category.ID,
				Name: category.CategoryName,
			}
		}

		tm := m.transactionsModel
		currentTx := tm.transactions[tm.cursor]
		tm.accountOptions = accountOptions
		tm.categoryOptions = categoryOptions
		tm.mode = transactionsModeFormEdit
		tm.formFieldCursor = txFormFieldAmount
		tm.formEditing = false
		tm.amountInput.SetValue(currentTx.Amount.String())
		tm.amountInput.Blur()
		tm.descriptionInput.SetValue(currentTx.TxDescription)
		tm.descriptionInput.Blur()
		tm.dateInput.SetValue(currentTx.TxDate.Format("2006-01-02"))
		tm.dateInput.Blur()

		if currentTx.Posted {
			tm.formPostedIndex = 0
		} else {
			tm.formPostedIndex = 1
		}

		tm.formAccountIndex = 0
		for i, option := range accountOptions {
			if option.ID == currentTx.AccountID {
				tm.formAccountIndex = i
				break
			}
		}
		tm.formCategoryIndex = 0
		for i, option := range categoryOptions {
			if option.ID == currentTx.CategoryID {
				tm.formCategoryIndex = i
				break
			}
		}
		tm.errorMsg = ""
		m.transactionsModel = tm
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" {
			return m, tea.Quit
		}

		if key == "q" {
			isEditing := false
			switch m.currentSection {
			case sectionAccounts:
				isEditing = m.accountsModel.IsEditing()
			case sectionTransactions:
				isEditing = m.transactionsModel.IsEditing()
				// case sectionCategories:
				// 	isEditing = m.categoriesModel.isEditing()
			}
			if !isEditing {
				return m, tea.Quit
			}
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
				var cmd tea.Cmd
				m.transactionsModel, cmd = m.transactionsModel.Update(msg)
				return m, cmd
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
	switch m.screen {
	case screenLogin:
		return m.loginView()
	case screenMain:
		return m.mainView()
	default:
		return ""
	}
}

func (m model) mainView() string {
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
		main = m.transactionsModel.View()
	}

	return sidebar + "\n---\n\n" + main
}

func (m model) navView() string {
	s := "BudgeTUI\n\n"

	for i, navItem := range m.navItems {
		cursor := " "
		if m.navCursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, navItem)
	}

	s += "\n(Press 'q' to quit, 'tab' to switch between nav bar and main content)\n"

	return s
}
