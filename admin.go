package main

import "net/http"

func (cfg *apiConfig) resetHandler(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		respondWithError(respWriter, http.StatusForbidden, "Dev feature")
		return
	}
	type respValues struct {
		HitsReset    bool `json:"hits_reset"`
		UsersCleared bool `json:"users_cleared"`
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ClearUsers(request.Context())
	if err != nil {
		respondWithError(respWriter, http.StatusInternalServerError, "Error while deleting users")
		return
	}
	respone := respValues{
		HitsReset:    true,
		UsersCleared: true,
	}
	respondWithJSON(respWriter, http.StatusOK, respone)
}
