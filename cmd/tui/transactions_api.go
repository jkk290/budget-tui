package main

import (
	"context"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type TransactionsAPI interface {
	ListTransactions(ctx context.Context) ([]Transaction, error)
}

type transactionsClient struct {
	client *Client
}

func (c *Client) Transactions() TransactionsAPI {
	return &transactionsClient{client: c}
}

func (t *transactionsClient) ListTransactions(ctx context.Context) ([]Transaction, error) {
	req, err := t.client.newRequest(ctx, http.MethodGet, "/transactions", nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

type transactionsLoadedMsg struct {
	transactions []Transaction
	err          error
}

func loadTransactionsCmd(api TransactionsAPI) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		transactions, err := api.ListTransactions(ctx)
		return transactionsLoadedMsg{
			transactions: transactions,
			err:          err,
		}
	}
}
