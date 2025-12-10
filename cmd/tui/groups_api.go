package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type GroupsAPI interface {
	ListGroups(ctx context.Context) ([]Group, error)
	CreateGroup(ctx context.Context, req CreateGroupRequest) (Group, error)
	UpdateGroup(ctx context.Context, id uuid.UUID, req UpdateGroupRequest) (Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
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

type CreateGroupRequest struct {
	Name string `json:"group_name"`
}

func (g *groupsClient) CreateGroup(ctx context.Context, req CreateGroupRequest) (Group, error) {
	httpReq, err := g.client.newJSONRequest(ctx, http.MethodPost, "/groups", req)
	if err != nil {
		return Group{}, err
	}

	res, err := g.client.httpClient.Do(httpReq)
	if err != nil {
		return Group{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return Group{}, fmt.Errorf("Failed creating group: %s", res.Status)
	}

	var group Group
	if err := json.NewDecoder(res.Body).Decode(&group); err != nil {
		return Group{}, err
	}

	return group, nil
}

type groupCreatedMsg struct {
	group Group
	err   error
}

func createGroupCmd(api GroupsAPI, req CreateGroupRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		group, err := api.CreateGroup(ctx, req)
		return groupCreatedMsg{
			group: group,
			err:   err,
		}
	}
}

type groupCreateSubmittedMsg struct {
	Name string
}

func submitCreateGroupMsg(name string) tea.Cmd {
	return func() tea.Msg {
		return groupCreateSubmittedMsg{
			Name: name,
		}
	}
}

type UpdateGroupRequest struct {
	Name string `json:"group_name"`
}

func (g *groupsClient) UpdateGroup(ctx context.Context, id uuid.UUID, req UpdateGroupRequest) (Group, error) {
	httpReq, err := g.client.newJSONRequest(ctx, http.MethodPut, "/groups/"+id.String(), req)
	if err != nil {
		return Group{}, err
	}

	res, err := g.client.httpClient.Do(httpReq)
	if err != nil {
		return Group{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, fmt.Errorf("Failed updating group: %s", res.Status)
	}

	var group Group
	if err := json.NewDecoder(res.Body).Decode(&group); err != nil {
		return Group{}, err
	}

	return group, nil
}

type groupUpdatedMsg struct {
	group Group
	err   error
}

func updateGroupCmd(api GroupsAPI, id uuid.UUID, req UpdateGroupRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		group, err := api.UpdateGroup(ctx, id, req)
		return groupUpdatedMsg{
			group: group,
			err:   err,
		}
	}
}

type groupUpdateSubmittedMsg struct {
	GroupID uuid.UUID
	Name    string
}

func submitUpdateGroupMsg(id uuid.UUID, name string) tea.Cmd {
	return func() tea.Msg {
		return groupUpdateSubmittedMsg{
			GroupID: id,
			Name:    name,
		}
	}
}

func (c *groupsClient) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	req, err := c.client.newRequest(ctx, http.MethodDelete, "/groups/"+id.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed deleting group: %s", res.Status)
	}

	return nil
}

type groupDeletedMsg struct {
	groupID uuid.UUID
	err     error
}

func deleteGroupCmd(api GroupsAPI, id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := api.DeleteGroup(ctx, id)
		return groupDeletedMsg{
			groupID: id,
			err:     err,
		}
	}
}

type groupDeleteSubmittedMsg struct {
	groupID uuid.UUID
}

func submitDeleteGroupMsg(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		return groupDeleteSubmittedMsg{
			groupID: id,
		}
	}
}

type groupsReloadRequestedMsg struct{}
