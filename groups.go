package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/budget-tui/internal/auth"
	"github.com/jkk290/budget-tui/internal/database"
)

type Group struct {
	ID        uuid.UUID `json:"id"`
	GroupName string    `json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createGroup(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		GroupName string `json:"group_name"`
	}

	type response struct {
		Group
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
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

	if params.GroupName == "" {
		respondWithError(w, http.StatusBadRequest, "Missing group name", errors.New("invalid parameters"))
		return
	}

	dbGroup, err := cfg.db.CreateGroup(req.Context(), database.CreateGroupParams{
		ID:        uuid.New(),
		GroupName: params.GroupName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create group", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Group: Group{
			ID:        dbGroup.ID,
			GroupName: dbGroup.GroupName,
			CreatedAt: dbGroup.CreatedAt,
			UpdatedAt: dbGroup.UpdatedAt,
			UserID:    dbGroup.UserID,
		},
	})
}

func (cfg *apiConfig) getGroups(w http.ResponseWriter, req *http.Request) {
	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbGroups, err := cfg.db.GetGroupsByUser(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get groups", err)
		return
	}

	groups := []Group{}
	for _, group := range dbGroups {
		groups = append(groups, Group{
			ID:        group.ID,
			GroupName: group.GroupName,
			CreatedAt: group.CreatedAt,
			UpdatedAt: group.UpdatedAt,
			UserID:    group.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, groups)

}

func (cfg *apiConfig) updateGroup(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		GroupName string `json:"group_name"`
	}

	type response struct {
		Group
	}

	groupIDString := req.PathValue("groupID")
	groupID, err := uuid.Parse(groupIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbGroup, err := cfg.db.GetGroupByID(req.Context(), groupID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get group", err)
		return
	}
	if dbGroup.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't update this group", errors.New("unauthorized"))
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.GroupName == "" {
		respondWithError(w, http.StatusBadRequest, "Missing group name", errors.New("invalid parameters"))
		return
	}

	updatedGroup, err := cfg.db.UpdateGroup(req.Context(), database.UpdateGroupParams{
		ID:        dbGroup.ID,
		GroupName: params.GroupName,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update group", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Group: Group{
			ID:        updatedGroup.ID,
			GroupName: updatedGroup.GroupName,
			CreatedAt: updatedGroup.CreatedAt,
			UpdatedAt: dbGroup.UpdatedAt,
			UserID:    dbGroup.UserID,
		},
	})
}

func (cfg *apiConfig) deleteGroup(w http.ResponseWriter, req *http.Request) {
	groupIDString := req.PathValue("groupID")
	groupID, err := uuid.Parse(groupIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	userID, err := checkToken(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbGroup, err := cfg.db.GetGroupByID(req.Context(), groupID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get group", err)
		return
	}
	if dbGroup.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this group", errors.New("unauthorized"))
		return
	}

	if err := cfg.db.DeleteGroup(req.Context(), dbGroup.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete group", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
