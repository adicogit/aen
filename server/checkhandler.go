package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	"github.com/gorilla/mux"
)

type checkIDList struct {
	ID []string `json:"id"`
}

func (server *Server) getChecks(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getChecks")
	w.Header().Set("Content-Type", "application/json")
	list, err := server.accountingManager.GetOpenCheckIDs()
	if err != nil {
		log.Log.Error("Failed to get open checks", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}
	result := checkIDList{
		ID: list,
	}
	log.Log.Debug("Exiting getChecks")
	json.NewEncoder(w).Encode(result)
}

func (server *Server) addCheck(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering addCheck")
	w.Header().Set("Content-Type", "application/json")
	var check payment.Check
	if err := json.NewDecoder(r.Body).Decode(&check); err != nil {
		log.Log.Error("Failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	if err := server.accountingManager.AddCheck(check); err != nil {
		log.Log.Error("Failed to add check", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Failed to add check: %v", err)})
		return
	}

	log.Log.Debug("Exiting addCheck")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse{Message: "Item created successfully"})
}

func (server *Server) getCheck(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getCheck")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	checkID := vars["checkID"]
	if checkID == "" {
		log.Log.Error("Check ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Check ID is required"})
		return
	}
	check, err := server.accountingManager.GetCheck(checkID)
	if err != nil {
		log.Log.Error("Failed to get check", "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}
	log.Log.Debug("Exiting getCheck")
	json.NewEncoder(w).Encode(check)
}

func (server *Server) payCheck(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering payCheck")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	checkID := vars["checkID"]
	if checkID == "" {
		log.Log.Error("Check ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Check ID is required"})
		return
	}
	if err := server.accountingManager.PayCheck(checkID); err != nil {
		log.Log.Error("Failed to pay check", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Failed to pay check: %v", err)})
		return
	}
	log.Log.Debug("Exiting payCheck")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: "Check paid successfully"})
}
