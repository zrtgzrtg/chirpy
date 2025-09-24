package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

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
	apiCfg = &apiConfig{atomic.Int32{}}
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
	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	server.ListenAndServe()
}
