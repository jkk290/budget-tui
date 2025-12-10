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

type BudgetCategoryResponse struct {
	CategoryID   uuid.UUID       `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Budget       decimal.Decimal `json:"budget"`
	TotalSpent   decimal.Decimal `json:"total_spent"`
	Remaining    decimal.Decimal `json:"remaining"`
	IsOverspent  bool            `json:"is_overspent"`
}

type BudgetGroupResponse struct {
	GroupID        uuid.UUID                `json:"group_id"`
	GroupName      string                   `json:"group_name"`
	Categories     []BudgetCategoryResponse `json:"categories"`
	TotalBudget    decimal.Decimal          `json:"total_budget"`
	TotalSpent     decimal.Decimal          `json:"total_spent"`
	TotalRemaining decimal.Decimal          `json:"total_remaining"`
}

type BudgetOverviewResponse struct {
	StartDate           time.Time                `json:"start_date"`
	EndDate             time.Time                `json:"end_date"`
	Groups              []BudgetGroupResponse    `json:"groups"`
	UngroupedCategories []BudgetCategoryResponse `json:"ungrouped_categories"`
	GrandTotalBudget    decimal.Decimal          `json:"grand_total_budget"`
	GrandTotalSpent     decimal.Decimal          `json:"grand_total_spent"`
	GrandTotalRemaining decimal.Decimal          `json:"grand_total_remaining"`
}

type BudgetAPI interface {
	GetBudgetOverview(ctx context.Context) (*BudgetOverviewResponse, error)
}

type budgetClient struct {
	client *Client
}

func (c *Client) Budget() BudgetAPI {
	return &budgetClient{client: c}
}

func (b *budgetClient) GetBudgetOverview(ctx context.Context) (*BudgetOverviewResponse, error) {
	req, err := b.client.newRequest(ctx, http.MethodGet, "/budget", nil)
	if err != nil {
		return &BudgetOverviewResponse{}, err
	}

	res, err := b.client.httpClient.Do(req)
	if err != nil {
		return &BudgetOverviewResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &BudgetOverviewResponse{}, fmt.Errorf("Failed getting budget overview: %s", res.Status)
	}

	var overview BudgetOverviewResponse
	if err := json.NewDecoder(res.Body).Decode(&overview); err != nil {
		return &BudgetOverviewResponse{}, err
	}

	return &overview, nil
}

type budgetReloadRequestedMsg struct{}

type budgetLoadedMsg struct {
	overview *BudgetOverviewResponse
	err      error
}

func loadBudgetCmd(api BudgetAPI) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		overview, err := api.GetBudgetOverview(ctx)
		return budgetLoadedMsg{
			overview: overview,
			err:      err,
		}
	}
}
