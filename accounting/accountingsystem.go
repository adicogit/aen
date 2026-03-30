package accounting

import (
	"sync"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
)

type accountingSystem struct {
	mu         sync.RWMutex
	openChecks map[string]payment.Check
}

var AccountingSystem *accountingSystem

func init() {
	log.Log.Debug("Entering AccountingSystem init")
	log.Log.Info("Create AccountingSystem")
	AccountingSystem = &accountingSystem{}
	log.Log.Debug("Exiting AccountingSystem init")
}
