package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/zrtgzrtg/chirpy/internal/database"
)

var apiCfg *apiConfig

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(int32(1))
		next.ServeHTTP(w, req)
	})
}
func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		// fallback for Boot.dev/local tests
		dbUrl = "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
	}
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiCfg = &apiConfig{atomic.Int32{}, *dbQueries, platform}
	serverMux := http.NewServeMux()
	server := &http.Server{
		Handler: serverMux,
		Addr:    ":8080",
	}
	appDir := "/home/zrtg/bootworkspace/chirpy/app"
	fServer := http.StripPrefix("/app", http.FileServer(http.Dir(appDir)))

	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(fServer))
	serverMux.Handle("/app/assets/logo.png", http.StripPrefix("/app", http.FileServer(http.Dir(appDir))))
	serverMux.HandleFunc("GET /api/healthz", handlerReady)
	serverMux.HandleFunc("GET /admin/metrics", handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", handlerReset)
	serverMux.HandleFunc("POST /api/users", handlerUser)
	serverMux.HandleFunc("POST /api/chirps", handlerPostChirp)
	server.ListenAndServe()
}
