package billiardroom

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/gamestation"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
)

func init() {
	log.SetLogLevel(slog.LevelDebug)
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "..", "config", "config.yml")
	config.Config.ReInitialize(path)
	Manager.loadFromConfig()
}

// TestBilliardRoom verify that billiard room works as expected
func TestBilliardRoom(t *testing.T) {
	numOfStation := Manager.GetNumberOfGamingStation()
	if numOfStation != 1 {
		t.Errorf("Billiard room manager initialization FAILED. Its station list has wrong number of elements: %d", numOfStation)
	}
	stationsIDs := Manager.GetGamingStationIDs()
	if stationsIDs[0] != "1" {
		t.Errorf("Billiard room manager initialization FAILED. Its first station is: %s", stationsIDs[0])
	}
	station, err := Manager.GetGamingStation(stationsIDs[0])
	if err != nil {
		t.Errorf("Billiard room manager initialization FAILED. Getting Station %s received the error: %s", stationsIDs[0], err)
	}
	if station.GetName() != "Postazione 1" {
		t.Errorf("Billiard room manager initialization FAILED. Expected station's name %s, but got %s", "Postazione 1", station.GetName())
	}
	numOfItems := Manager.GetNumberOfItems()
	if numOfItems != 1 {
		t.Errorf("Billiard room manager initialization FAILED. Its station list has wrong number of elements: %d", numOfItems)
	}
	itemIDs := Manager.GetItemIDs()
	if itemIDs[0] != "1" {
		t.Errorf("Billiard room manager initialization FAILED. Its first station is: %s", itemIDs[0])
	}
	item, err := Manager.GetItem(itemIDs[0])
	if err != nil {
		t.Errorf("Billiard room manager initialization FAILED. Getting Station %s received the error: %s", itemIDs[0], err)
	}
	if item.Name != "Acqua" {
		t.Errorf("Billiard room manager initialization FAILED. Expected station's name %s, but got %s", "Acqua", item.Name)
	}
}

// TestBilliardRoomAddItem verify that billiard room works as expected
func TestBilliardRoomAddItem(t *testing.T) {
	numOfItems := Manager.GetNumberOfItems()
	item := warehouse.Item{
		ID:            "1234567890",
		Name:          "CocaCola",
		PublicPrice:   350,
		IncomingPrice: 150,
	}
	Manager.AddItem(item)

	newNumberOfItems := Manager.GetNumberOfItems()
	if newNumberOfItems != numOfItems+1 {
		t.Errorf("Billiard room manager AddItem FAILED. Expected %d number of imtems, but got %d", numOfItems+1, newNumberOfItems)
	}
	newItem, err := Manager.GetItem(item.ID)
	if err != nil {
		t.Errorf("Billiard room manager AddItem FAILED. Got error retireving item with ID %s: %s", item.ID, err)
	}
	if newItem.ID != item.ID {
		t.Errorf("Billiard room manager AddItem FAILED. Expected item's ID %s, but got %s", item.ID, newItem.ID)
	}
}

// TestBilliardRoomAddStation verify that billiard room addStation
func TestBilliardRoomAddStation(t *testing.T) {
	numOfstations := Manager.GetNumberOfGamingStation()

	stationConfig := config.GameStationConfiguraiton{
		ID:      "test",
		Name:    "Postazione di test",
		Payment: config.Config.DefaultPayment,
	}
	station := gamestation.New(stationConfig)
	Manager.AddGamingStation(&station)

	newNumOfstations := Manager.GetNumberOfGamingStation()
	if newNumOfstations != numOfstations+1 {
		t.Errorf("Billiard room manager addStation FAILED. Expected %d number of imtems, but got %d", numOfstations+1, newNumOfstations)
	}
	newStation, err := Manager.GetGamingStation(station.GetID())
	if err != nil {
		t.Errorf("Billiard room manager addStation FAILED. Got error retireving item with ID %s: %s", station.GetID(), err)
	}
	if newStation.GetID() != station.GetID() {
		t.Errorf("Billiard room manager addStation FAILED. Expected item's ID %s, but got %s", station.GetID(), newStation.GetID())
	}
}
