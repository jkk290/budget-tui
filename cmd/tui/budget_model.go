package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type budgetMode int

const (
	budgetModeOverview budgetMode = iota
)

type budgetModel struct {
	mode     budgetMode
	overview *BudgetOverviewResponse
	errorMsg string
}

func initialBudgetModel() budgetModel {
	return budgetModel{
		mode:     budgetModeOverview,
		overview: nil,
		errorMsg: "",
	}
}

func (m budgetModel) Update(msg tea.Msg) (budgetModel, tea.Cmd) {
	switch msg := msg.(type) {
	case budgetLoadedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.overview = nil
			return m, nil
		}

		m.errorMsg = ""
		m.overview = msg.overview
		return m, nil
	}

	return m, nil
}

func (m budgetModel) View() string {
	if m.overview == nil {
		return "Loading budget data...\n"
	}
	s := fmt.Sprintf("Budget Overview - %s\n\n", m.overview.StartDate.Format("Jan 2006"))

	s += m.errorView()

	for _, group := range m.overview.Groups {
		s += fmt.Sprintf("=== %s ===\n", group.GroupName)
		s += fmt.Sprintf("Group Total - Budget: $%s | Spent: $%s | Remaining: $%s\n\n", group.TotalBudget.StringFixed(2), group.TotalSpent.StringFixed(2), group.TotalRemaining.StringFixed(2))

		for _, cat := range group.Categories {
			overspentTag := ""
			if cat.IsOverspent {
				overspentTag = " [OVERSPENT]"
			}
			s += fmt.Sprintf("  %s%s\n", cat.CategoryName, overspentTag)
			s += fmt.Sprintf("    Budget: $%s | Spent: $%s | Remaining: $%s\n\n", cat.Budget.StringFixed(2), cat.TotalSpent.StringFixed(2), cat.Remaining.StringFixed(2))
		}
		s += "\n"
	}

	if len(m.overview.UngroupedCategories) > 0 {
		s += "=== Ungrouped Categories ===\n\n"

		for _, cat := range m.overview.UngroupedCategories {
			overspentTag := ""
			if cat.IsOverspent {
				overspentTag = " [OVERSPENT]"
			}
			s += fmt.Sprintf("  %s%s\n", cat.CategoryName, overspentTag)
			s += fmt.Sprintf("    Budget: $%s | Spent: $%s | Remaining: $%s\n\n", cat.Budget.StringFixed(2), cat.TotalSpent.StringFixed(2), cat.Remaining.StringFixed(2))
		}
		s += "\n"
	}

	s += "===============================\n"
	s += fmt.Sprintf("TOTAL - Budget: $%s | Spent: $%s | Remaining: $%s\n", m.overview.GrandTotalBudget.StringFixed(2), m.overview.GrandTotalSpent.StringFixed(2), m.overview.GrandTotalRemaining.StringFixed(2))

	return s
}

func (m budgetModel) errorView() string {
	if m.errorMsg == "" {
		return ""
	}
	return fmt.Sprintf("Error: %s\n\n", m.errorMsg)
}
