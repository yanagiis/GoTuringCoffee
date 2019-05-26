//go:generate go-enum -fspi.go

package spiwrap

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/brian-armstrong/gpio"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

type SPI interface {
	io.Closer
	Open() error
	Tx(w, r []byte) error
}

type SPIPins struct {
	MISO uint `mapstructure:"miso"`
	MOSI uint `mapstructure:"mosi"`
	CS   uint `mapstructure:"cs"`
	CLK  uint `mapstructure:"clk"`
}

type Config struct {
	Path  string   `mapstructure:"path"`
	Speed int64    `mapstructure:"speed"`
	Mode  spi.Mode `mapstructure:"mode"`
	Bits  int      `mapstructure:"bits"`
}

type SPIDevice struct {
	Conf   Config
	closer spi.PortCloser
	conn   spi.Conn
}

type SPIGPIO struct {
	Conf   Config
	Pins   SPIPins
	opened bool
	miso   gpio.Pin
	mosi   gpio.Pin
	cs     gpio.Pin
	clk    gpio.Pin
}

func (s *SPIDevice) Open() error {
	var err error
	if s.closer == nil {
		s.closer, err = spireg.Open(s.Conf.Path)
		if err != nil {
			return err
		}
	}
	if s.conn == nil {
		if s.conn, err = s.closer.Connect(physic.Frequency(s.Conf.Speed)*physic.Hertz, s.Conf.Mode, s.Conf.Bits); err != nil {
			return err
		}
	}
	return nil
}

func (s *SPIDevice) IsOpen() bool {
	return s.conn != nil
}

func (s *SPIDevice) Close() error {
	err := s.closer.Close()
	if err != nil {
		return err
	}
	s.conn = nil
	s.closer = nil
	return nil
}

func (s *SPIDevice) Tx(w, r []byte) error {
	if s.conn == nil {
		fmt.Errorf("Not open")
		return errors.New("Not open")
	}
	return s.conn.Tx(w, r)
}

func (s *SPIGPIO) Open() error {
	s.mosi = gpio.NewOutput(s.Pins.MOSI, false)
	s.clk = gpio.NewOutput(s.Pins.CLK, false)
	s.cs = gpio.NewOutput(s.Pins.CS, true)
	s.miso = gpio.NewInput(s.Pins.MISO)
	s.opened = true
	return nil
}

func (s *SPIGPIO) IsOpen() bool {
	return s.opened
}

func (s *SPIGPIO) Close() error {
	s.mosi.Close()
	s.clk.Close()
	s.cs.Close()
	s.miso.Close()
	s.opened = false
	return nil
}

func (s *SPIGPIO) Tx(w, r []byte) error {
	s.cs.Low()
	for numByte, wb := range w {
		rb := byte(0x0)
		for mask := byte(0x80); mask != 0; mask >>= 1 {
			bit := wb & mask
			s.clk.High()

			if bit != 0 {
				s.mosi.High()
			} else {
				s.mosi.Low()
			}

			val, _ := s.miso.Read()
			if val != 0 {
				rb |= mask
			}

			time.Sleep(1 * time.Millisecond)
			s.clk.Low()
		}
		r[numByte] = rb
	}
	s.cs.High()
	return nil
}
