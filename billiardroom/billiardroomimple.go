package billiardroom

import (
	"fmt"
	"sync"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/gamestation"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/payment"
	"aen.it/poolmanager/warehouse"
)

type manager struct {
	mu           sync.RWMutex
	name         string
	background   string
	gameStations map[string]gamestation.GamingStation
	items        map[string]warehouse.Item
	itemQuantity map[string]int
	openChecks   map[string]payment.Check
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
	manage.name = config.BilliardRoomConfig.Name
	manage.gameStations = make(map[string]gamestation.GamingStation)
	for _, station := range config.BilliardRoomConfig.GamingStations {
		newGameStation := gamestation.New(station)
		manage.gameStations[station.ID] = &newGameStation
	}
	manage.items = make(map[string]warehouse.Item)
	manage.itemQuantity = make(map[string]int)
	manage.openChecks = make(map[string]payment.Check)
	for _, item := range config.BilliardRoomConfig.Items {
		newItem := warehouse.Item{
			ID:            item.ID,
			Name:          item.Name,
			PublicPrice:   item.PublicPrice,
			IncomingPrice: item.IncomingPrice,
		}
		manage.items[item.ID] = newItem
		manage.itemQuantity[item.ID] = 0
	}
	log.Log.Debug("Exiting loadFromConfig")
}

// Returns number of generic gaming station
func (manage *manager) GetNumberOfGamingStation() int {
	log.Log.Debug("Entering GetNumberOfGamingStation")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	log.Log.Debug("Exiting GetNumberOfGamingStation")
	return len(manage.gameStations)
}

// Returns list of Gaming Station's IDs
func (manage *manager) GetGamingStationIDs() []string {
	log.Log.Debug("Entering GetGamingStationIDs")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	ids := make([]string, 0, len(manage.gameStations))
	for id := range manage.gameStations {
		ids = append(ids, id)
	}
	log.Log.Debug("Exiting GetGamingStationIDs")
	return ids
}

// Returns required gaming station
func (manage *manager) GetGamingStation(id string) (gamestation.GamingStation, error) {
	log.Log.Debug("Entering GetGamingStation")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
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
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.gameStations[gs.GetID()]
	if ok {
		err := fmt.Errorf("gaming station with specified ID %s already exists", gs.GetID())
		log.Log.Debug("Exiting AddGamingStation")
		return err
	}
	manage.gameStations[gs.GetID()] = gs
	log.Log.Debug("Exiting AddGamingStation")
	return nil
}

// Returns number of available items
func (manage *manager) GetNumberOfItems() int {
	log.Log.Debug("Entering GetNumberOfItems")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	log.Log.Debug("Exiting GetNumberOfItems")
	return len(manage.items)
}

// Returns list of items's IDs
func (manage *manager) GetItemIDs() []string {
	log.Log.Debug("Entering GetItemIDs")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	ids := make([]string, 0, len(manage.items))
	for id := range manage.items {
		ids = append(ids, id)
	}
	log.Log.Debug("Exiting GetItemIDs")
	return ids
}

// Returns required item
func (manage *manager) GetItem(id string) (warehouse.Item, error) {
	log.Log.Debug("Entering GetItem")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
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
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.items[item.ID]
	if ok {
		err := fmt.Errorf("gaming station with specified ID %s already exists", item.ID)
		log.Log.Debug("Exiting AddItem")
		return err
	}
	manage.items[item.ID] = item
	manage.itemQuantity[item.ID] = 0
	log.Log.Debug("Exiting AddItem")
	return nil
}

// Add items with quantity to the warehouse
func (manage *manager) AddItems(item warehouse.Item, quantity int) {
	log.Log.Debug("Entering AddItems", "itemID", item.ID, "quantity", quantity)
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.items[item.ID]
	if !ok {
		// Item doesn't exist, create it
		manage.items[item.ID] = item
		manage.itemQuantity[item.ID] = quantity
	} else {
		// Item exists, update quantity and potentially update item details
		manage.items[item.ID] = item
		manage.itemQuantity[item.ID] += quantity
	}
	log.Log.Debug("Exiting AddItems", "itemID", item.ID, "newQuantity", manage.itemQuantity[item.ID])
}

// Update item properties (Name, PublicPrice, IncomingPrice)
func (manage *manager) UpdateItem(itemID string, name string, publicPrice int, incomingPrice int) error {
	log.Log.Debug("Entering UpdateItem", "itemID", itemID)
	manage.mu.Lock()
	defer manage.mu.Unlock()
	existingItem, ok := manage.items[itemID]
	if !ok {
		err := fmt.Errorf("item with specified ID %s does not exist", itemID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting UpdateItem")
		return err
	}

	// Update only the modifiable properties
	existingItem.Name = name
	existingItem.PublicPrice = publicPrice
	existingItem.IncomingPrice = incomingPrice

	manage.items[itemID] = existingItem
	log.Log.Debug("Exiting UpdateItem", "updated item", existingItem)
	return nil
}

// Delete an item from the warehouse
func (manage *manager) DeleteItem(itemID string) error {
	log.Log.Debug("Entering DeleteItem", "itemID", itemID)
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.items[itemID]
	if !ok {
		err := fmt.Errorf("item with specified ID %s does not exist", itemID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting DeleteItem")
		return err
	}

	delete(manage.items, itemID)
	delete(manage.itemQuantity, itemID)
	log.Log.Debug("Exiting DeleteItem", "deleted itemID", itemID)
	return nil
}

// Get quantity of a specific item
func (manage *manager) GetItemQuantity(itemID string) int {
	log.Log.Debug("Entering GetItemQuantity", "itemID", itemID)
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	quantity, ok := manage.itemQuantity[itemID]
	if !ok {
		log.Log.Debug("Exiting GetItemQuantity - item not found", "itemID", itemID)
		return 0
	}
	log.Log.Debug("Exiting GetItemQuantity", "itemID", itemID, "quantity", quantity)
	return quantity
}

// Get list of open checks' IDs
func (manage *manager) GetOpenCheckIDs() []string {
	log.Log.Debug("Entering GetOpenCheckIDs")
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	ids := make([]string, 0, len(manage.openChecks))
	for id := range manage.openChecks {
		ids = append(ids, id)
	}
	log.Log.Debug("Exiting GetOpenCheckIDs")
	return ids
}

// Add new check to the list
func (manage *manager) AddCheck(check payment.Check) error {
	log.Log.Debug("Entering AddCheck", "checkID", check.ID)
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.openChecks[check.ID]
	if ok {
		err := fmt.Errorf("check with specified ID %s already exists", check.ID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting AddCheck")
		return err
	}
	manage.openChecks[check.ID] = check
	log.Log.Debug("Exiting AddCheck", "added checkID", check.ID)
	return nil
}

// Get required check
func (manage *manager) GetCheck(checkID string) (payment.Check, error) {
	log.Log.Debug("Entering GetCheck", "checkID", checkID)
	manage.mu.RLock()
	defer manage.mu.RUnlock()
	check, ok := manage.openChecks[checkID]
	if !ok {
		err := fmt.Errorf("check with specified ID %s does not exist", checkID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting GetCheck")
		return payment.Check{}, err
	}
	log.Log.Debug("Exiting GetCheck", "found checkID", checkID)
	return check, nil
}

// Pay check
func (manage *manager) PayCheck(checkID string) error {
	log.Log.Debug("Entering PayCheck", "checkID", checkID)
	manage.mu.Lock()
	defer manage.mu.Unlock()
	_, ok := manage.openChecks[checkID]
	if !ok {
		err := fmt.Errorf("check with specified ID %s does not exist", checkID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting PayCheck")
		return err
	}
	delete(manage.openChecks, checkID)
	log.Log.Debug("Exiting PayCheck", "paid checkID", checkID)
	return nil
}
