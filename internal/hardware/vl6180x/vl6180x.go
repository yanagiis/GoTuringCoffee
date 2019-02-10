package vl6180x

import (
	"time"

	"github.com/rs/zerolog/log"

	"GoTuringCoffee/internal/hardware/i2cwrap"
)

const (
	// The fixed I2C addres
	VL6180X_DEFAULT_I2C_ADDR = 0x29

	///! Device model identification number
	VL6180X_REG_IDENTIFICATION_MODEL_ID = 0x000
	///! Interrupt configuration
	VL6180X_REG_SYSTEM_INTERRUPT_CONFIG = 0x014
	///! Interrupt clear bits
	VL6180X_REG_SYSTEM_INTERRUPT_CLEAR = 0x015
	///! Fresh out of reset bit
	VL6180X_REG_SYSTEM_FRESH_OUT_OF_RESET = 0x016
	///! Trigger Ranging
	VL6180X_REG_SYSRANGE_START = 0x018
	///! Trigger Lux Reading
	VL6180X_REG_SYSALS_START = 0x038
	///! Lux reading gain
	VL6180X_REG_SYSALS_ANALOGUE_GAIN = 0x03F
	///! Integration period for ALS mode, high byte
	VL6180X_REG_SYSALS_INTEGRATION_PERIOD_HI = 0x040
	///! Integration period for ALS mode, low byte
	VL6180X_REG_SYSALS_INTEGRATION_PERIOD_LO = 0x041
	///! Specific error codes
	VL6180X_REG_RESULT_RANGE_STATUS = 0x04d
	///! Interrupt status
	VL6180X_REG_RESULT_INTERRUPT_STATUS_GPIO = 0x04f
	///! Light reading value
	VL6180X_REG_RESULT_ALS_VAL = 0x050
	///! Ranging reading value
	VL6180X_REG_RESULT_RANGE_VAL = 0x062

	VL6180X_ALS_GAIN_1    = 0x06 ///< 1x gain
	VL6180X_ALS_GAIN_1_25 = 0x05 ///< 1.25x gain
	VL6180X_ALS_GAIN_1_67 = 0x04 ///< 1.67x gain
	VL6180X_ALS_GAIN_2_5  = 0x03 ///< 2.5x gain
	VL6180X_ALS_GAIN_5    = 0x02 ///< 5x gain
	VL6180X_ALS_GAIN_10   = 0x01 ///< 1=0x gain
	VL6180X_ALS_GAIN_20   = 0x00 ///< 2=0x gain
	VL6180X_ALS_GAIN_40   = 0x07 ///< 4=0x gain

	VL6180X_ERROR_NONE        = 0  ///< Success!
	VL6180X_ERROR_SYSERR_1    = 1  ///< System error
	VL6180X_ERROR_SYSERR_5    = 5  ///< Sysem error
	VL6180X_ERROR_ECEFAIL     = 6  ///< Early convergence estimate fail
	VL6180X_ERROR_NOCONVERGE  = 7  ///< No target detected
	VL6180X_ERROR_RANGEIGNORE = 8  ///< Ignore threshold check failed
	VL6180X_ERROR_SNR         = 11 ///< Ambient conditions too high
	VL6180X_ERROR_RAWUFLOW    = 12 ///< Raw range algo underflow
	VL6180X_ERROR_RAWOFLOW    = 13 ///< Raw range algo overflow
	VL6180X_ERROR_RANGEUFLOW  = 14 ///< Raw range algo underflow
	VL6180X_ERROR_RANGEOFLOW  = 15 ///< Raw range algo overflow
)

type Vl6180x struct {
	i2cAddr   int
	i2cDevice *i2cwrap.I2C
	ioTimeout time.Duration
	scaling   int

	opened bool
}

func New(i2cDevice *i2cwrap.I2C, i2cAddr int, scaling int) (*Vl6180x, error) {
	v := &Vl6180x{
		i2cAddr:   i2cAddr,
		i2cDevice: i2cDevice,
		scaling:   scaling,
		opened:    false,
	}
	return v, nil
}

func (v *Vl6180x) Open() error {
	if v.opened {
		return nil
	}

	err := v.i2cDevice.Open(v.i2cAddr)
	if err != nil {
		return err
	}

	v.LoadSettings()
	v.opened = true
	return nil
}

func (v *Vl6180x) ReadByte(reg uint16) byte {
	data, _, err := v.i2cDevice.ReadRegBytes(reg, 1)
	if err != nil {
		log.Error().Msgf("Can't read byte from %x", reg)
		log.Fatal().Err(err)
		return 0x00
	}

	return data[0]
}

func (v *Vl6180x) ReadBytes(reg uint16, n int) []byte {
	buf, _, err := v.i2cDevice.ReadRegBytes(reg, n)
	if err != nil {
		log.Error().Msgf("Can't read reg %x", reg)
		return []byte{}
	}

	return buf
}

func (v *Vl6180x) WriteU8(reg uint16, value uint8) error {
	return v.i2cDevice.WriteRegU8(reg, value)
}

func (v *Vl6180x) WriteU16(reg uint16, value uint16) error {
	err := v.i2cDevice.WriteRegU16(reg, value)
	if err != nil {
		return err
	}
	//log.Debug().Msgf("Write %x to 0x%0X", value, reg)

	return nil
}

func (v *Vl6180x) LoadSettings() {

	setup := v.ReadByte(VL6180X_REG_SYSTEM_FRESH_OUT_OF_RESET)

	if setup == 1 {
		log.Debug().Msg("Loading vl6180x settings")
		// private settings from page 24 of app note
		v.WriteU8(0x0207, 0x01)
		v.WriteU8(0x0208, 0x01)
		v.WriteU8(0x0096, 0x00)
		v.WriteU8(0x0097, 0xfd)
		v.WriteU8(0x00e3, 0x00)
		v.WriteU8(0x00e4, 0x04)
		v.WriteU8(0x00e5, 0x02)
		v.WriteU8(0x00e6, 0x01)
		v.WriteU8(0x00e7, 0x03)
		v.WriteU8(0x00f5, 0x02)
		v.WriteU8(0x00d9, 0x05)
		v.WriteU8(0x00db, 0xce)
		v.WriteU8(0x00dc, 0x03)
		v.WriteU8(0x00dd, 0xf8)
		v.WriteU8(0x009f, 0x00)
		v.WriteU8(0x00a3, 0x3c)
		v.WriteU8(0x00b7, 0x00)
		v.WriteU8(0x00bb, 0x3c)
		v.WriteU8(0x00b2, 0x09)
		v.WriteU8(0x00ca, 0x09)
		v.WriteU8(0x0198, 0x01)
		v.WriteU8(0x01b0, 0x17)
		v.WriteU8(0x01ad, 0x00)
		v.WriteU8(0x00ff, 0x05)
		v.WriteU8(0x0100, 0x05)
		v.WriteU8(0x0199, 0x05)
		v.WriteU8(0x01a6, 0x1b)
		v.WriteU8(0x01ac, 0x3e)
		v.WriteU8(0x01a7, 0x1f)
		v.WriteU8(0x0030, 0x00)

		// Recommended : Public registers - See data sheet for more detail
		v.WriteU8(0x0011, 0x10) // Enables polling for 'New Sample ready'
		// when measurement completes
		v.WriteU8(0x010a, 0x30) // Set the averaging sample period
		// (compromise between lower noise and
		// increased execution time)
		v.WriteU8(0x003f, 0x46) // Sets the light and dark gain (upper
		// nibble). Dark gain should not be
		// changed.
		v.WriteU8(0x0031, 0xFF) // sets the # of range measurements after
		// which auto calibration of system is
		// performed
		v.WriteU8(0x0040, 0x63) // Set ALS integration time to 100ms
		v.WriteU8(0x002e, 0x01) // perform a single temperature calibration
		// of the ranging sensor

		// Optional: Public registers - See data sheet for more detail
		v.WriteU8(0x001b, 0x09) // Set default ranging inter-measurement
		// period to 100ms
		v.WriteU8(0x003e, 0x31) // Set default ALS inter-measurement period
		// to 500ms
		v.WriteU8(0x0014, 0x24) //

		v.WriteU8(VL6180X_REG_SYSTEM_FRESH_OUT_OF_RESET, 0x00)
	}

	v.SetScaling(v.scaling)
	return
}

func (v *Vl6180x) StartRange() {
	// VL6180X_REG_SYSRANGE_START
	log.Debug().Msg("Start Range")
	v.WriteU8(0x018, 0x01)
}

func (v *Vl6180x) PollRange() {
	// wait for new measurement ready status
	/*
		for v.ReadRangeStatus() != 0x04 {
			log.Debug("Wait until the range status is 0x04 now is %x", v.ReadRangeStatus())
		}
	*/

	log.Debug().Msg("Polling range status")
	var status byte
	var rangeStatus byte

	for rangeStatus != 0x04 {
		status = v.ReadByte(0x04f)
		rangeStatus = status & 0x07
	}
}

func (v *Vl6180x) ClearInterrupt() {
	// VL6180X_REG_SYSTEM_INTERRUPT_CLEAR
	log.Debug().Msg("Clear interrupt")
	v.WriteU8(0x015, 0x07)
}

func (v *Vl6180x) SetScaling(newScaling int) {
	scalerValues := []uint16{0, 253, 127, 84}
	defaultCrosstalkValidHeight := 20
	if newScaling < 1 || newScaling > 3 {
		log.Warn().Msgf("Can't set scaling to %d", newScaling)
		return
	}

	ptpOffset := v.ReadByte(0x24)

	v.WriteU16(0x96, scalerValues[newScaling])
	v.WriteU8(0x24, byte(int(ptpOffset)/newScaling))
	v.WriteU8(0x21, byte(defaultCrosstalkValidHeight/newScaling))
	rce := v.ReadByte(0x2d)
	if newScaling == 1 {
		v.WriteU8(0x23, byte((rce&0xFE)|1))
	} else {
		v.WriteU8(0x23, byte(rce&0xFE))
	}
}

func (v *Vl6180x) ReadRange() uint8 {
	v.StartRange()
	v.PollRange()
	result := v.ReadByte(0x063)

	v.ClearInterrupt()

	return uint8(result)
}

func (v *Vl6180x) Close() error {
	v.i2cDevice.Close()
	v.opened = false
	return nil
}

func (v *Vl6180x) setTimeout(timeout time.Duration) {
	v.ioTimeout = timeout
}

func (v *Vl6180x) Init(i2c *i2cwrap.I2C) error {
	v.setTimeout(time.Millisecond * 1000)

	return nil
}
