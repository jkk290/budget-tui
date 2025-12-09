package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `json:"id"`
	GroupName string    `json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type groupsModel struct {
	groups []Group
}

func initialGroupsModel() groupsModel {
	return groupsModel{
		groups: []Group{},
	}
}

func (m groupsModel) Update(msg tea.Msg) (groupsModel, tea.Cmd) {
	return m, nil
}
