package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"aen.it/poolmanager/log"
)

func init() {
	log.SetLogLevel(slog.LevelDebug)
}

// TestConfig verify that New function works as expected
func TestConfig(t *testing.T) {
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "billiardRoomConfig.yml")
	ReInitializeConfig(path, BilliardRoomConfig)
	name := BilliardRoomConfig.Name
	if len(name) == 0 {
		t.Errorf("Config initialization FAILED. Its name is empty: %s", name)
	}
	costPerHour := BilliardRoomConfig.DefaultPayment.CostPerHour
	if costPerHour != 1200 {
		t.Errorf("Config initialization FAILED. Its cost per hour is wrong: %d", costPerHour)
	}
	stations := BilliardRoomConfig.GamingStations
	if len(stations) != 1 {
		t.Errorf("Config initialization FAILED. Its station list has wrong number of elements: %d", len(stations))
	}
}
