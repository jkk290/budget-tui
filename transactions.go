package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
)

type Transaction struct {
	ID            uuid.UUID `json:"id"`
	Amount        float32   `json:"amount"`
	TxDescription string    `json:"tx_description"`
	TxDate        time.Time `json:"tx_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Posted        bool      `json:"posted"`
	AccountID     uuid.UUID `json:"account_id"`
	CategoryID    uuid.UUID `json:"category_id"`
}

func (cfg *apiConfig) addTransaction(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Amount        float32   `json:"amount"`
		TxDescription string    `json:"tx_description"`
		TxDate        time.Time `json:"tx_date"`
		Posted        bool      `json:"posted"`
		AccountID     uuid.UUID `json:"account_id"`
		CategoryID    uuid.UUID `json:"category_id"`
	}

	type response struct {
		Transaction
	}

	_, err := checkToken(req.Header, cfg.jwtSecret)
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
	txCategoryID := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}
	if params.CategoryID != uuid.Nil {
		txCategoryID.UUID = params.CategoryID
		txCategoryID.Valid = true
	}

	dbTransaction, err := cfg.db.AddTransaction(req.Context(), database.AddTransactionParams{
		ID:            uuid.New(),
		Amount:        params.Amount,
		TxDescription: params.TxDescription,
		TxDate:        params.TxDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Posted:        params.Posted,
		AccountID:     params.AccountID,
		CategoryID:    txCategoryID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create transaction", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Transaction: Transaction{
			ID:            dbTransaction.ID,
			Amount:        dbTransaction.Amount,
			TxDescription: dbTransaction.TxDescription,
			TxDate:        dbTransaction.TxDate,
			CreatedAt:     dbTransaction.CreatedAt,
			UpdatedAt:     dbTransaction.UpdatedAt,
			Posted:        dbTransaction.Posted,
			AccountID:     dbTransaction.AccountID,
			CategoryID:    dbTransaction.CategoryID.UUID,
		},
	})
}
