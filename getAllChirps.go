package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chrips, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type response struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}
	resp := make([]response, 0, len(chrips))
	for _, chirp := range chrips {
		resp = append(resp, response{
			Id:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	dat, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
