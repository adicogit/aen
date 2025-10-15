package warehouse

// Record representing a generic consumable good
type Item struct {
	ID            string
	Name          string
	PublicPrice   int
	IncomingPrice int
}

// Record representing a generic item in the warehouse
type WarehouseItem struct {
	Item     Item
	Quantity int
}

// Warehouse interface
type Warehouse interface {
	// Add a given number of item to the warehouse
	AddItems(item Item, quantity int)
	// Remove a given number of item with specified ID from the warehouse. It returns an error if ther are not enough items in the warehouse
	RemoveItems(itemID string, quantity int) error
	// Return number of item's types present in the warehouse
	GetItemTypesCount() int
	// Return number of items with specified ID present in the warehouse
	GetItemsCount(itemID string) int
	// Return list of item's IDs present in the warehouse
	GetItemIDs() []string
	// Read an item from the warehouse without removing it
	GetItem(itemID string) (Item, error)
}
