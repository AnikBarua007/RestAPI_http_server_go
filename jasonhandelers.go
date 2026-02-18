package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"time"

	"github.com/AnikBarua007/http_server_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleruser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	type response struct {
		ID         string    `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	if err := decoder.Decode(&parameters); err != nil {
		respondWithError(w, 400, "invalid request body")
		return
	}
	if parameters.Email == "" {
		respondWithError(w, 400, "email is required")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), parameters.Email)
	if err != nil {
		respondWithError(w, 500, err.Error()) // temporary for debugging
		return
	}
	res := response{
		ID:         user.ID.String(),
		Created_At: user.CreatedAt,
		Updated_At: user.UpdatedAt,
		Email:      user.Email,
	}
	dat, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

//	func handlerValidate(w http.ResponseWriter, r *http.Request) {
//		type perameters struct {
//			Body string `json:"body"`
//		}
//		type response struct {
//			Cleaned_body string `json:"cleaned_body"`
//		}
//		decoder := json.NewDecoder(r.Body)
//		params := perameters{}
//		err := decoder.Decode(&params)
//		if err != nil {
//			respondWithError(w, 400, "Something went wrong")
//			return
//		}
//		const maxLength = 140
//		if len(params.Body) > 140 {
//			respondWithError(w, 400, "Chirp is too long")
//			return
//		}
//		mesage := strings.ToLower(params.Body)
//		words := strings.Split(mesage, " ")
//		badwords := []string{"kerfuffle", "sharbert", "fornax"}
//		for i, word := range words {
//			if slices.Contains(badwords, word) {
//				words[i] = "****"
//			}
//		}
//		message := strings.Join(words, " ")
//
//		respondWithJSON(w, 200, response{
//			Cleaned_body: message,
//		})
//		return
//	}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnvals struct {
		Error string `json:"error"`
	}
	respBody := returnvals{Error: msg}
	dat, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//type response struct {
	//	valid bool `json:"valid"`
	//}
	//respBody := response{valid: true}
	dat, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	parameters := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, "invalid request body")
		return
	}
	if parameters.Body == "" {
		respondWithError(w, 400, "body is required")
		return
	}
	const maxLength = 140
	if len(parameters.Body) > maxLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	words := strings.Split(parameters.Body, " ")
	badwords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if slices.Contains(badwords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	cleanBody := strings.Join(words, " ")

	userID, err := uuid.Parse(parameters.UserID)
	if err != nil {
		respondWithError(w, 400, "invalid user_id")
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	type response struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}
	res := response{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	dat, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
