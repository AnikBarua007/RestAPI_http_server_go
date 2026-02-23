package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnikBarua007/http_server_go/internal/auth"
	"github.com/AnikBarua007/http_server_go/internal/database"
)

func (cfg *apiConfig) handlerlogins(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID            string    `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		Refresh       string    `json:"refresh_token"`
		IS_CHIRPY_RED bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), parameters.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	matched, err := auth.CheckPasswordHash(parameters.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !matched {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create auth token")
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token")
		return
	}
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to persist refresh token")
		return
	}
	res := response{
		ID:            user.ID.String(),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		Refresh:       refreshToken,
		IS_CHIRPY_RED: user.IsChirpyRed,
	}
	w.Header().Set("Content-Type", "application/json")
	dat, _ := json.Marshal(res)
	w.Write(dat)
}
