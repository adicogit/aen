package server

import (
	"encoding/json"
	"net/http"

	"aen.it/poolmanager/log"
)

type statioIDList struct {
	ID []string `json:"id,omitempty"`
}

// return the Web UI configuration
func (server *Server) getGameStations(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getGameStations")
	w.Header().Set("Content-Type", "application/json")
	list := server.billiardManager.GetGamingStationIDs()
	result := statioIDList{
		ID: list,
	}

	log.Log.Debug("Exiting getGameStations", "result", result)
	json.NewEncoder(w).Encode(result)
}
