package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/magicznykacpur/chirpy/internal/auth"
)

type polkaRQ struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUpgradeWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != os.Getenv("POLKA_KEY") {
		writeError(err, "api key invalid", http.StatusUnauthorized, w)
		return
	}

	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(err, "couldn't read request bytes", http.StatusInternalServerError, w)
		return
	}

	var upgradeRQ polkaRQ
	err = json.Unmarshal(requestBytes, &upgradeRQ)
	if err != nil {
		writeError(err, "couldn't unmarshal request", http.StatusBadRequest, w)
		return
	}

	if upgradeRQ.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if upgradeRQ.Event == "user.upgraded" {
		userId, err := uuid.Parse(upgradeRQ.Data.UserId)
		if err != nil {
			writeError(err, "user id invalid", http.StatusBadRequest, w)
			return
		}

		user, err := cfg.db.GetUserById(r.Context(), userId)
		if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
			writeError(nil, "user not found", http.StatusNotFound, w)
			return
		}

		if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
			writeError(err, "cannot retrieve user", http.StatusInternalServerError, w)
			return
		}

		err = cfg.db.UpdateIsChirpyRed(r.Context(), user.ID)
		if err != nil {
			writeError(err, "cannot update user", http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
