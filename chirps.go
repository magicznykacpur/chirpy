package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/magicznykacpur/chirpy/internal/auth"
	"github.com/magicznykacpur/chirpy/internal/cleaner"
	"github.com/magicznykacpur/chirpy/internal/database"
)

type createChirpRQ struct {
	Body string `json:"body"`
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeError(err, "unauthorized", http.StatusUnauthorized, w)
		return
	}

	userId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		writeError(err, "cannot validate jwt token", http.StatusUnauthorized, w)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "request body invalid", http.StatusBadRequest, w)
		return
	}

	var createChirpRQ createChirpRQ
	err = json.Unmarshal(bytes, &createChirpRQ)
	if err != nil || createChirpRQ.Body == "" {
		writeError(err, "bad request, check if chirp contains body", http.StatusBadRequest, w)
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

	w.Header().Set("Content-Type", "application/json")
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
	authorParam := r.URL.Query().Get("author_id")
	var authorId uuid.UUID
	var err error

	if authorParam != "" {
		authorId, err = uuid.Parse(authorParam)
		if err != nil {
			writeError(err, "id malformed", http.StatusBadRequest, w)
			return
		}
	}

	var chirps []database.Chirp

	if authorId != (uuid.UUID{}) {
		chirps, err = cfg.db.GetChirpsByUser(r.Context(), authorId)
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
	}

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
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(err, "couldn't parse id", http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		writeError(nil, "chirp not found", http.StatusNotFound, w)
		return
	}

	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		writeError(err, "couldn't retrieve chirp", http.StatusInternalServerError, w)
		return
	}

	chirpRes := chirpRes{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}

	responseBytes, err := json.Marshal(chirpRes)
	if err != nil {
		writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeError(err, "couldn't get bearer token", http.StatusUnauthorized, w)
		return
	}

	userId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		writeError(err, "token invalid", http.StatusUnauthorized, w)
		return
	}

	defer r.Body.Close()
	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(err, "couldn't parse id", http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		writeError(nil, "chirp not found", http.StatusNotFound, w)
		return
	}
	if chirp.UserID != userId {
		writeError(err, "cannot delete other users chirps", http.StatusForbidden, w)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(),
		database.DeleteChirpParams{
			ID:     chirpId,
			UserID: userId,
		},
	)
	if err != nil {
		writeError(err, "couldn't delete chirp", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
