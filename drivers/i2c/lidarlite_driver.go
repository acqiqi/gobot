package i2c

import (
	"gobot.io/x/gobot"

	"time"
)

const lidarliteAddress = 0x62

type LIDARLiteDriver struct {
	name       string
	connector  I2cConnector
	connection I2cConnection
	I2cConfig
}

// NewLIDARLiteDriver creates a new driver with specified i2c interface
// Params:
//		conn I2cConnector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.Bus(int):	bus to use with this driver
//		i2c.Address(int):	address to use with this driver
//
func NewLIDARLiteDriver(a I2cConnector, options ...func(I2cConfig)) *LIDARLiteDriver {
	l := &LIDARLiteDriver{
		name:      gobot.DefaultName("LIDARLite"),
		connector: a,
		I2cConfig: NewI2cConfig(),
	}

	for _, option := range options {
		option(l)
	}

	// TODO: add commands to API
	return l
}

func (h *LIDARLiteDriver) Name() string                 { return h.name }
func (h *LIDARLiteDriver) SetName(n string)             { h.name = n }
func (h *LIDARLiteDriver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initialized the LIDAR
func (h *LIDARLiteDriver) Start() (err error) {
	bus := h.GetBus(h.connector.I2cGetDefaultBus())
	address := h.GetAddress(lidarliteAddress)

	h.connection, err = h.connector.I2cGetConnection(address, bus)
	if err != nil {
		return err
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *LIDARLiteDriver) Halt() (err error) { return }

// Distance returns the current distance in cm
func (h *LIDARLiteDriver) Distance() (distance int, err error) {
	if _, err = h.connection.Write([]byte{0x00, 0x04}); err != nil {
		return
	}
	time.Sleep(20 * time.Millisecond)

	if _, err = h.connection.Write([]byte{0x0F}); err != nil {
		return
	}

	upper := []byte{0}
	bytesRead, err := h.connection.Read(upper)
	if err != nil {
		return
	}

	if bytesRead != 1 {
		err = ErrNotEnoughBytes
		return
	}

	if _, err = h.connection.Write([]byte{0x10}); err != nil {
		return
	}

	lower := []byte{0}
	bytesRead, err = h.connection.Read(lower)
	if err != nil {
		return
	}

	if bytesRead != 1 {
		err = ErrNotEnoughBytes
		return
	}

	distance = ((int(upper[0]) & 0xff) << 8) | (int(lower[0]) & 0xff)

	return
}
