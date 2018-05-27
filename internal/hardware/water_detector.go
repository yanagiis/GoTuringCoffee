package hardware

import (
	"errors"
	"fmt"

	"github.com/yanagiis/periph/conn/gpio"
	"github.com/yanagiis/periph/conn/gpio/gpioreg"
)

type WaterDetector interface {
	Connect() error
	Disconnect() error
	IsWaterFull() (bool, error)
}

type GPIOWaterDetector struct {
	Pin  uint32
	gpio gpio.PinIO
}

func (w *GPIOWaterDetector) Connect() error {
	gpioName := fmt.Sprintf("GPIO%d", w.Pin)
	w.gpio = gpioreg.ByName(gpioName)
	if w.gpio == nil {
		return fmt.Errorf("Cannot open %s", gpioName)
	}
	if err := w.gpio.In(gpio.PullNoChange, gpio.NoEdge); err != nil {
		return err
	}
	return nil
}

func (w *GPIOWaterDetector) Disconnect() error {
	w.gpio = nil
	return nil
}

func (w *GPIOWaterDetector) IsWaterFull() (bool, error) {
	if w.gpio == nil {
		return false, errors.New("gpio not connected")
	}
	return w.gpio.Read() == gpio.High, nil
}
