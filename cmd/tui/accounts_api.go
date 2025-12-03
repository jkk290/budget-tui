package main

import (
	"context"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type AccountsAPI interface {
	ListAccounts(ctx context.Context) ([]Account, error)
	CreateAccount(ctx context.Context, req CreateAccountRequest) (Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, req UpdateAccountRequest) (Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
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

type accountsCreatedMsg struct {
	account Account
	err     error
}

type accountUpdatedMsg struct {
	account Account
	err     error
}

type accountDeleteMsg struct {
	accountID uuid.UUID
	err       error
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

func createAccountCmd(api AccountsAPI, req CreateAccountRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		account, err := api.CreateAccount(ctx, req)
		return accountsCreatedMsg{
			account: account,
			err:     err,
		}
	}
}

func updateAccountCmd(api AccountsAPI, id uuid.UUID, req UpdateAccountRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		account, err := api.UpdateAccount(ctx, id, req)
		return accountUpdatedMsg{
			account: account,
			err:     err,
		}
	}
}

func deleteAccountCmd(api AccountsAPI, id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := api.DeleteAccount(ctx, id)
		return accountDeleteMsg{
			accountID: id,
			err:       err,
		}
	}
}
