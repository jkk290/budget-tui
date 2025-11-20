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
	mux.HandleFunc("POST /api/v1/accounts", cfg.addAccount)

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
