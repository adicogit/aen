package payment

import (
	"testing"
	"time"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/warehouse"
)

// Helper function to create a test configuration
func createTestConfig() config.PaymentConfiguration {
	return config.PaymentConfiguration{
		CostPerHour:     6000, // 60 euros per hour in cents
		MinimumDuration: 15,   // 15 minutes minimum
	}
}

// Helper function to create a test item
func createTestItem(id, name string, price int) warehouse.Item {
	return warehouse.Item{
		ID:            id,
		Name:          name,
		PublicPrice:   price,
		IncomingPrice: price / 2,
	}
}

// Test New function creates a GamePayment with correct initial state
func TestNew(t *testing.T) {
	config := createTestConfig()
	gp := New(config, "Test Station")

	if gp.id == "" {
		t.Error("Expected non-empty ID")
	}

	if gp.status != Stopped {
		t.Errorf("Expected initial status to be Stopped, got %v", gp.status)
	}

	if gp.configuration.CostPerHour != config.CostPerHour {
		t.Errorf("Expected CostPerHour %d, got %d", config.CostPerHour, gp.configuration.CostPerHour)
	}

	if gp.previousDuration != 0 {
		t.Error("Expected previousDuration to be 0")
	}

	if len(gp.itemList) != 0 {
		t.Error("Expected empty itemList")
	}
}

// Test ConfigurePayment updates configuration
func TestConfigurePayment(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	newConfig := config.PaymentConfiguration{
		CostPerHour:     8000,
		MinimumDuration: 30,
	}

	gp.ConfigurePayment(newConfig)

	if gp.configuration.CostPerHour != newConfig.CostPerHour {
		t.Errorf("Expected CostPerHour %d, got %d", newConfig.CostPerHour, gp.configuration.CostPerHour)
	}

	if gp.configuration.MinimumDuration != newConfig.MinimumDuration {
		t.Errorf("Expected MinimumDuration %d, got %d", newConfig.MinimumDuration, gp.configuration.MinimumDuration)
	}
}

// Test StartCountingPayment from Stopped status
func TestStartCountingPayment_FromStopped(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	err := gp.StartCountingPayment()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if gp.status != Started {
		t.Errorf("Expected status to be Started, got %v", gp.status)
	}

	if gp.start.IsZero() {
		t.Error("Expected start time to be set")
	}
}

// Test StartCountingPayment from Suspended status
func TestStartCountingPayment_FromSuspended(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	gp.PauseCountingPayment()

	err := gp.StartCountingPayment()
	if err != nil {
		t.Errorf("Expected no error when starting from Suspended, got %v", err)
	}

	if gp.status != Started {
		t.Errorf("Expected status to be Started, got %v", gp.status)
	}
}

// Test StartCountingPayment fails when already Started
func TestStartCountingPayment_AlreadyStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()

	err := gp.StartCountingPayment()
	if err == nil {
		t.Error("Expected error when starting already started payment")
	}

	if gp.status != Started {
		t.Errorf("Expected status to remain Started, got %v", gp.status)
	}
}

// Test PauseCountingPayment from Started status
func TestPauseCountingPayment_FromStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)

	err := gp.PauseCountingPayment()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if gp.status != Suspended {
		t.Errorf("Expected status to be Suspended, got %v", gp.status)
	}

	if gp.previousDuration == 0 {
		t.Error("Expected previousDuration to be set")
	}
}

// Test PauseCountingPayment fails when not Started
func TestPauseCountingPayment_NotStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	err := gp.PauseCountingPayment()
	if err == nil {
		t.Error("Expected error when pausing non-started payment")
	}

	if gp.status != Stopped {
		t.Errorf("Expected status to remain Stopped, got %v", gp.status)
	}
}

// Test ClosePayment from Started status
func TestClosePayment_FromStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)

	err := gp.ClosePayment()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if gp.status != Stopped {
		t.Errorf("Expected status to be Stopped, got %v", gp.status)
	}

	if gp.check.Duration == 0 {
		t.Error("Expected check duration to be set")
	}

	if gp.check.Price == 0 {
		t.Error("Expected check price to be set")
	}
}

// Test ClosePayment from Suspended status
func TestClosePayment_FromSuspended(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)
	gp.PauseCountingPayment()

	err := gp.ClosePayment()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if gp.status != Stopped {
		t.Errorf("Expected status to be Stopped, got %v", gp.status)
	}
}

// Test ClosePayment fails when already Stopped
func TestClosePayment_AlreadyStopped(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	err := gp.ClosePayment()
	if err == nil {
		t.Error("Expected error when closing already stopped payment")
	}
}

// Test ClosePayment applies minimum duration
func TestClosePayment_MinimumDuration(t *testing.T) {
	config := createTestConfig()
	config.MinimumDuration = 15 // 15 minutes
	gp := New(config, "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond) // Less than minimum

	gp.ClosePayment()

	if gp.check.Duration != config.MinimumDuration {
		t.Errorf("Expected duration to be minimum %d minutes, got %d", config.MinimumDuration, gp.check.Duration)
	}

	expectedPrice := config.MinimumDuration * config.CostPerHour / 60
	if gp.check.Price != expectedPrice {
		t.Errorf("Expected price %d, got %d", expectedPrice, gp.check.Price)
	}
}

// Test GetCheck returns check after closing
func TestGetCheck_AfterClose(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)
	gp.ClosePayment()

	check, err := gp.GetCheck()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if check.Duration == 0 {
		t.Error("Expected non-zero duration in check")
	}

	if check.Price == 0 {
		t.Error("Expected non-zero price in check")
	}
}

// Test GetCheck fails when payment not closed
func TestGetCheck_NotClosed(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()

	_, err := gp.GetCheck()
	if err == nil {
		t.Error("Expected error when getting check for non-closed payment")
	}
}

// Test GetTemporaryCheck while Started
func TestGetTemporaryCheck_WhileStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)

	check := gp.GetTemporaryCheck()

	if check.Duration == 0 {
		t.Error("Expected non-zero duration in temporary check")
	}

	if check.Price == 0 {
		t.Error("Expected non-zero price in temporary check")
	}
}

// Test GetTemporaryCheck while Suspended
func TestGetTemporaryCheck_WhileSuspended(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)
	gp.PauseCountingPayment()

	check := gp.GetTemporaryCheck()

	if check.Duration == 0 {
		t.Error("Expected non-zero duration in temporary check")
	}
}

// Test GetTemporaryCheck when Stopped returns final check
func TestGetTemporaryCheck_WhenStopped(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	time.Sleep(100 * time.Millisecond)
	gp.ClosePayment()

	tempCheck := gp.GetTemporaryCheck()
	finalCheck, _ := gp.GetCheck()

	if tempCheck.Duration != finalCheck.Duration {
		t.Errorf("Expected temporary check duration %d to match final %d", tempCheck.Duration, finalCheck.Duration)
	}

	if tempCheck.Price != finalCheck.Price {
		t.Errorf("Expected temporary check price %d to match final %d", tempCheck.Price, finalCheck.Price)
	}
}

// Test GetPaymentStatus
func TestGetPaymentStatus(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	if gp.GetPaymentStatus() != Stopped {
		t.Errorf("Expected initial status Stopped, got %v", gp.GetPaymentStatus())
	}

	gp.StartCountingPayment()
	if gp.GetPaymentStatus() != Started {
		t.Errorf("Expected status Started, got %v", gp.GetPaymentStatus())
	}

	gp.PauseCountingPayment()
	if gp.GetPaymentStatus() != Suspended {
		t.Errorf("Expected status Suspended, got %v", gp.GetPaymentStatus())
	}

	gp.ClosePayment()
	if gp.GetPaymentStatus() != Stopped {
		t.Errorf("Expected status Stopped, got %v", gp.GetPaymentStatus())
	}
}

// Test AddConsumption while Started
func TestAddConsumption_WhileStarted(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()

	item := createTestItem("item1", "Test Item", 500)
	err := gp.AddConsumption(item)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(gp.itemList) != 1 {
		t.Errorf("Expected 1 item in list, got %d", len(gp.itemList))
	}

	if gp.itemList[0].ID != item.ID {
		t.Errorf("Expected item ID %s, got %s", item.ID, gp.itemList[0].ID)
	}
}

// Test AddConsumption while Suspended
func TestAddConsumption_WhileSuspended(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()
	gp.PauseCountingPayment()

	item := createTestItem("item1", "Test Item", 500)
	err := gp.AddConsumption(item)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(gp.itemList) != 1 {
		t.Errorf("Expected 1 item in list, got %d", len(gp.itemList))
	}
}

// Test AddConsumption fails when Stopped
func TestAddConsumption_WhenStopped(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	item := createTestItem("item1", "Test Item", 500)
	err := gp.AddConsumption(item)

	if err == nil {
		t.Error("Expected error when adding consumption to stopped payment")
	}

	if len(gp.itemList) != 0 {
		t.Errorf("Expected 0 items in list, got %d", len(gp.itemList))
	}
}

// Test AddConsumption multiple items
func TestAddConsumption_MultipleItems(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()

	item1 := createTestItem("item1", "Item 1", 500)
	item2 := createTestItem("item2", "Item 2", 300)
	item3 := createTestItem("item3", "Item 3", 700)

	gp.AddConsumption(item1)
	gp.AddConsumption(item2)
	gp.AddConsumption(item3)

	if len(gp.itemList) != 3 {
		t.Errorf("Expected 3 items in list, got %d", len(gp.itemList))
	}
}

// Test consumption items are included in final check
func TestClosePayment_WithConsumptions(t *testing.T) {
	config := createTestConfig()
	config.MinimumDuration = 15
	gp := New(config, "Test Station")
	gp.StartCountingPayment()

	item1 := createTestItem("item1", "Item 1", 500)
	item2 := createTestItem("item2", "Item 2", 300)
	gp.AddConsumption(item1)
	gp.AddConsumption(item2)

	time.Sleep(100 * time.Millisecond)
	gp.ClosePayment()

	check, _ := gp.GetCheck()

	if len(check.ItemList) != 2 {
		t.Errorf("Expected 2 items in check, got %d", len(check.ItemList))
	}

	expectedTimePrice := config.MinimumDuration * config.CostPerHour / 60
	expectedTotalPrice := expectedTimePrice + item1.PublicPrice + item2.PublicPrice

	if check.Price != expectedTotalPrice {
		t.Errorf("Expected total price %d, got %d", expectedTotalPrice, check.Price)
	}
}

// Test GetTemporaryCheck includes consumptions
func TestGetTemporaryCheck_WithConsumptions(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")
	gp.StartCountingPayment()

	item := createTestItem("item1", "Item 1", 500)
	gp.AddConsumption(item)

	time.Sleep(100 * time.Millisecond)
	check := gp.GetTemporaryCheck()

	if len(check.ItemList) != 1 {
		t.Errorf("Expected 1 item in temporary check, got %d", len(check.ItemList))
	}

	if check.Price < item.PublicPrice {
		t.Errorf("Expected price to include item price %d, got %d", item.PublicPrice, check.Price)
	}
}

// Test payment lifecycle: Start -> Pause -> Resume -> Close
func TestPaymentLifecycle_Complete(t *testing.T) {
	gp := New(createTestConfig(), "Test Station")

	// Start payment
	err := gp.StartCountingPayment()
	if err != nil {
		t.Fatalf("Failed to start payment: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	// Add consumption
	item := createTestItem("item1", "Test Item", 500)
	gp.AddConsumption(item)

	// Pause payment
	err = gp.PauseCountingPayment()
	if err != nil {
		t.Fatalf("Failed to pause payment: %v", err)
	}

	// Resume payment
	err = gp.StartCountingPayment()
	if err != nil {
		t.Fatalf("Failed to resume payment: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	// Close payment
	err = gp.ClosePayment()
	if err != nil {
		t.Fatalf("Failed to close payment: %v", err)
	}

	// Verify final check
	check, err := gp.GetCheck()
	if err != nil {
		t.Fatalf("Failed to get check: %v", err)
	}

	if check.Duration == 0 {
		t.Error("Expected non-zero duration")
	}

	if check.Price == 0 {
		t.Error("Expected non-zero price")
	}

	if len(check.ItemList) != 1 {
		t.Errorf("Expected 1 item in check, got %d", len(check.ItemList))
	}
}

// Made with Bob
