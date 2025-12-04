package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountsAPI interface {
	ListAccounts(ctx context.Context) ([]Account, error)
	CreateAccount(ctx context.Context, req CreateAccountRequest) (Account, error)
	// UpdateAccount(ctx context.Context, id uuid.UUID, req UpdateAccountRequest) (Account, error)
	// DeleteAccount(ctx context.Context, id uuid.UUID) error
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

type CreateAccountRequest struct {
	AccountName    string          `json:"account_name"`
	AccountType    string          `json:"account_type"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

func (a *accountsClient) CreateAccount(ctx context.Context, req CreateAccountRequest) (Account, error) {
	httpReq, err := a.client.newJSONRequest(ctx, http.MethodPost, "/accounts", req)
	if err != nil {
		return Account{}, err
	}

	resp, err := a.client.httpClient.Do(httpReq)
	if err != nil {
		return Account{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return Account{}, fmt.Errorf("Failed creating account: %s", resp.Status)
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

type accountCreateSubmittedMsg struct {
	Name        string
	Type        string
	BalanceText string
}

func submitNewAccountMsg(name, accountType, balance string) tea.Cmd {
	return func() tea.Msg {
		return accountCreateSubmittedMsg{
			Name:        name,
			Type:        accountType,
			BalanceText: balance,
		}
	}
}

type UpdateAccountRequest struct {
	AccountName string `json:"account_name"`
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

// func updateAccountCmd(api AccountsAPI, id uuid.UUID, req UpdateAccountRequest) tea.Cmd {
// 	return func() tea.Msg {
// 		ctx := context.Background()
// 		account, err := api.UpdateAccount(ctx, id, req)
// 		return accountUpdatedMsg{
// 			account: account,
// 			err:     err,
// 		}
// 	}
// }

// func deleteAccountCmd(api AccountsAPI, id uuid.UUID) tea.Cmd {
// 	return func() tea.Msg {
// 		ctx := context.Background()
// 		err := api.DeleteAccount(ctx, id)
// 		return accountDeleteMsg{
// 			accountID: id,
// 			err:       err,
// 		}
// 	}
// }
