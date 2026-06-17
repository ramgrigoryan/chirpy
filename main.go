package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ramgrigoryan/chirpy/internal/database"
)

const (
	port      = "8080"
	apiPrefix = "/app"
	rootDir   = "."
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("unable to connect to database")
		return
	}
	defer db.Close()

	dbQueries := database.New(db)

	serverMux := http.NewServeMux()

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	serverMux.Handle("/", cfg.middlewareMetricsInc(http.StripPrefix(apiPrefix, http.FileServer(http.Dir(rootDir)))))

	serverMux.HandleFunc("GET /api/healthz", HealthHandler)
	serverMux.HandleFunc("GET /admin/metrics", cfg.ReadHits)
	serverMux.HandleFunc("POST /admin/reset", cfg.Reset)
	serverMux.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	serverMux.HandleFunc("GET /api/chirps", cfg.GetChirps)
	serverMux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirp)
	serverMux.HandleFunc("POST /api/users", cfg.CreateUser)
	serverMux.HandleFunc("POST /api/login", cfg.AuthHandler)

	httpServer := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(httpServer.ListenAndServe())
}
