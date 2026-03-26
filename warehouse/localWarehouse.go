package warehouse

import (
	"fmt"
	"sync"

	"aen.it/poolmanager/log"
	"github.com/google/uuid"
)

type LocalWarehouse struct {
	id        string
	warehouse map[string]WarehouseItem
	mu        sync.RWMutex
}

// New function initialize LocalWarehouse
func NewLocalWarehouse() LocalWarehouse {
	log.Log.Debug("Entering NewLocalWarehouse")
	log.Log.Debug("Exiting NewLocalWarehouse")
	return LocalWarehouse{
		id:        uuid.New().String(),
		warehouse: make(map[string]WarehouseItem, 0),
	}
}

// Add a given number of item to the warehouse
func (lw *LocalWarehouse) AddItems(item Item, quantity int) {
	log.Log.Debug("Entering AddItems")
	lw.mu.Lock()
	defer lw.mu.Unlock()
	existingItem, ok := lw.warehouse[item.ID]
	if !ok {
		existingItem.Quantity = 0
		existingItem.Item = item
	}
	existingItem.Quantity += quantity
	lw.warehouse[item.ID] = existingItem
	log.Log.Debug("Exiting AddItems")
}

// Remove a given number of item with specified ID from the warehouse. It returns an error if ther are not enough items in the warehouse
func (lw *LocalWarehouse) RemoveItems(itemID string, quantity int) error {
	log.Log.Debug("Entering RemoveItems")
	lw.mu.Lock()
	defer lw.mu.Unlock()
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		err := fmt.Errorf("there are no items with specified ID: %s", itemID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting RemoveItems")
		return err
	}
	if existingItem.Quantity < quantity {
		err := fmt.Errorf("there are not enough items. %d is bigger than available quantity %d", quantity, existingItem.Quantity)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting RemoveItems")
		return err
	}
	existingItem.Quantity -= quantity
	lw.warehouse[itemID] = existingItem
	log.Log.Debug("Exiting RemoveItems")
	return nil
}

// Return number of item's types present in the warehouse
func (lw *LocalWarehouse) GetItemTypesCount() int {
	log.Log.Debug("Entering GetItemTypesCount")
	lw.mu.RLock()
	defer lw.mu.RUnlock()
	log.Log.Debug("Exiting GetItemTypesCount")
	return len(lw.warehouse)
}

// Return number of items with specified ID present in the warehouse
func (lw *LocalWarehouse) GetItemsCount(itemID string) int {
	log.Log.Debug("Entering GetItemsCount")
	lw.mu.RLock()
	defer lw.mu.RUnlock()
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		log.Log.Debug("Exiting GetItemsCount")
		return 0
	}
	log.Log.Debug("Exiting GetItemsCount")
	return existingItem.Quantity
}

// Return list of item's IDs present in the warehouse
func (lw *LocalWarehouse) GetItemIDs() []string {
	log.Log.Debug("Entering GetItemIDs")
	lw.mu.RLock()
	defer lw.mu.RUnlock()
	keys := make([]string, 0, len(lw.warehouse))
	for k := range lw.warehouse {
		keys = append(keys, k)
	}
	log.Log.Debug("Exiting GetItemIDs")
	return keys
}

// Read an item from the warehouse without removing it
func (lw *LocalWarehouse) GetItem(itemID string) (Item, error) {
	log.Log.Debug("Entering GetItem")
	lw.mu.RLock()
	defer lw.mu.RUnlock()
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		return Item{}, fmt.Errorf("there are no items with specified ID: %s", itemID)
	}
	log.Log.Debug("Exiting GetItem")
	return existingItem.Item, nil
}

// Update item properties (Name, PublicPrice, IncomingPrice). Returns error if item doesn't exist
func (lw *LocalWarehouse) UpdateItem(itemID string, name string, publicPrice int, incomingPrice int) error {
	log.Log.Debug("Entering UpdateItem", "itemID", itemID)
	lw.mu.Lock()
	defer lw.mu.Unlock()
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		err := fmt.Errorf("there are no items with specified ID: %s", itemID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting UpdateItem")
		return err
	}

	// Update only the modifiable properties
	existingItem.Item.Name = name
	existingItem.Item.PublicPrice = publicPrice
	existingItem.Item.IncomingPrice = incomingPrice

	lw.warehouse[itemID] = existingItem
	log.Log.Debug("Exiting UpdateItem", "updated item", existingItem)
	return nil
}

// Delete an item from the warehouse. Returns error if item doesn't exist
func (lw *LocalWarehouse) DeleteItem(itemID string) error {
	log.Log.Debug("Entering DeleteItem", "itemID", itemID)
	lw.mu.Lock()
	defer lw.mu.Unlock()
	_, ok := lw.warehouse[itemID]
	if !ok {
		err := fmt.Errorf("there are no items with specified ID: %s", itemID)
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting DeleteItem")
		return err
	}

	delete(lw.warehouse, itemID)
	log.Log.Debug("Exiting DeleteItem", "deleted itemID", itemID)
	return nil
}
