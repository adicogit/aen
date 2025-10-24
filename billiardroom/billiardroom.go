package billiardroom

import (
	"aen.it/poolmanager/gamestation"
	"aen.it/poolmanager/warehouse"
)

type BilliardRoom interface {
	// Returns number of generic gaming station
	GetNumberOfGamingStation() int
	// Returns list og Gaming Station's IDs
	GetGamingStationIDs() []string
	// Retturns required gaming station
	GetGamingStation(id string) (gamestation.GamingStation, error)
	// Add new Gaming station to the list
	AddGamingStation(gs gamestation.GamingStation) error
	// Returns number of available items
	GetNumberOfItems() int
	// Returns list og items's IDs
	GetItemIDs() []string
	// Retturns required item
	GetItem(id string) (warehouse.Item, error)
	// Add new item to the list
	AddItem(item warehouse.Item) error
}
