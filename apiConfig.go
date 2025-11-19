package main

import "github.com/jkk290/budget-tui/internal/database"

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
}
