package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loginResultMsg struct {
	token string
	err   error
}

func loginCmd(c *Client, username, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: username,
			Password: password,
		}

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return loginResultMsg{err: err}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/login", buf)
		if err != nil {
			return loginResultMsg{err: err}
		}
		req.Header.Set("content-type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return loginResultMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return loginResultMsg{err: fmt.Errorf("login failed: %s", resp.Status)}
		}

		var respBody struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return loginResultMsg{err: err}
		}

		return loginResultMsg{token: respBody.Token, err: nil}
	}
}

func (m model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "shift+tab":
			if m.loginUsername.Focused() {
				m.loginUsername.Blur()
				m.loginPassword.Focus()
			} else {
				m.loginPassword.Blur()
				m.loginUsername.Focus()
			}
			return m, nil

		case "enter":
			if m.loginPassword.Focused() {
				m.loginErr = ""
				cmd := loginCmd(m.client, m.loginUsername.Value(), m.loginPassword.Value())
				return m, cmd
			}
		}

		var cmds []tea.Cmd
		var cmd tea.Cmd
		m.loginUsername, cmd = m.loginUsername.Update(msg)
		cmds = append(cmds, cmd)
		m.loginPassword, cmd = m.loginPassword.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case loginResultMsg:
		if msg.err != nil {
			m.loginErr = msg.err.Error()
			return m, nil
		}

		m.jwt = msg.token
		m.client.SetJWT(m.jwt)
		m.accountsAPI = m.client.Accounts()
		m.transactionsAPI = m.client.Transactions()
		m.categoriesAPI = m.client.Categories()
		m.groupsAPI = m.client.Groups()
		var cmds []tea.Cmd
		cmds = append(cmds, loadBudgetCmd(m.budgetAPI))
		cmds = append(cmds, loadAccountsCmd(m.accountsAPI))
		cmds = append(cmds, loadTransactionsCmd(m.transactionsAPI))
		cmds = append(cmds, loadCategoriesCmd(m.categoriesAPI))
		cmds = append(cmds, loadGroupsCmd(m.groupsAPI))

		m.screen = screenMain
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) loginView() string {
	loginBoxWidth := m.width
	title := loginTitleStyle.
		Width(loginBoxWidth).
		Render("BudgeTUI Login")

	usernameInput := m.loginUsername.View()
	passwordInput := m.loginPassword.View()
	helpText := helpStyle.
		Width(loginBoxWidth).
		Render("(tab: switch between username/password • enter: on password to log in)")

	var errMsg string
	if m.loginErr != "" {
		errMsg = "Error: " + m.loginErr
	}
	loginErr := errorStyle.Render(errMsg)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		loginErr,
		usernameInput,
		passwordInput,
		"",
		helpText,
	)

	loginBox := loginStyle.
		Width(loginBoxWidth).
		Render(content)
	fullscreen := appBGStyle.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.Place(m.width, m.height,
				lipgloss.Center,
				lipgloss.Center,
				loginBox,
			),
		)
	return fullscreen
}
