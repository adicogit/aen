package main

import (
	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/server"
)

func main() {
	logLevel, err := log.ParseLevel(config.UIConfig.LogLevel)
	if err != nil {
		log.Log.Warn("Cannot set log level to requested velue", "requested value", config.UIConfig.LogLevel)
	} else {
		log.Log.Info("Log level changed to new value", "new value", config.UIConfig.LogLevel)
	}
	log.SetLogLevel(logLevel)
	server := server.Server{}
	server.Start()
}
