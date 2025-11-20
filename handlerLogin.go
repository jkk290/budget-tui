package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jkk290/budget-tui/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByUsername(req.Context(), params.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPw)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password", err)
		return
	}

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
		},
		Token: accessToken,
	})
}
