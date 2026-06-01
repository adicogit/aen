package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"aen.it/poolmanager/accounting"
	"aen.it/poolmanager/billiardroom"
	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
)

// Server is the implementaiton of the http server
type Server struct {
	billiardManager   billiardroom.BilliardRoom
	accountingManager accounting.Accounting
	mu                sync.Mutex
}

func commonHeader(next http.Handler) http.Handler {
	log.Log.Debug("Entering Server.commonHeader")
	log.Log.Debug("Exiting Server.commonHeader")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	})
}

// Start starts the http server
func (server *Server) Start() error {
	log.Log.Debug("Entering Server.Start")
	var err error
	server.billiardManager = billiardroom.Manager
	server.accountingManager = accounting.AccountingSystem

	if len(server.accountingManager.GetCurrentAccountingDay()) == 0 {
		today := time.Now().Format("20060102")
		if err := server.accountingManager.SetAccountingDay(today); err != nil {
			log.Log.Error("Failed to initialize accounting day to today's date", "date", today, "error", err)
		} else {
			log.Log.Info("Successfully initialized accounting day to today's date", "date", today)
		}
	}

	router := mux.NewRouter().StrictSlash(true)

	// API to handle UI configuration
	router.HandleFunc("/api/v1/uiconfig", server.getUIConfig).Methods("GET")
	router.HandleFunc("/api/v1/uiconfig", server.setUIConfig).Methods("PUT")

	// API to handle game stations
	router.HandleFunc("/api/v1/gamestations", server.getGameStations).Methods("GET")
	//router.HandleFunc("/api/v1/gamestations", server.createGameStation).Methods("POST")
	router.HandleFunc("/api/v1/gamestations/{gsID}", server.getGameStation).Methods("GET")
	router.HandleFunc("/api/v1/gamestations/{gsID}/action", server.actionGameStation).Methods("POST")
	router.HandleFunc("/api/v1/gamestations/{gsID}/status", server.getGameStationStatus).Methods("GET")
	router.HandleFunc("/api/v1/gamestations/{gsID}/consumption", server.addGameStationConsumption).Methods("POST")

	// API to handle accounting day
	router.HandleFunc("/api/v1/accountingday", server.getAccountingDay).Methods("GET")
	router.HandleFunc("/api/v1/accountingday", server.setAccountingDay).Methods("PUT")

	// API to handle checks
	router.HandleFunc("/api/v1/checks", server.getChecks).Methods("GET")
	router.HandleFunc("/api/v1/checks", server.addCheck).Methods("POST")
	router.HandleFunc("/api/v1/checks/{checkID}", server.getCheck).Methods("GET")
	router.HandleFunc("/api/v1/checks/{checkID}", server.payCheck).Methods("PUT")

	// API to handle warehouse items
	router.HandleFunc("/api/v1/warehouseitems", server.getWarehouseItems).Methods("GET")
	router.HandleFunc("/api/v1/warehouseitems", server.createWarehouseItems).Methods("POST")
	router.HandleFunc("/api/v1/warehouseitems/{itemID}", server.getWarehouseItem).Methods("GET")
	router.HandleFunc("/api/v1/warehouseitems/{itemID}", server.modifyWarehouseItems).Methods("PUT")
	router.HandleFunc("/api/v1/warehouseitems/{itemID}", server.deleteWarehouseItems).Methods("DELETE")

	WebuiHandler.SetStaticPath(config.UIConfig.WebsiteDir)
	router.PathPrefix("/").Handler(WebuiHandler)

	// Set common header for all requests
	router.Use(commonHeader)

	portNumber := fmt.Sprintf(":%d", config.UIConfig.PortNumber)
	if len(config.UIConfig.CertFile) == 0 {
		log.Log.Debug("Start listening http", "port number", config.UIConfig.PortNumber)
		err = http.ListenAndServe(portNumber, router)
	} else {
		log.Log.Debug("Start listening https", "port number", config.UIConfig.PortNumber)
		err = http.ListenAndServeTLS(portNumber, config.UIConfig.CertFile, config.UIConfig.KeyFile, router)
	}
	log.Log.Debug("Listen and Serve passed with", "result", err)
	server.Shutdown()
	log.Log.Debug("Exiting Server.Start")
	return err
}

func (server *Server) Shutdown() {
	log.Log.Info("Entering Server.Shutdown")
	//server.accountingManager.Close()
	log.Log.Info("Exiting Server.Shutdown")
}
