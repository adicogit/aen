package accounting

import (
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	persistency "aen.it/poolmanager/perisistency"
)

const (
	ACCOUNTING_DAYS_KEY  = "Accounting_Days"
	OPENED_CHECKS_KEY    = "Opened_Checks"
	CLOSED_CHECKS_KEY    = "Closed_Checks"
	CURRENT_INCOMING_KEY = "Current_Incoming"
	CURRENT_EXPECTED_KEY = "Current_Expected"
)

type closedCheck struct {
	Check payment.Check
	Payed int
}

type accountingSystem struct {
	mu                   sync.RWMutex
	openChecks           map[string]payment.Check
	closedChecks         map[string]closedCheck
	currentAccountingDay string
	currentIncoming      int
	currentExpected      int
}

var AccountingSystem *accountingSystem

func init() {
	log.Log.Debug("Entering AccountingSystem init")
	log.Log.Info("Create AccountingSystem")
	AccountingSystem = &accountingSystem{}
	log.Log.Debug("Exiting AccountingSystem init")
}

// Set the current accounting day. If it already exists it mus be read from persisted data, ortherwise its status must be initialized
func (accounting *accountingSystem) SetAccountingDay(date string) error {
	log.Log.Debug("Entering SetAccountingDay", "date", date)
	// Check if date is a valid date in the format aaaammdd
	_, err := time.Parse("19701608", date)
	if err != nil {
		error := fmt.Errorf("Invalid date used in SetAccountingDay. Reason: %s", err)
		log.Log.Error("Trying to use invalid date as accounting date", "date", date, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return error
	}

	currentAccountingDay := accounting.currentAccountingDay
	// Add date as available accounting date
	var days []string
	err = persistency.BoltDBPersistency.ReadData(ACCOUNTING_DAYS_KEY, days)
	if err != nil {
		error := fmt.Errorf("Unable to read data from DB. Reason: %s", err)
		log.Log.Error("Error in trying to read data from DB", "key", ACCOUNTING_DAYS_KEY, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return error
	}

	// Set date as current accounting day
	accounting.currentAccountingDay = date

	if !slices.Contains(days, date) {
		days = append(days, date)
		err = persistency.BoltDBPersistency.WriteData(ACCOUNTING_DAYS_KEY, days)
		if err != nil {
			// Restore current accounting day to saved value
			accounting.currentAccountingDay = currentAccountingDay
			error := fmt.Errorf("Unable to persist updated list of accounting days. Reason: %s", err)
			log.Log.Error("Error in trying to persist updated list of accounting days", "key", ACCOUNTING_DAYS_KEY, "data", days, "error", error)
			log.Log.Debug("Exiting SetAccountingDay")
			return error
		}
		accounting.closedChecks = make(map[string]closedCheck)
		accounting.openChecks = make(map[string]payment.Check)
		accounting.currentExpected = 0
		accounting.currentIncoming = 0
	} else {
		err = accounting.restore()
		if err != nil {
			accounting.currentAccountingDay = currentAccountingDay
			error := fmt.Errorf("Unable to read accounting day %s. Reason: %s", accounting.currentAccountingDay, err)
			log.Log.Error("Error in trying to read current accounting day", "accountingDay", date, "error", error)
			log.Log.Debug("Exiting SetAccountingDay")
			return error
		}
	}

	log.Log.Debug("Exiting SetAccountingDay")
	return nil
}

func (accounting *accountingSystem) GetCurrentAccountingDay() string {
	log.Log.Debug("Entering GetCurrentAccountingDay")
	log.Log.Debug("Exiting GetCurrentAccountingDay", "currentAccountingDay", accounting.currentAccountingDay)
	return accounting.currentAccountingDay
}

func (accounting *accountingSystem) GetAccountingDays() ([]string, error) {
	log.Log.Debug("Entering GetAccountingDays")
	var result []string

	// Read data from DB
	err := persistency.BoltDBPersistency.ReadData(ACCOUNTING_DAYS_KEY, result)
	if err != nil {
		error := fmt.Errorf("Unable to read data from DB. Reason: %s", err)
		log.Log.Error("Error in trying to read data from DB", "key", ACCOUNTING_DAYS_KEY, "error", error)
		log.Log.Debug("Exiting GetAccountingDays")
		return result, error
	}

	log.Log.Debug("Exiting GetAccountingDays", "result", result)
	return result, nil
}

func (accounting *accountingSystem) GetCurrentIncoming() (int, error) {
	log.Log.Debug("Entering GetCurrentIncoming")

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before reading any information about it")
		log.Log.Error("Error in getting current incoming", "error", error)
		log.Log.Debug("Exiting GetCurrentIncoming")
		return 0, error
	}

	result := accounting.currentIncoming
	log.Log.Debug("Exiting GetCurrentIncoming", "currentIncoming", result)
	return result, nil
}

func (accounting *accountingSystem) GetCurrentExpectedIncoming() (int, error) {
	log.Log.Debug("Entering GetCurrentExpectedIncoming")

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before reading any information about it")
		log.Log.Error("Error in getting current expected incoming", "error", error)
		log.Log.Debug("Exiting GetCurrentExpectedIncoming")
		return 0, error
	}

	result := accounting.currentExpected
	log.Log.Debug("Exiting GetCurrentExpectedIncoming", "currentExpectedIncoming", result)
	return result, nil
}

func (accounting *accountingSystem) CloseCurrentAccountingDay() error {
	log.Log.Debug("Entering CloseCurrentAccountingDay")

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before performing any operation it")
		log.Log.Error("Error in closing current acconting day", "error", error)
		log.Log.Debug("Exiting GetOpenChCloseCurrentAccountingDayeckIDs")
		return error
	}

	// Store current accounting day
	err := accounting.persist()
	if err != nil {
		error := fmt.Errorf("Unable to persist current accounting. Reason: %s", err)
		log.Log.Error("Error in closing current acconting day", "error", error)
		log.Log.Debug("Exiting CloseCurrentAccountingDay")
		return error
	}

	// Reset current accounting data
	accounting.currentAccountingDay = ""
	accounting.currentExpected = 0
	accounting.currentIncoming = 0
	accounting.closedChecks = make(map[string]closedCheck)
	accounting.openChecks = make(map[string]payment.Check)

	log.Log.Debug("Exiting CloseCurrentAccountingDay")
	return nil
}

func (accounting *accountingSystem) GetOpenCheckIDs() ([]string, error) {
	log.Log.Debug("Entering GetOpenCheckIDs")

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before reading any information about it")
		log.Log.Error("Error in getting list of opened check IDs", "error", error)
		log.Log.Debug("Exiting GetOpenCheckIDs")
		return []string{}, error
	}

	result := slices.Collect(maps.Keys(accounting.openChecks))
	log.Log.Debug("Exiting GetOpenCheckIDs", "result", result)
	return result, nil
}

func (accounting *accountingSystem) GetCheck(checkID string) (payment.Check, error) {
	log.Log.Debug("Entering GetCheck", "checkID", checkID)

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before reading any information about it")
		log.Log.Error("Error in getting specified check", "checkID", checkID, "error", error)
		log.Log.Debug("Exiting GetOpenCheckIDs")
		return payment.Check{}, error
	}

	result, contained := accounting.openChecks[checkID]
	if !contained {
		closedCheck, contained := accounting.closedChecks[checkID]
		if !contained {
			error := fmt.Errorf("Unable to find out requested data. Check ID: %s", checkID)
			log.Log.Error("Error in looking for check IDs", "checkID", checkID, "error", error)
			log.Log.Debug("Exiting GetOpenCheckIDs")
			return payment.Check{}, error
		}
		result = closedCheck.Check
	}

	log.Log.Debug("Exiting GetOpenCheckIDs", "result", result)
	return result, nil
}

func (accounting *accountingSystem) AddCheck(check payment.Check) error {
	log.Log.Debug("Entering AddCheck", "check", check)

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		error := fmt.Errorf("Current Accounting Day must be set before setting any information about it")
		log.Log.Error("Error in adding new check", "check", check, "error", error)
		log.Log.Debug("Exiting AddCheck")
		return error
	}

	// Add received check to list of opened checks
	accounting.openChecks[check.ID] = check

	// persist the new status
	err := accounting.persist()
	if err != nil {
		// Error in serializing the new check. Must remove it from the list of open check aborting the operation
		delete(accounting.openChecks, check.ID)
		error := fmt.Errorf("Unable to add new check with ID %s. Reason: %s", check.ID, err)
		log.Log.Error("Error in adding new check", "check", check, "error", error)
		log.Log.Debug("Exiting AddCheck")
		return error
	}

	log.Log.Debug("Exiting AddCheck")
	return nil
}

func (accounting *accountingSystem) PayCheck(checkID string, payment int) error {
	log.Log.Debug("Entering PayCheck", "checkID", checkID, "payment", payment)

	// Store current status to restore it in case of errors
	originalCurrentExpected := accounting.currentExpected
	originalCurrentIncoming := accounting.currentIncoming

	check := accounting.openChecks[checkID]

	// Update current expected cash flow
	accounting.currentExpected += check.Price

	// Update current cash flow
	accounting.currentIncoming += payment

	// Update list of opened checks
	delete(accounting.openChecks, checkID)

	// update list of closed checks
	closedCheck := closedCheck{
		Check: check,
		Payed: payment,
	}
	accounting.closedChecks[checkID] = closedCheck

	err := accounting.persist()
	if err != nil {
		// Restore updated information
		accounting.currentExpected = originalCurrentExpected
		accounting.currentIncoming = originalCurrentIncoming
		accounting.openChecks[checkID] = check
		delete(accounting.closedChecks, checkID)

		error := fmt.Errorf("Unable to persist updated current accounting day. Reason: %s", err)
		log.Log.Error("Error in paying check", "checkID", checkID, "payment", payment, "error", error)
		log.Log.Debug("Exiting PayCheck")
	}

	log.Log.Debug("Exiting PayCheck")
	return nil
}

func (accounting *accountingSystem) restore() error {
	log.Log.Debug("Entering restore")
	var operationError error

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		operationError = fmt.Errorf("Current Accounting Day must be set before restoring it")
		log.Log.Error("Error in restoring current accounting", "error", operationError)
	} else {
		key := ACCOUNTING_DAYS_KEY + persistency.BUCKET_SEPARATOR + accounting.currentAccountingDay
		err := persistency.BoltDBPersistency.ReadData(key, accounting)
		if err != nil {
			operationError = fmt.Errorf("Unable to read current accounting. Reason: %s", err)
			log.Log.Error("Error in restoring current accounting", "error", operationError)
		}
	}

	log.Log.Debug("Exiting restore")
	return operationError
}

func (accounting *accountingSystem) persist() error {
	log.Log.Debug("Entering persist")
	var operationError error

	// Check if currentAccountingDay has been set
	if len(accounting.currentAccountingDay) == 0 {
		operationError = fmt.Errorf("Current Accounting Day must be set before persisting it")
		log.Log.Error("Error in persisting current accounting", "error", operationError)
	} else {
		key := ACCOUNTING_DAYS_KEY + persistency.BUCKET_SEPARATOR + accounting.currentAccountingDay
		err := persistency.BoltDBPersistency.WriteData(key, accounting)
		if err != nil {
			operationError = fmt.Errorf("Unable to persist current accounting. Reason: %s", err)
			log.Log.Error("Error in persisting current accounting", "error", operationError)
		}
	}

	log.Log.Debug("Exiting persist")
	return operationError
}
