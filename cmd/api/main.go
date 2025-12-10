package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jkk290/budget-tui/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set in .env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in .env")
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET not set in .env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to postgres database: %v", err)
	}
	defer db.Close()

	cfg := &apiConfig{
		db:        database.New(db),
		jwtSecret: tokenSecret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/hello", handlerHello)

	mux.HandleFunc("POST /api/v1/users", cfg.createUser)
	mux.HandleFunc("POST /api/v1/login", cfg.handlerLogin)

	mux.HandleFunc("GET /api/v1/accounts", cfg.getAccounts)
	mux.HandleFunc("POST /api/v1/accounts", cfg.addAccount)
	mux.HandleFunc("PUT /api/v1/accounts/{accountID}", cfg.updateAccountInfo)
	mux.HandleFunc("DELETE /api/v1/accounts/{accountID}", cfg.deleteAccount)
	mux.HandleFunc("GET /api/v1/accounts/{accountID}/transactions", cfg.getAccountTransactions)

	mux.HandleFunc("GET /api/v1/groups", cfg.getGroups)
	mux.HandleFunc("POST /api/v1/groups", cfg.createGroup)
	mux.HandleFunc("PUT /api/v1/groups/{groupID}", cfg.updateGroup)
	mux.HandleFunc("DELETE /api/v1/groups/{groupID}", cfg.deleteGroup)

	mux.HandleFunc("GET /api/v1/categories", cfg.getCategories)
	mux.HandleFunc("POST /api/v1/categories", cfg.createCategory)
	mux.HandleFunc("PUT /api/v1/categories/{categoryID}", cfg.updateCategory)
	mux.HandleFunc("DELETE /api/v1/categories/{categoryID}", cfg.deleteCategory)
	mux.HandleFunc("GET /api/v1/categories/{categoryID}/transactions", cfg.getCategoryTransactions)

	mux.HandleFunc("GET /api/v1/transactions", cfg.getUserTransactions)
	mux.HandleFunc("POST /api/v1/transactions", cfg.addTransaction)
	mux.HandleFunc("PUT /api/v1/transactions/{transactionID}", cfg.updateTransaction)
	mux.HandleFunc("DELETE /api/v1/transactions/{transactionID}", cfg.deleteTransaction)

	mux.HandleFunc("GET /api/v1/budget", cfg.handlerGetBudgetOverview)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerHello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
