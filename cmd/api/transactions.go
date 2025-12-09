package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/database"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	Amount        decimal.Decimal `json:"amount"`
	TxDescription string          `json:"tx_description"`
	TxDate        time.Time       `json:"tx_date"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Posted        bool            `json:"posted"`
	AccountID     uuid.UUID       `json:"account_id"`
	CategoryID    uuid.UUID       `json:"category_id"`
}

func (cfg *apiConfig) addTransaction(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Amount        decimal.Decimal `json:"amount"`
		TxDescription string          `json:"tx_description"`
		TxDate        time.Time       `json:"tx_date"`
		Posted        bool            `json:"posted"`
		AccountID     uuid.UUID       `json:"account_id"`
		CategoryID    uuid.UUID       `json:"category_id"`
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

	if params.Amount.Equal(decimal.NewFromInt(0)) || params.TxDescription == "" || params.TxDate.IsZero() || params.AccountID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Missing amount, description, date, and/or account", errors.New("invalid parameters"))
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

	// dbAmountFloat, err := strconv.ParseFloat(dbTransaction.Amount, 64)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't parse amount", err)
	// 	return
	// }

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

func (cfg *apiConfig) getUserTransactions(w http.ResponseWriter, req *http.Request) {
	type userTransaction struct {
		ID            uuid.UUID       `json:"id"`
		Amount        decimal.Decimal `json:"amount"`
		TxDescription string          `json:"tx_description"`
		TxDate        time.Time       `json:"tx_date"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		Posted        bool            `json:"posted"`
		AccountID     uuid.UUID       `json:"account_id"`
		CategoryID    uuid.UUID       `json:"category_id"`
		AccountName   string          `json:"account_name"`
		CategoryName  string          `json:"category_name"`
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbTransactions, err := cfg.db.GetUserTransactions(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user transactions", err)
	}

	var transactions []userTransaction
	for _, tx := range dbTransactions {
		transactions = append(transactions, userTransaction{
			ID:            tx.ID,
			Amount:        tx.Amount,
			TxDescription: tx.TxDescription,
			TxDate:        tx.TxDate,
			CreatedAt:     tx.CreatedAt,
			UpdatedAt:     tx.UpdatedAt,
			Posted:        tx.Posted,
			AccountID:     tx.AccountID,
			CategoryID:    tx.CategoryID.UUID,
			AccountName:   tx.AccountName,
			CategoryName:  tx.CategoryName,
		})
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

func (cfg *apiConfig) updateTransaction(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Amount        decimal.Decimal `json:"amount"`
		TxDescription string          `json:"tx_description"`
		TxDate        time.Time       `json:"tx_date"`
		Posted        bool            `json:"posted"`
		AccountID     uuid.UUID       `json:"account_id"`
		CategoryID    uuid.UUID       `json:"category_id"`
	}

	type response struct {
		Transaction
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	transactionIDString := req.PathValue("transactionID")
	transactionID, err := uuid.Parse(transactionIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	dbTransactionUserID, err := cfg.db.GetTransactionUserID(req.Context(), transactionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get transaction", err)
		return
	}
	if dbTransactionUserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Can't update transaction", errors.New("unauthorized"))
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Amount.Equal(decimal.NewFromInt(0)) || params.TxDescription == "" || params.TxDate.IsZero() || params.AccountID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Missing amount, description, date, and/or account", errors.New("invalid parameters"))
		return
	}

	updatedCatergoryID := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}
	if params.CategoryID != uuid.Nil {
		updatedCatergoryID.UUID = params.CategoryID
		updatedCatergoryID.Valid = true
	}

	updatedTransaction, err := cfg.db.UpdateTransaction(req.Context(), database.UpdateTransactionParams{
		ID:            transactionID,
		Amount:        params.Amount,
		TxDescription: params.TxDescription,
		TxDate:        params.TxDate,
		Posted:        params.Posted,
		AccountID:     params.AccountID,
		CategoryID:    updatedCatergoryID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update transaction", err)
		return
	}

	// updatedAmountFloat, err := strconv.ParseFloat(updatedTransaction.Amount, 64)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't parse amount", err)
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, response{
		Transaction: Transaction{
			ID:            updatedTransaction.ID,
			Amount:        updatedTransaction.Amount,
			TxDescription: updatedTransaction.TxDescription,
			TxDate:        updatedTransaction.TxDate,
			Posted:        updatedTransaction.Posted,
			CreatedAt:     updatedTransaction.CreatedAt,
			UpdatedAt:     updatedTransaction.UpdatedAt,
			AccountID:     updatedTransaction.AccountID,
			CategoryID:    updatedTransaction.CategoryID.UUID,
		},
	})
}

func (cfg *apiConfig) deleteTransaction(w http.ResponseWriter, req *http.Request) {
	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	transactionIDString := req.PathValue("transactionID")
	transactionID, err := uuid.Parse(transactionIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	dbTransactionUserID, err := cfg.db.GetTransactionUserID(req.Context(), transactionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get transaction", err)
		return
	}
	if dbTransactionUserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Can't delete transaction", errors.New("unauthorized"))
		return
	}

	if err := cfg.db.DeleteTransaction(req.Context(), transactionID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete transaction", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) getAccountTransactions(w http.ResponseWriter, req *http.Request) {
	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	accountIDString := req.PathValue("accountID")
	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	dbAccount, err := cfg.db.GetAccountByID(req.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get account", err)
		return
	}
	if dbAccount.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can't view account transactions", errors.New("unauthorized"))
		return
	}

	dbTransactions, err := cfg.db.GetTransactionsByAccount(req.Context(), dbAccount.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get transactions", err)
		return
	}

	transactions := []Transaction{}
	for _, transaction := range dbTransactions {
		// amountFloat, err := strconv.ParseFloat(transaction.Amount, 64)
		// if err != nil {
		// 	log.Printf("Error parsing transaction amount: %v", transaction.ID)
		// 	continue
		// }
		transactions = append(transactions, Transaction{
			ID:            transaction.ID,
			Amount:        transaction.Amount,
			TxDescription: transaction.TxDescription,
			TxDate:        transaction.TxDate,
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
			Posted:        transaction.Posted,
			AccountID:     transaction.AccountID,
			CategoryID:    transaction.CategoryID.UUID,
		})
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

func (cfg *apiConfig) getCategoryTransactions(w http.ResponseWriter, req *http.Request) {
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
		respondWithError(w, http.StatusForbidden, "Can't view category transactions", errors.New("unauthorized"))
		return
	}

	catId := uuid.NullUUID{
		UUID:  dbCategory.ID,
		Valid: true,
	}

	dbTransactions, err := cfg.db.GetTransactionsByCategory(req.Context(), catId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get transactions", err)
		return
	}

	transactions := []Transaction{}
	for _, transaction := range dbTransactions {
		transactions = append(transactions, Transaction{
			ID:            transaction.ID,
			Amount:        transaction.Amount,
			TxDescription: transaction.TxDescription,
			TxDate:        transaction.TxDate,
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
			Posted:        transaction.Posted,
			AccountID:     transaction.AccountID,
			CategoryID:    transaction.CategoryID.UUID,
		})
	}

	respondWithJSON(w, http.StatusOK, transactions)
}
