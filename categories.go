package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
)

type Category struct {
	ID           uuid.UUID `json:"id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Budget       float32   `json:"budget"`
	UserID       uuid.UUID `json:"user_id"`
	GroupID      uuid.UUID `json:"group_id"`
}

func (cfg *apiConfig) createCategory(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		CategoryName string    `json:"category_name"`
		Budget       float32   `json:"budget"`
		GroupID      uuid.UUID `json:"groupID"`
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
		Budget:       params.Budget,
		UserID:       userID,
		GroupID:      categoryGroup,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create group", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Category: Category{
			ID:           dbCategory.ID,
			CategoryName: dbCategory.CategoryName,
			CreatedAt:    dbCategory.CreatedAt,
			UpdatedAt:    dbCategory.UpdatedAt,
			Budget:       dbCategory.Budget,
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
		categories = append(categories, Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
			Budget:       category.Budget,
			UserID:       category.UserID,
			GroupID:      category.GroupID.UUID,
		})
	}

	respondWithJSON(w, http.StatusOK, categories)
}
