package config

import (
	"os"
	"path/filepath"

	"aen.it/poolmanager/log"
)

// Define the payment configuration
type PaymentConfiguration struct {
	// Specify minimum duration to be payed
	MinimumDuration int `yaml:"minimumDuration"`
	// Specify cost for any hour
	CostPerHour int `yaml:"costPerHour"`
}

// Define the game station configuraiton
type GameStationConfiguraiton struct {
	// Specify the Game Station's name
	Name string `yaml:"name"`
	// Specify the Game Station's ID
	ID string `yaml:"id"`
	// Specify the Game Station's payment model
	Payment PaymentConfiguration `yaml:"payment"`
}

// Define the item configuraiton
type ItemConfiguraiton struct {
	// Specify the Item's name
	Name string `yaml:"name"`
	// Specify the Item's ID
	ID string `yaml:"id"`
	// Specify the Item's price for the public
	PublicPrice int `yaml:"publicPrice"`
	// Specify the price payed for this Item
	IncomingPrice int `yaml:"incomingPrice"`
}

type billiardRoomConfig struct {
	DefaultPayment PaymentConfiguration       `yaml:"defaultPayment"`
	GamingStations []GameStationConfiguraiton `yaml:"gamingStations"`
	Items          []ItemConfiguraiton        `yaml:"items"`
	Name           string
}

var BilliardRoomConfig *billiardRoomConfig

func init() {
	log.Log.Debug("Entering BilliardRoomConfig init")
	//create new configuragion object
	BilliardRoomConfig = &billiardRoomConfig{}
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "config", "billiardRoomConfig.yml")
	//load configuration from default path
	loadConfig(path, BilliardRoomConfig)
	log.Log.Debug("Exiting BilliardRoomConfig init")
}
