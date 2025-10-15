package warehouse

import (
	"testing"

	"github.com/google/uuid"
)

// TestLocalWarehouseInitialization verify that New function works as expected
func TestLocalWarehouseInitialization(t *testing.T) {
	warehouse := NewLocalWahouse()
	if err := uuid.Validate(warehouse.id); err != nil {
		t.Errorf("New created warehouse does not have a valid UUID: %s", warehouse.id)
	}
	itemsLen := warehouse.GetItemTypesCount()
	if itemsLen != 0 {
		t.Errorf("New created warehouse does not have an empty list of items. Len is: %d", itemsLen)
	}
}

// TestLocalWarehouseAddItem verify that AddItems function works as expected
func TestLocalWarehouseAddItem(t *testing.T) {
	item := Item{
		ID:            uuid.NewString(),
		Name:          "CocaCola",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse := NewLocalWahouse()
	warehouse.AddItems(item, 10)
	itemsLen := len(warehouse.GetItemIDs())
	if itemsLen != 1 {
		t.Errorf("Warehouse does not have correct number of different items: %d", itemsLen)
	}
	numberOfItem := warehouse.GetItemsCount(item.ID)
	if numberOfItem != 10 {
		t.Errorf("Warehouse does not have correct number of same items: %d", numberOfItem)
	}

	item = Item{
		ID:            uuid.NewString(),
		Name:          "Pizza",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse.AddItems(item, 2)
	itemsLen = len(warehouse.GetItemIDs())
	if itemsLen != 2 {
		t.Errorf("Warehouse does not have correct number of different items: %d", itemsLen)
	}
	numberOfItem = warehouse.GetItemsCount(item.ID)
	if numberOfItem != 2 {
		t.Errorf("Warehouse does not have correct number of same items: %d", numberOfItem)
	}
}

// TestLocalWarehouseRemoveItem verify that RemoveItems function works as expected
func TestLocalWarehouseRemoveItem(t *testing.T) {
	item := Item{
		ID:            uuid.NewString(),
		Name:          "CocaCola",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse := NewLocalWahouse()
	warehouse.AddItems(item, 10)
	error := warehouse.RemoveItems(item.ID, 3)
	if error != nil {
		t.Errorf("Warehouse generated unexpected error in removing an items %s", error)
	}
	itemsLen := len(warehouse.GetItemIDs())
	if itemsLen != 1 {
		t.Errorf("Warehouse does not have correct number of different items: %d", itemsLen)
	}
	numberOfItem := warehouse.GetItemsCount(item.ID)
	if numberOfItem != 7 {
		t.Errorf("Warehouse does not have correct number of same items: %d", numberOfItem)
	}
	error = warehouse.RemoveItems(item.ID, 10)
	if error == nil {
		t.Errorf("Warehouse did not generate an error in removing more items then present %s", error)
	}
	numberOfItem = warehouse.GetItemsCount(item.ID)
	if numberOfItem != 7 {
		t.Errorf("Warehouse does not have correct number of same items after failed removal: %d", numberOfItem)
	}
}

// TestLocalWarehouseItemTypesCount verify that GetItemTypesCount function works as expected
func TestLocalWarehouseItemTypesCount(t *testing.T) {
	warehouse := NewLocalWahouse()

	item := Item{
		ID:            uuid.NewString(),
		Name:          "CocaCola",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse.AddItems(item, 10)

	item = Item{
		ID:            uuid.NewString(),
		Name:          "Pizza",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse.AddItems(item, 2)

	item = Item{
		ID:            uuid.NewString(),
		Name:          "Pasta",
		PublicPrice:   200,
		IncomingPrice: 200,
	}
	warehouse.AddItems(item, 2)

	itemsLen := len(warehouse.GetItemIDs())
	if itemsLen != 3 {
		t.Errorf("Warehouse does not have correct number of different items: %d", itemsLen)
	}
	numberOfItemTypes := warehouse.GetItemTypesCount()
	if numberOfItemTypes != 3 {
		t.Errorf("Warehouse does not have correct number of different items: %d", numberOfItemTypes)
	}

	idList := warehouse.GetItemIDs()
	item, error := warehouse.GetItem(idList[0])
	if error != nil {
		t.Errorf("Warehouse does not have ID %s retrieved from GetItemIDs.", idList[0])
	}
	warehouse.AddItems(item, 5)
	numberOfItemTypes = warehouse.GetItemTypesCount()
	if numberOfItemTypes != 3 {
		t.Errorf("Warehouse does not have correct number of different items %d, after modifying quantity for existing", numberOfItemTypes)
	}
}
