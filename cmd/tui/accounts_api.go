package main

import (
	"context"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type AccountsAPI interface {
	ListAccounts(ctx context.Context) ([]Account, error)
}

type accountsClient struct {
	client *Client
}

func (c *Client) Accounts() AccountsAPI {
	return &accountsClient{client: c}
}

func (a *accountsClient) ListAccounts(ctx context.Context) ([]Account, error) {
	req, err := a.client.newRequest(ctx, http.MethodGet, "/accounts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var accounts []Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

type accountsLoadedMsg struct {
	accounts []Account
	err      error
}

func loadAccountsCmd(api AccountsAPI) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		accounts, err := api.ListAccounts(ctx)
		return accountsLoadedMsg{
			accounts: accounts,
			err:      err,
		}
	}
}
