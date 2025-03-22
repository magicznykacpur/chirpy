package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (cfg *apiConfig) middlewareServerHitsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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
