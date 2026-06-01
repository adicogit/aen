package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"aen.it/poolmanager/accounting"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	"github.com/gorilla/mux"
)

type checkListResponse struct {
	Open            []payment.Check          `json:"open"`
	Closed          []accounting.ClosedCheck `json:"closed"`
	CurrentIncoming int                      `json:"currentIncoming"`
	CurrentExpected int                      `json:"currentExpected"`
}

func (server *Server) getChecks(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getChecks")
	w.Header().Set("Content-Type", "application/json")

	openChecks, err := server.accountingManager.GetOpenChecks()
	if err != nil {
		log.Log.Error("Failed to get open checks", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	closedChecks, err := server.accountingManager.GetClosedChecks()
	if err != nil {
		log.Log.Error("Failed to get closed checks", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	currentIncoming, err := server.accountingManager.GetCurrentIncoming()
	if err != nil {
		log.Log.Error("Failed to get current incoming", "error", err)
		currentIncoming = 0
	}

	currentExpected, err := server.accountingManager.GetCurrentExpectedIncoming()
	if err != nil {
		log.Log.Error("Failed to get current expected incoming", "error", err)
		currentExpected = 0
	}

	result := checkListResponse{
		Open:            openChecks,
		Closed:          closedChecks,
		CurrentIncoming: currentIncoming,
		CurrentExpected: currentExpected,
	}

	log.Log.Debug("Exiting getChecks", "openCount", len(openChecks), "closedCount", len(closedChecks))
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log.Error("Failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Failed to read request body"})
		return
	}

	var payment int
	// Try parsing body as struct first: {"payment": 100}
	var req struct {
		Payment int `json:"payment"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err == nil && req.Payment > 0 {
		payment = req.Payment
	} else {
		// Fallback to parsing body as raw integer: 100
		if err := json.Unmarshal(bodyBytes, &payment); err != nil || payment <= 0 {
			log.Log.Error("Failed to decode payment amount", "body", string(bodyBytes), "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: "Invalid payment amount"})
			return
		}
	}

	if err := server.accountingManager.PayCheck(checkID, payment); err != nil {
		log.Log.Error("Failed to pay check", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Failed to pay check: %v", err)})
		return
	}
	log.Log.Debug("Exiting payCheck")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: "Check paid successfully"})
}
