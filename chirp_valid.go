package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/magicznykacpur/chirpy/internal/cleaner"
)

type chirpValidRQ struct {
	Body string `json:"body"`
}

type chirpValidRes struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}


func handlerChirpValid(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "request body invalid", http.StatusBadRequest, w)
		return
	}

	var chirpValid chirpValidRQ
	err = json.Unmarshal(bytes, &chirpValid)
	if err != nil {
		writeError(err, "bad request", http.StatusBadRequest, w)
		return
	}

	if len(chirpValid.Body) < 140 {
		response := chirpValidRes{Valid: true}
		cleaned := cleaner.CleanBodyBy(chirpValid.Body, "kerfuffle")
		cleaned = cleaner.CleanBodyBy(cleaned, "sharbert")
		cleaned = cleaner.CleanBodyBy(cleaned, "fornax")
		response.CleanedBody = cleaned

		bytes, err := json.Marshal(response)
		if err != nil {
			writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	} else {
		response := chirpValidRes{Valid: false}
		bytes, err := json.Marshal(response)
		if err != nil {
			writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
	}
}
