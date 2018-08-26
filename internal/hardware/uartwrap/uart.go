package uartwrap

import (
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tarm/serial"
)

type UART interface {
	io.ReadWriter
	io.Closer
	Open() error
	IsOpen() bool
}

type Config struct {
	Path        string        `mapstructure:"path"`
	Baudrate    uint32        `mapstructure:"baudrate"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}

type UARTDevice struct {
	Conf Config
	uart *serial.Port
}

func (u *UARTDevice) Open() error {
	var err error

	config := &serial.Config{
		Name:        u.Conf.Path,
		Baud:        int(u.Conf.Baudrate),
		ReadTimeout: time.Second * u.Conf.ReadTimeout,
	}
	u.uart, err = serial.OpenPort(config)

	log.Debug().Msgf("%v\n", config)

	return err
}

func (u *UARTDevice) IsOpen() bool {
	return u.uart != nil
}

func (u *UARTDevice) Read(p []byte) (int, error) {
	n, err := u.uart.Read(p)
	return n, err
}

func (u *UARTDevice) Write(p []byte) (int, error) {
	return u.uart.Write(p)
}

func (u *UARTDevice) Close() error {
	if err := u.uart.Close(); err != nil {
		return err
	}
	u.uart = nil
	return nil
}
