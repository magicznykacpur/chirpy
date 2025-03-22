package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/magicznykacpur/chirpy/internal/database"
)

type createUserRQ struct {
	Email string `json:"email"`
}

type createUserRes struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "couldn't read request bytes", http.StatusInternalServerError, w)
		return
	}

	var userRQ createUserRQ
	err = json.Unmarshal(requestBytes, &userRQ)
	if err != nil {
		writeError(err, "couldn't unmarshall user request", http.StatusBadRequest, w)
		return
	}

	if userRQ.Email == "" {
		writeError(nil, "user email cannot be empty", http.StatusBadRequest, w)
		return
	}

	user, err := cfg.db.Createuser(r.Context(), userRQ.Email)
	if err != nil {
		writeError(err, "couldn't create user", http.StatusInternalServerError, w)
		return
	}

	userRes := createUserRes{Id: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	responseBytes, err := json.Marshal(userRes)
	if err != nil {
		writeError(err, "couldn't marshall user response", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

type getUsersRes struct {
	Users []database.User `json:"users"`
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		writeError(err, "couldn't retrieve users", http.StatusInternalServerError, w)
	}

	usersRes := getUsersRes{Users: users}
	responseBytes, err := json.Marshal(usersRes)
	if err != nil {
		writeError(err, "couldn't marshall users response", http.StatusInternalServerError, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
