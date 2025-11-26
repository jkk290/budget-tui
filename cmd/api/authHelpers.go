package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/auth"
)

func checkToken(h http.Header, secret string) (uuid.UUID, error) {

	accessToken, err := auth.GetBearerToken(h)
	if err != nil {
		return uuid.Nil, err
	}

	return auth.ValidateJWT(accessToken, secret)
}
