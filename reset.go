package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") == "dev" {
		cfg.fileserverHits.Swap(0)
		cfg.db.DeleteUsers(r.Context())
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
