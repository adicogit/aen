package config

import (
	"os"

	"aen.it/poolmanager/log"
	"gopkg.in/yaml.v2"
)

// loadConfig allow to fill configuration obkect with information from file
func loadConfig(configPath string, config any) {
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
}

// ReInitialize allor to reload configuraiton specifying different
func ReInitializeConfig(configPath string, config any) {
	log.Log.Debug("Entering ReInitialize")
	loadConfig(configPath, config)
	log.Log.Debug("Exiting ReInitialize")
}
