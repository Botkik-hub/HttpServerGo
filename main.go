package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serveMux := http.ServeMux{}
	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /healthz", readinessCheck)
	serveMux.HandleFunc("GET /metrics", cfg.metricsCheck)
	serveMux.HandleFunc("POST /reset", cfg.clearMetrics)
	server.ListenAndServe()
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
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
		<h1>Hits: %v</h1>
	</body>
</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) clearMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
