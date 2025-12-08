package main

import (
	"context"
	"encoding/json"
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
}
