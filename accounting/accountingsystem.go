package accounting

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"slices"
	"sync"
	"time"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	persistency "aen.it/poolmanager/perisistency"
)

const (
	ACCOUNTING_DAYS_KEY = "Accounting_Days"
)

type accountingSystem struct {
	mu                   sync.RWMutex
	openChecks           map[string]payment.Check
	currentAccountingDay string
}

var AccountingSystem *accountingSystem

func init() {
	log.Log.Debug("Entering AccountingSystem init")
	log.Log.Info("Create AccountingSystem")
	AccountingSystem = &accountingSystem{}
	log.Log.Debug("Exiting AccountingSystem init")
}

func (accounting *accountingSystem) SetAccountingDay(date string) error {
	log.Log.Debug("Entering SetAccountingDay", "date", date)
	// Check if date is a valid date in the format aaaammdd
	_, err := time.Parse("19701608", date)
	if err != nil {
		error := fmt.Errorf("Invalid date used in SetAccountingDay. Error: %s", err)
		log.Log.Error("Trying to use invalid date as accounting date", "date", date, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return error
	}

	// Add date as available accounting date
	var days []string
	data, err := persistency.BoltDBPersistency.ReadData(ACCOUNTING_DAYS_KEY)
	if err != nil {
		error := fmt.Errorf("Unable to read data from DB. Error: %s", err)
		log.Log.Error("Error in trying to read data from DB", "key", ACCOUNTING_DAYS_KEY, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return error
	}

	err = convertFromByteArray(data, days)
	if err != nil {
		error := fmt.Errorf("Unable to convert data read from DB. Error: %s", err)
		log.Log.Error("Error in trying to convert data read from DB", "data", data, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return error
	}

	if !slices.Contains(days, date) {
		days = append(days, date)
		data, err = convertToByteArray(date)
		if err != nil {
			error := fmt.Errorf("Unable to convert date to byte array. Error: %s", err)
			log.Log.Error("Error in trying to convert data to write to DB", "data", date, "error", error)
			log.Log.Debug("Exiting SetAccountingDay")
			return error
		}

		err = persistency.BoltDBPersistency.WriteData(ACCOUNTING_DAYS_KEY, data)
		if err != nil {
			error := fmt.Errorf("Unable to write data to DB. Error: %s", err)
			log.Log.Error("Error in trying to write to DB", "data", data, "error", error)
			log.Log.Debug("Exiting SetAccountingDay")
			return error
		}
	}

	// Set date as current accounting day
	accounting.currentAccountingDay = date

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
	data, err := persistency.BoltDBPersistency.ReadData(ACCOUNTING_DAYS_KEY)
	if err != nil {
		error := fmt.Errorf("Unable to read data from DB. Error: %s", err)
		log.Log.Error("Error in trying to read data from DB", "key", ACCOUNTING_DAYS_KEY, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return result, error
	}

	err = convertFromByteArray(data, result)
	if err != nil {
		error := fmt.Errorf("Unable to convert data read from DB. Error: %s", err)
		log.Log.Error("Error in trying to convert data read from DB", "data", data, "error", error)
		log.Log.Debug("Exiting SetAccountingDay")
		return result, error
	}

	log.Log.Debug("Exiting GetAccountingDays", "result", result)
	return result, nil
}

func convertToByteArray(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&data)
	if err != nil {
		log.Log.Error("Unable to convert data into byte array", "data", data, "error", err)
	}
	return buffer.Bytes(), err
}

func convertFromByteArray(dataToRead []byte, data interface{}) error {
	buffer := bytes.NewBuffer(dataToRead)

	dec := gob.NewDecoder(buffer)
	err := dec.Decode(&data)
	if err != nil {
		log.Log.Error("Unable to convert data into object", "data", dataToRead, "error", err)
	}
	return err
}
