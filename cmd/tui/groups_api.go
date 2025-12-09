package main

import (
	"context"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type GroupsAPI interface {
	ListGroups(ctx context.Context) ([]Group, error)
}

type groupsClient struct {
	client *Client
}

func (c *Client) Groups() GroupsAPI {
	return &groupsClient{client: c}
}

func (c *groupsClient) ListGroups(ctx context.Context) ([]Group, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, "/groups", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var groups []Group
	if err := json.NewDecoder(res.Body).Decode(&groups); err != nil {
		return nil, err
	}

	return groups, nil
}

type groupsLoadedMsg struct {
	groups []Group
	err    error
}

func loadGroupsCmd(api GroupsAPI) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		groups, err := api.ListGroups(ctx)
		return groupsLoadedMsg{
			groups: groups,
			err:    err,
		}
	}
}
