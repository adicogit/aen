package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"aen.it/poolmanager/billiardroom"
	"aen.it/poolmanager/log"
)

// Server is the implementaiton of the http server
type Server struct {
	billiardManager billiardroom.BilliardRoom
	mu              sync.Mutex
}

func commonHeader(next http.Handler) http.Handler {
	log.Log.Debug("Entering Server.commonHeader")
	log.Log.Debug("Exiting Server.commonHeader")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Security-Policy", "default-src 'self';")

		next.ServeHTTP(w, r)
	})
}

// Start starts the http server
func (server *Server) Start() error {
	log.Log.Debug("Entering Server.commonHeader")
	router := mux.NewRouter().StrictSlash(true)

	// Set common header for all requests
	router.Use(commonHeader)

	log.Log.Debug("Exiting Server.commonHeader")
	return nil
}
