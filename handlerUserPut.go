package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnikBarua007/http_server_go/internal/auth"
	"github.com/AnikBarua007/http_server_go/internal/database"
)

func (cfg *apiConfig) handlerUserPut(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID           string    `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Is_Chirp_Red bool      `json:"is_chirp_red"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	hashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		ID:             userID,
		Email:          parameters.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	respondWithJSON(w, http.StatusUnauthorized, response{
		ID:           updatedUser.ID.String(),
		CreatedAt:    updatedUser.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
		Email:        updatedUser.Email,
		Is_Chirp_Red: updatedUser.IsChirpyRed,
	})
}
