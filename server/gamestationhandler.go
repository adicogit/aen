package server

import (
	"encoding/json"
	"net/http"

	"aen.it/poolmanager/log"
	"github.com/gorilla/mux"
)

type stationIDList struct {
	ID []string `json:"id"`
}

type deviceListProp struct {
}

type gemestationProp struct {
	ID         string           `json:"id"`
	IconPath   string           `json:"iconPath"`
	Name       string           `json:"name"`
	Status     int              `json:"status"`
	DeviceList []deviceListProp `json:"deviceList"`
}

// return the game station's ID list
func (server *Server) getGameStations(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getGameStations")
	w.Header().Set("Content-Type", "application/json")
	list := server.billiardManager.GetGamingStationIDs()
	result := stationIDList{
		ID: list,
	}

	log.Log.Debug("Exiting getGameStations", "result", result)
	json.NewEncoder(w).Encode(result)
}

// return the game station's parameters
func (server *Server) getGameStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["gsID"]
	log.Log.Debug("Entering getGameStation", "id", id, "vars", vars)

	gs, err := server.billiardManager.GetGamingStation(id)
	if err != nil {
		log.Log.Error("Unable to get game station", "is", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result := gemestationProp{
		ID:       gs.GetID(),
		Name:     gs.GetName(),
		Status:   int(gs.GetStatus()),
		IconPath: gs.GetIconPath(),
	}
	log.Log.Debug("Exiting getGameStation", "result", result)
	json.NewEncoder(w).Encode(result)
}
