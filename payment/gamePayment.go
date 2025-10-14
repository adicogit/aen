package payment

import (
	"errors"
	"time"
)

type GamePayment struct {
	configuraiton    PaymentConfiguration
	status           PaymentStatus
	start            time.Time
	previousDuration time.Duration
	check            Check
}

// New function initialize GamePayment passing a PaymentConfiguration
func New(config PaymentConfiguration) GamePayment {
	return GamePayment{
		configuraiton:    config,
		status:           Stopped,
		start:            time.Time{},
		previousDuration: time.Duration(0),
		check:            Check{},
	}
}

// Set the confguration for the payment system
func (gp *GamePayment) ConfigurePayment(config PaymentConfiguration) {
	gp.configuraiton = config
}

// Start new payment by starting counting the time. It returns an error if the payment is neither in a stopped nor in suspended status
func (gp *GamePayment) StartCountingPayment() error {
	// Check if payment already started, if it is the case return an error
	if gp.status == Started {
		err := errors.New("can not start new payment because its current status already is Started, but it must be either Stopped or Suspended")
		return err
	}
	// Start a new payment
	gp.start = time.Now()
	gp.status = Started
	gp.check = Check{}

	return nil
}

// Stop current payment and calculate the billing. It returns an error if the payment is neither in a started nor in suspended status
func (gp *GamePayment) ClosePayment() error {
	// Check if payment already stopped, if it is the case return an error
	if gp.status == Stopped {
		err := errors.New("can not stop current payment because its current status already is Stopped, but it must be either Started or Paused")
		return err
	}
	duration := gp.previousDuration
	if gp.status == Started {
		duration += time.Since(gp.start)
	}
	if duration.Minutes() < float64(gp.configuraiton.MinimumDuration) {
		duration = time.Duration(float64(gp.configuraiton.MinimumDuration) * float64(time.Minute))
	}
	gp.check = Check{
		Duration: int(duration.Minutes()),
		Price:    float32(duration.Minutes()) * float32(gp.configuraiton.CostPerHour) / 60,
	}
	gp.status = Stopped

	return nil
}

// Pause current payment, to calculate billing it must be stopped. It returns an error if the payment is not in started status
func (gp *GamePayment) PauseCountingPayment() error {
	if gp.status != Started {
		err := errors.New("can not pause current payment because its current status is not Started")
		return err
	}
	gp.previousDuration = time.Since(gp.start)
	gp.status = Suspended

	return nil
}

// Return the bill calculated for the payment. It returns an error if the payment has not been closed
func (gp *GamePayment) GetCheck() (Check, error) {
	if gp.status != Stopped {
		err := errors.New("can not return check for this mayment becaus it has not been closed. Please close it")
		return Check{}, err
	}
	return gp.check, nil
}

// Return the payment status
func (gp *GamePayment) GetPaymentStatus() PaymentStatus {
	return gp.status
}
