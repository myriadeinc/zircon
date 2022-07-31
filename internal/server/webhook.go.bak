package server

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"fmt"
)

func NewJobsPushWebhook(w http.ResponseWriter, req *http.Request) {
	sessions := GetSessionHandler()
	sessions.BroadcastNewJobs()
	fmt.Fprintf(w, "OK\n")
}
func ListenWebhook() {
	endpoint := "0.0.0.0:4990"
	path := "/new"
	http.HandleFunc(path, NewJobsPushWebhook)
	log.Info().Msgf("Listening on %s%s", endpoint, path)
	http.ListenAndServe(endpoint, nil)
}
