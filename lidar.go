package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"

	"time"
)

var _ gobot.Driver = (*LIDARLiteDriver)(nil)

const lidarliteAddress = 0x62

type LIDARLiteDriver struct {
	name       string
	connection i2c.I2c
}

// NewLIDARLiteDriver creates a new driver with specified name and i2c interface
func NewLIDARLiteDriver(a i2c.I2c, name string) *LIDARLiteDriver {
	return &LIDARLiteDriver{
		name:       name,
		connection: a,
	}
}

func (h *LIDARLiteDriver) Name() string                 { return h.name }
func (h *LIDARLiteDriver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the LIDAR
func (h *LIDARLiteDriver) Start() (errs []error) {
	if err := h.connection.I2cStart(lidarliteAddress); err != nil {
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *LIDARLiteDriver) Halt() (errs []error) { return }

func (h *LIDARLiteDriver) readByteReg(address byte) (b byte, err error) {
	if err = h.connection.I2cWrite(lidarliteAddress, []byte{address}); err != nil {
		return
	}

	<-time.After(20 * time.Millisecond)
	ret, err := h.connection.I2cRead(lidarliteAddress, 1)
	if err != nil {
		return 0, err
	}
	if len(ret) != 1 {
		err = i2c.ErrNotEnoughBytes
	} else {
		b = ret[0]
	}

	return
}

// Distance returns the current distance
func (h *LIDARLiteDriver) Distance() (distance int, err error) {
	if err = h.connection.I2cWrite(lidarliteAddress, []byte{0x00, 0x04}); err != nil {
		return
	}
	<-time.After(20 * time.Millisecond)

	dist := []byte{}
	for _, b := range []byte{0x0f, 0x10} {
		r, err := h.readByteReg(b)
		if err != nil {
			return 0, err
		}
		dist = append(dist, r)
	}

	if len(dist) == 2 {
		distance = (int(dist[1]) + int(dist[0])*256)
	} else {
		err = i2c.ErrNotEnoughBytes
	}
	return
}
