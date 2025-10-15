package payment

import "aen.it/poolmanager/warehouse"

// Represents the payment status
type PaymentStatus int

// Enumeration for the payment status
const (
	Started PaymentStatus = iota
	Stopped
	Suspended
)

// Define the payment configuration
type PaymentConfiguration struct {
	// Specify minimum duration to be payed
	MinimumDuration int `json:"minimumDuration"`
	// Specify cost for any hour
	CostPerHour int `json:"costPerHour"`
}

// define information for the Check related to current payment
type Check struct {
	// Duration for the current check
	Duration int
	// Price for the current check
	Price int
	// List of cosnumed items
	ItemList []warehouse.Item
}

// Payment interface
type Payment interface {
	// Set the confguration for the payment system
	ConfigurePayment(config PaymentConfiguration)
	// Start new payment by starting counting the time. It returns an error if the payment is neither in a stopped nor in suspended status
	StartCountingPayment() error
	// Stop current payment and calculate the billing. It returns an error if the payment is neither in a started nor in suspended status
	ClosePayment() error
	// Pause current payment, to calculate billing it must be stopped. It returns an error if the payment is not in started status
	PauseCountingPayment() error
	// Return the bill calculated for the payment. It returns an error if the payment has not been closed
	GetCheck() (Check, error)
	// Return the payment status
	GetPaymentStatus() PaymentStatus
	// Add a consumprion to current payment. If payment is n stoppes status it reruns error
	AddConsumption(item warehouse.Item) error
}
