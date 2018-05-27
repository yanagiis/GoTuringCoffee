//go:generate go-enum -fspi.go

package spiwrap

import (
	"errors"
	"fmt"
	"io"

	"github.com/yanagiis/periph/conn/spi"
	"github.com/yanagiis/periph/conn/spi/spireg"
)

type SPI interface {
	io.Closer
	Open() error
	Tx(w, r []byte) error
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

func (s *SPIDevice) Open() error {
	var err error
	if s.closer == nil {
		s.closer, err = spireg.Open(s.Conf.Path)
		if err != nil {
			return err
		}
	}
	if s.conn == nil {
		if s.conn, err = s.closer.Connect(s.Conf.Speed, s.Conf.Mode, s.Conf.Bits); err != nil {
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
