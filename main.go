package main

import (
	"os"
	"os/signal"
	"syscall"

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

	srv := server.Server{}

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Log.Error("Server stopped with error", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-sigs
	log.Log.Info("Received signal to terminate", "signal", sig)
	srv.Shutdown()
	log.Log.Info("Application terminated gracefully")
}
