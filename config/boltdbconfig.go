package config

import (
	"os"
	"path/filepath"

	"aen.it/poolmanager/log"
)

type boltDBConfig struct {
	path   string
	DBpath string `yaml:"dbpath"`
}

var BoltDBConfig *boltDBConfig

func init() {
	log.Log.Debug("Entering BoltDBConfig init")
	//create new UIconfiguragion object
	BoltDBConfig = &boltDBConfig{}
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "config", "boltdbconfig.yml")
	BoltDBConfig.SetConfigFilePath(path)
	//load configuration from default path
	BoltDBConfig.LoadConfig()
	log.Log.Debug("Exiting AccountingConfig init")
}

// load configuration from its own config file
func (conf *boltDBConfig) LoadConfig() error {
	log.Log.Debug("Entering LoadConfig")
	err := loadConfigFomFile(conf.path, conf)
	log.Log.Debug("Exiting LoadConfig", "errror", err)
	return err
}

// seve configuraiton to its own config file
func (conf *boltDBConfig) SaveConfig() error {
	log.Log.Debug("Entering SaveConfig")
	err := saveConfig(conf.path, conf)
	log.Log.Debug("Exiting SaveConfig", "errror", err)
	return nil
}

// Set config file path for this configuration
func (conf *boltDBConfig) SetConfigFilePath(configFilePath string) {
	log.Log.Debug("Entering SetConfigFilePath", "configFilePath", configFilePath)
	conf.path = configFilePath
	log.Log.Debug("Exiting SetConfigFilePath")
}

// Get config file path for this configuration
func (conf *boltDBConfig) GetConfigFilePath() string {
	log.Log.Debug("Entering SetConfigFilePath")
	log.Log.Debug("Exiting GetConfigFilePath", "path", conf.path)
	return conf.path
}
