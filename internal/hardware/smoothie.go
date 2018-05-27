package hardware

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

var initCmds = [...]string{"G28", "G21", "G90", "M83"}

type SmoothiePort interface {
	io.ReadWriter
	io.Closer
	Open() error
	IsOpen() bool
}

// Smoothie setting
type Smoothie struct {
	port SmoothiePort
	io   *bufio.ReadWriter
}

// ConnectSmoothie is used to check extruder is alive or not.
func NewSmoothie(port SmoothiePort) *Smoothie {
	return &Smoothie{
		port: port,
		io:   bufio.NewReadWriter(bufio.NewReaderSize(port, 60), bufio.NewWriter(port)),
	}
}

func (s *Smoothie) Connect() error {
	if err := s.port.Open(); err != nil {
		return err
	}

	s.io.Flush()
	if s.Writeline("G") {
		line, err := s.Readline()
		if err != nil {
			return err
		}
		if strings.Compare(line, "Ok") == 0 {
			return nil
		}
	}

	for _, cmd := range initCmds {
		var line string
		var err error

		if s.Writeline(cmd) {
			return errors.New("initial failed")
		}
		if line, err = s.Readline(); err != nil {
			return err
		}
		if strings.Compare(line, "ok") != 0 {
			return errors.New("initial failed")
		}
	}

	return errors.New("no response")
}

// Disconnect extruder
func (s *Smoothie) Disconnect() error {
	s.io.Flush()
	if err := s.port.Close(); err != nil {
		return err
	}
	return nil
}

// Writeline is used to write a line to extruder
func (s *Smoothie) Writeline(msg string) bool {
	var buffer bytes.Buffer

	buffer.WriteString(msg)
	buffer.WriteByte('\n')

	if _, err := buffer.WriteTo(s.io); err != nil {
		return false
	}
	return true
}

// Readline is used to read a line from extruder
func (s *Smoothie) Readline() (string, error) {
	line, isPrefix, err := s.io.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("Line is too long")
	}
	return string(line), err
}
