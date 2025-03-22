package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)



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
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValid)

	server := http.Server{Handler: &mux, Addr: ":" + port}

	fmt.Printf("starting server on %v\n", server.Addr)
	server.ListenAndServe()
}

