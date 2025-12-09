package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type categoriesMode int

const (
	categoriesModeList categoriesMode = iota
	categoriesModeDetails
	categoriesModeFormNew
	categoriesModeFormEdit
	categoriesModeDelete
)

const (
	catFormFieldName = iota
	catFormFieldBudget
	catFormFieldGroup
	catFormFieldSave
)

const (
	catConfirmYes = iota
	catConfirmCancel
)

type Category struct {
	ID           uuid.UUID       `json:"id"`
	CategoryName string          `json:"category_name"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Budget       decimal.Decimal `json:"budget"`
	UserID       uuid.UUID       `json:"user_id"`
	GroupID      uuid.UUID       `json:"group_id"`
	GroupName    string          `json:"group_name"`
}

type catGroupOption struct {
	ID   uuid.UUID
	Name string
}

type categoriesModel struct {
	mode            categoriesMode
	categories      []Category
	catTxs          []Transaction
	cursor          int
	formEditing     bool
	formFieldCursor int
	nameInput       textinput.Model
	budgetInput     textinput.Model
	formGroupIndex  int
	groupOptions    []catGroupOption
	confirmCursor   int
	errorMsg        string
}

func initialCategoriesModel() categoriesModel {
	catName := textinput.New()
	catName.CharLimit = 64
	catName.Blur()

	catBudget := textinput.New()
	catBudget.CharLimit = 64
	catBudget.Blur()

	return categoriesModel{
		categories:      []Category{},
		cursor:          0,
		formEditing:     false,
		formFieldCursor: catFormFieldName,
		nameInput:       catName,
		budgetInput:     catBudget,
		formGroupIndex:  0,
		groupOptions:    []catGroupOption{},
		confirmCursor:   catConfirmCancel,
	}
}

func (m categoriesModel) Update(msg tea.Msg) (categoriesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case categoryCreatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.mode = categoriesModeList
		return m, func() tea.Msg {
			return categoriesReloadRequestedMsg{}
		}

	case loadCategoryTxsMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.errorMsg = ""
		m.catTxs = msg.catTxs
		m.mode = categoriesModeDetails
		return m, nil

	case categoryUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.mode = categoriesModeList
		return m, func() tea.Msg {
			return categoriesReloadRequestedMsg{}
		}

	case categoryDeletedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}

		filtered := m.categories[:0]
		for _, category := range m.categories {
			if category.ID != msg.categoryID {
				filtered = append(filtered, category)
			}
		}
		m.categories = filtered

		if len(m.categories) == 1 {
			m.cursor = 0
			m.mode = categoriesModeList
		} else {
			prevCursor := m.cursor
			if prevCursor >= len(m.categories) {
				m.cursor = len(m.categories) - 1
			} else {
				m.cursor = prevCursor
			}
			m.mode = categoriesModeList
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		switch m.mode {
		case categoriesModeList:
			switch key {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.categories)-1 {
					m.cursor++
				}
			case "n":
				return m, func() tea.Msg {
					return categoriesNewRequestedMsg{}
				}
			case "d":
				if len(m.categories) > 0 {
					m.mode = categoriesModeDelete
					m.confirmCursor = catConfirmCancel
				}
			case "enter":
				if len(m.categories) > 0 {
					return m, submitLoadCategoryTxsMsg(m.categories[m.cursor].ID)
				}

			}
		case categoriesModeDetails:
			switch key {
			case "esc":
				m.mode = categoriesModeList
			case "e":
				return m, func() tea.Msg {
					return categoriesEditRequestedMsg{}
				}
			case "d":
				m.mode = categoriesModeDelete
				m.confirmCursor = catConfirmCancel
			}
		case categoriesModeFormNew, categoriesModeFormEdit:
			if m.formEditing {
				if key == "esc" {
					m.formEditing = false
					m.nameInput.Blur()
					m.budgetInput.Blur()
					return m, nil
				}

				switch m.formFieldCursor {
				case catFormFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				case catFormFieldBudget:
					var cmd tea.Cmd
					m.budgetInput, cmd = m.budgetInput.Update(msg)
					return m, cmd
				}
			}

			switch key {
			case "esc":
				m.mode = categoriesModeList
				m.errorMsg = ""
				return m, nil
			case "up", "k":
				if m.formFieldCursor > catFormFieldName {
					m.formFieldCursor--
				}
			case "down", "j":
				if m.formFieldCursor < catFormFieldSave {
					m.formFieldCursor++
				}
			case "enter":
				switch m.formFieldCursor {
				case catFormFieldName, catFormFieldBudget:
					m.formEditing = true
					m.nameInput.Blur()
					m.budgetInput.Blur()

					switch m.formFieldCursor {
					case catFormFieldName:
						m.nameInput.Focus()
					case catFormFieldBudget:
						m.budgetInput.Focus()
					}
				case catFormFieldSave:
					switch m.mode {
					case categoriesModeFormNew:
						name := m.nameInput.Value()
						budget := m.budgetInput.Value()
						group := m.groupOptions[m.formGroupIndex].ID
						return m, submitCreateCategoryMsg(name, budget, group)
					case categoriesModeFormEdit:
						name := m.nameInput.Value()
						budget := m.budgetInput.Value()
						group := m.groupOptions[m.formGroupIndex].ID
						return m, submitUpdateCategoryMsg(m.categories[m.cursor].ID, name, budget, group)
					}
				}
			default:
				switch m.formFieldCursor {
				case catFormFieldName:
					var cmd tea.Cmd
					m.nameInput, cmd = m.nameInput.Update(msg)
					return m, cmd
				case catFormFieldBudget:
					var cmd tea.Cmd
					m.budgetInput, cmd = m.budgetInput.Update(msg)
					return m, cmd
				case catFormFieldGroup:
					if key == "left" || key == "h" {
						if m.formGroupIndex > 0 {
							m.formGroupIndex--
						}
					}
					if key == "right" || key == "l" {
						if m.formGroupIndex < len(m.groupOptions)-1 {
							m.formGroupIndex++
						}
					}
				}
			}
		case categoriesModeDelete:
			switch key {
			case "esc":
				m.mode = categoriesModeList
			case "up", "k":
				if m.confirmCursor > confirmYes {
					m.confirmCursor--
				}
			case "down", "j":
				if m.confirmCursor < confirmCancel {
					m.confirmCursor++
				}
			case "enter":
				switch m.confirmCursor {
				case confirmYes:
					return m, submitDeleteCategoryMsg(m.categories[m.cursor].ID)
				case confirmCancel:
					m.mode = categoriesModeList
				}
			}
		}
	}
	return m, nil
}

func (m categoriesModel) View() string {
	switch m.mode {
	case categoriesModeList:
		s := "Categories\n\n"
		s += m.errorView()
		for i, category := range m.categories {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s Name: %s | Budget: $%s\n", cursor, category.CategoryName, category.Budget)
		}
		s += "\n(Use 'j'/'k' to move, 'enter' to view details, 'n' to create a new category, 'd' to delete category)\n"
		return s
	case categoriesModeDetails:
		cat := m.categories[m.cursor]
		s := "Category Details\n\n"
		s += fmt.Sprintf("Name: %s\n", cat.CategoryName)
		s += fmt.Sprintf("Budget: $%s\n", cat.Budget.String())
		s += fmt.Sprintf("Group: %s\n\n", cat.GroupName)

		s += "Transactions\n\n"
		for _, transaction := range m.catTxs {

			dateStr := transaction.TxDate.Format("2006-01-02")
			s += fmt.Sprintf("%s | %s | %s\n", dateStr, transaction.TxDescription, transaction.Amount)
		}

		s += "\n(Press 'esc' to go back, 'e' to edit, 'd' to delete)\n"

		return s
	case categoriesModeFormNew, categoriesModeFormEdit:
		s := "New Category\n\n"
		if m.mode == categoriesModeFormEdit {
			s = "Edit Category\n\n"
		}

		s += m.errorView()

		currentRow := func(field int) string {
			if m.formFieldCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s Name: %s\n", currentRow(catFormFieldName), m.nameInput.View())
		s += fmt.Sprintf("%s Budget: %s\n", currentRow(catFormFieldBudget), m.budgetInput.View())
		s += fmt.Sprintf("%s Group ('h'/'l' to change): %v\n", currentRow(catFormFieldGroup), m.groupOptions[m.formGroupIndex].Name)

		s += fmt.Sprintf("%s [ Save ]\n", currentRow(catFormFieldSave))
		s += "\n(Use 'j'/'k' to move, 'enter' to edit field, 'esc' to stop editing, 'esc' again to cancel)\n"

		return s
	case categoriesModeDelete:
		s := "Delete Category\n\n"
		s += fmt.Sprintf("Are you sure you want to delete category '%s'?\n", m.categories[m.cursor].CategoryName)

		currentRow := func(field int) string {
			if m.confirmCursor == field {
				return ">"
			}
			return " "
		}

		s += fmt.Sprintf("%s [ Yes ]\n", currentRow(catConfirmYes))
		s += fmt.Sprintf("%s [ Cancel ]\n", currentRow(catConfirmCancel))

		s += "\n(Use 'j'/'k' to move, 'enter' to select, 'esc' to cancel)"

		return s
	}

	return "Unknown category mode"
}

func (m categoriesModel) errorView() string {
	if m.errorMsg == "" {
		return ""
	}
	return fmt.Sprintf("Error: %s\n\n", m.errorMsg)
}

func (m categoriesModel) IsEditing() bool {
	return m.formEditing
}
