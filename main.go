package main

import (
	"log"
	"net/http"
)

const (
	port      = "8080"
	apiPrefix = "/app"
	rootDir   = "."
)

func main() {

	serverMux := http.NewServeMux()

	cfg := &apiConfig{}

	serverMux.Handle("/", cfg.middlewareMetricsInc(http.StripPrefix(apiPrefix, http.FileServer(http.Dir(rootDir)))))

	serverMux.HandleFunc("GET /api/healthz", HealthHandler)
	serverMux.HandleFunc("GET /admin/metrics", cfg.ReadHits)
	serverMux.HandleFunc("POST /admin/reset", cfg.Reset)
	serverMux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	httpServer := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(httpServer.ListenAndServe())
}
