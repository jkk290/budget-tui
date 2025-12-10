package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type groupsMode int

const (
	groupsModeList groupsMode = iota
	groupsModeDetails
	groupsModeFormNew
	groupsModeFormEdit
	groupsModeDelete
)

const (
	groupFormFieldName = iota
	groupFormFieldSave
)

const (
	groupConfirmYes = iota
	groupConfirmCancel
)

type Group struct {
	ID        uuid.UUID `json:"id"`
	GroupName string    `json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type groupsModel struct {
	mode            groupsMode
	groups          []Group
	cursor          int
	formEditing     bool
	formFieldCursor int
	nameInput       textinput.Model
	confirmCursor   int
	errorMsg        string
}

func initialGroupsModel() groupsModel {
	groupName := textinput.New()
	groupName.CharLimit = 64
	groupName.Blur()

	return groupsModel{
		groups:          []Group{},
		cursor:          0,
		formEditing:     false,
		formFieldCursor: groupFormFieldName,
		nameInput:       groupName,
		confirmCursor:   groupConfirmCancel,
	}
}

func (m groupsModel) Update(msg tea.Msg) (groupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case groupCreatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.mode = groupsModeList
		return m, func() tea.Msg {
			return groupsReloadRequestedMsg{}
		}

	case groupUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.mode = groupsModeList
		return m, func() tea.Msg {
			return groupsReloadRequestedMsg{}
		}

	case groupDeletedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}

		filtered := m.groups[:0]
		for _, group := range m.groups {
			if group.ID != msg.groupID {
				filtered = append(filtered, group)
			}
		}
		m.groups = filtered

		if len(m.groups) == 1 {
			m.cursor = 0
			m.mode = groupsModeList
		} else {
			prevCursor := m.cursor
			if prevCursor >= len(m.groups) {
				m.cursor = len(m.groups) - 1
			} else {
				m.cursor = prevCursor
			}
			m.mode = groupsModeList
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		switch m.mode {
		case groupsModeList:
			switch key {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.groups)-1 {
					m.cursor++
				}
			case "n":
				m.mode = groupsModeFormNew

				m.formFieldCursor = groupFormFieldName
				m.formEditing = false

				m.nameInput.SetValue("")
				m.nameInput.Blur()
				m.errorMsg = ""
			case "d":
				if len(m.groups) > 0 {
					m.mode = groupsModeDelete
					m.confirmCursor = groupConfirmCancel
				}
			case "enter":
				if len(m.groups) > 0 {
					m.mode = groupsModeDetails
				}

			}
		case groupsModeDetails:
			switch key {
			case "esc":
				m.mode = groupsModeList
			case "e":
				m.mode = groupsModeFormEdit
				m.formFieldCursor = groupFormFieldName
				m.formEditing = false

				m.nameInput.SetValue(m.groups[m.cursor].GroupName)
				m.nameInput.Blur()
			case "d":
				m.mode = groupsModeDelete
				m.confirmCursor = groupConfirmCancel
			}
		case groupsModeFormNew, groupsModeFormEdit:
			if m.formEditing {
				if key == "esc" {
					m.formEditing = false
					m.nameInput.Blur()
					return m, nil
				}

				switch m.formFieldCursor {
				case groupFormFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				}
			}

			switch key {
			case "esc":
				m.mode = groupsModeList
				m.errorMsg = ""
				return m, nil
			case "up", "k":
				if m.formFieldCursor > groupFormFieldName {
					m.formFieldCursor--
				}
			case "down", "j":
				if m.formFieldCursor < groupFormFieldSave {
					m.formFieldCursor++
				}
			case "enter":
				switch m.formFieldCursor {
				case groupFormFieldName:
					m.formEditing = true
					m.nameInput.Blur()

					switch m.formFieldCursor {
					case groupFormFieldName:
						m.nameInput.Focus()
					}
				case groupFormFieldSave:
					switch m.mode {
					case groupsModeFormNew:
						name := m.nameInput.Value()
						return m, submitCreateGroupMsg(name)
					case groupsModeFormEdit:
						name := m.nameInput.Value()
						return m, submitUpdateGroupMsg(m.groups[m.cursor].ID, name)
					}
				}
			default:
				switch m.formFieldCursor {
				case groupFormFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				}
			}
		case groupsModeDelete:
			switch key {
			case "esc":
				m.mode = groupsModeList
			case "up", "k":
				if m.confirmCursor > groupConfirmYes {
					m.confirmCursor--
				}
			case "down", "j":
				if m.confirmCursor < groupConfirmCancel {
					m.confirmCursor++
				}
			case "enter":
				switch m.confirmCursor {
				case groupConfirmYes:
					return m, submitDeleteGroupMsg(m.groups[m.cursor].ID)
				case groupConfirmCancel:
					m.mode = groupsModeList
				}
			}
		}
	}
	return m, nil
}

func (m groupsModel) View() string {
	switch m.mode {
	case groupsModeList:
		s := "Groups\n\n"
		s += m.errorView()
		for i, group := range m.groups {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s Name: %s\n", cursor, group.GroupName)
		}
		s += "\n(Use 'j'/'k' to move, 'enter' to view details, 'n' to create a new group, 'd' to delete group)\n"
		return s
	case groupsModeDetails:
		group := m.groups[m.cursor]
		s := "Group Details\n\n"
		s += fmt.Sprintf("Name: %s\n", group.GroupName)

		s += "\n(Press 'esc' to go back, 'e' to edit, 'd' to delete)\n"

		return s
	case groupsModeFormNew, groupsModeFormEdit:
		s := "New Group\n\n"
		if m.mode == groupsModeFormEdit {
			s = "Edit Group\n\n"
		}

		s += m.errorView()

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Name: %s\n", currentRow(groupFormFieldName), m.nameInput.View())

		s += fmt.Sprintf("%s [ Save ]\n", currentRow(groupFormFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"

		return s
	case groupsModeDelete:
		s := "Delete Group\n\n"
		s += fmt.Sprintf("Are you sure you want to delete group '%s'?\n", m.groups[m.cursor].GroupName)

		currentRow := func(field int) string {
			if m.confirmCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s [ Yes ]\n", currentRow(groupConfirmYes))
		s += fmt.Sprintf("%s [ Cancel ]\n", currentRow(groupConfirmCancel))

		s += "\n(Use 'j'/'k' to move, 'enter' to select, 'esc' to cancel)"

		return s
	}

	return "Unknown group mode"
}

func (m groupsModel) errorView() string {
	if m.errorMsg == "" {
		return ""
	}
	return fmt.Sprintf("Error: %s\n\n", m.errorMsg)
}

func (m groupsModel) IsEditing() bool {
	return m.formEditing
}
