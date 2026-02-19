package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerlogins(w http.ResponseWriter, r *http.Request) {
	type parmas struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	looger := decoder.Decode(&parmas{})

}
