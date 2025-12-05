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

	res, err := a.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var accounts []Account
	if err := json.NewDecoder(res.Body).Decode(&accounts); err != nil {
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

	res, err := a.client.httpClient.Do(httpReq)
	if err != nil {
		return Account{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return Account{}, fmt.Errorf("Failed creating account: %s", res.Status)
	}

	var account Account
	if err := json.NewDecoder(res.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

type accountCreatedMsg struct {
	account Account
	err     error
}

func createAccountCmd(api AccountsAPI, req CreateAccountRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		account, err := api.CreateAccount(ctx, req)
		return accountCreatedMsg{
			account: account,
			err:     err,
		}
	}
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

func (a *accountsClient) UpdateAccount(ctx context.Context, Id uuid.UUID, req UpdateAccountRequest) (Account, error) {
	httpReq, err := a.client.newJSONRequest(ctx, http.MethodPut, "/accounts/"+Id.String(), req)
	if err != nil {
		return Account{}, err
	}

	res, err := a.client.httpClient.Do(httpReq)
	if err != nil {
		return Account{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Account{}, fmt.Errorf("Failed updating account: %s", res.Status)
	}

	var account Account
	if err := json.NewDecoder(res.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

type accountUpdatedMsg struct {
	account Account
	err     error
}

func updateAccountCmd(api AccountsAPI, Id uuid.UUID, req UpdateAccountRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		account, err := api.UpdateAccount(ctx, Id, req)
		return accountUpdatedMsg{
			account: account,
			err:     err,
		}
	}
}

type accountUpdatedSubmittedMsg struct {
	AccountId uuid.UUID
	Name      string
}

func submitUpdateAccountMsg(Id uuid.UUID, name string) tea.Cmd {
	return func() tea.Msg {
		return accountUpdatedSubmittedMsg{
			AccountId: Id,
			Name:      name,
		}
	}
}

func (a *accountsClient) DeleteAccount(ctx context.Context, Id uuid.UUID) error {
	req, err := a.client.newRequest(ctx, http.MethodDelete, "/accounts/"+Id.String(), nil)
	if err != nil {
		return err
	}

	res, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed deleting account: %s", res.Status)
	}

	return nil
}

type accountDeletedMsg struct {
	accountID uuid.UUID
	err       error
}

func deleteAccountCmd(api AccountsAPI, id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := api.DeleteAccount(ctx, id)
		return accountDeletedMsg{
			accountID: id,
			err:       err,
		}
	}
}

type accountDeleteSubmittedMsg struct {
	AccountId uuid.UUID
}

func submitDeleteAccountMsg(Id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return accountDeleteSubmittedMsg{
			AccountId: Id,
		}
	}
}

type accountsReloadRequestedMsg struct{}
