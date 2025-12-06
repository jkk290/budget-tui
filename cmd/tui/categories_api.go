package main

import (
	"context"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type CategoriesAPI interface {
	ListCategories(ctx context.Context) ([]Category, error)
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
