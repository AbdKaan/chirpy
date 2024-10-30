package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Api endpoint to check if the api is up and running
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlContent := fmt.Sprintf(`
	<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
	</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(htmlContent))
}

func (cfg *apiConfig) handlerResetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(200)
	w.Write([]byte("Config resetted."))
}
