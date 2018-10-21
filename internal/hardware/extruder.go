package hardware

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type ExtruderPort interface {
	io.ReadWriter
	io.Closer
	Open(ctx context.Context) error
	IsOpen() bool
}

// Extruder setting
type Extruder struct {
	port ExtruderPort
	io   *bufio.ReadWriter
}

// ConnectExtruder is used to check extruder is alive or not.
func NewExtruder(port ExtruderPort) *Extruder {
	return &Extruder{
		port: port,
		io:   bufio.NewReadWriter(bufio.NewReader(port), bufio.NewWriter(port)),
	}
}

func (e *Extruder) Connect(ctx context.Context) error {
	log.Info().Msg("Connecting to Extruder")
	if err := e.port.Open(ctx); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	e.io.Flush()
	if e.Writeline("") == nil {
		line, err := e.Readline()
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		log.Debug().Msg(line)
		if strings.Compare(line, "ok") == 0 {
			log.Info().Msg("Connect to Extruder successful")
			return nil
		}
	}
	return errors.New("no response")
}

// Disconnect extruder
func (e *Extruder) Disconnect() error {
	if err := e.io.Flush(); err != nil {
		return err
	}
	if err := e.port.Close(); err != nil {
		return err
	}
	return nil
}

// Writeline is used to write a line to extruder
func (e *Extruder) Writeline(msg string) error {
	var buffer bytes.Buffer
	var err error

	if _, err = buffer.WriteString(msg); err != nil {
		return err
	}
	buffer.WriteByte('\r')
	buffer.WriteByte('\n')

	if _, err := buffer.WriteTo(e.io); err != nil {
		return err
	}
	e.io.Flush()
	return nil
}

// Readline is used to read a line from extruder
func (e *Extruder) Readline() (string, error) {
	line, isPrefix, err := e.io.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("Line is too long")
	}
	return string(line), err
}
