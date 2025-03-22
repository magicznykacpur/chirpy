package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/magicznykacpur/chirpy/internal/database"
)

func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("couldn't open database: %v\n", err)
		os.Exit(1)
	}

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, db: database.New(db)}
	mux := http.ServeMux{}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareServerHitsInc(fileServerHandler))
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValid)

	server := http.Server{Handler: &mux, Addr: ":" + port}

	fmt.Printf("starting server on %v\n", server.Addr)
	server.ListenAndServe()
}
