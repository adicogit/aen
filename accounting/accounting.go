package accounting

import "aen.it/poolmanager/payment"

type Accounting interface {
	// Get list of already available accounting days
	GetAccountingDays() ([]string, error)
	// Set accounting day: date represents the accounting day.
	SetAccountingDay(date string) error
	// Get accounting day. Returns a string representing the current accounting day or empty string if it is not set
	GetCurrentAccountingDay() string
	// Get current incoming amount
	GetCurrentIncoming() (int, error)
	// Get current expected incoming amount
	GetCurrentExpectedIncoming() (int, error)
	// Close current accounting day
	CloseCurrentAccountingDay() error
	// Get list of open checks' IDs for currently set accounting day. Returns an eror no accounting day has been set.
	GetOpenCheckIDs() ([]string, error)
	// Get required check for currently set accounting day. Returns an eror no accounting day has been set.
	GetCheck(checkID string) (payment.Check, error)
	// Add new check to the list for currently set accounting day. Returns an eror no accounting day has been set.
	AddCheck(check payment.Check) error
	// Pay check for currently set accounting day. Returns an eror no accounting day has been set.
	PayCheck(checkID string, payment int) error
}
