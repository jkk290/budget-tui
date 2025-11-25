package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
)

type Category struct {
	ID           uuid.UUID `json:"id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Budget       float64   `json:"budget"`
	UserID       uuid.UUID `json:"user_id"`
	GroupID      uuid.UUID `json:"group_id"`
}

func (cfg *apiConfig) createCategory(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		CategoryName string    `json:"category_name"`
		Budget       float64   `json:"budget"`
		GroupID      uuid.UUID `json:"group_id"`
	}

	type response struct {
		Category
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	categoryGroup := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}

	if params.GroupID != uuid.Nil {
		categoryGroup.UUID = params.GroupID
		categoryGroup.Valid = true
	}

	dbCategory, err := cfg.db.CreateCategory(req.Context(), database.CreateCategoryParams{
		ID:           uuid.New(),
		CategoryName: params.CategoryName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Budget:       strconv.FormatFloat(params.Budget, 'f', 2, 64),
		UserID:       userID,
		GroupID:      categoryGroup,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create group", err)
		return
	}

	dbBudgetFloat, err := strconv.ParseFloat(dbCategory.Budget, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse budget", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Category: Category{
			ID:           dbCategory.ID,
			CategoryName: dbCategory.CategoryName,
			CreatedAt:    dbCategory.CreatedAt,
			UpdatedAt:    dbCategory.UpdatedAt,
			Budget:       dbBudgetFloat,
			UserID:       dbCategory.UserID,
			GroupID:      dbCategory.GroupID.UUID,
		},
	})

}

func (cfg *apiConfig) getCategories(w http.ResponseWriter, req *http.Request) {
	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbCategories, err := cfg.db.GetCategoriesByUser(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get categories", err)
		return
	}

	categories := []Category{}
	for _, category := range dbCategories {
		budgetFloat, err := strconv.ParseFloat(category.Budget, 64)
		if err != nil {
			log.Printf("Error parsing category budget: %v", category.ID)
			continue
		}

		categories = append(categories, Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
			Budget:       budgetFloat,
			UserID:       category.UserID,
			GroupID:      category.GroupID.UUID,
		})
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (cfg *apiConfig) updateCategory(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		CategoryName string    `json:"category_name"`
		Budget       string    `json:"budget"`
		GroupID      uuid.UUID `json:"group_id"`
	}

	type response struct {
		Category
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	categoryIDString := req.PathValue("categoryID")
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	dbCategory, err := cfg.db.GetCategoryByID(req.Context(), categoryID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get category", err)
		return
	}
	if dbCategory.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can't update category", errors.New("Unauthorized"))
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	updatedName := dbCategory.CategoryName
	updatedBudget := dbCategory.Budget
	updatedGroup := dbCategory.GroupID

	if params.CategoryName != updatedName {
		updatedName = params.CategoryName
	}

	if params.Budget != updatedBudget {
		updatedBudget = params.Budget
	}

	if params.GroupID != updatedGroup.UUID {
		updatedGroup.UUID = params.GroupID
		updatedGroup.Valid = true
	}

	updatedCategory, err := cfg.db.UpdateCategory(req.Context(), database.UpdateCategoryParams{
		ID:           dbCategory.ID,
		CategoryName: updatedName,
		Budget:       updatedBudget,
		GroupID:      updatedGroup,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update category", err)
		return
	}

	updatedBudgetFloat, err := strconv.ParseFloat(dbCategory.Budget, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse budget", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Category: Category{
			ID:           updatedCategory.ID,
			CategoryName: updatedCategory.CategoryName,
			CreatedAt:    updatedCategory.CreatedAt,
			UpdatedAt:    updatedCategory.UpdatedAt,
			Budget:       updatedBudgetFloat,
			UserID:       updatedCategory.UserID,
			GroupID:      updatedCategory.GroupID.UUID,
		},
	})

}

func (cfg *apiConfig) deleteCategory(w http.ResponseWriter, req *http.Request) {
	categoryIDString := req.PathValue("categoryID")
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbCategory, err := cfg.db.GetCategoryByID(req.Context(), categoryID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get category", err)
		return
	}
	if dbCategory.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can't delete category", errors.New("Unauthorized"))
		return
	}

	if err := cfg.db.DeleteCategory(req.Context(), dbCategory.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete category", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
