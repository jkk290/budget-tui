package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	navGroups
	navAccounts
	navTransactions
)

type section int

const (
	sectionBudget section = iota
	sectionCategories
	sectionGroups
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
	groupsAPI       GroupsAPI

	navItems  []string
	navCursor int

	currentSection section

	budgetModel       budgetModel
	budgetAPI         BudgetAPI
	categoriesModel   categoriesModel
	groupsModel       groupsModel
	accountsModel     accountsModel
	transactionsModel transactionsModel

	focus  focus
	width  int
	height int
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

		navItems:          []string{"Budget", "Categories", "Category Groups", "Accounts", "Transactions"},
		navCursor:         0,
		currentSection:    sectionBudget,
		budgetModel:       initialBudgetModel(),
		budgetAPI:         client.Budget(),
		accountsModel:     initialAccountModel(),
		accountsAPI:       client.Accounts(),
		transactionsModel: initialTransactionsModel(),
		transactionsAPI:   client.Transactions(),
		categoriesModel:   initialCategoriesModel(),
		categoriesAPI:     client.Categories(),
		groupsModel:       initialGroupsModel(),
		groupsAPI:         client.Groups(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
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

	// Budget
	case budgetReloadRequestedMsg:
		cmd := loadBudgetCmd(m.budgetAPI)
		return m, cmd

	case budgetLoadedMsg:
		var cmd tea.Cmd
		m.budgetModel, cmd = m.budgetModel.Update(msg)
		return m, cmd

	// Categories
	case categoriesReloadRequestedMsg:
		cmd := loadCategoriesCmd(m.categoriesAPI)
		return m, cmd

	case categoriesLoadedMsg:
		if msg.err != nil {
			var cmd tea.Cmd
			m.categoriesModel, cmd = m.categoriesModel.Update(categoriesLoadedMsg{
				err: msg.err,
			})
			return m, cmd
		}

		m.categoriesModel.categories = msg.categories
		m.categoriesModel.cursor = 0
		return m, nil

	case loadCategoryTxsSubmittedMsg:
		return m, loadCategoryTxsCmd(m.categoriesAPI, msg.categoryID)

	case loadCategoryTxsMsg:
		var cmd tea.Cmd
		m.categoriesModel, cmd = m.categoriesModel.Update(msg)
		return m, cmd

	case categoryCreateSubmittedMsg:
		budgetDecimal, err := decimal.NewFromString(msg.BudgetText)
		if err != nil {
			var cmd tea.Cmd
			m.categoriesModel, cmd = m.categoriesModel.Update(categoryCreatedMsg{
				err: fmt.Errorf("invalid budget: %w", err),
			})
			return m, cmd
		}
		req := CreateCategoryRequest{
			Name:    msg.Name,
			Budget:  budgetDecimal,
			GroupID: msg.GroupID,
		}
		return m, createCategoryCmd(m.categoriesAPI, req)

	case categoryCreatedMsg:
		var cmd tea.Cmd
		m.categoriesModel, cmd = m.categoriesModel.Update(msg)
		return m, cmd

	case categoryUpdateSubmittedMsg:
		budgetDecimal, err := decimal.NewFromString(msg.BudgetText)
		if err != nil {
			var cmd tea.Cmd
			m.categoriesModel, cmd = m.categoriesModel.Update(categoryUpdatedMsg{
				err: fmt.Errorf("invalid budget: %w", err),
			})
			return m, cmd
		}

		req := UpdateCategoryRequest{
			Name:    msg.Name,
			Budget:  budgetDecimal,
			GroupID: msg.GroupID,
		}
		return m, updateCategoryCmd(m.categoriesAPI, msg.CategoryID, req)

	case categoryUpdatedMsg:
		var cmd tea.Cmd
		m.categoriesModel, cmd = m.categoriesModel.Update(msg)
		return m, cmd

	case categoryDeleteSubmittedMsg:
		return m, deleteCategoryCmd(m.categoriesAPI, msg.categoryID)

	case categoryDeletedMsg:
		var cmd tea.Cmd
		m.categoriesModel, cmd = m.categoriesModel.Update(msg)
		return m, cmd

	case categoriesNewRequestedMsg:
		groupOptions := make([]catGroupOption, len(m.groupsModel.groups))
		for i, group := range m.groupsModel.groups {
			groupOptions[i] = catGroupOption{
				ID:   group.ID,
				Name: group.GroupName,
			}
		}

		groupOptions = append(groupOptions, catGroupOption{
			ID:   uuid.Nil,
			Name: "None",
		})

		cm := m.categoriesModel
		cm.groupOptions = groupOptions
		cm.mode = categoriesModeFormNew
		cm.formFieldCursor = catFormFieldName
		cm.formEditing = false
		cm.nameInput.SetValue("")
		cm.nameInput.Blur()
		cm.budgetInput.SetValue("")
		cm.budgetInput.Blur()
		cm.formGroupIndex = 0
		cm.errorMsg = ""
		m.categoriesModel = cm
		return m, nil

	case categoriesEditRequestedMsg:
		groupOptions := make([]catGroupOption, len(m.groupsModel.groups))
		for i, group := range m.groupsModel.groups {
			groupOptions[i] = catGroupOption{
				ID:   group.ID,
				Name: group.GroupName,
			}
		}

		groupOptions = append(groupOptions, catGroupOption{
			ID:   uuid.Nil,
			Name: "None",
		})

		cm := m.categoriesModel
		currentCat := cm.categories[cm.cursor]
		cm.groupOptions = groupOptions
		cm.mode = categoriesModeFormEdit
		cm.formFieldCursor = catFormFieldName
		cm.formEditing = false
		cm.nameInput.SetValue(currentCat.CategoryName)
		cm.nameInput.Blur()
		cm.budgetInput.SetValue(currentCat.Budget.String())
		cm.budgetInput.Blur()
		cm.formGroupIndex = 0
		for i, option := range groupOptions {
			if option.ID == currentCat.GroupID {
				cm.formGroupIndex = i
				break
			}
		}

		cm.errorMsg = ""
		m.categoriesModel = cm
		return m, nil

	// Groups
	case groupsReloadRequestedMsg:
		cmd := loadGroupsCmd(m.groupsAPI)
		return m, cmd

	case groupsLoadedMsg:
		if msg.err != nil {
			var cmd tea.Cmd
			m.groupsModel, cmd = m.groupsModel.Update(groupsLoadedMsg{
				err: msg.err,
			})
			return m, cmd
		}

		m.groupsModel.groups = msg.groups
		m.groupsModel.cursor = 0
		return m, nil

	case groupCreateSubmittedMsg:
		req := CreateGroupRequest{
			Name: msg.Name,
		}
		return m, createGroupCmd(m.groupsAPI, req)

	case groupCreatedMsg:
		var cmd tea.Cmd
		m.groupsModel, cmd = m.groupsModel.Update(msg)
		return m, cmd

	case groupUpdateSubmittedMsg:
		req := UpdateGroupRequest{
			Name: msg.Name,
		}
		return m, updateGroupCmd(m.groupsAPI, msg.GroupID, req)

	case groupUpdatedMsg:
		var cmd tea.Cmd
		m.groupsModel, cmd = m.groupsModel.Update(msg)
		return m, cmd

	case groupDeleteSubmittedMsg:
		return m, deleteGroupCmd(m.groupsAPI, msg.groupID)

	case groupDeletedMsg:
		var cmd tea.Cmd
		m.groupsModel, cmd = m.groupsModel.Update(msg)
		return m, cmd

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
				if m.currentSection == sectionBudget {
					return m, loadBudgetCmd(m.budgetAPI)
				}
			}
		case focusMain:
			switch m.currentSection {
			case sectionAccounts:
				var cmd tea.Cmd
				m.accountsModel, cmd = m.accountsModel.Update(msg)
				return m, cmd
			case sectionBudget:
				var cmd tea.Cmd
				m.budgetModel, cmd = m.budgetModel.Update(msg)
				return m, cmd
			case sectionCategories:
				var cmd tea.Cmd
				m.categoriesModel, cmd = m.categoriesModel.Update(msg)
				return m, cmd
			case sectionGroups:
				var cmd tea.Cmd
				m.groupsModel, cmd = m.groupsModel.Update(msg)
				return m, cmd
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

	content := m.mainContentView()
	main := mainStyle.Render(content)

	if m.width > 0 {
		remainingWidth := m.width - lipgloss.Width(sidebar)
		if remainingWidth < 0 {
			remainingWidth = 0
		}
		main = mainStyle.
			Width(remainingWidth).
			Render(content)
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		main,
	)

	view := row
	if m.width > 0 && m.height > 0 {
		view = appStyle.
			Width(m.width).
			Height(m.height).
			Render(row)
	} else {
		view = appStyle.Render(row)
	}

	return view
}

func (m model) mainContentView() string {
	switch m.currentSection {
	case sectionBudget:
		return m.budgetModel.View()
	case sectionCategories:
		return m.categoriesModel.View()
	case sectionGroups:
		return m.groupsModel.View()
	case sectionAccounts:
		return m.accountsModel.View()
	case sectionTransactions:
		return m.transactionsModel.View()
	default:
		return ""
	}
}

func (m model) navView() string {
	title := sidebarTitleStyle.Render("BudgeTUI")

	var items []string
	for i, navItem := range m.navItems {
		style := sidebarItemStyle
		prefix := "  "

		if m.navCursor == i {
			style = sidebarItemActiveStyle
			prefix = "▶ "
		}

		items = append(items, style.Render(prefix+navItem))
	}

	navList := lipgloss.JoinVertical(lipgloss.Left, items...)

	helpText := helpStyle.Render("q: quit • tab: switch focus")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		navList,
		"",
		helpText,
	)

	return sidebarStyle.Render(content)
}
