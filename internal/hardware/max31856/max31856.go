//go:generate go-enum -fmax31856.go

package max31856

import (
	"errors"

	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
)

// TC Type
// ENUM(
// B, E, J, K, N, R, S, T
// )
type Type byte

// MAX31856 mode
// ENUM(
// Manual, Automatic
// )
type Mode byte

// Sample average
// ENUM(
// Avg1, Avg2, Avg4, Avg8, Avg16
// )
type Sample byte

const (
	addrWriteMAX31856 byte = 0x80
	addrCR0                = 0x0
	addrCR1                = 0x1
	addrMASK               = 0x2
	addrCJHF               = 0x3
	addrCJLF               = 0x4
	addrLTHFTH             = 0x5
	addrLTHFTL             = 0x6
	addrLTLFTH             = 0x7
	addrLTLFTL             = 0x8
	addrCJTO               = 0x9
	addrCJTH               = 0xA
	addrCJTL               = 0xB
	addrLTCBH              = 0xC
	addrLTCBM              = 0xD
	addrLTCBL              = 0xE
	addrSR                 = 0xF
)

const (
	resolutionTC = 0.0078125
	resolutionCJ = 0.015625
)

// MAX31856 TC sensor
type MAX31856 struct {
	spi       spiwrap.SPI
	conf      Config
	connected bool
}

// Config is used to setting max31856 sensor
type Config struct {
	TC   Type
	Avg  Sample
	Mode Mode
}

func New(spi spiwrap.SPI, conf Config) *MAX31856 {
	return &MAX31856{
		spi:       spi,
		conf:      conf,
		connected: false,
	}
}

func (m *MAX31856) Connect() error {
	var err error
	var mode Mode
	var t Type
	var sample Sample

	if m.connected {
		return nil
	}

	if err = m.spi.Open(); err != nil {
		return err
	}

	if err = m.setMode(m.conf.Mode); err != nil {
		return err
	}

	if err = m.setSampleAvg(m.conf.Avg); err != nil {
		return err
	}

	if err = m.setTCType(m.conf.TC); err != nil {
		return err
	}

	mode, err = m.getMode()
	if err != nil {
		return err
	}
	if mode != m.conf.Mode {
		return errors.New("max31856 set mode failed")
	}

	t, err = m.getTCType()
	if err != nil {
		return err
	}
	if t != m.conf.TC {
		return errors.New("max31856 set tc-type failed")
	}

	sample, err = m.getSampleAvg()
	if err != nil {
		return err
	}
	if sample != m.conf.Avg {
		return errors.New("max31856 set sample avg failed")
	}

	m.connected = true
	return nil
}

func (m *MAX31856) Disconnect() error {
	if err := m.spi.Close(); err != nil {
		return err
	}
	m.connected = false
	return nil
}

// GetTemperature from max31856
func (m *MAX31856) GetTemperature() (float64, error) {
	if !m.connected {
		return 0, errors.New("sensor is connected yet")
	}

	buf := make([]byte, 4) // [t2, t1, t0, fault]
	if err := m.readReg(addrLTCBH, buf); err != nil {
		return 0, err
	}

	adcValue := ((int32(buf[0]) << 16) | (int32(buf[1]) << 8) | int32(buf[2])) >> 5

	if (buf[0] & 0x80) != 0 {
		adcValue -= 0x80000
	}

	if adcValue == 0 {
		return 0, errors.New("get zero ADC Value")
	}

	return float64(adcValue) * resolutionTC, nil
}

func (m *MAX31856) getMode() (Mode, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrCR0, buf)
	return Mode(buf[0] >> 7), err
}

func (m *MAX31856) setMode(mode Mode) error {
	buf := make([]byte, 1)
	if err := m.readReg(addrCR0, buf); err != nil {
		return err
	}
	buf[0] = (buf[0] & 0x7f) | (byte(mode) << 7)
	return m.writeReg(addrCR0, buf)
}

func (m *MAX31856) getSampleAvg() (Sample, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrCR1, buf)
	return Sample((buf[0] % 0x70) >> 4), err
}

func (m *MAX31856) setSampleAvg(avg Sample) error {
	buf := make([]byte, 1)
	if err := m.readReg(addrCR0, buf); err != nil {
		return err
	}
	buf[0] = (buf[0] & 0x8f) | (byte(avg) << 4)
	return m.writeReg(addrCR0, buf)
}

func (m *MAX31856) getTCType() (Type, error) {
	buf := make([]byte, 1)
	err := m.readReg(addrCR1, buf)
	return Type((buf[0] & 0x0f)), err
}

func (m *MAX31856) setTCType(t Type) error {
	buf := make([]byte, 1)
	if err := m.readReg(addrCR1, buf); err != nil {
		return err
	}
	buf[0] = (buf[0] & 0xf0) | byte(t)
	return m.writeReg(addrCR1, buf)
}

func (m *MAX31856) readReg(addr byte, buf []byte) error {
	wbuf := append([]byte{addr}, buf...)
	rbuf := make([]byte, len(wbuf))
	if err := m.spi.Tx(wbuf, rbuf); err != nil {
		return err
	}
	copy(buf, rbuf[1:])
	return nil
}

func (m *MAX31856) writeReg(addr byte, buf []byte) error {
	wbuf := append([]byte{addr | addrWriteMAX31856}, buf...)
	rbuf := make([]byte, len(wbuf))
	return m.spi.Tx(wbuf, rbuf)
}
