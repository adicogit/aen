package accounting

import (
	"fmt"
	"sync"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	persistency "aen.it/poolmanager/perisistency"
)

type accountingSystem struct {
	mu                   sync.RWMutex
	db                   persistency.Persistency
	currentAccountingDay string
	openChecks           map[string]payment.Check
}

var AccountingSystem *accountingSystem

func init() {
	log.Log.Debug("Entering AccountingSystem init")
	log.Log.Info("Create AccountingSystem")
	AccountingSystem = &accountingSystem{}
	AccountingSystem.db = persistency.BoltDBPersistency
	AccountingSystem.openChecks = make(map[string]payment.Check)
	log.Log.Debug("Exiting AccountingSystem init")
}

func (system *accountingSystem) SetAccountingDay(date string) error {
	log.Log.Debug("Entering SetAccountingDay", "day", date)
	system.mu.Lock()
	defer system.mu.Unlock()
	system.currentAccountingDay = date
	log.Log.Debug("Exiting SetAccountingDay")
	return nil
}

func (system *accountingSystem) GetCurrentAccountingDay() string {
	log.Log.Debug("Entering GetCurrentAccountingDay")
	system.mu.RLock()
	defer system.mu.RUnlock()
	log.Log.Debug("Exiting GetCurrentAccountingDay")
	return system.currentAccountingDay
}

func (system *accountingSystem) GetAccountingDays() []string {
	log.Log.Debug("Entering GetAccountingDays")
	// For now return empty, we can implement it later with DB
	log.Log.Debug("Exiting GetAccountingDays")
	return []string{}
}

func (system *accountingSystem) CloseCurrentAccountingDay() error {
	log.Log.Debug("Entering CloseCurrentAccountingDay")
	system.mu.Lock()
	defer system.mu.Unlock()
	system.currentAccountingDay = ""
	// Should probably clear openChecks too or save them
	system.openChecks = make(map[string]payment.Check)
	log.Log.Debug("Exiting CloseCurrentAccountingDay")
	return nil
}

func (system *accountingSystem) GetOpenCheckIDs() ([]string, error) {
	log.Log.Debug("Entering GetOpenCheckIDs")
	system.mu.RLock()
	defer system.mu.RUnlock()
	if system.currentAccountingDay == "" {
		return nil, fmt.Errorf("no accounting day set")
	}
	ids := make([]string, 0, len(system.openChecks))
	for id := range system.openChecks {
		ids = append(ids, id)
	}
	log.Log.Debug("Exiting GetOpenCheckIDs")
	return ids, nil
}

func (system *accountingSystem) GetCheck(checkID string) (payment.Check, error) {
	log.Log.Debug("Entering GetCheck", "checkID", checkID)
	system.mu.RLock()
	defer system.mu.RUnlock()
	if system.currentAccountingDay == "" {
		return payment.Check{}, fmt.Errorf("no accounting day set")
	}
	check, ok := system.openChecks[checkID]
	if !ok {
		return payment.Check{}, fmt.Errorf("check with ID %s not found", checkID)
	}
	log.Log.Debug("Exiting GetCheck")
	return check, nil
}

func (system *accountingSystem) AddCheck(check payment.Check) error {
	log.Log.Debug("Entering AddCheck", "checkID", check.ID)
	system.mu.Lock()
	defer system.mu.Unlock()
	if system.currentAccountingDay == "" {
		return fmt.Errorf("no accounting day set")
	}
	if _, ok := system.openChecks[check.ID]; ok {
		return fmt.Errorf("check with ID %s already exists", check.ID)
	}
	system.openChecks[check.ID] = check
	log.Log.Debug("Exiting AddCheck")
	return nil
}

func (system *accountingSystem) PayCheck(checkID string) error {
	log.Log.Debug("Entering PayCheck", "checkID", checkID)
	system.mu.Lock()
	defer system.mu.Unlock()
	if system.currentAccountingDay == "" {
		return fmt.Errorf("no accounting day set")
	}
	if _, ok := system.openChecks[checkID]; !ok {
		return fmt.Errorf("check with ID %s not found", checkID)
	}
	delete(system.openChecks, checkID)
	log.Log.Debug("Exiting PayCheck")
	return nil
}

func (system *accountingSystem) Close() {
	log.Log.Debug("Entering accountingSystem.Close")
	system.mu.Lock()
	defer system.mu.Unlock()
	system.db.CloseDB()
	log.Log.Debug("Exiting accountingSystem.Close")
}
