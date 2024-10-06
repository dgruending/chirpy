package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/dgruending/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	serverSecret   string
}

func main() {
	// Setup environment variable with godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
		serverSecret:   os.Getenv("SERVER_SECRET"),
	}
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("GET /api/healthz", handlerReady)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.hitHandler)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serverMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
	serverMux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}

func handlerReady(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte(http.StatusText(http.StatusOK)))
}
