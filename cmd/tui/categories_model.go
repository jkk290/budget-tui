package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Category struct {
	ID           uuid.UUID       `json:"id"`
	CategoryName string          `json:"category_name"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Budget       decimal.Decimal `json:"budget"`
	UserID       uuid.UUID       `json:"user_id"`
	GroupID      uuid.UUID       `json:"group_id"`
}

type categoriesModel struct {
	categories []Category
}

func initialCategoriesModel() categoriesModel {
	return categoriesModel{
		categories: []Category{},
	}
}

func (m categoriesModel) Update(msg tea.Msg) (categoriesModel, tea.Cmd) {
	return m, nil
}
