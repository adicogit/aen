package payment

import (
	"log/slog"
	"testing"
	"time"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
)

func init() {
	//	log.SetLogLevel(slog.LevelDebug)
	log.SetLogLevel(slog.LevelInfo)
}

// TestGamePaymentInitialization verify that New function works as expected
func TestGamePaymentInitialization(t *testing.T) {
	config := config.PaymentConfiguration{
		MinimumDuration: 30,
		CostPerHour:     10,
	}
	gamePayment := New(config)
	// Verify MinimumDuration
	if gamePayment.configuration.MinimumDuration != config.MinimumDuration {
		t.Errorf("GamePayment initialization FAILED. Its MinimumDuration value is %d but it should be %d", gamePayment.configuration.MinimumDuration, config.MinimumDuration)
	}
	// Verify CostPerHour
	if gamePayment.configuration.CostPerHour != config.CostPerHour {
		t.Errorf("GamePayment initialization FAILED. Its CostPerHour value is %d but it should be %d", gamePayment.configuration.CostPerHour, config.CostPerHour)
	}
	// Verify start
	zeroTime := time.Time{}
	if gamePayment.start.Compare(zeroTime) != 0 {
		t.Errorf("GamePayment initialization FAILED. Its start value is %s but it should be %s", gamePayment.start.String(), zeroTime.String())
	}
	// Verify previousDuration
	zeroDuration := time.Duration(0)
	if gamePayment.previousDuration != zeroDuration {
		t.Errorf("GamePayment initialization FAILED. Its previousDuration value is %d but it should be %d", int(gamePayment.previousDuration.Minutes()), int(zeroDuration.Minutes()))
	}
}

// TestGamePaymentStart verify start function works as expected
func TestGamePaymentStart(t *testing.T) {
	config := config.PaymentConfiguration{
		MinimumDuration: 30,
		CostPerHour:     10,
	}
	gamePayment := New(config)
	err := gamePayment.StartCountingPayment()
	if err != nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. First start returned error %s but it should be OK", err)
	}
	status := gamePayment.GetPaymentStatus()
	if status != Started {
		t.Errorf("GamePayment StartCountingPayment FAILED. Expecting status %d but it is %d", Started, status)
	}
	_, err = gamePayment.GetCheck()
	if err == nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. Without closing the game after a start should generate an error but it did not")
	}
	err = gamePayment.StartCountingPayment()
	if err == nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. Second start is OK but error is expected")
	}
}

// TestGamePaymentPause verify pause works as expected
func TestGamePaymentPause(t *testing.T) {
	config := config.PaymentConfiguration{
		MinimumDuration: 0,
		CostPerHour:     10,
	}
	gamePayment := New(config)
	gamePayment.StartCountingPayment()
	//Wait for a little
	time.Sleep(5 * time.Second)
	err := gamePayment.PauseCountingPayment()
	if err != nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. First pause returned error %s but it should be OK", err)
	}
	status := gamePayment.GetPaymentStatus()
	if status != Suspended {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Expecting status %d but it is %d", Suspended, status)
	}
	_, err = gamePayment.GetCheck()
	if err == nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Without closing the game after a pause should generate an error but it did not")
	}
	err = gamePayment.PauseCountingPayment()
	if err == nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Second pause is OK but error is expected")
	}
	err = gamePayment.StartCountingPayment()
	if err != nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Starting counting after pause it returned error %s but it should be OK", err)
	}
}

// TestGamePaymentPause verify pause works as expected
func TestGamePaymentClosure(t *testing.T) {
	config := config.PaymentConfiguration{
		MinimumDuration: 15,
		CostPerHour:     10,
	}
	gamePayment := New(config)
	err := gamePayment.StartCountingPayment()
	if err != nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. Trying to start new game returned an error %s but it should be OK", err)
	}
	item := warehouse.Item{
		ID:            "",
		Name:          "acqua",
		PublicPrice:   50,
		IncomingPrice: 10,
	}
	err = gamePayment.AddConsumption(item)
	if err != nil {
		t.Errorf("GamePayment AddConsumption FAILED. Returned error %s rying to add new consumption to a started payment, but it should be OK", err)
	}
	gamePayment.start = time.Now().Add(-10 * time.Minute)
	err = gamePayment.PauseCountingPayment()
	if err != nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Trying to pause a started returned an error %s but it should be OK", err)
	}
	err = gamePayment.AddConsumption(item)
	if err != nil {
		t.Errorf("GamePayment AddConsumption FAILED. Returned error %s rying to add new consumption to a started payment, but it should be OK", err)
	}
	// verify payment without pause and with less than minimum duration
	err = gamePayment.ClosePayment()
	if err != nil {
		t.Errorf("GamePayment ClosePayment FAILED. First close returned error %s but it should be OK", err)
	}
	status := gamePayment.GetPaymentStatus()
	if status != Stopped {
		t.Errorf("GamePayment ClosePayment FAILED. Expecting status %d but it is %d", Stopped, status)
	}
	err = gamePayment.AddConsumption(item)
	if err == nil {
		t.Errorf("GamePayment AddConsumption FAILED. Did not return any error, but it should. you cannot add new consumption to a stopped payment")
	}
	check, err := gamePayment.GetCheck()
	if err != nil {
		t.Errorf("GamePayment ClosePayment FAILED. Closing the game should not generate an error in getting the check")
	}
	if check.Duration != config.MinimumDuration {
		t.Errorf("GamePayment GetCheck FAILED. Check's duration is %d but it is expeted to be %d", check.Duration, config.MinimumDuration)
	}
	expectedPrice := config.MinimumDuration * config.CostPerHour / 60
	expectedPrice += item.PublicPrice * 2
	if check.Price != expectedPrice {
		t.Errorf("GamePayment GetCheck FAILED. Check's price is %.2f but it is expeted to be %.2f", float32(check.Price/100), float32(expectedPrice/100))
	}
	err = gamePayment.ClosePayment()
	if err == nil {
		t.Errorf("GamePayment ClosePayment FAILED. Second close is OK but error is expected")
	}

	// verify payment after pause
	err = gamePayment.StartCountingPayment()
	if err != nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. Trying to start new game after a close returned an error %s but it should be OK", err)
	}
	gamePayment.start = time.Now().Add(-15 * time.Minute)
	err = gamePayment.PauseCountingPayment()
	if err != nil {
		t.Errorf("GamePayment PauseCountingPayment FAILED. Trying to pause a started game %s but it should be OK", err)
	}
	gamePayment.StartCountingPayment()
	if err != nil {
		t.Errorf("GamePayment StartCountingPayment FAILED. Trying to start a paused game %s but it should be OK", err)
	}
	gamePayment.start = time.Now().Add(-10 * time.Minute)
	// Game duration is expeted to be 25 minutes
	gamePayment.ClosePayment()
	if err != nil {
		t.Errorf("GamePayment ClosePayment FAILED. First close returned error %s but it should be OK", err)
	}
	status = gamePayment.GetPaymentStatus()
	if status != Stopped {
		t.Errorf("GamePayment ClosePayment FAILED. Expecting status %d but it is %d", Stopped, status)
	}
	check, err = gamePayment.GetCheck()
	if err != nil {
		t.Errorf("GamePayment ClosePayment FAILED. Closing the game should not generate an error in getting the check")
	}
	if check.Duration != 25 {
		t.Errorf("GamePayment GetCheck FAILED. Check's duration is %d but it is expeted to be %d", check.Duration, 25)
	}
	expectedPrice = 25 * config.CostPerHour / 60
	if check.Price != expectedPrice {
		t.Errorf("GamePayment GetCheck FAILED. Check's price is %.2f but it is expeted to be %.2f", float32(check.Price), float32(expectedPrice))
	}
	err = gamePayment.ClosePayment()
	if err == nil {
		t.Errorf("GamePayment ClosePayment FAILED. Second close is OK but error is expected")
	}
}
