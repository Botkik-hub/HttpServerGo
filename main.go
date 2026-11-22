package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/botkik-hub/httpServerGo/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print("Server cannot connect to the database, cannot run the server.")
		return
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
	}
	serveMux := http.ServeMux{}
	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", cfg.readinessCheck)
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsCheck)
	serveMux.HandleFunc("POST /admin/reset", cfg.clearMetrics)
	serveMux.HandleFunc("POST /api/users", cfg.createUser)
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.GetAllChirp)
	server.ListenAndServe()
}
