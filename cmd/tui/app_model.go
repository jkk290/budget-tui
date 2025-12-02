package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	navItems  []string
	navCursor int

	currentSection section

	// budgetModel budgetModel
	// categoriesModel categoriesModel
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

	case accountsLoadedMsg:
		if msg.err != nil {
			// implement error message
			// m.accountsModel.err = msg.err
			return m, nil
		}

		m.accountsModel.accounts = msg.accounts
		m.accountsModel.cursor = 0
		return m, nil

	case transactionsLoadedMsg:
		if msg.err != nil {
			// implement error message
			// m.transactionsModel.err = msg.err
			return m, nil
		}

		m.transactionsModel.transactions = msg.transactions
		m.transactionsModel.cursor = 0

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
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
