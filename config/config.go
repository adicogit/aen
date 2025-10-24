package config

import "aen.it/poolmanager/log"

// Define the payment configuration
type PaymentConfiguration struct {
	// Specify minimum duration to be payed
	MinimumDuration int `json:"minimumDuration"`
	// Specify cost for any hour
	CostPerHour int `json:"costPerHour"`
}

// Define the game station configuraiton
type GameStationConfiguraiton struct {
	Name    string               `json:"name"`
	ID      string               `json:"id"`
	Payment PaymentConfiguration `json:"payment"`
}

// Define the item configuraiton
type ItemConfiguraiton struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	PublicPrice   int    `json:"publicPrice"`
	IncomingPrice int    `json:"incomingPrice"`
}

type configInfo struct {
	DefaultPayment PaymentConfiguration       `json:"defaultPayment"`
	GamingStations []GameStationConfiguraiton `json:"gamingStations"`
	Items          []ItemConfiguraiton        `json:"items"`
	Name           string
}

var Config *configInfo

func init() {
	log.Log.Debug("Entering Config init")
	//create new configuragion object
	Config = &configInfo{}
	//load configuration from default path
	loadConfig("/etc/baas/baas.yml", Config)
	log.Log.Debug("Exiting Config init")
}

//loadConfig allow to fill configuration obkect with information from file
func loadConfig(configPath string, config *configInfo) {
	log.Log.Debug("Entering loadConfig")
	if len(configPath) == 0 {
		log.Log.Info("Used empty config path")
	}
	config.Name = "Prova"
	defaultPayment := PaymentConfiguration{
		CostPerHour:     500,
		MinimumDuration: 30,
	}
	config.DefaultPayment = defaultPayment
	config.GamingStations = make([]GameStationConfiguraiton, 1)

	currentGamingStation := GameStationConfiguraiton{}
	currentGamingStation.Name = "Postazione 1"
	currentGamingStation.ID = "1"
	currentGamingStation.Payment = defaultPayment
	config.GamingStations[0] = currentGamingStation

	config.Items = make([]ItemConfiguraiton, 1)

	currentItem := ItemConfiguraiton{}
	currentItem.Name = "Acqua"
	currentItem.ID = "1"
	currentItem.IncomingPrice = 50
	currentItem.PublicPrice = 100
	config.Items[0] = currentItem

	log.Log.Debug("Exiting loadConfig")
}

//ReInitialize allor to reload configuraiton specifying different
func (config *configInfo) ReInitialize(configPath string) {
	log.Log.Debug("Entering ReInitialize")
	loadConfig(configPath, config)
	log.Log.Debug("Exiting ReInitialize")
}
