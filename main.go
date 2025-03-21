package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
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

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	mux := http.ServeMux{}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareServerHitsInc(fileServerHandler))
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := http.Server{Handler: &mux, Addr: ":" + port}

	fmt.Printf("starting server on %v\n", server.Addr)
	server.ListenAndServe()
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
