package main

import (
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	serverMux.HandleFunc("/healthz", handlerReady)

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
