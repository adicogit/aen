package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
	"github.com/gorilla/mux"
)

type stationIDList struct {
	ID []string `json:"id"`
}

type deviceListProp struct {
}

type gamestationProp struct {
	ID         string           `json:"id"`
	IconPath   string           `json:"iconPath"`
	Name       string           `json:"name"`
	Status     int              `json:"status"`
	Cost       int              `json:"cost"`
	DeviceList []deviceListProp `json:"deviceList"`
}

type gamestationStatus struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Cost   int    `json:"cost"`
}

type gamestationAction struct {
	Action string `json:"action"`
}

type gamestationActionResult struct {
	Action string `json:"action"`
	Result string `json:"result"`
}

type gamestationConsumption struct {
	ItemIDs []string `json:"itemIDs"`
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
	result := gamestationProp{
		ID:       gs.GetID(),
		Name:     gs.GetName(),
		Status:   int(gs.GetStatus()),
		Cost:     gs.GetTemporaryCheck().Price,
		IconPath: gs.GetIconPath(),
	}
	log.Log.Debug("Exiting getGameStation", "result", result)
	json.NewEncoder(w).Encode(result)
}

// run an action on the game station's parameters
func (server *Server) actionGameStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["gsID"]

	var request gamestationAction
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Log.Error("Unable to decode request body", "body", r.Body, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Log.Debug("Entering actionGameStation", "id", id, "vars", vars, "action", request.Action)

	gs, err := server.billiardManager.GetGamingStation(id)
	if err != nil {
		log.Log.Error("Unable to get game station", "is", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch action := strings.ToLower(request.Action); action {
	case "start":
		err = gs.StartMatch()
	case "stop":
		err = gs.CloseMatch()
	case "suspend":
		err = gs.PauseMatch()
	default:
		log.Log.Error("Invalid actions has been specified", "action", action)
		http.Error(w, "Invalid actions has been specified", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Log.Error("Unable to run requested action", "action", request.Action, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := gamestationActionResult{
		Action: request.Action,
		Result: "success",
	}
	log.Log.Debug("Exiting actionGameStation", "result", result)
	json.NewEncoder(w).Encode(result)
}

// return the game station's parameters status
func (server *Server) getGameStationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["gsID"]
	log.Log.Debug("Entering getGameStationStatus", "id", id, "vars", vars)

	gs, err := server.billiardManager.GetGamingStation(id)
	if err != nil {
		log.Log.Error("Unable to get game station", "is", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	check := gs.GetTemporaryCheck()
	result := gamestationStatus{
		ID:     gs.GetID(),
		Status: int(gs.GetStatus()),
		Cost:   check.Price,
	}
	log.Log.Debug("Exiting getGameStationStatus", "result", result)
	json.NewEncoder(w).Encode(result)
}

// Add new consumption to the specified ganme station
func (server *Server) addGameStationConsumption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["gsID"]

	var request gamestationConsumption
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Log.Error("Unable to decode request body", "body", r.Body, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Log.Debug("Entering addGameStationConsumption", "gsID", id, "request", request)

	gs, err := server.billiardManager.GetGamingStation(id)
	if err != nil {
		log.Log.Error("Unable to get game station", "is", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Transaction: First validate all items exist before adding any
	items := make([]warehouse.Item, 0, len(request.ItemIDs))
	for _, itemID := range request.ItemIDs {
		item, err := server.billiardManager.GetItem(itemID)
		if err != nil {
			log.Log.Error("Unable to get item from warehouse", "itemID", itemID, "error", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		items = append(items, item)
	}

	// All items validated successfully, now add them all
	for i, item := range items {
		err = gs.AddConsumption(item)
		if err != nil {
			log.Log.Error("Error in adding consumption to game station", "gsID", id, "itemID", request.ItemIDs[i], "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	check := gs.GetTemporaryCheck()
	result := gamestationStatus{
		ID:     gs.GetID(),
		Status: int(gs.GetStatus()),
		Cost:   check.Price,
	}

	log.Log.Debug("Exiting addGameStationConsumption", "result", result)
	json.NewEncoder(w).Encode(result)
}
