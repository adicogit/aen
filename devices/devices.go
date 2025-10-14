package devices

// Represents the device status
type DeviceStatus int

// Enumeration for the payment status
const (
	On DeviceStatus = iota
	Off
	Unresponsive
)

// Represents the device type
type DeviceType int

// Enumeration for the payment status
const (
	Bulbs DeviceType = iota
	Camera
	BallBox
)

// Devices interface
type Device interface {
	// Return device type
	GetType() DeviceType
	// Return device status
	GetStatus() DeviceStatus
	// Turn ON the devices
	TurnOn() error
	// Turn OFF the devices
	TurnOff() error
}

//webcam protocol: rtsp
//https://github.com/bluenviron/gortsplib
