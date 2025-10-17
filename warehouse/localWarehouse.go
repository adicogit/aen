package warehouse

import (
	"fmt"

	"aen.it/poolmanager/log"
	"github.com/google/uuid"
)

type LocalWahouse struct {
	id        string
	warehouse map[string]WarehouseItem
}

// New function initialize LocalWahouse
func NewLocalWahouse() LocalWahouse {
	log.Log.Debug("Entering NewLocalWahouse")
	log.Log.Debug("Exiting NewLocalWahouse")
	return LocalWahouse{
		id:        uuid.New().String(),
		warehouse: make(map[string]WarehouseItem, 0),
	}
}

// Add a given number of item to the warehouse
func (lw *LocalWahouse) AddItems(item Item, quantity int) {
	log.Log.Debug("Entering AddItems")
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
func (lw *LocalWahouse) RemoveItems(itemID string, quantity int) error {
	log.Log.Debug("Entering RemoveItems")
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
func (lw *LocalWahouse) GetItemTypesCount() int {
	log.Log.Debug("Entering GetItemTypesCount")
	log.Log.Debug("Exiting GetItemTypesCount")
	return len(lw.warehouse)
}

// Return number of items with specified ID present in the warehouse
func (lw *LocalWahouse) GetItemsCount(itemID string) int {
	log.Log.Debug("Entering GetItemsCount")
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		log.Log.Debug("Exiting GetItemsCount")
		return 0
	}
	log.Log.Debug("Exiting GetItemsCount")
	return existingItem.Quantity
}

// Return list of item's IDs present in the warehouse
func (lw *LocalWahouse) GetItemIDs() []string {
	log.Log.Debug("Exiting GetItemIDs")
	keys := make([]string, 0, len(lw.warehouse))
	for k := range lw.warehouse {
		keys = append(keys, k)
	}
	log.Log.Debug("Exiting GetItemIDs")
	return keys
}

// Read an item from the warehouse without removing it
func (lw *LocalWahouse) GetItem(itemID string) (Item, error) {
	log.Log.Debug("Entering GetItem")
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		return Item{}, fmt.Errorf("there are no items with specified ID: %s", itemID)
	}
	log.Log.Debug("Exiting GetItem")
	return existingItem.Item, nil
}
