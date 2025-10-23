package payment

import (
	"errors"
	"time"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
	"github.com/google/uuid"
)

type GamePayment struct {
	id               string
	configuration    config.PaymentConfiguration
	status           PaymentStatus
	start            time.Time
	previousDuration time.Duration
	check            Check
	itemList         []warehouse.Item
}

// New function initialize GamePayment passing a PaymentConfiguration
func New(config config.PaymentConfiguration) GamePayment {
	log.Log.Debug("Entering New for GamePayment")
	log.Log.Debug("Exiting  New for GamePayment")
	return GamePayment{
		id:               uuid.New().String(),
		configuration:    config,
		status:           Stopped,
		start:            time.Time{},
		previousDuration: time.Duration(0),
		check:            Check{},
		itemList:         make([]warehouse.Item, 0),
	}
}

// Set the confguration for the payment system
func (gp *GamePayment) ConfigurePayment(config config.PaymentConfiguration) {
	log.Log.Debug("Entering ConfigurePayment")
	gp.configuration = config
	log.Log.Debug("Entering ConfigurePayment")
}

// Start new payment by starting counting the time. It returns an error if the payment is neither in a stopped nor in suspended status
func (gp *GamePayment) StartCountingPayment() error {
	log.Log.Debug("Entering StartCountingPayment")
	// Check if payment already started, if it is the case return an error
	if gp.status == Started {
		err := errors.New("can not start new payment because its current status already is Started, but it must be either Stopped or Suspended")
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting StartCountingPayment")
		return err
	}
	// Start a new payment
	gp.start = time.Now()
	gp.status = Started
	gp.check = Check{}
	gp.itemList = make([]warehouse.Item, 0)

	log.Log.Info("New payment has been started")
	log.Log.Debug("Exiting StartCountingPayment")
	return nil
}

// Stop current payment and calculate the billing. It returns an error if the payment is neither in a started nor in suspended status
func (gp *GamePayment) ClosePayment() error {
	log.Log.Debug("Entering ClosePayment")
	// Check if payment already stopped, if it is the case return an error
	if gp.status == Stopped {
		err := errors.New("can not stop current payment because its current status already is Stopped, but it must be either Started or Paused")
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting ClosePayment")
		return err
	}
	duration := gp.previousDuration
	if gp.status == Started {
		duration += time.Since(gp.start)
	}
	if duration.Minutes() < float64(gp.configuration.MinimumDuration) {
		duration = time.Duration(float64(gp.configuration.MinimumDuration) * float64(time.Minute))
	}

	gp.check = Check{
		Duration: int(duration.Minutes()),
		Price:    int(duration.Minutes()) * gp.configuration.CostPerHour / 60,
		ItemList: make([]warehouse.Item, len(gp.itemList)),
	}
	copy(gp.check.ItemList, gp.itemList)
	for _, item := range gp.check.ItemList {
		gp.check.Price += item.PublicPrice
	}
	gp.status = Stopped

	log.Log.Info("Payment has been closed and related check has been created")
	log.Log.Debug("Exiting ClosePayment")
	return nil
}

// Pause current payment, to calculate billing it must be stopped. It returns an error if the payment is not in started status
func (gp *GamePayment) PauseCountingPayment() error {
	log.Log.Debug("Entering PauseCountingPayment")
	if gp.status != Started {
		err := errors.New("can not pause current payment because its current status is not Started")
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting PauseCountingPayment")
		return err
	}
	gp.previousDuration = time.Since(gp.start)
	gp.status = Suspended

	log.Log.Info("Payment has been paused")
	log.Log.Debug("Exiting ClosePayment")
	return nil
}

// Return the bill calculated for the payment. It returns an error if the payment has not been closed
func (gp *GamePayment) GetCheck() (Check, error) {
	log.Log.Debug("Entering GetCheck")
	if gp.status != Stopped {
		err := errors.New("can not return check for this mayment becaus it has not been closed. Please close it")
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting GetCheck")
		return Check{}, err
	}

	log.Log.Info("Check for this payment has been returned")
	log.Log.Debug("Exiting GetCheck")
	return gp.check, nil
}

// Return the payment status
func (gp *GamePayment) GetPaymentStatus() PaymentStatus {
	log.Log.Debug("Entering GetPaymentStatus")
	log.Log.Debug("Exiting GetPaymentStatus")
	return gp.status
}

// Add a consumprion to current payment. If payment is n stoppes status it reruns error
func (gp *GamePayment) AddConsumption(item warehouse.Item) error {
	log.Log.Debug("Entering AddConsumption")
	if gp.status == Stopped {
		err := errors.New("can not add new consumption to closed payment")
		log.Log.Error(err.Error())
		log.Log.Debug("Exiting AddConsumption")
		return err
	}
	gp.itemList = append(gp.itemList, item)
	log.Log.Info("New item has been consumed")
	log.Log.Debug("Exiting AddConsumption")
	return nil
}
