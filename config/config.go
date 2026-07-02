package config

import (
	"os"

	"aen.it/poolmanager/log"
	"gopkg.in/yaml.v2"
)

type Config interface {
	// load configuration from its own config file
	LoadConfig() error
	// seve configuraiton to its own config file
	SaveConfig() error
	// Set config file path for this configuration
	SetConfigFilePath(configFilePath string)
	// Get config file path for this configuration
	GetConfigFilePath() string
}

// loadConfig allow to fill configuration obkect with information from file
func loadConfigFomFile(configPath string, config Config) error {
	var err error
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
				log.Log.Info("New config succesfully loaded from specified configPath", "configPath", configPath, "config", BilliardRoomConfig)
			}
		} else {
			log.Log.Info("Unable to open config file specified configPath", "configPath", configPath)
		}
	} else {
		log.Log.Error("Unable to find specified config path", "configPath", configPath, "error", err)
	}

	log.Log.Debug("Exiting loadConfig")
	return err
}

// SaveConfig allow to persist current configuration
func saveConfig(configPath string, config Config) error {
	var err error
	log.Log.Debug("Entering saveConfig")
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
		file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err == nil {
			defer file.Close()

			// Init new YAML decode
			d := yaml.NewEncoder(file)

			// Start YAML decoding from file
			err = d.Encode(config)
			if err != nil {
				log.Log.Info("Unable to save current configuration to specified configPath", "configPath", configPath, "error", err)
			} else {
				log.Log.Info("Current configuration succesfully saved to specified configPath", "configPath", configPath, "config", BilliardRoomConfig)
			}
		} else {
			log.Log.Info("Unable to open config file specified configPath", "configPath", configPath)
		}
	} else {
		log.Log.Error("Unable to find specified config path", "configPath", configPath, "error", err)
	}

	log.Log.Debug("Exiting saveConfig")
	return err
}
