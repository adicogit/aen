package config

import (
	"os"
	"path/filepath"

	"aen.it/poolmanager/log"
)

type uiAspect struct {
	BackgroundImage string `yaml:"background_image,omitempty"`
}

type uiConfig struct {
	PortNumber        int      `yaml:"portNumber,omitempty"`
	CertFile          string   `yaml:"crt_file,omitempty"`
	KeyFile           string   `yaml:"key_file,omitempty"`
	WebsiteDir        string   `yaml:"websiteDir,omitempty"`
	PortalFrontendURL string   `yaml:"portal_frontend_url,omitempty"`
	LogLevel          string   `yaml:"log_level,omitempty"`
	UIAspect          uiAspect `yaml:"ui_aspect,omitempty"`
	path              string
}

var UIConfig *uiConfig

func init() {
	log.Log.Debug("Entering UIConfig init")
	//create new UIconfiguragion object
	UIConfig = &uiConfig{}
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "config", "uiconfig.yml")
	UIConfig.SetConfigFilePath(path)
	//load configuration from default path
	UIConfig.LoadConfig()
	log.Log.Debug("Exiting Config init")
}

// load configuration from its own config file
func (conf *uiConfig) LoadConfig() error {
	log.Log.Debug("Entering LoadConfig")
	err := loadConfigFomFile(conf.path, conf)
	log.Log.Debug("Exiting LoadConfig", "errror", err)
	return err
}

// seve configuraiton to its own config file
func (conf *uiConfig) SaveConfig() error {
	log.Log.Debug("Entering SaveConfig")
	err := saveConfig(conf.path, conf)
	log.Log.Debug("Exiting SaveConfig", "errror", err)
	return nil
}

// Set config file path for this configuration
func (conf *uiConfig) SetConfigFilePath(configFilePath string) {
	log.Log.Debug("Entering SetConfigFilePath", "configFilePath", configFilePath)
	conf.path = configFilePath
	log.Log.Debug("Exiting SetConfigFilePath")
}

// Get config file path for this configuration
func (conf *uiConfig) GetConfigFilePath() string {
	log.Log.Debug("Entering SetConfigFilePath")
	log.Log.Debug("Exiting GetConfigFilePath", "path", conf.path)
	return conf.path
}
