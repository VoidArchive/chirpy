package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		errorResp := map[string]string{"error": "Invalid JSON"}
		dat, _ := json.Marshal(errorResp)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		errorResp := map[string]string{"error": "Chirp is too long"}
		dat, _ := json.Marshal(errorResp)
		w.Write(dat)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	validResp := map[string]bool{"valid": true}
	dat, _ := json.Marshal(validResp)
	w.Write(dat)
}
