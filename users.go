package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/magicznykacpur/chirpy/internal/auth"
	"github.com/magicznykacpur/chirpy/internal/database"
)

type userRQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userRes struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "couldn't read request bytes", http.StatusInternalServerError, w)
		return
	}

	var userRQ userRQ
	err = json.Unmarshal(requestBytes, &userRQ)
	if err != nil {
		writeError(err, "couldn't unmarshall user request", http.StatusBadRequest, w)
		return
	}

	if userRQ.Email == "" {
		writeError(nil, "user email cannot be empty", http.StatusBadRequest, w)
		return
	}

	if userRQ.Password == "" {
		writeError(nil, "user password cannot be empty", http.StatusBadRequest, w)
		return
	}

	hashedPassword, err := auth.HashPassword(userRQ.Password)
	if err != nil {
		writeError(err, "couldn't hash password", http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.db.Createuser(r.Context(), database.CreateuserParams{Email: userRQ.Email, HashedPassword: hashedPassword})
	if err != nil {
		writeError(err, "couldn't create user", http.StatusInternalServerError, w)
		return
	}

	response := userRes{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeError(err, "couldn't marshall user response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "couldn't read request bytes", http.StatusInternalServerError, w)
		return
	}

	var userRQ userRQ
	err = json.Unmarshal(requestBytes, &userRQ)
	if err != nil {
		writeError(err, "couldn't unmarshall user request", http.StatusBadRequest, w)
		return
	}

	if userRQ.Email == "" {
		writeError(nil, "user email cannot be empty", http.StatusBadRequest, w)
		return
	}

	if userRQ.Password == "" {
		writeError(nil, "user password cannot be empty", http.StatusBadRequest, w)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), userRQ.Email)
	if err != nil {
		writeError(nil, "incorrect email or password", http.StatusUnauthorized, w)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, userRQ.Password)
	if err != nil {
		writeError(nil, "incorrect email or password", http.StatusUnauthorized, w)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		writeError(err, "couldn't create a token", http.StatusInternalServerError, w)
		return
	}

	randomString, err := auth.MakeRefreshToken()
	if err != nil {
		writeError(err, "couldn't generate random string", http.StatusInternalServerError, w)
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(),
		database.CreateRefreshTokenParams{
			Token:     randomString,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		},
	)
	if err != nil {
		writeError(err, "couldn't create refresh token", http.StatusInternalServerError, w)
		return
	}

	response := userRes{
		Id:           user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed.Bool,
		Token:        token,
		RefreshToken: refreshToken.Token,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerUpdateEmailAndPassword(w http.ResponseWriter, r *http.Request) {
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
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "couldn't read request bytes", http.StatusInternalServerError, w)
		return
	}

	var userRQ userRQ
	err = json.Unmarshal(requestBytes, &userRQ)
	if err != nil {
		writeError(err, "couldn't unmarshall request bytes", http.StatusInternalServerError, w)
		return
	}

	hashedPassword, err := auth.HashPassword(userRQ.Password)
	if err != nil {
		writeError(err, "couldn't hash password", http.StatusInternalServerError, w)
		return
	}

	err = cfg.db.UpdateUserEmailAndPassword(r.Context(),
		database.UpdateUserEmailAndPasswordParams{
			Email:          userRQ.Email,
			HashedPassword: hashedPassword,
			ID:             userId,
		},
	)
	if err != nil {
		writeError(err, "couldn't update users data", http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), userRQ.Email)
	if err != nil {
		writeError(err, "couldn't retrieve user", http.StatusNotFound, w)
		return
	}

	userRes := userRes{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}

	responseBytes, err := json.Marshal(userRes)
	if err != nil {
		writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		writeError(err, "couldn't retrieve users", http.StatusInternalServerError, w)
	}

	response := []userRes{}
	for _, user := range users {
		response = append(response,
			userRes{
				Id:          user.ID.String(),
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
				Email:       user.Email,
				IsChirpyRed: user.IsChirpyRed.Bool,
			},
		)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeError(err, "couldn't marshall users response", http.StatusInternalServerError, w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
