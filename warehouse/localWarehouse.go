package warehouse

import (
	"fmt"

	"github.com/google/uuid"
)

type LocalWahouse struct {
	id        string
	warehouse map[string]WarehouseItem
}

// New function initialize LocalWahouse
func NewLocalWahouse() LocalWahouse {
	return LocalWahouse{
		id:        uuid.New().String(),
		warehouse: make(map[string]WarehouseItem, 0),
	}
}

// Add a given number of item to the warehouse
func (lw *LocalWahouse) AddItems(item Item, quantity int) {
	existingItem, ok := lw.warehouse[item.ID]
	if !ok {
		existingItem.Quantity = 0
		existingItem.Item = item
	}
	existingItem.Quantity += quantity
	lw.warehouse[item.ID] = existingItem
}

// Remove a given number of item with specified ID from the warehouse. It returns an error if ther are not enough items in the warehouse
func (lw *LocalWahouse) RemoveItems(itemID string, quantity int) error {
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		return fmt.Errorf("There are no items with specified ID: %s", itemID)
	}
	if existingItem.Quantity < quantity {
		return fmt.Errorf("There are not enough items. %d is bigger than available quoantity %d", quantity, existingItem.Quantity)
	}
	existingItem.Quantity -= quantity
	lw.warehouse[itemID] = existingItem
	return nil
}

// Return number of item's types present in the warehouse
func (lw *LocalWahouse) GetItemTypesCount() int {
	return len(lw.warehouse)
}

// Return number of items with specified ID present in the warehouse
func (lw *LocalWahouse) GetItemsCount(itemID string) int {
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		return 0
	}
	return existingItem.Quantity
}

// Return list of item's IDs present in the warehouse
func (lw *LocalWahouse) GetItemIDs() []string {
	keys := make([]string, 0, len(lw.warehouse))
	for k := range lw.warehouse {
		keys = append(keys, k)
	}
	return keys
}

// Read an item from the warehouse without removing it
func (lw *LocalWahouse) GetItem(itemID string) (Item, error) {
	existingItem, ok := lw.warehouse[itemID]
	if !ok {
		return Item{}, fmt.Errorf("There are no items with specified ID: %s", itemID)
	}
	return existingItem.Item, nil
}
