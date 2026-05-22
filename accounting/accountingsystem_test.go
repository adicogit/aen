package accounting

import (
	"os"
	"testing"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/payment"
	"aen.it/poolmanager/perisistency"
)

func TestAccountingSystem(t *testing.T) {
	// Setup temporary DB
	tmpDB := "/tmp/test_accounting.db"
	config.BoltDBConfig.DBpath = tmpDB
	defer os.Remove(tmpDB)

	// Ensure DB is closed after test
	defer persistency.BoltDBPersistency.CloseDB()

	// Initialize the system (it's already initialized by init, but we might want to reset it)
	AccountingSystem.CurrentAccountingDay = ""
	AccountingSystem.OpenChecks = make(map[string]payment.Check)
	AccountingSystem.ClosedChecks = make(map[string]closedCheck)
	AccountingSystem.CurrentIncoming = 0
	AccountingSystem.CurrentExpected = 0

	date := "20231027"

	// Test SetAccountingDay
	err := AccountingSystem.SetAccountingDay(date)
	if err != nil {
		t.Fatalf("Failed to set accounting day: %v", err)
	}

	if AccountingSystem.GetCurrentAccountingDay() != date {
		t.Errorf("Expected accounting day %s, got %s", date, AccountingSystem.GetCurrentAccountingDay())
	}

	// Test AddCheck
	check := payment.Check{
		ID:    "check1",
		Price: 100,
	}
	err = AccountingSystem.AddCheck(check)
	if err != nil {
		t.Fatalf("Failed to add check: %v", err)
	}

	ids, err := AccountingSystem.GetOpenCheckIDs()
	if err != nil {
		t.Fatalf("Failed to get open check IDs: %v", err)
	}
	if len(ids) != 1 || ids[0] != "check1" {
		t.Errorf("Expected open checks [check1], got %v", ids)
	}

	// Test GetCheck
	retrieved, err := AccountingSystem.GetCheck("check1")
	if err != nil {
		t.Fatalf("Failed to get check: %v", err)
	}
	if retrieved.ID != "check1" || retrieved.Price != 100 {
		t.Errorf("Retrieved check mismatch: %+v", retrieved)
	}

	// Test PayCheck
	err = AccountingSystem.PayCheck("check1", 100)
	if err != nil {
		t.Fatalf("Failed to pay check: %v", err)
	}

	incoming, _ := AccountingSystem.GetCurrentIncoming()
	if incoming != 100 {
		t.Errorf("Expected incoming 100, got %d", incoming)
	}

	expected, _ := AccountingSystem.GetCurrentExpectedIncoming()
	if expected != 100 {
		t.Errorf("Expected expected incoming 100, got %d", expected)
	}

	ids, _ = AccountingSystem.GetOpenCheckIDs()
	if len(ids) != 0 {
		t.Errorf("Expected 0 open checks, got %d", len(ids))
	}

	// Test CloseCurrentAccountingDay
	err = AccountingSystem.CloseCurrentAccountingDay()
	if err != nil {
		t.Fatalf("Failed to close accounting day: %v", err)
	}

	if AccountingSystem.GetCurrentAccountingDay() != "" {
		t.Errorf("Expected empty accounting day after close, got %s", AccountingSystem.GetCurrentAccountingDay())
	}

	// Test Restore (SetAccountingDay again with same date)
	err = AccountingSystem.SetAccountingDay(date)
	if err != nil {
		t.Fatalf("Failed to Re-set accounting day: %v", err)
	}

	incoming, _ = AccountingSystem.GetCurrentIncoming()
	if incoming != 100 {
		t.Errorf("Expected restored incoming 100, got %d", incoming)
	}
}
