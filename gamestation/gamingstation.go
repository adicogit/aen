package gamestation

import (
	"aen.it/poolmanager/devices"
	"aen.it/poolmanager/payment"
	"aen.it/poolmanager/warehouse"
)

// GamingStation interface
type GamingStation interface {
	// Start the match on the current GamingStation
	StartMatch() error
	// Pause the match on the current GamingStation
	PauseMatch() error
	// Close the match on the current GamingStation
	CloseMatch() error
	// Add a consumption on the current GamingStation
	AddConsumption(item warehouse.Item) error
	// Return the bill calculated for the payment. It returns an error if the payment has not been closed
	GetCheck() (payment.Check, error)
	// Add new device. If device is already present it will not be added
	AddDevice(device devices.Device)
	// Return list of devices associated to this GameStation
	GetDevicesList() []devices.Device
	// Set the name for current GameStation
	SetName(name string)
	// Get the name for current GameStation
	GetName() string
}
