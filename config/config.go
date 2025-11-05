package config

import (
	"os"
	"path/filepath"

	"aen.it/poolmanager/log"
	"gopkg.in/yaml.v2"
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

type configInfo struct {
	DefaultPayment PaymentConfiguration       `yaml:"defaultPayment"`
	GamingStations []GameStationConfiguraiton `yaml:"gamingStations"`
	Items          []ItemConfiguraiton        `yaml:"items"`
	Name           string
}

var Config *configInfo

func init() {
	log.Log.Debug("Entering Config init")
	//create new configuragion object
	Config = &configInfo{}
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "config", "config.yml")
	//load configuration from default path
	loadConfig(path, Config)
	log.Log.Debug("Exiting Config init")
}

// loadConfig allow to fill configuration obkect with information from file
func loadConfig(configPath string, config *configInfo) {
	log.Log.Debug("Entering loadConfig")
	if len(configPath) == 0 {
		log.Log.Info("Used empty config path")
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Log.Error("Error getting current dir", "error", err)
	}
	log.Log.Debug("Current dir", "dir name", dir)
	// Validate the config path
	s, err := os.Stat(configPath)
	if err == nil && !s.IsDir() {
		// Open config file
		file, err := os.Open(configPath)
		if err == nil {
			defer file.Close()

			// Init new YAML decode
			d := yaml.NewDecoder(file)

			// Start YAML decoding from file
			err = d.Decode(config)
			if err != nil {
				log.Log.Info("Unable to load new config from specified configPath", "configPath", configPath, "error", err)
			} else {
				log.Log.Info("New config succesfully loaded from specified configPath", "configPath", configPath, "config", Config)
			}
		} else {
			log.Log.Info("Unable to open config file specified configPath", "configPath", configPath)
		}
	} else {
		log.Log.Error("Unable to find specified config path", "configPath", configPath, "error", err)
	}

	log.Log.Debug("Exiting loadConfig")
}

// ReInitialize allor to reload configuraiton specifying different
func (config *configInfo) ReInitialize(configPath string) {
	log.Log.Debug("Entering ReInitialize")
	loadConfig(configPath, Config)
	log.Log.Debug("Exiting ReInitialize")
}
