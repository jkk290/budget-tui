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

type CategoriesAPI interface {
	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategoryTransactions(ctx context.Context, id uuid.UUID) ([]Transaction, error)
}

type categoriesClient struct {
	client *Client
}

func (c *Client) Categories() CategoriesAPI {
	return &categoriesClient{client: c}
}

func (c *categoriesClient) ListCategories(ctx context.Context) ([]Category, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, "/categories", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var categories []Category
	if err := json.NewDecoder(res.Body).Decode(&categories); err != nil {
		return nil, err
	}

	return categories, nil
}

type categoriesLoadedMsg struct {
	categories []Category
	err        error
}

func loadCategoriesCmd(api CategoriesAPI) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		categories, err := api.ListCategories(ctx)
		return categoriesLoadedMsg{
			categories: categories,
			err:        err,
		}
	}
}

func (c *categoriesClient) ListCategoryTransactions(ctx context.Context, id uuid.UUID) ([]Transaction, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, "/categories/"+id.String()+"/transactions", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var catTxs []Transaction
	if err := json.NewDecoder(res.Body).Decode(&catTxs); err != nil {
		return nil, err
	}

	return catTxs, nil
}

type loadCategoryTxsMsg struct {
	catTxs []Transaction
	err    error
}

func loadCategoryTxsCmd(api CategoriesAPI, id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		catTxs, err := api.ListCategoryTransactions(ctx, id)
		return loadCategoryTxsMsg{
			catTxs: catTxs,
			err:    err,
		}
	}
}

type loadCategoryTxsSubmittedMsg struct {
	categoryID uuid.UUID
}

func submitLoadCategoryTxsMsg(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return loadCategoryTxsSubmittedMsg{
			categoryID: id,
		}
	}
}

type categoriesNewRequestedMsg struct{}
type categoriesEditRequestedMsg struct{}

type CreateCategoryRequest struct {
	Name    string          `json:"category_name"`
	Budget  decimal.Decimal `json:"budget"`
	GroupID uuid.UUID       `json:"group_id"`
}

func (c *categoriesClient) CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error) {
	httpReq, err := c.client.newJSONRequest(ctx, http.MethodPost, "/categories", req)
	if err != nil {
		return Category{}, err
	}

	res, err := c.client.httpClient.Do(httpReq)
	if err != nil {
		return Category{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return Category{}, fmt.Errorf("Failed to create category: %s", res.Status)
	}

	var category Category
	if err := json.NewDecoder(res.Body).Decode(&category); err != nil {
		return Category{}, err
	}

	return category, nil
}

type categoryCreatedMsg struct {
	category Category
	err      error
}

func createCategoryCmd(api CategoriesAPI, req CreateCategoryRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		category, err := api.CreateCategory(ctx, req)
		return categoryCreatedMsg{
			category: category,
			err:      err,
		}
	}
}

type categoryCreateSubmittedMsg struct {
	Name       string
	BudgetText string
	GroupID    uuid.UUID
}

func submitCreateCategoryMsg(name, budget string, groupID uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return categoryCreateSubmittedMsg{
			Name:       name,
			BudgetText: budget,
			GroupID:    groupID,
		}
	}
}

type UpdateCategoryRequest struct {
	Name    string          `json:"category_name"`
	Budget  decimal.Decimal `json:"budget"`
	GroupID uuid.UUID       `json:"group_id"`
}

func (c *categoriesClient) UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (Category, error) {
	httpReq, err := c.client.newJSONRequest(ctx, http.MethodPut, "/categories/"+id.String(), req)
	if err != nil {
		return Category{}, err
	}

	res, err := c.client.httpClient.Do(httpReq)
	if err != nil {
		return Category{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Category{}, fmt.Errorf("Failed updating category: %s", res.Status)
	}

	var category Category
	if err := json.NewDecoder(res.Body).Decode(&category); err != nil {
		return Category{}, err
	}

	return category, nil
}

type categoryUpdatedMsg struct {
	category Category
	err      error
}

func updateCategoryCmd(api CategoriesAPI, id uuid.UUID, req UpdateCategoryRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		category, err := api.UpdateCategory(ctx, id, req)
		return categoryUpdatedMsg{
			category: category,
			err:      err,
		}
	}
}

type categoryUpdateSubmittedMsg struct {
	CategoryID uuid.UUID
	Name       string
	BudgetText string
	GroupID    uuid.UUID
}

func submitUpdateCategoryMsg(id uuid.UUID, name, budget string, groupID uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return categoryUpdateSubmittedMsg{
			CategoryID: id,
			Name:       name,
			BudgetText: budget,
			GroupID:    groupID,
		}
	}
}

func (c *categoriesClient) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	req, err := c.client.newRequest(ctx, http.MethodDelete, "/categories/"+id.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed deleting category: %s", res.Status)
	}

	return nil
}

type categoryDeletedMsg struct {
	categoryID uuid.UUID
	err        error
}

func deleteCategoryCmd(api CategoriesAPI, id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := api.DeleteCategory(ctx, id)
		return categoryDeletedMsg{
			categoryID: id,
			err:        err,
		}
	}
}

type categoryDeleteSubmittedMsg struct {
	categoryID uuid.UUID
}

func submitDeleteCategoryMsg(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return categoryDeleteSubmittedMsg{
			categoryID: id,
		}
	}
}

type categoriesReloadRequestedMsg struct{}
