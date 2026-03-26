package billiardroom

import (
	"aen.it/poolmanager/gamestation"
	"aen.it/poolmanager/warehouse"
)

type BilliardRoom interface {
	// Returns number of generic gaming station
	GetNumberOfGamingStation() int
	// Returns list of Gaming Station's IDs
	GetGamingStationIDs() []string
	// Returns required gaming station
	GetGamingStation(id string) (gamestation.GamingStation, error)
	// Add new Gaming station to the list
	AddGamingStation(gs gamestation.GamingStation) error
	// Returns number of available items
	GetNumberOfItems() int
	// Returns list of items's IDs
	GetItemIDs() []string
	// Returns required item
	GetItem(id string) (warehouse.Item, error)
	// Add new item to the list
	AddItem(item warehouse.Item) error
	// Add items with quantity to the warehouse
	AddItems(item warehouse.Item, quantity int)
	// Update item properties (Name, PublicPrice, IncomingPrice)
	UpdateItem(itemID string, name string, publicPrice int, incomingPrice int) error
	// Delete an item from the warehouse
	DeleteItem(itemID string) error
}
