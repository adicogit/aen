package config

import (
	"os"
	"path/filepath"

	"aen.it/poolmanager/log"
)

type uiConfig struct {
	PortNumber        int    `yaml:"portNumber,omitempty"`
	CertFile          string `yaml:"crt_file,omitempty"`
	KeyFile           string `yaml:"key_file,omitempty"`
	WebsiteDir        string `yaml:"websiteDir,omitempty"`
	PortalFrontendURL string `yaml:"portal_frontend_url,omitempty"`
	LogLevel          string `yaml:"log_level,omitempty"`
}

var UIConfig *uiConfig

func init() {
	log.Log.Debug("Entering UIConfig init")
	//create new UIconfiguragion object
	UIConfig = &uiConfig{}
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "config", "uiconfig.yml")
	//load configuration from default path
	loadConfig(path, UIConfig)
	log.Log.Debug("Exiting Config init")
}
