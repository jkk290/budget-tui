package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
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

func (cfg *apiConfig) handlerGetBudgetOverview(w http.ResponseWriter, req *http.Request) {
	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()

	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	rows, err := cfg.db.GetUserBudgetOverviewForMonth(req.Context(), database.GetUserBudgetOverviewForMonthParams{
		UserID:   userID,
		TxDate:   startDate,
		TxDate_2: endDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get budget overview", err)
		return
	}

	groupsMap := make(map[uuid.UUID]*BudgetGroupResponse)
	var ungroupedCategories []BudgetCategoryResponse

	for _, row := range rows {
		remaining := row.Budget.Sub(row.TotalSpent)
		isOverspent := remaining.IsNegative()

		category := BudgetCategoryResponse{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Budget:       row.Budget,
			TotalSpent:   row.TotalSpent,
			Remaining:    remaining,
			IsOverspent:  isOverspent,
		}

		if !row.GroupID.Valid {
			ungroupedCategories = append(ungroupedCategories, category)
			continue
		}

		groupID := row.GroupID.UUID
		if _, exists := groupsMap[groupID]; !exists {
			groupsMap[groupID] = &BudgetGroupResponse{
				GroupID:        groupID,
				GroupName:      row.GroupName.String,
				Categories:     []BudgetCategoryResponse{},
				TotalBudget:    decimal.Zero,
				TotalSpent:     decimal.Zero,
				TotalRemaining: decimal.Zero,
			}
		}

		group := groupsMap[groupID]
		group.Categories = append(group.Categories, category)
		group.TotalBudget = group.TotalBudget.Add(row.Budget)
		group.TotalSpent = group.TotalSpent.Add(row.TotalSpent)
		group.TotalRemaining = group.TotalRemaining.Add(remaining)
	}

	groups := make([]BudgetGroupResponse, 0, len(groupsMap))
	for _, group := range groupsMap {
		groups = append(groups, *group)
	}

	grandTotalBudget := decimal.Zero
	grandTotalSpent := decimal.Zero

	for _, group := range groups {
		grandTotalBudget = grandTotalBudget.Add(group.TotalBudget)
		grandTotalSpent = grandTotalSpent.Add(group.TotalSpent)
	}

	for _, cat := range ungroupedCategories {
		grandTotalBudget = grandTotalBudget.Add(cat.Budget)
		grandTotalSpent = grandTotalSpent.Add(cat.TotalSpent)
	}

	grandTotalRemaining := grandTotalBudget.Sub(grandTotalSpent)

	response := BudgetOverviewResponse{
		StartDate:           startDate,
		EndDate:             endDate,
		Groups:              groups,
		UngroupedCategories: ungroupedCategories,
		GrandTotalBudget:    grandTotalBudget,
		GrandTotalSpent:     grandTotalSpent,
		GrandTotalRemaining: grandTotalRemaining,
	}
	respondWithJSON(w, http.StatusOK, response)
}
