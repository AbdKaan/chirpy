package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	handler := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	handler.Handle("/app/", fsHandler)

	handler.HandleFunc("GET /api/healthz", handlerReadiness)
	handler.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	handler.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	handler.HandleFunc("POST /admin/reset", apiCfg.handlerResetConfig)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
