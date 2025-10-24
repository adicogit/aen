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
func (gs *GameStation) StartMatch() error {
	log.Log.Debug("Entering StartMatch")
	var err error
	containsError := false
	deviceLen := len(gs.devices)
	log.Log.Info("Turning ON devices associated to ", gs.id, gs.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gs.devices[i].TurnOn()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Starting counting payment for ", gs.id, gs.name)
		err = gs.payment.StartCountingPayment()
	}
	log.Log.Debug("Exiting StartMatch")
	return err
}

// Pause the match on the current GamingStation
func (gs *GameStation) PauseMatch() error {
	log.Log.Debug("Entering PauseMatch")
	var err error
	containsError := false
	deviceLen := len(gs.devices)
	log.Log.Info("Turning OFF devices associated to ", gs.id, gs.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gs.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Pausing counting payment for ", gs.id, gs.name)
		err = gs.payment.PauseCountingPayment()
	}
	log.Log.Debug("Exiting PauseMatch")
	return err
}

// Close the match on the current GamingStation
func (gs *GameStation) CloseMatch() error {
	log.Log.Debug("Entering CloseMatch")
	var err error
	containsError := false
	deviceLen := len(gs.devices)
	log.Log.Info("Turning OFF devices associated to ", gs.id, gs.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gs.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Closing counting payment for ", gs.id, gs.name)
		err = gs.payment.ClosePayment()
	}
	log.Log.Debug("Exiting CloseMatch")
	return err
}

// Add a consumption on the current GamingStation
func (gs *GameStation) AddConsumption(item warehouse.Item) error {
	log.Log.Debug("Entering AddConsumption")
	log.Log.Debug("Exiting AddConsumption")
	return gs.payment.AddConsumption(item)
}

// Return the bill calculated for the payment. It returns an error if the payment has not been closed
func (gs *GameStation) GetCheck() (payment.Check, error) {
	log.Log.Debug("Entering GetCheck")
	log.Log.Debug("Exiting GetCheck")
	return gs.payment.GetCheck()
}

// Add new device. If device is already present it will not be added
func (gs *GameStation) AddDevice(device devices.Device) {
	log.Log.Debug("Entering AddConsumption")
	log.Log.Info("Adding new device to this game station", gs.id, gs.name, device.GetID(), device.GetType())
	gs.devices = append(gs.devices, device)
	log.Log.Debug("Exiting AddConsumption")
}

// Return list of devices associated to this GameStation
func (gs *GameStation) GetDevicesList() []devices.Device {
	log.Log.Debug("Entering GetDevicesList")
	list := gs.devices
	log.Log.Debug("Exiting GetDevicesList")
	return list
}

// Set the name for current GameStation
func (gs *GameStation) SetName(name string) {
	log.Log.Debug("Entering SetName")
	gs.name = name
	log.Log.Debug("Exiting SetName")
}

// Get the name for current GameStation
func (gs *GameStation) GetName() string {
	log.Log.Debug("Entering GetName")
	log.Log.Debug("Exiting GetName")
	return gs.name
}
