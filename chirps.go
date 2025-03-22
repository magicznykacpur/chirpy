package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/magicznykacpur/chirpy/internal/cleaner"
	"github.com/magicznykacpur/chirpy/internal/database"
)

type createChirpRQ struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type chirpRes struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "request body invalid", http.StatusBadRequest, w)
		return
	}

	var createChirpRQ createChirpRQ
	err = json.Unmarshal(bytes, &createChirpRQ)
	if err != nil || createChirpRQ.Body == "" || createChirpRQ.UserId == "" {
		writeError(err, "bad request, check if chirp contains user id and body", http.StatusBadRequest, w)
		return
	}

	userId, err := uuid.Parse(createChirpRQ.UserId)
	if err != nil {
		writeError(err, "cannot parse user id", http.StatusBadRequest, w)
		return
	}

	if len(createChirpRQ.Body) > 140 {
		writeError(nil, "chirp body too long, max 140 characters", http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(),
		database.CreateChirpParams{
			Body:   cleanChirpBody(createChirpRQ.Body),
			UserID: userId,
		},
	)
	if err != nil {
		writeError(err, "couldn't create chirp", http.StatusInternalServerError, w)
		return
	}

	response := chirpRes{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func cleanChirpBody(body string) string {
	cleaned := body
	for _, badWord := range badWords {
		cleaned = cleaner.CleanBodyBy(cleaned, badWord)
	}
	return cleaned
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		writeError(err, "couldn't retrieve chirps", http.StatusInternalServerError, w)
		return
	}

	chirpResList := []chirpRes{}
	for _, chirp := range chirps {
		chirpResList = append(chirpResList,
			chirpRes{
				Id:        chirp.ID.String(),
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID.String(),
			},
		)
	}

	responseBytes, err := json.Marshal(chirpResList)
	if err != nil {
		writeError(err, "couldn't marshall response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	w.Write(responseBytes)
}
