package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnikBarua007/http_server_go/internal/auth"
)

func (cfg *apiConfig) handlerlogins(w http.ResponseWriter, r *http.Request) {
	type parmas struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := parmas{}
	err := decoder.Decode(&parameters)
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
	res := response{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	dat, _ := json.Marshal(res)
	w.Write(dat)

}
