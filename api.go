package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/magicznykacpur/chirpy/internal/cleaner"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareServerHitsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	adminHTML, err := os.ReadFile("admin.html")

	if err != nil {
		fmt.Println("couldn't read file admin.html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	} else {
		serverHits := strings.ReplaceAll(string(adminHTML), "%d", fmt.Sprintf("%d", cfg.fileserverHits.Load()))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverHits))
	}

}

type chirpValidRQ struct {
	Body string `json:"body"`
}

type chirpValidRes struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type errorRes struct {
	Message string `json:"error"`
}

func handlerChirpValid(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		marshalError(err, "request body invalid", http.StatusBadRequest, w)
		return
	}

	var chirpValid chirpValidRQ
	err = json.Unmarshal(bytes, &chirpValid)
	if err != nil {
		marshalError(err, "bad request", http.StatusBadRequest, w)
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
			marshalError(err, "couldn't marshal response", http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	} else {
		response := chirpValidRes{Valid: false}
		bytes, err := json.Marshal(response)
		if err != nil {
			marshalError(err, "couldn't marshal response", http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
	}
}

func marshalError(err error, message string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	response := errorRes{}
	if err == nil {
		response.Message = message
	} else {
		response.Message = fmt.Sprintf("%s: %v", message, err)
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("marshalling error"))
	} else {
		w.Write(bytes)
	}
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
