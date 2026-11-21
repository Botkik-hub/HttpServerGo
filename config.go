package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/botkik-hub/httpServerGo/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("charset", "utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("charset", "utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(
		`
<html>
  <body>
      <h1>Welcome, Chirpy Admin</h1>
	      <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) clearMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	if err := cfg.dbQueries.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot clear users: %s", err))
		return
	}
	if err := cfg.dbQueries.DeleteAllChirps(r.Context()); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot clear users: %s", err))
		return
	}
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		User string `json:"email"`
	}
	type respOk struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Updated   time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	got := params{}
	if err := decoder.Decode(&got); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot decrypt request body: %s", err))
		return
	}
	if got.User == "" {
		respondWithError(w, 400, "Cannot create user with empty email")
		return
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), got.User)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot inser user: %s", err))
		return
	}
	resp := respOk{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		Updated:   user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, resp)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type respOk struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Updated   time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	got := params{}
	if err := decoder.Decode(&got); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot decrypt request body: %s", err))
		return
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(),
		database.CreateChirpParams{
			Body:   got.Body,
			UserID: got.UserId,
		})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot inser chirp: %s", err))
		return
	}
	resp := respOk{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		Updated:   chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJSON(w, 201, resp)
}
