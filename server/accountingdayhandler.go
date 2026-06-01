package server

import (
	"encoding/json"
	"net/http"

	"aen.it/poolmanager/log"
)

type accountingDayResponse struct {
	CurrentAccountingDay string `json:"currentAccountingDay"`
}

type setAccountingDayRequest struct {
	Date string `json:"date"`
}

func (server *Server) getAccountingDay(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getAccountingDay")
	w.Header().Set("Content-Type", "application/json")

	day := server.accountingManager.GetCurrentAccountingDay()
	log.Log.Debug("Exiting getAccountingDay", "currentAccountingDay", day)
	json.NewEncoder(w).Encode(accountingDayResponse{CurrentAccountingDay: day})
}

func (server *Server) setAccountingDay(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering setAccountingDay")
	w.Header().Set("Content-Type", "application/json")

	var req setAccountingDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Log.Error("Failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request body"})
		return
	}

	if err := server.accountingManager.SetAccountingDay(req.Date); err != nil {
		log.Log.Error("Failed to set accounting day", "date", req.Date, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	log.Log.Debug("Exiting setAccountingDay")
	json.NewEncoder(w).Encode(successResponse{Message: "Accounting day set successfully"})
}
