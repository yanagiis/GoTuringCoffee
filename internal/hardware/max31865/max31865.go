//go:generate go-enum -fmax31865.go

package max31865

import (
	"errors"
	"fmt"

	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
)

// RTD wire type
// ENUM(
// 2, 3, 4
// )
type Wire byte

// MAX31865 Mode
// ENUM(
// Manual, Automatic
// )
type Mode byte

const (
	addrWriteMAX31865 byte = 0x80
	addrCR                 = 0x0
	addrRTDH               = 0x1
	addrRTDL               = 0x2
	addrHFTH               = 0x3
	addrHFTL               = 0x4
	addrLFTH               = 0x5
	addrLFTL               = 0x6
	addrFault              = 0x7
)

// Config is used to setting max31865 sensor
type Config struct {
	Wire Wire `mapstructure:"wire"`
	Mode Mode `mapstructure:"mode"`
}

// MAX31865 RTD sensor
type MAX31865 struct {
	spi       spiwrap.SPI
	conf      Config
	connected bool
}

func New(spi spiwrap.SPI, conf Config) *MAX31865 {
	return &MAX31865{
		spi:       spi,
		conf:      conf,
		connected: false,
	}
}

func (m *MAX31865) Connect() error {
	var err error
	var mode Mode
	var wire Wire

	if m.connected {
		return nil
	}

	if err = m.spi.Open(); err != nil {
		return err
	}

	err = m.setMode(m.conf.Mode)
	if err != nil {
		return err
	}
	err = m.setWire(m.conf.Wire)
	if err != nil {
		return err
	}

	mode, err = m.getMode()
	if err != nil {
		return err
	}
	if mode != m.conf.Mode {
		return errors.New("max31865 set mode failed")
	}

	wire, err = m.getWire()
	if err != nil {
		return err
	}
	if wire != m.conf.Wire {
		return errors.New("max31865 set wire failed")
	}

	buf := make([]byte, 1)
	err = m.readReg(addrCR, buf)
	if err != nil {
		return err
	}

	buf[0] = (buf[0] & 0x7f) | (1 << 7)
	err = m.writeReg(addrCR, buf)
	if err != nil {
		return err
	}

	m.connected = true
	return nil
}

func (m *MAX31865) Disconnect() error {
	if err := m.spi.Close(); err != nil {
		return err
	}
	m.connected = false
	return nil
}

// GetTemperature get temperature
//
// Callendar-Van Dusen equation:
// R(T) = R0(1 + aT + bT2 + c(T - 100)T3)
// where:
//     T = temperature (C)
//     R(T) = resistance at T
//     R0 = resistance at T = 0C
//     IEC 751 specifies α = 0.00385055 and the following
//     Callendar-Van Dusen coefficient values:
//         a = 3.90830 x 10^-3
//         b = -5.77500 x 10^-7
//         c = -4.18301 x 10^-12 for -200C < T < 0C, 0 for 0C < T < +850C
//
// Linearizing Temperature Data
// For a temperature range of -100C to +100C, a good
// approximation of temperature can be made by simply
// using the RTD data as shown below:
// Temperature (C) ≈ (ADC code/32) – 256
func (m *MAX31865) GetTemperature() (float64, error) {

	if !m.connected {
		return 0, errors.New("sensor is not connected yet")
	}

	buf := make([]byte, 2)
	err := m.readReg(addrRTDH, buf)
	if err != nil {
		return 0, err
	}

	if (buf[1] & 0x1) != 0 {
		var value byte
		value, err = m.getError()
		if err != nil {
			return 0, err
		}

		err = m.clearError()
		if err != nil {
			return 0, err
		}

		return 0, fmt.Errorf("max31865 error %02x", value)
	}

	adcValue := ((int32(buf[0]) << 8) | (int32(buf[1]))) >> 1
	if adcValue == 0 {
		return 0, errors.New("get zero adcValue")
	}

	return float64(adcValue)/32 - 256, nil
}

func (m *MAX31865) getMode() (Mode, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrCR, buf)
	return Mode((buf[0] & 0x40) >> 6), err
}

func (m *MAX31865) setMode(mode Mode) error {
	buf := make([]byte, 1)
	err := m.readReg(addrCR, buf)
	if err != nil {
		return err
	}
	buf[0] = (buf[0] & 0xbf) | (byte(mode) << 6)
	return m.writeReg(addrCR, buf)
}

func (m *MAX31865) getWire() (Wire, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrCR, buf)
	return Wire((buf[0] & 0x10) >> 4), err
}

func (m *MAX31865) setWire(wire Wire) error {
	buf := make([]byte, 1)
	err := m.readReg(addrCR, buf)
	if err != nil {
		return err
	}
	buf[0] = (buf[0] & 0xef) | (byte(wire) << 4)
	return m.writeReg(addrCR, buf)
}

func (m *MAX31865) getError() (byte, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrFault, buf)
	return buf[0], err
}

func (m *MAX31865) clearError() error {
	buf := make([]byte, 1)
	err := m.readReg(addrCR, buf)
	if err != nil {
		return err
	}
	buf[0] |= 0x2
	return m.writeReg(addrCR, buf)
}

func (m *MAX31865) readReg(addr byte, buf []byte) error {
	wbuf := append([]byte{addr}, buf...)
	rbuf := make([]byte, len(wbuf))
	if err := m.spi.Tx(wbuf, rbuf); err != nil {
		return err
	}
	copy(buf, rbuf[1:])
	return nil
}

func (m *MAX31865) writeReg(addr byte, buf []byte) error {
	wbuf := append([]byte{addr | addrWriteMAX31865}, buf...)
	rbuf := make([]byte, len(wbuf))
	return m.spi.Tx(wbuf, rbuf)
}
