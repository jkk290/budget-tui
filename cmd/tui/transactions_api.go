package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionsAPI interface {
	ListTransactions(ctx context.Context) ([]Transaction, error)
	CreateTransaction(ctx context.Context, req CreateTransactionRequest) (Transaction, error)
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

type transactionsNewRequestedMsg struct{}

type CreateTransactionRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	TxDescription string          `json:"tx_description"`
	TxDate        time.Time       `json:"tx_date"`
	Posted        bool            `json:"posted"`
	AccountID     uuid.UUID       `json:"account_id"`
	CategoryID    uuid.UUID       `json:"category_id"`
}

func (t *transactionsClient) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (Transaction, error) {
	httpReq, err := t.client.newJSONRequest(ctx, http.MethodPost, "/transactions", req)
	if err != nil {
		return Transaction{}, err
	}

	res, err := t.client.httpClient.Do(httpReq)
	if err != nil {
		return Transaction{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return Transaction{}, fmt.Errorf("Failed creating transaction: %s", res.Status)
	}

	var transaction Transaction
	if err := json.NewDecoder(res.Body).Decode(&transaction); err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

type transactionCreatedMsg struct {
	transaction Transaction
	err         error
}

func createTransactionCmd(api TransactionsAPI, req CreateTransactionRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		transaction, err := api.CreateTransaction(ctx, req)
		return transactionCreatedMsg{
			transaction: transaction,
			err:         err,
		}
	}
}

type transactionCreateSubmittedMsg struct {
	AmountText    string
	TxDescription string
	TxDate        string
	Posted        bool
	AccountID     string
	CategoryID    string
}

func submitNewTransactionMsg(amountText, txDescription, txDate, accountID, categoryID string, posted bool) tea.Cmd {
	return func() tea.Msg {
		return transactionCreateSubmittedMsg{
			AmountText:    amountText,
			TxDescription: txDescription,
			TxDate:        txDate,
			Posted:        posted,
			AccountID:     accountID,
			CategoryID:    categoryID,
		}
	}
}
