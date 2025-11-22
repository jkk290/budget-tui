package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
)

type Account struct {
	ID          uuid.UUID `json:"id"`
	AccountName string    `json:"account_name"`
	AccountType string    `json:"account_type"`
	Balance     float32   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addAccount(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		AccountName string  `json:"account_name"`
		AccountType string  `json:"account_type"`
		Balance     float32 `json:"balance"`
	}

	type response struct {
		Account
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

	account, err := cfg.db.AddAccount(req.Context(), database.AddAccountParams{
		ID:          uuid.New(),
		AccountName: params.AccountName,
		AccountType: params.AccountType,
		Balance:     params.Balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add account", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Account: Account{
			ID:          account.ID,
			AccountName: account.AccountName,
			AccountType: account.AccountType,
			Balance:     account.Balance,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
			UserID:      account.UserID,
		},
	})
}

func (cfg *apiConfig) getAccounts(w http.ResponseWriter, req *http.Request) {

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbAccounts, err := cfg.db.GetAccountsByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve accounts", err)
		return
	}

	accounts := []Account{}
	for _, account := range dbAccounts {
		accounts = append(accounts, Account{
			ID:          account.ID,
			AccountName: account.AccountName,
			AccountType: account.AccountType,
			Balance:     account.Balance,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
			UserID:      account.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, accounts)
}

func (cfg *apiConfig) updateAccountInfo(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		AccountName string `json:"account_name"`
	}

	type response struct {
		Account
	}

	accountIDString := req.PathValue("accountID")
	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbAccount, err := cfg.db.GetAccountByID(req.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get account", err)
		return
	}
	if dbAccount.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't update this account's info", errors.New("Unauthorized"))
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	updatedAccount, err := cfg.db.UpdateAccountInfo(req.Context(), database.UpdateAccountInfoParams{
		ID:          dbAccount.ID,
		AccountName: params.AccountName,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update the account balance", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Account: Account{
			ID:          updatedAccount.ID,
			AccountName: updatedAccount.AccountName,
			AccountType: updatedAccount.AccountType,
			Balance:     updatedAccount.Balance,
			CreatedAt:   updatedAccount.CreatedAt,
			UpdatedAt:   updatedAccount.UpdatedAt,
			UserID:      updatedAccount.UserID,
		},
	})
}

func (cfg *apiConfig) deleteAccount(w http.ResponseWriter, req *http.Request) {
	accountIDString := req.PathValue("accountID")
	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbAccount, err := cfg.db.GetAccountByID(req.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get account", err)
		return
	}
	if dbAccount.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this account", errors.New("Unauthorized"))
		return
	}

	if err := cfg.db.DeleteAccount(req.Context(), dbAccount.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete account", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
