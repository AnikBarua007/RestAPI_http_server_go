package main

import (
	"encoding/json"
	"net/http"

	"github.com/AnikBarua007/http_server_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request) {
	type params struct {
		EventID string `json:"event"`
		Data    struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.polkaapi {
		respondWithError(w, http.StatusUnauthorized, "api key not allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	if err := decoder.Decode(&parameters); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if parameters.EventID != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := cfg.dbQueries.Upgrade_chirpy(r.Context(), parameters.Data.UserID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}
	//type webhookResponse struct {
	//	Status  string `json:"status"`
	//	Message string `json:"message"`
	//	UserID  string `json:"user_id,omitempty"`
	//}
	//resp := webhookResponse{
	//	Status:  "ok",
	//	Message: "user upgraded",
	//	UserID:  parameters.Data.UserID.String(),
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(resp)
	respondWithError(w, http.StatusOK, "success here upgrade user")
	w.Write([]byte("ok"))
}
