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

func (cfg *apiConfig) hitHandler(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHitsHandler(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	respWriter.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}
