package billiardroom

import (
	"fmt"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/gamestation"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
)

type manager struct {
	name         string
	backgound    string
	gameStations map[string]gamestation.GamingStation
	items        map[string]warehouse.Item
}

var Manager *manager

func init() {
	log.Log.Debug("Entering manager init")
	log.Log.Info("Create billiard room manager")
	Manager = &manager{}
	Manager.loadFromConfig()
	log.Log.Debug("Exiting manager init")
}

func (manage *manager) loadFromConfig() {
	log.Log.Debug("Entering loadFromConfig")
	manage.name = config.Config.Name
	manage.gameStations = make(map[string]gamestation.GamingStation)
	for _, station := range config.Config.GamingStations {
		newGameStation := gamestation.New(station)
		manage.gameStations[station.ID] = &newGameStation
	}
	manage.items = make(map[string]warehouse.Item)
	for _, item := range config.Config.Items {
		newItem := warehouse.Item{
			ID:            item.ID,
			Name:          item.Name,
			PublicPrice:   item.PublicPrice,
			IncomingPrice: item.IncomingPrice,
		}
		manage.items[item.ID] = newItem
	}
	log.Log.Debug("Exiting loadFromConfig")
}

// Returns number of generic gaming station
func (manage *manager) GetNumberOfGamingStation() int {
	log.Log.Debug("Entering GetNumberOfGamingStation")
	log.Log.Debug("Exiting GetNumberOfGamingStation")
	return len(manage.gameStations)
}

// Returns list og Gaming Station's IDs
func (manage *manager) GetGamingStationIDs() []string {
	log.Log.Debug("Entering GetGamingStationIDs")
	ids := make([]string, len(manage.gameStations))
	i := 0
	for id := range manage.gameStations {
		ids[i] = id
		i++
	}
	log.Log.Debug("Exiting GetGamingStationIDs")
	return ids
}

// Retturns required gaming station
func (manage *manager) GetGamingStation(id string) (gamestation.GamingStation, error) {
	log.Log.Debug("Entering GetGamingStation")
	station, ok := manage.gameStations[id]
	if !ok {
		err := fmt.Errorf("gaming station with specified ID %s does not exist", id)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting GetGamingStation")
		return nil, err
	}
	log.Log.Debug("Exiting GetGamingStation")
	return station, nil
}

// Add new Gaming station to the list
func (manage *manager) AddGamingStation(gs gamestation.GamingStation) error {
	log.Log.Debug("Entering AddGamingStation")
	log.Log.Debug("Exiting AddGamingStation")
	return nil
}

// Returns number of available items
func (manage *manager) GetNumberOfItems() int {
	log.Log.Debug("Entering GetNumberOfItems")
	log.Log.Debug("Exiting GetNumberOfItems")
	return len(manage.items)
}

// Returns list og items's IDs
func (manage *manager) GetItemIDs() []string {
	log.Log.Debug("Entering GetItemIDs")
	ids := make([]string, len(manage.items))
	i := 0
	for id := range manage.items {
		ids[i] = id
		i++
	}
	log.Log.Debug("Exiting GetItemIDs")
	return ids
}

// Retturns required item
func (manage *manager) GetItem(id string) (warehouse.Item, error) {
	log.Log.Debug("Entering GetGaminGetItemgStation")
	item, ok := manage.items[id]
	if !ok {
		err := fmt.Errorf("item with specified ID %s does not exist", id)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting GetItem")
		return warehouse.Item{}, err
	}
	log.Log.Debug("Exiting GetItem")
	return item, nil
}

// Add new item to the list
func (manage *manager) AddItem(item warehouse.Item) error {
	log.Log.Debug("Entering AddItem")
	log.Log.Debug("Exiting AddItem")
	return nil
}
