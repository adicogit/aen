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
	log.Log.Debug("Entering StartMatch")
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	log.Log.Info("Turning ON devices associated to ", gp.id, gp.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOn()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Starting counting payment for ", gp.id, gp.name)
		err = gp.payment.StartCountingPayment()
	}
	log.Log.Debug("Exiting StartMatch")
	return err
}

// Pause the match on the current GamingStation
func (gp *GameStation) PauseMatch() error {
	log.Log.Debug("Entering PauseMatch")
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	log.Log.Info("Turning OFF devices associated to ", gp.id, gp.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Pausing counting payment for ", gp.id, gp.name)
		err = gp.payment.PauseCountingPayment()
	}
	log.Log.Debug("Exiting PauseMatch")
	return err
}

// Close the match on the current GamingStation
func (gp *GameStation) CloseMatch() error {
	log.Log.Debug("Entering CloseMatch")
	var err error
	containsError := false
	deviceLen := len(gp.devices)
	log.Log.Info("Turning OFF devices associated to ", gp.id, gp.name)
	for i := 0; i < deviceLen && !containsError; i++ {
		err = gp.devices[i].TurnOff()
		containsError = err != nil
	}
	if !containsError {
		log.Log.Info("Closing counting payment for ", gp.id, gp.name)
		err = gp.payment.ClosePayment()
	}
	log.Log.Debug("Exiting CloseMatch")
	return err
}

// Add a consumption on the current GamingStation
func (gp *GameStation) AddConsumption(item warehouse.Item) error {
	log.Log.Debug("Entering AddConsumption")
	log.Log.Debug("Exiting AddConsumption")
	return gp.payment.AddConsumption(item)
}

// Return the bill calculated for the payment. It returns an error if the payment has not been closed
func (gp *GameStation) GetCheck() (payment.Check, error) {
	log.Log.Debug("Entering GetCheck")
	log.Log.Debug("Exiting GetCheck")
	return gp.payment.GetCheck()
}

// Add new device. If device is already present it will not be added
func (gp *GameStation) AddDevice(device devices.Device) {
	log.Log.Debug("Entering AddConsumption")
	log.Log.Info("Adding new device to this game station", gp.id, gp.name, device.GetID(), device.GetType())
	gp.devices = append(gp.devices, device)
	log.Log.Debug("Exiting AddConsumption")
}

// Return list of devices associated to this GameStation
func (gp *GameStation) GetDevicesList() []devices.Device {
	log.Log.Debug("Entering GetDevicesList")
	list := gp.devices
	log.Log.Debug("Exiting GetDevicesList")
	return list
}

// Set the name for current GameStation
func (gp *GameStation) SetName(name string) {
	log.Log.Debug("Entering SetName")
	gp.name = name
	log.Log.Debug("Exiting SetName")
}

// Get the name for current GameStation
func (gp *GameStation) GetName() string {
	log.Log.Debug("Entering GetName")
	log.Log.Debug("Exiting GetName")
	return gp.name
}
