package gamestation

import (
	"aen.it/poolmanager/config"
	"aen.it/poolmanager/devices"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	"aen.it/poolmanager/warehouse"
)

type GameStation struct {
	name    string
	id      string
	payment payment.GamePayment
	devices []devices.Device
}

func init() {

}

// New function initialize simpleStation passing a name
func New(config config.GameStationConfiguraiton) GameStation {
	log.Log.Debug("Entering New for GameStation")
	log.Log.Debug("Exiting  New for GameStation")
	return GameStation{
		name:    config.Name,
		id:      config.ID,
		payment: payment.New(config.Payment),
		devices: make([]devices.Device, 0),
	}
}

// Start the match on the current GamingStation
func (gp *GameStation) StartMatch() error {
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOn()
		containsError = err != nil
	}
	if !containsError {
		err = gp.payment.StartCountingPayment()
	}
	return err
}

// Pause the match on the current GamingStation
func (gp *GameStation) PauseMatch() error {
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		err = gp.payment.PauseCountingPayment()
	}
	return err
}

// Close the match on the current GamingStation
func (gp *GameStation) CloseMatch() error {
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		err = gp.payment.ClosePayment()
	}
	return err
}

// Add a consumption on the current GamingStation
func (gp *GameStation) AddConsumption(item warehouse.Item) error {
	return gp.payment.AddConsumption(item)
}

// Return the bill calculated for the payment. It returns an error if the payment has not been closed
func (gp *GameStation) GetCheck() (payment.Check, error) {
	return gp.payment.GetCheck()
}

// Add new device. If device is already present it will not be added
func (gp *GameStation) AddDevice(device devices.Device) {
	gp.devices = append(gp.devices, device)
}

/*
	// Return list of devices associated to this GameStation
	GetDevicesList() []devices.Device
	// Set the name for current GameStation
	SetName(name string)
	// Get the name for current GameStation
	GetName() string
*/
