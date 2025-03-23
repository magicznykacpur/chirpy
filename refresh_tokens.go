package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/magicznykacpur/chirpy/internal/auth"
)

type tokenRes struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeError(err, "couldn't get bearer token", http.StatusUnauthorized, w)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		writeError(err, "couldn't get refresh token", http.StatusUnauthorized, w)
		return
	}

	if time.Until(refreshToken.ExpiresAt) < 0 {
		writeError(err, "token expired", http.StatusUnauthorized, w)
		return
	}

	if refreshToken.RevokedAt != (sql.NullTime{}) {
		writeError(err, "token revoked", http.StatusUnauthorized, w)
		return
	}

	jwtToken, err := auth.MakeJWT(refreshToken.UserID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		writeError(err, "couldn't create a jwt token", http.StatusInternalServerError, w)
		return
	}

	response := tokenRes{Token: jwtToken}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeError(err, "couldn't marshal response", http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeError(err, "couldn't get bearer token", http.StatusUnauthorized, w)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		writeError(err, "couldn't revoke refresh token", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
