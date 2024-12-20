package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/AbdKaan/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_ID   string    `json:"user_id"`
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	platform := os.Getenv("PLATFORM")

	secret := os.Getenv("SECRET")

	polkaKey := os.Getenv("POLKA_KEY")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
		polkaKey:       polkaKey,
	}

	handler := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	handler.Handle("/app/", fsHandler)

	handler.HandleFunc("GET /api/healthz", handlerReadiness)

	handler.HandleFunc("GET /api/chirps", apiCfg.handlerGetPosts)
	handler.HandleFunc("POST /api/chirps", apiCfg.handlerCreatePost)
	handler.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeletePost)
	handler.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetPost)

	handler.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	handler.HandleFunc("PUT /api/users", apiCfg.handlerUpdateEmailAndPassword)

	handler.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	handler.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	handler.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	handler.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	handler.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// webhooks
	handler.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeRedChirpy)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
