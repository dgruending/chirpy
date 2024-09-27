package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("GET /healthz", handlerReady)
	serverMux.HandleFunc("GET /metrics", apiCfg.hitHandler)
	serverMux.HandleFunc("POST /reset", apiCfg.resetHitsHandler)

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
