package barista

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
)

type Controller interface {
	Connect() error
	Disconnect() error
	Do(p *lib.Point) error
}

type SEController struct {
	Smoothie *hardware.Smoothie
	Extruder *hardware.Extruder
}

func (se *SEController) Connect() (err error) {
	if err = se.Extruder.Connect(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			se.Extruder.Disconnect()
		}
	}()

	if err = se.Smoothie.Connect(); err != nil {
		return err
	}
	return nil
}

func (se *SEController) Disconnect() error {
	se.Extruder.Disconnect()
	se.Smoothie.Disconnect()
	return nil
}

func (se *SEController) Do(p *lib.Point) error {
	gcode, gerr := pointToGCode(p)
	hcode, herr := pointToHCode(p)

	if gerr == nil {
		se.Smoothie.Writeline(gcode)
	}
	if herr == nil {
		se.Extruder.Writeline(hcode)
	}

	if gerr == nil {
		resp, err := se.Smoothie.Readline()
		if strings.Compare(resp, "ok") != 0 {
			log.Error().Err(err)
		}
	}
	if herr == nil {
		resp, err := se.Extruder.Readline()
		if strings.Compare(resp, "ok") != 0 {
			log.Error().Err(err)
		}
	}

	return nil
}

func pointToGCode(p *lib.Point) (string, error) {
	if p.X == nil && p.Y == nil && p.Z == nil {
		return "", errors.New("no x, y, and z")
	}

	var buffer bytes.Buffer
	buffer.WriteString("G1")
	if p.X != nil {
		buffer.WriteString(fmt.Sprintf(" X%0.5f", *p.X))
	}
	if p.Y != nil {
		buffer.WriteString(fmt.Sprintf(" Y%0.5f", *p.Y))
	}
	if p.Z != nil {
		buffer.WriteString(fmt.Sprintf(" Z%0.5f", *p.Z))
	}
	buffer.WriteString(fmt.Sprintf(" F%0.5f", *p.F))
	return buffer.String(), nil
}

func pointToHCode(p *lib.Point) (string, error) {
	if p.Time == nil {
		return "", errors.New("no time")
	}

	var buffer bytes.Buffer
	buffer.WriteString("H")
	if p.E1 != nil && *p.E1 != 0 {
		buffer.WriteString(fmt.Sprintf(" E0 %05f", *p.E1))
	}
	if p.E2 != nil && *p.E2 != 0 {
		buffer.WriteString(fmt.Sprintf(" E1 %05f", *p.E2))
	}
	buffer.WriteString(fmt.Sprintf(" T%0.5f", *p.Time))
	return buffer.String(), nil
}
